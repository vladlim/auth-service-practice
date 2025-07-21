package models

type Student struct {
	UserID         string `db:"user_id"`
	GroupID        string `db:"group_id"`
	UniversityID   string `db:"university_id"`
	EnrollmentYear int    `db:"enrollment_year"`
	User           User   `db:"-"`
}
