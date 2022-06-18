package controller

import (
	"bytes"
	"encoding/json"
	"hexagony/auth/domain"
	"hexagony/auth/domain/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthenticate(t *testing.T) {
	mockAuthUseCase := new(mocks.AuthUseCase)

	token := `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJIZXhhZ29ueSIsInN1YiI
	6Imh0dHBzOi8vZ2l0aHViLmNvbS9jeXJ1emluL2hleGFnb255IiwiYXVkIjpbIkNsZWFuIEFyY2hpd
	GVjdHVyZSJdLCJleHAiOjg3MjkxNTU1MTYsImlkIjoiMTBkMTQ5NWQtODUyOS00OWJmLThlZjUtNzl
	kNzE2ZmQyMTgwIiwibmFtZSI6IkN5cm8gRHViZXV4IiwiZW1haWwiOiJ4b3J5Y3hAZ21haWwuY29tI
	n0.Nvc7gd4JoBUwjhehZ5vVub295k2846JsIsVdRGKvfvc`

	authToken := &domain.AuthToken{Token: token}

	mockAuthUseCase.
		On("Authenticate",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(authToken, nil)

	handler := AuthHandler{
		authUseCase: mockAuthUseCase,
	}

	router := chi.NewRouter()

	credentials := domain.Auth{Email: "xorycx@gmail.com", Password: "12345678"}

	payload, err := json.Marshal(credentials)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/auth", bytes.NewBuffer(payload))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	router.HandleFunc("/auth", handler.Authenticate)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	mockAuthUseCase.AssertExpectations(t)
}

func TestAuthenticateFail(t *testing.T) {
	mockAuthUseCase := new(mocks.AuthUseCase)

	mockAuthUseCase.
		On("Authenticate",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(nil, domain.ErrAuth)

	handler := AuthHandler{
		authUseCase: mockAuthUseCase,
	}

	router := chi.NewRouter()

	credentials := domain.Auth{Email: "xorycx@gmail.com", Password: "12345678"}

	payload, err := json.Marshal(credentials)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/auth", bytes.NewBuffer(payload))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	router.HandleFunc("/auth", handler.Authenticate)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

	mockAuthUseCase.AssertExpectations(t)
}
