package storage

const (
	GetTeacherByIDQuery = `
        SELECT 
            t.user_id, t.university_id, t.degree,
            u.id, u.username, u.email, u.first_name, u.last_name, u.created_at
        FROM teachers t
        JOIN users u ON t.user_id = u.id
        WHERE t.user_id = $1
    `
)
