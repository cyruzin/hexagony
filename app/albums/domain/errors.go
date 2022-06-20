package domain

import "errors"

var (
	ErrFindAll   = errors.New("failed to list the albums")
	ErrFindByID  = errors.New("failed to get the album")
	ErrAdd       = errors.New("failed to insert the album")
	ErrUpdate    = errors.New("failed to update the album")
	ErrDelete    = errors.New("failed to delete the album")
	ErrUUIDParse = errors.New("failed to parse the UUID")

	ErrResourceNotFound = errors.New("the resource you requested could not be found")
)
