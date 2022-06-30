package usecase

import (
	"context"
	"hexagony/internal/albums/domain"

	"github.com/google/uuid"
)

type albumUseCase struct {
	albumRepository domain.AlbumRepository
}

func NewAlbumUseCase(ar domain.AlbumRepository) domain.AlbumUseCase {
	return &albumUseCase{albumRepository: ar}
}

func (s *albumUseCase) FindAll(ctx context.Context) ([]*domain.Album, error) {
	album, err := s.albumRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return album, nil
}

func (s *albumUseCase) FindByID(ctx context.Context, uuid uuid.UUID) (*domain.Album, error) {
	album, err := s.albumRepository.FindByID(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return album, nil
}

func (s *albumUseCase) Add(ctx context.Context, album *domain.Album) error {
	if err := s.albumRepository.Add(ctx, album); err != nil {
		return err
	}
	return nil
}

func (s *albumUseCase) Update(ctx context.Context, uuid uuid.UUID, album *domain.Album) error {
	if err := s.albumRepository.Update(ctx, uuid, album); err != nil {
		return err
	}
	return nil
}

func (s *albumUseCase) Delete(ctx context.Context, uuid uuid.UUID) error {
	if err := s.albumRepository.Delete(ctx, uuid); err != nil {
		return err
	}
	return nil
}
