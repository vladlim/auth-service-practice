package models

type Teacher struct {
	UserID       string `db:"user_id"`
	UniversityID string `db:"university_id"`
	Degree       string `db:"degree"`
	User         User   `db:"-"`
}
