package mariadb

const (
	sqlFindAll = "SELECT * FROM users ORDER BY updated_at DESC LIMIT 10"

	sqlFindByID = "SELECT * FROM users WHERE uuid=$1"

	sqlAdd = `
	INSERT INTO 
	users (uuid, name, email, password, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6)
	`

	sqlUpdate = `
	UPDATE users 
	SET name=$1, email=$2, password=$3, updated_at=$4
	WHERE uuid=$5
	`

	sqlDelete = "DELETE FROM users WHERE uuid=$1"
)
