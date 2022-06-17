package domain

import "errors"

var (
	ErrFindAll   = errors.New("couldn't list the users")
	ErrFindByID  = errors.New("couldn't get the user")
	ErrAdd       = errors.New("couldn't insert the user")
	ErrUpdate    = errors.New("couldn't update the user")
	ErrDelete    = errors.New("couldn't delete the user")
	ErrUUIDParse = errors.New("couldn't parse the UUID")

	ErrResourceNotFound = errors.New("the resource you requested could not be found")
)
