package queries

const (
	SqlAlbumsFindAll = "SELECT * FROM albums ORDER BY updated_at DESC LIMIT 10"

	SqlAlbumsFindByID = "SELECT * FROM albums WHERE uuid=$1"

	SqlAlbumsAdd = `
	INSERT INTO 
	albums (uuid, name, length, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5)
	`

	SqlAlbumsUpdate = `
	UPDATE albums 
	SET name=$1, length=$2, updated_at=$3
	WHERE uuid=$4
	`

	SqlAlbumsDelete = "DELETE FROM albums WHERE uuid=$1"
)
