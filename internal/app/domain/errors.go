package domain

import "errors"

var (
	ErrFindAll   = errors.New("couldn't list the albums")
	ErrFindByID  = errors.New("couldn't get the album")
	ErrAdd       = errors.New("couldn't insert the album")
	ErrUpdate    = errors.New("couldn't update the album")
	ErrDelete    = errors.New("couldn't delete the album")
	ErrUUIDParse = errors.New("couldn't parse the UUID")

	ErrResourceNotFound = errors.New("the resource you requested could not be found")
)
