package domain

// Diagram describes the details of the database
type Diagram struct {
	Title         string
	TableList     []Table
	ReferenceList []Reference
}

// Table describes the details of a table.
type Table struct {
	Name       string   `json:"name"`
	ColumnList []Column `json:"columns"`
	Color      string   `json:"color"`
}

// Column describes the details of a column.
type Column struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	IsPrimaryKey bool   `json:"is_primary_key"`
	IsForeignKey bool   `json:"is_foreign_key"`
	IsExtraField bool   `json:"is_extra_field"`
}

// Reference describes the references for a table.
type Reference struct {
	FromTableName   string
	FromTableColumn string
	ToTableName     string
	TypeOfReference string
}
