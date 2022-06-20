package domain

import "errors"

var (
	ErrFindAll   = errors.New("failed to list the users")
	ErrFindByID  = errors.New("failed to get the user")
	ErrAdd       = errors.New("failed to insert the user")
	ErrUpdate    = errors.New("failed to update the user")
	ErrDelete    = errors.New("failed to delete the user")
	ErrUUIDParse = errors.New("failed to parse the UUID")

	ErrResourceNotFound = errors.New("the resource you requested could not be found")
	ErrHashPassword     = errors.New("failed to hash the password")
)
