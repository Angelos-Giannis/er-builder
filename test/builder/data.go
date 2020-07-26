package test

import (
	"github.com/eujoy/erbuilder/internal/config"
	"github.com/eujoy/erbuilder/internal/domain"
	"github.com/urfave/cli/v2"
)

// DataBuilder describes the data builder being used for the tests.
type DataBuilder struct{}

// NewDataBuilder creates a new test data builder.
func NewDataBuilder() *DataBuilder {
	return &DataBuilder{}
}

// GetWriterTestDiagram returns a diagram to be used by the writer test.
func (d *DataBuilder) GetWriterTestDiagram() domain.Diagram {
	return domain.Diagram{
		Title: "example_db",
		TableList: []domain.Table{
			domain.Table{
				Name: "user",
				ColumnList: []domain.Column{
					createColumn("first_name", "varchar", false, false, false),
					createColumn("lastname", "varchar", false, false, false),
					createColumn("id", "integer", true, false, false),
				},
			},
			domain.Table{
				Name: "phone_number",
				ColumnList: []domain.Column{
					createColumn("user_id", "integer", false, true, false),
					createColumn("mobile", "varchar", false, false, false),
					createColumn("landline", "varchar", false, false, false),
					createColumn("id", "integer", true, false, false),
				},
			},
			domain.Table{
				Name: "address",
				ColumnList: []domain.Column{
					createColumn("id", "integer", true, false, false),
					createColumn("user_id", "integer", false, true, false),
					createColumn("street", "varchar", false, false, false),
					createColumn("number", "varchar", false, false, false),
					createColumn("zip_code", "varchar", false, false, false),
					createColumn("city_id", "integer", false, true, false),
				},
			},
			domain.Table{
				Name: "city",
				ColumnList: []domain.Column{
					createColumn("id", "integer", true, false, false),
					createColumn("name", "varchar", false, false, false),
				},
			},
		},
		ReferenceList: []domain.Reference{
			createReference("phone_number", "user_id", "user", "*--*"),
			createReference("address", "user_id", "user", "*--*"),
			createReference("address", "city_id", "city", "*--*"),
		},
	}
}

// GetOptionsForOptionTest returns a list of options for the option_test file.
func (d *DataBuilder) GetOptionsForOptionTest(cfg config.Config) domain.Options {
	return domain.Options{
		Directory:      "./../../../test",
		IDField:        "id",
		FileList:       cli.StringSlice{},
		OutputFilename: "test-example-er-diagram",
		OutputPath:     "./../../../test",
		Tag:            "db",
		Title:          "example_db",
		ColumnNameCase: "snake_case",
		TableNameCase:  "snake_case",
		Config:         cfg,
	}
}

// createColumn creates and returns a domain.Column instance with the provided values.
func createColumn(name, colType string, isPrimaryKey, isForeignKey, isExtraField bool) domain.Column {
	return domain.Column{
		Name:         name,
		Type:         colType,
		IsPrimaryKey: isPrimaryKey,
		IsForeignKey: isForeignKey,
		IsExtraField: isExtraField,
	}
}

// createReference creates and returns a domain.Reference instance with the provided values.
func createReference(fromTableName, fromTableColumn, toTableName, typeOfReference string) domain.Reference {
	return domain.Reference{
		FromTableName:   fromTableName,
		FromTableColumn: fromTableColumn,
		ToTableName:     toTableName,
		TypeOfReference: typeOfReference,
	}
}
