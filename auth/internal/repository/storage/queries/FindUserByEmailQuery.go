package storage

const (
	FindUserByEmailQuery = `
		SELECT id, password_hash
		FROM users
		WHERE email = $1
	`
)
