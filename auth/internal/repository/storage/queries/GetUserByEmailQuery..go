package storage

const (
	GetUserByEmailQuery = `
	SELECT id, username, email, first_name, last_name, created_at
	FROM users
	WHERE email = $1`
)
