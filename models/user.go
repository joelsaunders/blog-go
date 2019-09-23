package models

type User struct {
	ID       int    `db:"id"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

type NewUser struct {
	Email    string `db:"email"`
	Password string `db:"password"`
}
