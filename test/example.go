package test

// User example test struct.
type User struct {
	ID        int
	FirstName string `db:"first_name"`
	LastName  string `db:"lastname"`
}

// PhoneNumber example test struct.
type PhoneNumber struct {
	ID       int
	UserID   int    `db:"user_id"`
	Mobile   string `db:"mobile"`
	Landline string `db:"landline"`
}

// Address example test struct.
type Address struct {
	ID      int    `db:"id"`
	UserID  int    `db:"user_id"`
	Street  string `db:"street"`
	Number  string `db:"number"`
	ZipCode string `db:"zip_code"`
	CityID  int    `db:"city_id"`
}

// City example test struct.
type City struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}
