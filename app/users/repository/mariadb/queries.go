package mariadb

const (
	sqlFindAll = "SELECT * FROM users ORDER BY uuid DESC LIMIT 10"

	sqlFindByID = "SELECT * FROM users WHERE uuid=?"

	sqlAdd = `
	INSERT INTO 
	users (uuid, name, email, password, created_at, updated_at) 
	VALUES (?, ?, ?, ?, ?, ?)
	`

	sqlUpdate = `
	UPDATE users 
	SET name=?, email=?, password=?, updated_at=?
	WHERE uuid=?
	`

	sqlDelete = "DELETE FROM users WHERE uuid=?"
)
