package test

import (
	"github.com/eujoy/erbuilder/internal/domain"
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
					domain.Column{
						Name:         "first_name",
						Type:         "varchar",
						IsPrimaryKey: false,
						IsForeignKey: false,
						IsExtraField: false,
					},
					domain.Column{
						Name:         "lastname",
						Type:         "varchar",
						IsPrimaryKey: false,
						IsForeignKey: false,
						IsExtraField: false,
					},
					domain.Column{
						Name:         "id",
						Type:         "integer",
						IsPrimaryKey: true,
						IsForeignKey: false,
						IsExtraField: false,
					},
				},
			},
			domain.Table{
				Name: "phone_number",
				ColumnList: []domain.Column{
					domain.Column{
						Name:         "user_id",
						Type:         "integer",
						IsPrimaryKey: false,
						IsForeignKey: true,
						IsExtraField: false,
					},
					domain.Column{
						Name:         "mobile",
						Type:         "varchar",
						IsPrimaryKey: false,
						IsForeignKey: false,
						IsExtraField: false,
					},
					domain.Column{
						Name:         "landline",
						Type:         "varchar",
						IsPrimaryKey: false,
						IsForeignKey: false,
						IsExtraField: false,
					},
					domain.Column{
						Name:         "id",
						Type:         "integer",
						IsPrimaryKey: true,
						IsForeignKey: false,
						IsExtraField: false,
					},
				},
			},
			domain.Table{
				Name: "address",
				ColumnList: []domain.Column{
					domain.Column{
						Name:         "id",
						Type:         "integer",
						IsPrimaryKey: true,
						IsForeignKey: false,
						IsExtraField: false,
					},
					domain.Column{
						Name:         "user_id",
						Type:         "integer",
						IsPrimaryKey: false,
						IsForeignKey: true,
						IsExtraField: false,
					},
					domain.Column{
						Name:         "street",
						Type:         "varchar",
						IsPrimaryKey: false,
						IsForeignKey: false,
						IsExtraField: false,
					},
					domain.Column{
						Name:         "number",
						Type:         "varchar",
						IsPrimaryKey: false,
						IsForeignKey: false,
						IsExtraField: false,
					},
					domain.Column{
						Name:         "zip_code",
						Type:         "varchar",
						IsPrimaryKey: false,
						IsForeignKey: false,
						IsExtraField: false,
					},
					domain.Column{
						Name:         "city_id",
						Type:         "integer",
						IsPrimaryKey: false,
						IsForeignKey: true,
						IsExtraField: false,
					},
				},
			},
			domain.Table{
				Name: "city",
				ColumnList: []domain.Column{
					domain.Column{
						Name:         "id",
						Type:         "integer",
						IsPrimaryKey: true,
						IsForeignKey: false,
						IsExtraField: false,
					},
					domain.Column{
						Name:         "name",
						Type:         "varchar",
						IsPrimaryKey: false,
						IsForeignKey: false,
						IsExtraField: false,
					},
				},
			},
		},
		ReferenceList: []domain.Reference{
			domain.Reference{
				FromTableName:   "phone_number",
				FromTableColumn: "user_id",
				ToTableName:     "user",
				TypeOfReference: "*--*",
			},
			domain.Reference{
				FromTableName:   "address",
				FromTableColumn: "user_id",
				ToTableName:     "user",
				TypeOfReference: "*--*",
			},
			domain.Reference{
				FromTableName:   "address",
				FromTableColumn: "city_id",
				ToTableName:     "city",
				TypeOfReference: "*--*",
			},
		},
	}
}
