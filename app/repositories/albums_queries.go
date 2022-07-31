package postgres

const (
	sqlAlbumsFindAll = "SELECT * FROM albums ORDER BY updated_at DESC LIMIT 10"

	sqlAlbumsFindByID = "SELECT * FROM albums WHERE uuid=$1"

	sqlAlbumsAdd = `
	INSERT INTO 
	albums (uuid, name, length, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5)
	`

	sqlAlbumsUpdate = `
	UPDATE albums 
	SET name=$1, length=$2, updated_at=$3
	WHERE uuid=$4
	`

	sqlAlbumsDelete = "DELETE FROM albums WHERE uuid=$1"
)
