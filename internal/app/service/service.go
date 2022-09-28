package service

import (
	"encoding/json"
	"errors"
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

	externalSurvey "github.com/AlecAivazis/survey/v2"
	"github.com/eujoy/erbuilder/internal/domain"
	"gopkg.in/go-playground/colors.v1"
)

const (
	columnTypeInteger = "integer"
	columnTypeVarchar = "varchar"
	columnTypeOther   = "other"

	addMoreTable   = "Table"
	addMoreColumn  = "Column"
	addMoreNothing = "Nothing"
)

var tableNameQuestion = []*externalSurvey.Question{
	{
		Name:     "table_name",
		Prompt:   &externalSurvey.Input{Message: "What is the name table name?"},
		Validate: externalSurvey.Required,
	},
	{
		Name: "table_color",
		Prompt: &externalSurvey.Input{
			Message: "What should be the background color for the table?",
			Default: "#ebe486",
		},
		Validate: func(val interface{}) error {
			_, err := colors.ParseHEX(fmt.Sprintf("%v", val))
			fmt.Println(err)
			if err != nil && err != colors.ErrBadColor {
				return errors.New("the provided value should a hexadecimal color value")
			}
			return nil
		},
	},
}

var columnDefinitionQuestion = []*externalSurvey.Question{
	{
		Name:      "column_name",
		Prompt:    &externalSurvey.Input{Message: "What is the name of the column?"},
		Validate:  externalSurvey.Required,
		Transform: externalSurvey.Title,
	},
	{
		Name: "column_type",
		Prompt: &externalSurvey.Select{
			Message: "Choose the type of the column:",
			Options: []string{columnTypeInteger, columnTypeVarchar, columnTypeOther},
			Default: "varchar",
		},
	},
	{
		Name: "is_primary_key",
		Prompt: &externalSurvey.Confirm{
			Message: "Is this field a primary key?",
			Default: false,
		},
	},
	{
		Name: "is_foreign_key",
		Prompt: &externalSurvey.Confirm{
			Message: "Is this field a foreign key?",
			Default: false,
		},
	},
	{
		Name: "add_more",
		Prompt: &externalSurvey.Select{
			Message: "Do you want to add some more:",
			Options: []string{addMoreTable, addMoreColumn, addMoreNothing},
			Default: addMoreColumn,
		},
	},
}

type survey interface {
	AskTableDetails(questions []*externalSurvey.Question) (domain.TableAnswer, error)
	AskColumnDetails(questions []*externalSurvey.Question) (domain.ColumnAnswer, error)
}

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
	survey  survey
	util    util
	writer  writer
}

// New creates and returns a new service.
func New(options domain.Options, survey survey, util util, writer writer) *Service {
	return &Service{
		options: options,
		survey:  survey,
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

	if s.options.ExtraTablesSurvey {
		extraTables, err := s.Build()
		if err != nil {
			return err
		}
		diagram.TableList = append(diagram.TableList, extraTables...)
	}

	if s.options.ExtraTablesDefinition != "" {
		var extraTables []domain.Table
		err := json.Unmarshal([]byte(s.options.ExtraTablesDefinition), &extraTables)
		if err != nil {
			return err
		}
		diagram.TableList = append(diagram.TableList, extraTables...)
	}

	err := s.writer.WriteFile(diagram)
	if err != nil {
		return err
	}

	return nil
}

// Build performs the action to build extra details for the cli tool.
func (s *Service) Build() ([]domain.Table, error) {
	var tableList []domain.Table
	isCompleted := false

	for {
		var columnList []domain.Column

		// perform the questions for table details.
		tAnswers, err := s.survey.AskTableDetails(tableNameQuestion)
		if err != nil {
			return []domain.Table{}, err
		}

		for {
			// perform the questions for table details.
			cAnswers, err := s.survey.AskColumnDetails(columnDefinitionQuestion)
			if err != nil {
				return []domain.Table{}, err
			}

			cl := domain.Column{
				Name:         s.util.GetCaseOfString(cAnswers.Name, s.options.ColumnNameCase),
				Type:         cAnswers.Type,
				IsPrimaryKey: cAnswers.IsPrimaryKey,
				IsForeignKey: cAnswers.IsForeignKey,
				IsExtraField: false,
			}

			columnList = append(columnList, cl)

			if cAnswers.AddMore != addMoreColumn {
				isCompleted = cAnswers.AddMore == addMoreNothing
				break
			}
		}

		t := domain.Table{
			Name:       s.util.GetCaseOfString(tAnswers.Name, s.options.TableNameCase),
			ColumnList: columnList,
			Color:      tAnswers.Color,
		}

		tableList = append(tableList, t)

		if isCompleted {
			break
		}
	}

	return tableList, nil
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

	columnList := s.getTagFieldsFromStruct(tagRegexp, structDecl.Fields.List)
	if len(columnList) == 0 {
		return tableDetails, false
	}

	tableDetails = domain.Table{
		Name:       s.util.GetCaseOfString(structName, s.options.TableNameCase),
		ColumnList: columnList,
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
			IsPrimaryKey: s.options.IDField == columnName,
			IsForeignKey: false,
			IsExtraField: false,
		}
		columns = append(columns, newCol)
		pkFound = pkFound || newCol.IsPrimaryKey
	}

	if len(columns) == 0 {
		return []domain.Column{}
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
				filesToParse = append(filesToParse, defineFilesToParse(fmt.Sprintf("%v/%v", directory, f.Name()), []string{})...)
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
