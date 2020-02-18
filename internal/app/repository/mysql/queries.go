package mysql

const (
	sqlFindAll = "SELECT * FROM albums"
	sqlFindByID = "SELECT * FROM albums WHERE uuid=?"
	sqlAdd = `
	INSERT INTO 
	albums (uuid, name, length, created_at, updated_at) 
	VALUES (?, ?, ?, ?, ?)
	`
	sqlUpdate = `
	UPDATE albums 
	SET name=?, length=?, updated_at=?
	WHERE uuid=?
	`
	sqlDelete = "DELETE FROM albums WHERE uuid=?"
)