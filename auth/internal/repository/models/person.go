package models

type Person struct {
	ID       int    `db:"id"`
	Username string `db:"username"`
}
