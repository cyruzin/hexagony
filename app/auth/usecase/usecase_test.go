package usecase

import (
	"context"
	"errors"
	"hexagony/app/auth/domain/mocks"
	domainUsers "hexagony/app/users/domain"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthenticate(t *testing.T) {
	mockAuthRepo := new(mocks.AuthRepository)

	mockUser := &domainUsers.User{
		UUID:      uuid.New(),
		Name:      "Cyro Dubeux",
		Email:     "xorycx@gmail.com",
		Password:  "$2a$10$Vm8jmbPV5NMgoCag3O/iM.LTfMs6rmmwgDwRUw9m8QGFyis7EA/Gy",
		Token:     "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("success", func(t *testing.T) {
		mockAuthRepo.On("Authenticate",
			mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("string")).
			Return(mockUser, nil).
			Once()

		a := NewAuthUsecase(mockAuthRepo)
		_, err := a.Authenticate(context.TODO(), "xorycx@gmail.com", "12345678")

		assert.NoError(t, err)

		mockAuthRepo.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockAuthRepo.On("Authenticate",
			mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string")).
			Return(nil, errors.New("Unexpected error")).
			Once()

		a := NewAuthUsecase(mockAuthRepo)
		token, err := a.Authenticate(context.TODO(), "xorycx@gmail.com", "12345678")

		assert.Nil(t, token)
		assert.NotNil(t, err)

		mockAuthRepo.AssertExpectations(t)
	})
}
