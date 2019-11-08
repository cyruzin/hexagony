package albums

import (
	"context"

	"github.com/google/uuid"
)

type service struct {
	albumRepository Repository
}

// NewService creates a service with the necessary dependencies.
func NewService(r Repository) Service {
	return &service{r}
}

// FindAll finds the latest albums.
func (s *service) FindAll(ctx context.Context) ([]*Album, error) {
	album, err := s.albumRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return album, nil
}

// FindByID finds an album by ID.
func (s *service) FindByID(ctx context.Context, uuid uuid.UUID) (*Album, error) {
	album, err := s.albumRepository.FindByID(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return album, nil
}

// Add adds the given album to the database.
func (s *service) Add(ctx context.Context, album *Album) error {
	if err := s.albumRepository.Add(ctx, album); err != nil {
		return err
	}
	return nil
}

// Update updates an album by ID.
func (s *service) Update(ctx context.Context, uuid uuid.UUID, album *Album) error {
	if err := s.albumRepository.Update(ctx, uuid, album); err != nil {
		return err
	}
	return nil
}

// Delete deletes an album by ID.
func (s *service) Delete(ctx context.Context, uuid uuid.UUID) error {
	if err := s.albumRepository.Delete(ctx, uuid); err != nil {
		return err
	}
	return nil
}
