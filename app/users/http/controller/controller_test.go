package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"hexagony/app/users/domain"
	"hexagony/app/users/domain/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewUserHandler(t *testing.T) {
	router := chi.NewRouter()

	mockUserUseCase := new(mocks.UserUseCase)

	NewUserHandler(router, mockUserUseCase)
}

func TestFindAll(t *testing.T) {
	now := time.Now()
	mockUserUseCase := new(mocks.UserUseCase)
	mockUserList := make([]*domain.User, 0)

	mockUser := domain.User{
		UUID:      uuid.New(),
		Name:      "Cyro Dubeux",
		Email:     "xorycX@gmail.com",
		Password:  "12345678",
		CreatedAt: now,
		UpdatedAt: now,
	}
	mockUserList = append(mockUserList, &mockUser)

	mockUserUseCase.
		On("FindAll", mock.Anything).
		Return(mockUserList, nil)

	handler := UserHandler{
		userUseCase: mockUserUseCase,
	}

	router := chi.NewRouter()

	req, err := http.NewRequest(http.MethodGet, "/user", nil)
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	router.HandleFunc("/user", handler.FindAll)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	mockUserUseCase.AssertExpectations(t)
}

func TestFindAllFail(t *testing.T) {
	mockUserUseCase := new(mocks.UserUseCase)

	mockUserUseCase.
		On("FindAll", mock.Anything).
		Return(nil, domain.ErrFindAll)

	handler := UserHandler{
		userUseCase: mockUserUseCase,
	}

	router := chi.NewRouter()

	req, err := http.NewRequest(http.MethodGet, "/user", nil)
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	router.HandleFunc("/user", handler.FindAll)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

	mockUserUseCase.AssertExpectations(t)
}

func TestFetchByID(t *testing.T) {
	now := time.Now()
	newUUID := uuid.New()
	mockUserUseCase := new(mocks.UserUseCase)

	mockUser := &domain.User{
		UUID:      newUUID,
		Name:      "Cyro Dubeux",
		Email:     "xorycx@gmail.com",
		Password:  "12345678",
		CreatedAt: now,
		UpdatedAt: now,
	}

	mockUserUseCase.
		On("FindByID", mock.Anything, mock.Anything).
		Return(mockUser, nil)

	handler := UserHandler{
		userUseCase: mockUserUseCase,
	}

	router := chi.NewRouter()

	req, err := http.NewRequest(http.MethodGet, "/user/"+newUUID.String(), nil)
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	router.HandleFunc("/user/{uuid}", handler.FindByID)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	mockUserUseCase.AssertExpectations(t)
}

func TestFetchByIDFail(t *testing.T) {
	newUUID := uuid.New()
	mockUserUseCase := new(mocks.UserUseCase)

	mockUserUseCase.
		On("FindByID", mock.Anything, mock.Anything).
		Return(nil, domain.ErrFindByID)

	handler := UserHandler{
		userUseCase: mockUserUseCase,
	}

	router := chi.NewRouter()

	req, err := http.NewRequest(http.MethodGet, "/user/"+newUUID.String(), nil)
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	router.HandleFunc("/user/{uuid}", handler.FindByID)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

	mockUserUseCase.AssertExpectations(t)

	// err uuid parsing

	mockUserUseCase.
		On("FindByID", mock.Anything, mock.Anything).
		Return(nil, domain.ErrUUIDParse)

	req, err = http.NewRequest(http.MethodGet, "/user/{uuid}", nil)
	assert.NoError(t, err)

	rec = httptest.NewRecorder()

	router.HandleFunc("/user/{uuid}", handler.FindByID)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

	mockUserUseCase.AssertExpectations(t)
}

func TestAdd(t *testing.T) {
	now := time.Now()
	newUUID := uuid.New()
	mockUserUseCase := new(mocks.UserUseCase)

	mockUser := &domain.User{
		UUID:      newUUID,
		Name:      "Cyro Dubeux",
		Email:     "xorycx@gmail.com",
		Password:  "12345678",
		CreatedAt: now,
		UpdatedAt: now,
	}

	mockUserUseCase.
		On("Add", mock.Anything, mock.Anything).
		Return(nil)

	handler := UserHandler{
		userUseCase: mockUserUseCase,
	}

	router := chi.NewRouter()

	payload, err := json.Marshal(mockUser)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(payload))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	router.HandleFunc("/user", handler.Add)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	mockUserUseCase.AssertExpectations(t)
}

func TestAddFail(t *testing.T) {
	mockUserUseCase := new(mocks.UserUseCase)

	mockUserUseCase.
		On("Add", mock.Anything, mock.Anything).
		Return(domain.ErrAdd)

	handler := UserHandler{
		userUseCase: mockUserUseCase,
	}

	router := chi.NewRouter()

	payload := []byte(`{
		"name": "Cyro Dubeux",
		"email": "xorycx@gmail.com", 
		"password": "12345678"
		}`)

	req, err := http.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(payload))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	router.HandleFunc("/user", handler.Add)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

	mockUserUseCase.AssertExpectations(t)

	// err decoding json

	mockUserUseCase.
		On("Add", mock.Anything, mock.Anything).
		Return(domain.ErrAdd)

	mockUser2 := []byte(`{"id:""1"}`)

	req, err = http.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(mockUser2))
	assert.NoError(t, err)

	rec = httptest.NewRecorder()

	router.HandleFunc("/user", handler.Add)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

	mockUserUseCase.AssertExpectations(t)

	// // validation errors

	mockUserUseCase.
		On("Add", mock.Anything, mock.Anything).
		Return(errors.New("the length field is required"))

	mockUser3 := []byte(`{"name":"Cyro Dubeux"}`)

	req, err = http.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(mockUser3))
	assert.NoError(t, err)

	rec = httptest.NewRecorder()

	router.HandleFunc("/user", handler.Add)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	mockUserUseCase.AssertExpectations(t)
}

func TestUpdate(t *testing.T) {
	now := time.Now()
	newUUID := uuid.New()
	mockUserUseCase := new(mocks.UserUseCase)

	mockUser := &domain.User{
		Name:      "Cyro Dubeux",
		Email:     "xorycx@gmail.com",
		Password:  "12345678",
		UpdatedAt: now,
	}

	mockUserUseCase.
		On("Update", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	handler := UserHandler{
		userUseCase: mockUserUseCase,
	}

	router := chi.NewRouter()

	payload, err := json.Marshal(mockUser)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPut, "/user/"+newUUID.String(), bytes.NewBuffer(payload))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	router.HandleFunc("/user/{uuid}", handler.Update)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	mockUserUseCase.AssertExpectations(t)
}

func TestUpdateFail(t *testing.T) {
	now := time.Now()
	newUUID := uuid.New()
	mockUserUseCase := new(mocks.UserUseCase)

	mockUser := &domain.User{
		Name:      "Cyro Dubeux",
		Email:     "xorycx@gmail.com",
		Password:  "12345678",
		UpdatedAt: now,
	}

	mockUserUseCase.
		On("Update", mock.Anything, mock.Anything, mock.Anything).
		Return(domain.ErrUpdate)

	handler := UserHandler{
		userUseCase: mockUserUseCase,
	}

	router := chi.NewRouter()

	payload, err := json.Marshal(mockUser)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPut, "/user/"+newUUID.String(), bytes.NewBuffer(payload))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	router.HandleFunc("/user/{uuid}", handler.Update)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

	mockUserUseCase.AssertExpectations(t)

	// err uuid parsing

	mockUserUseCase.
		On("Update", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, domain.ErrUUIDParse)

	req, err = http.NewRequest(http.MethodPut, "/user/{uuid}", bytes.NewBuffer(payload))
	assert.NoError(t, err)

	rec = httptest.NewRecorder()

	router.HandleFunc("/user/{uuid}", handler.Update)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

	mockUserUseCase.AssertExpectations(t)

	// err decoding json

	mockUserUseCase.
		On("Update", mock.Anything, mock.Anything, mock.Anything).
		Return(domain.ErrUpdate)

	mockUser2 := []byte(`{"id:""1"}`)

	req, err = http.NewRequest(http.MethodPut, "/user/"+newUUID.String(), bytes.NewBuffer(mockUser2))
	assert.NoError(t, err)

	rec = httptest.NewRecorder()

	router.HandleFunc("/user/{uuid}", handler.Update)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

	mockUserUseCase.AssertExpectations(t)

	// validation errors

	mockUserUseCase.
		On("Update", mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("the length field is required"))

	mockUser3 := []byte(`{"name":"Cyro Dubeux"}`)

	req, err = http.NewRequest(http.MethodPut, "/user/"+newUUID.String(), bytes.NewBuffer(mockUser3))
	assert.NoError(t, err)

	rec = httptest.NewRecorder()

	router.HandleFunc("/user/{uuid}", handler.Update)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	mockUserUseCase.AssertExpectations(t)
}

func TestDelete(t *testing.T) {
	newUUID := uuid.New()
	mockUserUseCase := new(mocks.UserUseCase)

	mockUserUseCase.
		On("Delete", mock.Anything, mock.Anything).
		Return(nil)

	handler := UserHandler{
		userUseCase: mockUserUseCase,
	}

	router := chi.NewRouter()

	req, err := http.NewRequest(http.MethodDelete, "/user/"+newUUID.String(), nil)
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	router.HandleFunc("/user/{uuid}", handler.Delete)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	mockUserUseCase.AssertExpectations(t)
}

func TestDeleteFail(t *testing.T) {
	newUUID := uuid.New()
	mockUserUseCase := new(mocks.UserUseCase)

	mockUserUseCase.
		On("Delete", mock.Anything, mock.Anything).
		Return(domain.ErrDelete)

	handler := UserHandler{
		userUseCase: mockUserUseCase,
	}

	router := chi.NewRouter()

	req, err := http.NewRequest(http.MethodDelete, "/user/"+newUUID.String(), nil)
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	router.HandleFunc("/user/{uuid}", handler.Delete)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

	mockUserUseCase.AssertExpectations(t)

	// err uuid parsing

	mockUserUseCase.
		On("Delete", mock.Anything, mock.Anything).
		Return(nil, domain.ErrUUIDParse)

	req, err = http.NewRequest(http.MethodDelete, "/user/{uuid}", nil)
	assert.NoError(t, err)

	rec = httptest.NewRecorder()

	router.HandleFunc("/user/{uuid}", handler.Delete)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

	mockUserUseCase.AssertExpectations(t)

	// err decoding json

	mockUserUseCase.
		On("Delete", mock.Anything, mock.Anything).
		Return(domain.ErrDelete)

	mockUser2 := []byte(`{"id":"1"}`)

	req, err = http.NewRequest(http.MethodDelete, "/user/"+newUUID.String(), bytes.NewBuffer(mockUser2))
	assert.NoError(t, err)

	rec = httptest.NewRecorder()

	router.HandleFunc("/user/{uuid}", handler.Delete)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

	mockUserUseCase.AssertExpectations(t)
}
