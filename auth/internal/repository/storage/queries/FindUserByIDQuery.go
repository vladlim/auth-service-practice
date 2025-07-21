package storage

const (
	FindUserByIDQuery = `
	SELECT EXISTS
	(SELECT 1 FROM users 
	WHERE id = $1)`
)
