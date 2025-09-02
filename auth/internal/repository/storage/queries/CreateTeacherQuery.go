package storage

const (
	CreateTeacherQuery = `
        INSERT INTO teachers (user_id, university_id, degree)
        VALUES ($1, $2, $3)
    `
)
