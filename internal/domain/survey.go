package domain

// TableAnswer captures the answers for the table definition questions.
type TableAnswer struct {
	Name  string `survey:"table_name"`
	Color string `survey:"table_color"`
}

// ColumnAnswer captures the answers for the column related questions.
type ColumnAnswer struct {
	Name         string `survey:"column_name"`
	Type         string `survey:"column_type"`
	IsPrimaryKey bool   `survey:"is_primary_key"`
	IsForeignKey bool   `survey:"is_foreign_key"`
	AddMore      string `survey:"add_more"`
}
