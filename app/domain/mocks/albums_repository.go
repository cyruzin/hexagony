package mocks

import (
	"context"
	"hexagony/app/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type AlbumRepository struct {
	mock.Mock
}

func (m *AlbumRepository) FindAll(ctx context.Context) ([]*domain.Albums, error) {
	args := m.Called(ctx)

	var album []*domain.Albums

	if rf, ok := args.Get(0).(func(context.Context) []*domain.Albums); ok {
		album = rf(ctx)
	} else {
		if args.Get(0) != nil {
			album = args.Get(0).([]*domain.Albums)
		}
	}

	var err error
	if rf, ok := args.Get(1).(func(context.Context) error); ok {
		err = rf(ctx)
	} else {
		err = args.Error(1)
	}

	return album, err
}

func (m *AlbumRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Albums, error) {
	args := m.Called(ctx, id)

	var album *domain.Albums

	if rf, ok := args.Get(0).(func(context.Context, uuid.UUID) *domain.Albums); ok {
		album = rf(ctx, id)
	} else {
		if args.Get(0) != nil {
			album = args.Get(0).(*domain.Albums)
		}
	}

	var err error

	if rf, ok := args.Get(1).(func(context.Context, uuid.UUID) error); ok {
		err = rf(ctx, id)
	} else {
		err = args.Error(1)
	}

	return album, err
}

func (m *AlbumRepository) Add(ctx context.Context, album *domain.Albums) error {
	args := m.Called(ctx, album)

	var err error

	if rf, ok := args.Get(0).(func(context.Context, *domain.Albums) error); ok {
		err = rf(ctx, album)
	} else {
		err = args.Error(0)
	}

	return err
}

func (m *AlbumRepository) Update(ctx context.Context, id uuid.UUID, album *domain.Albums) error {
	args := m.Called(ctx, id, album)

	var err error

	if rf, ok := args.Get(0).(func(context.Context, uuid.UUID, *domain.Albums) error); ok {
		err = rf(ctx, id, album)
	} else {
		err = args.Error(0)
	}

	return err
}

func (m *AlbumRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)

	var err error

	if rf, ok := args.Get(0).(func(context.Context, uuid.UUID) error); ok {
		err = rf(ctx, id)
	} else {
		err = args.Error(0)
	}

	return err
}
