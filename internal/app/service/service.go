package service

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/Angelos-Giannis/er-builder/internal/domain"
)

// Service describes the service flow.
type Service struct {
	options domain.Options
}

type column struct {
	fieldName interface{}
	fieldType interface{}
	// fieldLabel interface{}
}

// New creates and returns a new service.
func New(options domain.Options) *Service {
	return &Service{
		options: options,
	}
}

// Generate performs the action to generate the .er file based on the provided input.
func (s *Service) Generate() error {
	if s.options.Directory == "" && len(s.options.FileList.Value()) == 0 {
		return errors.New("Need to provide at least one of 'directory' or 'file_list'")
	}

	filesToParse := defineFilesToParse(s.options.Directory, s.options.FileList.Value())
	tagRegexp := getTagRegexp(s.options.Tag)

	tableMapping := make(map[string][]column)

	for _, fl := range filesToParse {
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, fl, nil, parser.ParseComments)
		if err != nil {
			return err
		}

		for i := 0; i < len(node.Decls); i++ {
			if reflect.TypeOf(node.Decls[i]) != reflect.TypeOf(&ast.GenDecl{}) {
				continue
			}
			typeDecl := node.Decls[i].(*ast.GenDecl)

			for j := 0; j < len(typeDecl.Specs); j++ {
				if reflect.TypeOf(typeDecl.Specs[j]) != reflect.TypeOf(&ast.TypeSpec{}) {
					continue
				}

				if reflect.TypeOf(typeDecl.Specs[j].(*ast.TypeSpec).Type) != reflect.TypeOf(&ast.StructType{}) {
					continue
				}

				structDecl := typeDecl.Specs[j].(*ast.TypeSpec).Type.(*ast.StructType)
				structName := fmt.Sprintf("%v", typeDecl.Specs[j].(*ast.TypeSpec).Name)

				tableMapping[structName] = getTagFieldsFromStruct(tagRegexp, structDecl.Fields.List)
			}
		}
	}

	err := s.writeToFile(tableMapping)
	if err != nil {
		return err
	}

	return nil

}

// writeToFile creates the .er file and writes all the details in there.
func (s *Service) writeToFile(tableMapping map[string][]column) error {
	outputFile, err := os.Create(fmt.Sprintf("%v/%v.er", s.options.OutputPath, s.options.OutputFilename))
	if err != nil {
		return err
	}
	defer outputFile.Close()

	if s.options.Title != "" {
		_, err = outputFile.WriteString(fmt.Sprintf("title {label: \"%v\"}\n", s.options.Title))
		if err != nil {
			return err
		}
	}

	var foreignKeyConnections []string

	keys := make([]string, 0, len(tableMapping))
	for k := range tableMapping {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	_, err = outputFile.WriteString("# Definition of tables.\n")
	if err != nil {
		return err
	}
	for _, tb := range keys {
		checkingTable := strings.ToLower(fmt.Sprintf("%v", tb))
		if len(tableMapping[tb]) == 0 {
			continue
		}

		_, err = outputFile.WriteString(fmt.Sprintf("[%v]\n", tb))
		if err != nil {
			return err
		}

		if s.options.IDField != "" {
			_, err = outputFile.WriteString(fmt.Sprintf("\t*%v\n", s.options.IDField))
			if err != nil {
				return err
			}
		}

		for _, col := range tableMapping[tb] {
			fkPrefix := ""
			for intTB := range tableMapping {
				fld := strings.ToLower(fmt.Sprintf("%v", col.fieldName))
				currentTB := strings.ToLower(fmt.Sprintf("%v", intTB))

				if strings.Contains(fld, currentTB) {
					foreignKeyConnections = append(foreignKeyConnections, fmt.Sprintf("%v *--* %v {label: \"%v\"}", checkingTable, currentTB, fld))
					fkPrefix = "+"
				}
			}

			_, err = outputFile.WriteString(fmt.Sprintf("\t%v%v {label: \"%v\"}\n", fkPrefix, col.fieldName, col.fieldType))
			if err != nil {
				return err
			}
		}

		for _, c := range s.options.CommonFields.Value() {
			_, err = outputFile.WriteString(fmt.Sprintf("\t%v%v\n", "", c))
			if err != nil {
				return err
			}
		}
		_, err = outputFile.WriteString("\n")
		if err != nil {
			return err
		}
	}

	_, err = outputFile.WriteString("\n")
	if err != nil {
		return err
	}
	_, err = outputFile.WriteString("# Definition of foreign keys.\n")
	if err != nil {
		return err
	}

	for _, fk := range foreignKeyConnections {
		_, err = outputFile.WriteString(fmt.Sprintf("%v\n", fk))
		if err != nil {
			return err
		}
	}

	return nil
}

// defineFilesToParse prepares and returns the list of files that the service need to parse.
func defineFilesToParse(directory string, filesList []string) []string {
	filesToParse := filesList
	if directory != "" {
		directory = strings.TrimRight(directory, "/")
		filesToParse = []string{}
		files, err := ioutil.ReadDir(directory)
		if err != nil {
			log.Fatal(err)
		}
		for _, f := range files {
			if f.IsDir() {
				continue
			}

			fullFilePath := fmt.Sprintf("%v/%v", directory, f.Name())
			extension := filepath.Ext(fullFilePath)
			if extension != ".go" {
				continue
			}
			filesToParse = append(filesToParse, fullFilePath)
		}
	}

	return filesToParse
}

// getTagRegex prepares and returns the regexp for getting the value of a tag.
func getTagRegexp(tag string) *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf("%v:\"(.*?)\"", tag))
}

// getTagFieldsFromStruct retrieves and returns the values that exist on a respective tag.
func getTagFieldsFromStruct(tagRegexp *regexp.Regexp, fields []*ast.Field) []column {
	var cols []column
	for _, field := range fields {
		match := tagRegexp.FindStringSubmatch(field.Tag.Value)

		if len(match) == 0 {
			continue
		}

		newCol := column{
			fieldName: match[1],
			fieldType: field.Type,
		}
		cols = append(cols, newCol)
	}

	return cols
}
