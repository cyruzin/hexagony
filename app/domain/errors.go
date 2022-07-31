package domain

import "errors"

var (
	ErrAuth             = errors.New("authentication failed")
	ErrAuthEmptyClaim   = errors.New("claim is empty")
	ErrAuthSign         = errors.New("failed to sign the key")
	ErrAuthUserNotFound = errors.New("user not found")
	ErrAuthPassword     = errors.New("wrong password")
)

var (
	ErrAlbumsFindAll   = errors.New("failed to list the albums")
	ErrAlbumsFindByID  = errors.New("failed to get the album")
	ErrAlbumsAdd       = errors.New("failed to insert the album")
	ErrAlbumsUpdate    = errors.New("failed to update the album")
	ErrAlbumsDelete    = errors.New("failed to delete the album")
	ErrAlbumsUUIDParse = errors.New("failed to parse the UUID")
)

var (
	ErrUsersFindAll        = errors.New("failed to list the users")
	ErrUsersFindByID       = errors.New("failed to get the user")
	ErrUsersAdd            = errors.New("failed to insert the user")
	ErrUsersUpdate         = errors.New("failed to update the user")
	ErrUsersDelete         = errors.New("failed to delete the user")
	ErrUsersUUIDParse      = errors.New("failed to parse the UUID")
	ErrUsersHashPassword   = errors.New("failed to hash the password")
	ErrUsersDuplicateEmail = errors.New("this email already exists")
)

var (
	ErrResourceNotFound = errors.New("the resource you requested could not be found")
)
