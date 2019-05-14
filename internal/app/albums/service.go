package albums

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	// ErrInvalidID is returned when the provided ID is invalid.
	ErrInvalidID = errors.New("invalid id provided")

	// ErrDuplicate is used when a album already exists.
	ErrDuplicate = errors.New("album already exists")

	// ErrInvalidPrice is returned when the provided price is invalid.
	ErrInvalidPrice = errors.New("invalid price provided")

	// ErrNotFound is returned when the album is not found.
	ErrNotFound = errors.New("could not find the album")
)

// Service provides albums operations.
type Service interface {
	AddAlbum(ctx context.Context, album *Album) error
	DeleteAlbum(ctx context.Context, uuid uuid.UUID) error
}

// Repository provides access to albums repository.
type Repository interface {
	Add(context.Context, *Album) error
	Find(context.Context, uuid.UUID) (*Album, error)
	FindAll(context.Context) ([]*Album, error)
	Update(context.Context, uuid.UUID, *Album) error
	Delete(context.Context, uuid.UUID) error
}

type service struct {
	aR Repository
}

// NewService creates a service with the necessary dependencies.
func NewService(r Repository) Service {
	return &service{r}
}

// AddAlbum adds the given album to the database.
func (s *service) AddAlbum(ctx context.Context, album *Album) error {
	if err := s.aR.Add(ctx, album); err != nil {
		return err
	}
	return nil
}

// DeleteAlbum deletes the given album from the database.
func (s *service) DeleteAlbum(ctx context.Context, uuid uuid.UUID) error {
	if err := s.aR.Delete(ctx, uuid); err != nil {
		return err
	}
	return nil
}
