package usecase

import (
	"context"
	"errors"
	"hexagony/app/users/domain"
	"hexagony/app/users/domain/mocks"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFindAll(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	mockUser := &domain.User{
		UUID:      uuid.New(),
		Name:      "Cyro Dubeux",
		Email:     "xorycx@gmailcom",
		Password:  "12345678",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockListUsers := make([]*domain.User, 0)
	mockListUsers = append(mockListUsers, mockUser)

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("FindAll",
			mock.AnythingOfType("*context.emptyCtx")).
			Return(mockListUsers, nil).Once()

		a := NewUserUseCase(mockUserRepo)
		list, err := a.FindAll(context.TODO())

		assert.Equal(t, "Cyro Dubeux", list[0].Name)
		assert.NoError(t, err)
		assert.Len(t, list, len(mockListUsers))

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockUserRepo.On("FindAll",
			mock.AnythingOfType("*context.emptyCtx")).
			Return(nil, errors.New("Unexpected error")).Once()

		a := NewUserUseCase(mockUserRepo)
		_, err := a.FindAll(context.TODO())

		assert.NotNil(t, err)

		mockUserRepo.AssertExpectations(t)
	})
}

func TestFindByID(t *testing.T) {
	newUUID := uuid.New()
	mockUserRepo := new(mocks.UserRepository)
	mockUser := &domain.User{
		UUID:      newUUID,
		Name:      "Cyro Dubeux",
		Email:     "xorycx@gmailcom",
		Password:  "12345678",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("FindByID",
			mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("uuid.UUID")).
			Return(mockUser, nil).Once()

		a := NewUserUseCase(mockUserRepo)
		user, err := a.FindByID(context.TODO(), newUUID)

		assert.Equal(t, "Cyro Dubeux", user.Name)
		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("failure", func(t *testing.T) {
		mockUserRepo.On("FindByID",
			mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("uuid.UUID")).
			Return(nil, errors.New("Unexpected error")).Once()

		a := NewUserUseCase(mockUserRepo)
		_, err := a.FindByID(context.TODO(), newUUID)

		assert.NotNil(t, err)

		mockUserRepo.AssertExpectations(t)
	})
}

func TestAdd(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	mockUser := &domain.User{
		UUID:      uuid.New(),
		Name:      "Cyro Dubeux",
		Email:     "xorycx@gmailcom",
		Password:  "12345678",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("Add",
			mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("*domain.User")).
			Return(nil).Once()

		u := NewUserUseCase(mockUserRepo)
		err := u.Add(context.TODO(), mockUser)

		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("failure", func(t *testing.T) {
		mockUserRepo.On("Add",
			mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("*domain.User")).
			Return(errors.New("Unexpected error")).Once()

		u := NewUserUseCase(mockUserRepo)
		err := u.Add(context.TODO(), mockUser)

		assert.NotNil(t, err)

		mockUserRepo.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	newUUID := uuid.New()
	mockUserRepo := new(mocks.UserRepository)
	mockUser := &domain.User{
		Name:      "Cyro Dubeux",
		Email:     "xorycx@gmailcom",
		Password:  "12345678",
		UpdatedAt: time.Now(),
	}

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("Update",
			mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("*domain.UserUpdate")).
			Return(nil).Once()

		a := NewUserUseCase(mockUserRepo)
		err := a.Update(context.TODO(), newUUID, mockUser)

		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("failure", func(t *testing.T) {
		mockUserRepo.On("Update",
			mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("*domain.UserUpdate")).
			Return(errors.New("Unexpected error")).Once()

		a := NewUserUseCase(mockUserRepo)
		err := a.Update(context.TODO(), newUUID, mockUser)

		assert.NotNil(t, err)

		mockUserRepo.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	newUUID := uuid.New()
	mockUserRepo := new(mocks.UserRepository)

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("Delete",
			mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("uuid.UUID")).
			Return(nil).Once()

		u := NewUserUseCase(mockUserRepo)
		err := u.Delete(context.TODO(), newUUID)

		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("failure", func(t *testing.T) {
		mockUserRepo.On("Delete",
			mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("uuid.UUID")).
			Return(errors.New("Unexpected error")).Once()

		a := NewUserUseCase(mockUserRepo)
		err := a.Delete(context.TODO(), newUUID)

		assert.NotNil(t, err)

		mockUserRepo.AssertExpectations(t)
	})
}
