package postgres

const (
	sqlFindAll = "SELECT * FROM albums ORDER BY updated_at DESC LIMIT 10"

	sqlFindByID = "SELECT * FROM albums WHERE uuid=$1"

	sqlAdd = `
	INSERT INTO 
	albums (uuid, name, length, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5)
	`

	sqlUpdate = `
	UPDATE albums 
	SET name=$1, length=$2, updated_at=$3
	WHERE uuid=$4
	`

	sqlDelete = "DELETE FROM albums WHERE uuid=$1"
)
