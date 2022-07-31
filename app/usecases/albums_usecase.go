package usecase

import (
	"context"
	"hexagony/app/domain"

	"github.com/google/uuid"
)

type albumsUseCase struct {
	albumRepository domain.AlbumsRepository
}

func NewAlbumsUseCase(ar domain.AlbumsRepository) domain.AlbumsUseCase {
	return &albumsUseCase{albumRepository: ar}
}

func (s *albumsUseCase) FindAll(ctx context.Context) ([]*domain.Albums, error) {
	album, err := s.albumRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return album, nil
}

func (s *albumsUseCase) FindByID(ctx context.Context, uuid uuid.UUID) (*domain.Albums, error) {
	album, err := s.albumRepository.FindByID(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return album, nil
}

func (s *albumsUseCase) Add(ctx context.Context, album *domain.Albums) error {
	if err := s.albumRepository.Add(ctx, album); err != nil {
		return err
	}
	return nil
}

func (s *albumsUseCase) Update(ctx context.Context, uuid uuid.UUID, album *domain.Albums) error {
	if err := s.albumRepository.Update(ctx, uuid, album); err != nil {
		return err
	}
	return nil
}

func (s *albumsUseCase) Delete(ctx context.Context, uuid uuid.UUID) error {
	if err := s.albumRepository.Delete(ctx, uuid); err != nil {
		return err
	}
	return nil
}
