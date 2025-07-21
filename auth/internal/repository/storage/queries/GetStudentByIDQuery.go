package storage

const (
	GetStudentByIDQuery = `
    SELECT 
        s.user_id, s.group_id, s.university_id, s.enrollment_year,
        u.id, u.username, u.email, u.first_name, u.last_name, u.created_at
    FROM students s
    JOIN users u ON s.user_id = u.id
    WHERE s.user_id = $1
`
)
