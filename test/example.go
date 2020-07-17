package example

// User example test struct.
type User struct {
	ID        int    ``
	FirstName string `db:"first_name"`
	LastName  string `db:"lastname"`
}

// PhoneNumber example test struct.
type PhoneNumber struct {
	ID       int    ``
	UserID   int    `db:"user_id"`
	Mobile   string `db:"mobile"`
	Landline string `db:"landline"`
}
