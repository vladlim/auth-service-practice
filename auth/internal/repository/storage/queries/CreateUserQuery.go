package storage

const (
	CreateUserQuery = `
        INSERT INTO users (
            username,
            email,
            password_hash,
            first_name,
            last_name
        ) VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `
)
