package storage

const (
	AddUserRoleQuery = `
        INSERT INTO user_roles (user_id, role_id)
        SELECT $1, id FROM roles WHERE name = $2
    `
)
