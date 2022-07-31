package postgres

const (
	sqlUsersFindAll = "SELECT uuid,name,email,created_at,updated_at FROM users ORDER BY updated_at DESC LIMIT 10"

	sqlUsersFindByID = "SELECT uuid,name,email,created_at,updated_at FROM users WHERE uuid=$1"

	sqlUsersAdd = `
	INSERT INTO 
	users (uuid, name, email, password, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6)
	`

	sqlUsersUpdate = `
	UPDATE users 
	SET name=$1, email=$2, updated_at=$3
	WHERE uuid=$4
	`

	sqlUsersDelete = "DELETE FROM users WHERE uuid=$1"

	sqlUsersCheckDuplicate = "SELECT email FROM users WHERE email=$1"
)
