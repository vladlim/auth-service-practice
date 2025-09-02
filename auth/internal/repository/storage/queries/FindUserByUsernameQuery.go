package storage

const (
	FindUserByUsernameQuery = `
		SELECT id, password_hash
		FROM users
		WHERE username = $1
	`
)
