package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/Angelos-Giannis/erd-builder/internal/config"
	"github.com/Angelos-Giannis/erd-builder/internal/domain"
	"github.com/urfave/cli/v2"
)

const (
	// define the configuration file for the system.
	configurationFile = "configuration.yaml"
)

func main() {
	cfg, err := config.New(configurationFile)
	if err != nil {
		fmt.Printf("Error parsing configuration : %v\n", err)
		os.Exit(1)
	}

	var app = cli.NewApp()
	info(app, cfg)

	var options domain.Options

	app.Commands = []*cli.Command{
		&cli.Command{
			Name:    "generate",
			Aliases: []string{"g"},
			Usage:   "Generate the .er file based on the provided structures.",
			Flags: []cli.Flag{
				options.GetCommonField(),
				options.GetDirectoryFlag(),
				options.GetFileList(),
				options.GetIDField(),
				options.GetOutputFilename(),
				options.GetOutputPath(),
				options.GetTag(),
				options.GetTitle(),
			},
			Action: func(c *cli.Context) error {
				if options.Directory == "" && len(options.FileList.Value()) == 0 {
					return errors.New("Need to provide at least one of 'directory' or 'file_list'")
				}

				filesToParse := options.FileList.Value()
				if options.Directory != "" {
					options.Directory = strings.TrimRight(options.Directory, "/")
					filesToParse = []string{}
					files, err := ioutil.ReadDir(options.Directory)
					if err != nil {
						log.Fatal(err)
					}
					for _, f := range files {
						filesToParse = append(filesToParse, fmt.Sprintf("%v/%v", options.Directory, f.Name()))
					}
				}

				re := regexp.MustCompile(fmt.Sprintf("%v:\"(.*?)\"", options.Tag))
				type column struct {
					fieldName  interface{}
					fieldType  interface{}
					fieldLabel interface{}
				}
				tableMapping := make(map[*ast.Ident][]column)

				for _, fl := range filesToParse {
					fset := token.NewFileSet()
					node, err := parser.ParseFile(fset, fl, nil, parser.ParseComments)
					if err != nil {
						log.Fatal(err)
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

							tableName := typeDecl.Specs[j].(*ast.TypeSpec).Name
							if _, exists := tableMapping[tableName]; !exists {
								tableMapping[tableName] = []column{}
							}

							fields := structDecl.Fields.List

							for _, field := range fields {
								match := re.FindStringSubmatch(field.Tag.Value)

								if len(match) == 0 {
									continue
								}

								newCol := column{
									fieldName: match[1],
									fieldType: field.Type,
								}
								tableMapping[tableName] = append(tableMapping[tableName], newCol)
							}
						}
					}
				}

				outputFile, err := os.Create(fmt.Sprintf("%v/%v.er", options.OutputPath, options.OutputFilename))
				if err != nil {
					panic(err)
				}
				defer outputFile.Close()

				if options.Title != "" {
					outputFile.WriteString(fmt.Sprintf("title {label: \"%v\"}\n", options.Title))
				}

				var foreignKeyConnections []string

				outputFile.WriteString("# Definition of tables.\n")
				for tb := range tableMapping {
					checkingTable := strings.ToLower(fmt.Sprintf("%v", tb))
					if len(tableMapping[tb]) == 0 {
						continue
					}

					outputFile.WriteString(fmt.Sprintf("[%v]\n", tb))

					if options.IDField != "" {
						outputFile.WriteString(fmt.Sprintf("\t*%v\n", options.IDField))
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

						outputFile.WriteString(fmt.Sprintf("\t%v%v {label: \"%v\"}\n", fkPrefix, col.fieldName, col.fieldType))
					}

					for _, c := range options.CommonFields.Value() {
						outputFile.WriteString(fmt.Sprintf("\t%v%v\n", "", c))
					}
					outputFile.WriteString("\n")
				}

				outputFile.WriteString("\n")
				outputFile.WriteString("# Definition of foreign keys.\n")

				for _, fk := range foreignKeyConnections {
					outputFile.WriteString(fmt.Sprintf("%v\n", fk))
				}

				return nil
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// info sets up the information of the tool.
func info(app *cli.App, cfg config.Config) {
	var appAuthors []*cli.Author
	for _, author := range cfg.Application.Authors {
		newAuthor := cli.Author{
			Name:  author.Name,
			Email: author.Email,
		}
		appAuthors = append(appAuthors, &newAuthor)
	}

	app.Authors = appAuthors
	app.Name = cfg.Application.Name
	app.Usage = cfg.Application.Usage
	app.Version = cfg.Application.Version
}
