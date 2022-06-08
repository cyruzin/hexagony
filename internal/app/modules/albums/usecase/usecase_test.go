package usecase

import (
	"context"
	"errors"
	"hexagony/internal/app/domain"
	"hexagony/internal/app/domain/mocks"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFindAll(t *testing.T) {
	mockAlbumRepo := new(mocks.AlbumRepository)
	mockAlbum := &domain.Album{
		UUID:      uuid.New(),
		Name:      "St. Anger",
		Length:    60,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockListAlbum := make([]*domain.Album, 0)
	mockListAlbum = append(mockListAlbum, mockAlbum)

	t.Run("success", func(t *testing.T) {
		mockAlbumRepo.On("FindAll",
			mock.AnythingOfType("*context.emptyCtx")).
			Return(mockListAlbum, nil).Once()

		a := NewAlbumUseCase(mockAlbumRepo)
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

		a := NewAlbumUseCase(mockAlbumRepo)
		_, err := a.FindAll(context.TODO())

		assert.NotNil(t, err)

		mockAlbumRepo.AssertExpectations(t)
	})
}

func TestFindByID(t *testing.T) {
	newUUID := uuid.New()
	mockAlbumRepo := new(mocks.AlbumRepository)
	mockAlbum := &domain.Album{
		UUID:      newUUID,
		Name:      "St. Anger",
		Length:    60,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("success", func(t *testing.T) {
		mockAlbumRepo.On("FindByID",
			mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("uuid.UUID")).
			Return(mockAlbum, nil).Once()

		a := NewAlbumUseCase(mockAlbumRepo)
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

		a := NewAlbumUseCase(mockAlbumRepo)
		_, err := a.FindByID(context.TODO(), newUUID)

		assert.NotNil(t, err)

		mockAlbumRepo.AssertExpectations(t)
	})
}

func TestAdd(t *testing.T) {
	mockAlbumRepo := new(mocks.AlbumRepository)
	mockAlbum := &domain.Album{
		UUID:      uuid.New(),
		Name:      "St. Anger",
		Length:    60,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("success", func(t *testing.T) {
		mockAlbumRepo.On("Add",
			mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("*domain.Album")).
			Return(nil).Once()

		u := NewAlbumUseCase(mockAlbumRepo)
		err := u.Add(context.TODO(), mockAlbum)

		assert.NoError(t, err)
		mockAlbumRepo.AssertExpectations(t)
	})

	t.Run("failure", func(t *testing.T) {
		mockAlbumRepo.On("Add",
			mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("*domain.Album")).
			Return(errors.New("Unexpected error")).Once()

		u := NewAlbumUseCase(mockAlbumRepo)
		err := u.Add(context.TODO(), mockAlbum)

		assert.NotNil(t, err)

		mockAlbumRepo.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	newUUID := uuid.New()
	mockAlbumRepo := new(mocks.AlbumRepository)
	mockAlbum := &domain.Album{
		UUID:      newUUID,
		Name:      "St. Anger",
		Length:    60,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("success", func(t *testing.T) {
		mockAlbumRepo.On("Update",
			mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("*domain.Album")).
			Return(nil).Once()

		a := NewAlbumUseCase(mockAlbumRepo)
		err := a.Update(context.TODO(), newUUID, mockAlbum)

		assert.NoError(t, err)
		mockAlbumRepo.AssertExpectations(t)
	})

	t.Run("failure", func(t *testing.T) {
		mockAlbumRepo.On("Update",
			mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("*domain.Album")).
			Return(errors.New("Unexpected error")).Once()

		a := NewAlbumUseCase(mockAlbumRepo)
		err := a.Update(context.TODO(), newUUID, mockAlbum)

		assert.NotNil(t, err)

		mockAlbumRepo.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	newUUID := uuid.New()
	mockAlbumRepo := new(mocks.AlbumRepository)

	t.Run("success", func(t *testing.T) {
		mockAlbumRepo.On("Delete",
			mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("uuid.UUID")).
			Return(nil).Once()

		u := NewAlbumUseCase(mockAlbumRepo)
		err := u.Delete(context.TODO(), newUUID)

		assert.NoError(t, err)
		mockAlbumRepo.AssertExpectations(t)
	})

	t.Run("failure", func(t *testing.T) {
		mockAlbumRepo.On("Delete",
			mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("uuid.UUID")).
			Return(errors.New("Unexpected error")).Once()

		a := NewAlbumUseCase(mockAlbumRepo)
		err := a.Delete(context.TODO(), newUUID)

		assert.NotNil(t, err)

		mockAlbumRepo.AssertExpectations(t)
	})
}
