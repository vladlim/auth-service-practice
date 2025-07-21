package storage

const (
	GetUserByIDQuery = `
	SELECT id, username, email, first_name, last_name, created_at
	FROM users
	WHERE id = $1`
)
