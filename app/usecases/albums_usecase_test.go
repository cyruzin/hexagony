package usecase

import (
	"context"
	"errors"
	"hexagony/app/domain"
	"hexagony/app/domain/mocks"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAlbumsFindAll(t *testing.T) {
	mockAlbumRepo := new(mocks.AlbumRepository)
	mockAlbum := &domain.Albums{
		UUID:      uuid.New(),
		Name:      "St. Anger",
		Length:    75,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockListAlbum := make([]*domain.Albums, 0)
	mockListAlbum = append(mockListAlbum, mockAlbum)

	t.Run("success", func(t *testing.T) {
		mockAlbumRepo.On("FindAll",
			mock.AnythingOfType("*context.emptyCtx")).
			Return(mockListAlbum, nil).Once()

		a := NewAlbumsUseCase(mockAlbumRepo)
		list, err := a.FindAll(context.TODO())

		assert.Equal(t, "St. Anger", list[0].Name)
		assert.NoError(t, err)
		assert.Len(t, list, len(mockListAlbum))

		mockAlbumRepo.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockAlbumRepo.On("FindAll",
			mock.AnythingOfType("*context.emptyCtx")).
			Return(nil, errors.New("Unexpected error")).Once()

		a := NewAlbumsUseCase(mockAlbumRepo)
		_, err := a.FindAll(context.TODO())

		assert.NotNil(t, err)

		mockAlbumRepo.AssertExpectations(t)
	})
}

func TestAlbumsFindByID(t *testing.T) {
	newUUID := uuid.New()
	mockAlbumRepo := new(mocks.AlbumRepository)
	mockAlbum := &domain.Albums{
		UUID:      newUUID,
		Name:      "St. Anger",
		Length:    75,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("success", func(t *testing.T) {
		mockAlbumRepo.On("FindByID",
			mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("uuid.UUID")).
			Return(mockAlbum, nil).Once()

		a := NewAlbumsUseCase(mockAlbumRepo)
		album, err := a.FindByID(context.TODO(), newUUID)

		assert.Equal(t, "St. Anger", album.Name)
		assert.NoError(t, err)
		mockAlbumRepo.AssertExpectations(t)
	})

	t.Run("failure", func(t *testing.T) {
		mockAlbumRepo.On("FindByID",
			mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("uuid.UUID")).
			Return(nil, errors.New("Unexpected error")).Once()

		a := NewAlbumsUseCase(mockAlbumRepo)
		_, err := a.FindByID(context.TODO(), newUUID)

		assert.NotNil(t, err)

		mockAlbumRepo.AssertExpectations(t)
	})
}

func TestAlbumsAdd(t *testing.T) {
	mockAlbumRepo := new(mocks.AlbumRepository)
	mockAlbum := &domain.Albums{
		UUID:      uuid.New(),
		Name:      "St. Anger",
		Length:    75,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("success", func(t *testing.T) {
		mockAlbumRepo.On("Add",
			mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("*domain.Albums")).
			Return(nil).Once()

		u := NewAlbumsUseCase(mockAlbumRepo)
		err := u.Add(context.TODO(), mockAlbum)

		assert.NoError(t, err)
		mockAlbumRepo.AssertExpectations(t)
	})

	t.Run("failure", func(t *testing.T) {
		mockAlbumRepo.On("Add",
			mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("*domain.Albums")).
			Return(errors.New("Unexpected error")).Once()

		u := NewAlbumsUseCase(mockAlbumRepo)
		err := u.Add(context.TODO(), mockAlbum)

		assert.NotNil(t, err)

		mockAlbumRepo.AssertExpectations(t)
	})
}

func TestAlbumsUpdate(t *testing.T) {
	newUUID := uuid.New()
	mockAlbumRepo := new(mocks.AlbumRepository)
	mockAlbum := &domain.Albums{
		UUID:      newUUID,
		Name:      "St. Anger",
		Length:    75,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("success", func(t *testing.T) {
		mockAlbumRepo.On("Update",
			mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("*domain.Albums")).
			Return(nil).Once()

		a := NewAlbumsUseCase(mockAlbumRepo)
		err := a.Update(context.TODO(), newUUID, mockAlbum)

		assert.NoError(t, err)
		mockAlbumRepo.AssertExpectations(t)
	})

	t.Run("failure", func(t *testing.T) {
		mockAlbumRepo.On("Update",
			mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("*domain.Albums")).
			Return(errors.New("Unexpected error")).Once()

		a := NewAlbumsUseCase(mockAlbumRepo)
		err := a.Update(context.TODO(), newUUID, mockAlbum)

		assert.NotNil(t, err)

		mockAlbumRepo.AssertExpectations(t)
	})
}

func TestAlbumsDelete(t *testing.T) {
	newUUID := uuid.New()
	mockAlbumRepo := new(mocks.AlbumRepository)

	t.Run("success", func(t *testing.T) {
		mockAlbumRepo.On("Delete",
			mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("uuid.UUID")).
			Return(nil).Once()

		u := NewAlbumsUseCase(mockAlbumRepo)
		err := u.Delete(context.TODO(), newUUID)

		assert.NoError(t, err)
		mockAlbumRepo.AssertExpectations(t)
	})

	t.Run("failure", func(t *testing.T) {
		mockAlbumRepo.On("Delete",
			mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("uuid.UUID")).
			Return(errors.New("Unexpected error")).Once()

		a := NewAlbumsUseCase(mockAlbumRepo)
		err := a.Delete(context.TODO(), newUUID)

		assert.NotNil(t, err)

		mockAlbumRepo.AssertExpectations(t)
	})
}
