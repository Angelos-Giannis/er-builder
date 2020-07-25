package domain

// Diagram describes the details of the database
type Diagram struct {
	Title         string
	TableList     []Table
	ReferenceList []Reference
}

// Table describes the details of a table.
type Table struct {
	Name       string
	ColumnList []Column
}

// Column describes the details of a column.
type Column struct {
	Name         string
	Type         string
	IsPrimaryKey bool
	IsForeignKey bool
	IsExtraField bool
}

// Reference describes the references for a table.
type Reference struct {
	FromTableName   string
	FromTableColumn string
	ToTableName     string
	TypeOfReference string
}
