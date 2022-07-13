package postgres

const (
	sqlFindAll = "SELECT uuid,name,email,created_at,updated_at FROM users ORDER BY updated_at DESC LIMIT 10"

	sqlFindByID = "SELECT uuid,name,email,created_at,updated_at FROM users WHERE uuid=$1"

	sqlAdd = `
	INSERT INTO 
	users (uuid, name, email, password, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6)
	`

	sqlUpdate = `
	UPDATE users 
	SET name=$1, email=$2, updated_at=$3
	WHERE uuid=$4
	`

	sqlDelete = "DELETE FROM users WHERE uuid=$1"

	sqlCheckDuplicate = "SELECT email FROM users WHERE email=$1"
)
