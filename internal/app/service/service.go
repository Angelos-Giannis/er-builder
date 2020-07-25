package service

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	"github.com/eujoy/erbuilder/internal/domain"
)

type util interface {
	GetCaseOfString(initialValue, convertToCase string) string
	GetValueCount(isPlural bool, initialValue string) string
	GetDBDataTypeFromCodeDataType(dataType string) string
}

type writer interface {
	WriteFile(diagram domain.Diagram) error
}

// Service describes the service flow.
type Service struct {
	options domain.Options
	util    util
	writer  writer
}

// New creates and returns a new service.
func New(options domain.Options, util util, writer writer) *Service {
	return &Service{
		options: options,
		util:    util,
		writer:  writer,
	}
}

// Generate performs the action to generate the .er file based on the provided input.
func (s *Service) Generate() error {
	filesToParse := defineFilesToParse(s.options.Directory, s.options.FileList.Value())
	diagram := domain.Diagram{Title: s.options.Title}

	for _, fl := range filesToParse {
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, fl, nil, parser.ParseComments)
		if err != nil {
			return err
		}

		diagram.TableList = append(diagram.TableList, s.getAllTables(node.Decls)...)
	}

	s.enrichForeignKeyReferences(&diagram)

	err := s.writer.WriteFile(diagram)
	if err != nil {
		return err
	}

	return nil

}

// getAllTables retrieves and returns all the table definitions with their columns.
func (s *Service) getAllTables(declarations []ast.Decl) []domain.Table {
	var tableList []domain.Table
	for i := 0; i < len(declarations); i++ {
		if reflect.TypeOf(declarations[i]) != reflect.TypeOf(&ast.GenDecl{}) {
			continue
		}
		typeDecl := declarations[i].(*ast.GenDecl)

		tableDefinition, found := s.getTableDefinition(typeDecl.Specs)
		if found {
			tableList = append(tableList, tableDefinition)
		}
	}

	return tableList
}

// getTableDefinition retrieves the definition of a table alongside with it's columns and returns it.
func (s *Service) getTableDefinition(typeSpec []ast.Spec) (domain.Table, bool) {
	var tableDetails domain.Table
	tagRegexp := getTagRegexp(s.options.Tag)

	if reflect.TypeOf(typeSpec[0]) != reflect.TypeOf(&ast.TypeSpec{}) {
		return tableDetails, false
	}

	if reflect.TypeOf(typeSpec[0].(*ast.TypeSpec).Type) != reflect.TypeOf(&ast.StructType{}) {
		return tableDetails, false
	}

	structDecl := typeSpec[0].(*ast.TypeSpec).Type.(*ast.StructType)
	structName := fmt.Sprintf("%v", typeSpec[0].(*ast.TypeSpec).Name)

	tableDetails = domain.Table{
		Name:       s.util.GetCaseOfString(structName, s.options.TableNameCase),
		ColumnList: s.getTagFieldsFromStruct(tagRegexp, structDecl.Fields.List),
	}

	return tableDetails, true
}

// getTagFieldsFromStruct retrieves and returns the values that exist on a respective tag.
func (s *Service) getTagFieldsFromStruct(tagRegexp *regexp.Regexp, fields []*ast.Field) []domain.Column {
	var columns []domain.Column
	pkFound := false
	for _, field := range fields {
		columnName := ""
		if field.Tag != nil && len(field.Tag.Value) > 0 {
			match := tagRegexp.FindStringSubmatch(field.Tag.Value)
			if len(match) < 2 {
				continue
			}

			columnName = match[1]
		}

		if len(columnName) == 0 {
			continue
		}

		newCol := domain.Column{
			Name:         columnName,
			Type:         s.util.GetDBDataTypeFromCodeDataType(fmt.Sprintf("%v", field.Type)),
			IsPrimaryKey: (s.options.IDField == columnName),
			IsForeignKey: false,
			IsExtraField: false,
		}
		columns = append(columns, newCol)
		pkFound = (pkFound || newCol.IsPrimaryKey)
	}

	if !pkFound {
		pkColumn := domain.Column{
			Name:         s.options.IDField,
			Type:         "integer",
			IsPrimaryKey: true,
			IsForeignKey: false,
			IsExtraField: false,
		}
		columns = append(columns, pkColumn)
	}

	if len(s.options.CommonFields.Value()) > 0 {
		for _, commonField := range s.options.CommonFields.Value() {
			extraColumn := domain.Column{
				Name:         commonField,
				Type:         "",
				IsPrimaryKey: false,
				IsForeignKey: false,
				IsExtraField: true,
			}
			columns = append(columns, extraColumn)
		}
	}

	return columns
}

// enrichForeignKeyReferences enriches the references between tables in the diagram.
func (s *Service) enrichForeignKeyReferences(diagram *domain.Diagram) {
	for idx := range diagram.TableList {
		diagram.ReferenceList = append(diagram.ReferenceList, s.getReferencesToTable(diagram, diagram.TableList[idx].Name)...)
	}
}

// getReferencesToTable finds and returns a list of all the references to a table.
func (s *Service) getReferencesToTable(diagram *domain.Diagram, searchForTable string) []domain.Reference {
	var referenceList []domain.Reference
	for idxTb := range diagram.TableList {
		if diagram.TableList[idxTb].Name == searchForTable {
			continue
		}

		for idxCol := range diagram.TableList[idxTb].ColumnList {
			if strings.Contains(diagram.TableList[idxTb].ColumnList[idxCol].Name, searchForTable) || strings.Contains(diagram.TableList[idxTb].ColumnList[idxCol].Name, s.util.GetValueCount(s.options.TableNamePlural, searchForTable)) {
				newReference := domain.Reference{
					FromTableName:   diagram.TableList[idxTb].Name,
					FromTableColumn: diagram.TableList[idxTb].ColumnList[idxCol].Name,
					ToTableName:     searchForTable,
					TypeOfReference: "*--*",
				}
				referenceList = append(referenceList, newReference)
				diagram.TableList[idxTb].ColumnList[idxCol].IsForeignKey = true
			}
		}
	}

	return referenceList
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
