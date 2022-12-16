package query

const (
	SqlUsersFindAll = "SELECT uuid,name,email,created_at,updated_at FROM users ORDER BY updated_at DESC LIMIT 10"

	SqlUsersFindByID = "SELECT uuid,name,email,created_at,updated_at FROM users WHERE uuid=$1"

	SqlUsersAdd = `
	INSERT INTO 
	users (uuid, name, email, password, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6)
	`

	SqlUsersUpdate = `
	UPDATE users 
	SET name=$1, email=$2, updated_at=$3
	WHERE uuid=$4
	`

	SqlUsersDelete = "DELETE FROM users WHERE uuid=$1"

	SqlUsersCheckDuplicate = "SELECT email FROM users WHERE email=$1"
)
