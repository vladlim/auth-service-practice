package storage

const (
	CreateStudentQuery = `
        INSERT INTO students (user_id, group_id, university_id, enrollment_year)
        VALUES ($1, $2, $3, $4)
    `
)
