package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"hexagony/internal/albums/domain"
	"hexagony/internal/albums/domain/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewAlbumHandler(t *testing.T) {
	router := chi.NewRouter()

	mockAlbumUseCase := new(mocks.AlbumUseCase)

	NewAlbumHandler(router, mockAlbumUseCase)
}

func TestFindAll(t *testing.T) {
	now := time.Now()
	mockAlbumUseCase := new(mocks.AlbumUseCase)
	mockAlbumList := make([]*domain.Album, 0)

	mockAlbum := domain.Album{UUID: uuid.New(), Name: "St. Anger", Length: 75, CreatedAt: now, UpdatedAt: now}
	mockAlbumList = append(mockAlbumList, &mockAlbum)

	mockAlbumUseCase.
		On("FindAll", mock.Anything).
		Return(mockAlbumList, nil)

	handler := AlbumHandler{
		albumUseCase: mockAlbumUseCase,
	}

	router := chi.NewRouter()

	req, err := http.NewRequest(http.MethodGet, "/album", nil)
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	router.HandleFunc("/album", handler.FindAll)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	mockAlbumUseCase.AssertExpectations(t)
}

func TestFindAllFail(t *testing.T) {
	mockAlbumUseCase := new(mocks.AlbumUseCase)

	mockAlbumUseCase.
		On("FindAll", mock.Anything).
		Return(nil, domain.ErrFindAll)

	handler := AlbumHandler{
		albumUseCase: mockAlbumUseCase,
	}

	router := chi.NewRouter()

	req, err := http.NewRequest(http.MethodGet, "/album", nil)
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	router.HandleFunc("/album", handler.FindAll)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	mockAlbumUseCase.AssertExpectations(t)
}

func TestFetchByID(t *testing.T) {
	now := time.Now()
	newUUID := uuid.New()
	mockAlbumUseCase := new(mocks.AlbumUseCase)

	mockAlbum := &domain.Album{UUID: newUUID, Name: "St. Anger", Length: 75, CreatedAt: now, UpdatedAt: now}

	mockAlbumUseCase.
		On("FindByID", mock.Anything, mock.Anything).
		Return(mockAlbum, nil)

	handler := AlbumHandler{
		albumUseCase: mockAlbumUseCase,
	}

	router := chi.NewRouter()

	req, err := http.NewRequest(http.MethodGet, "/album/"+newUUID.String(), nil)
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	router.HandleFunc("/album/{uuid}", handler.FindByID)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	mockAlbumUseCase.AssertExpectations(t)
}

func TestFetchByIDFail(t *testing.T) {
	newUUID := uuid.New()
	mockAlbumUseCase := new(mocks.AlbumUseCase)

	mockAlbumUseCase.
		On("FindByID", mock.Anything, mock.Anything).
		Return(nil, domain.ErrFindByID)

	handler := AlbumHandler{
		albumUseCase: mockAlbumUseCase,
	}

	router := chi.NewRouter()

	req, err := http.NewRequest(http.MethodGet, "/album/"+newUUID.String(), nil)
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	router.HandleFunc("/album/{uuid}", handler.FindByID)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	mockAlbumUseCase.AssertExpectations(t)

	// err uuid parsing

	mockAlbumUseCase.
		On("FindByID", mock.Anything, mock.Anything).
		Return(nil, domain.ErrUUIDParse)

	req, err = http.NewRequest(http.MethodGet, "/album/{uuid}", nil)
	assert.NoError(t, err)

	rec = httptest.NewRecorder()

	router.HandleFunc("/album/{uuid}", handler.FindByID)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	mockAlbumUseCase.AssertExpectations(t)
}

func TestFetchByIDFailResource(t *testing.T) {
	newUUID := uuid.New()
	mockAlbumUseCase := new(mocks.AlbumUseCase)

	mockAlbumUseCase.
		On("FindByID", mock.Anything, mock.Anything).
		Return(nil, domain.ErrResourceNotFound)

	handler := AlbumHandler{
		albumUseCase: mockAlbumUseCase,
	}

	router := chi.NewRouter()

	req, err := http.NewRequest(http.MethodGet, "/album/"+newUUID.String(), nil)
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	router.HandleFunc("/album/{uuid}", handler.FindByID)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)

	mockAlbumUseCase.AssertExpectations(t)
}

func TestAdd(t *testing.T) {
	now := time.Now()
	newUUID := uuid.New()
	mockAlbumUseCase := new(mocks.AlbumUseCase)

	mockAlbum := &domain.Album{UUID: newUUID, Name: "St. Anger", Length: 75, CreatedAt: now, UpdatedAt: now}

	mockAlbumUseCase.
		On("Add", mock.Anything, mock.Anything).
		Return(nil)

	handler := AlbumHandler{
		albumUseCase: mockAlbumUseCase,
	}

	router := chi.NewRouter()

	payload, err := json.Marshal(mockAlbum)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/album", bytes.NewBuffer(payload))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	router.HandleFunc("/album", handler.Add)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	mockAlbumUseCase.AssertExpectations(t)
}

func TestAddFail(t *testing.T) {
	now := time.Now()
	newUUID := uuid.New()
	mockAlbumUseCase := new(mocks.AlbumUseCase)

	mockAlbum := &domain.Album{UUID: newUUID, Name: "St. Anger", Length: 75, CreatedAt: now, UpdatedAt: now}

	mockAlbumUseCase.
		On("Add", mock.Anything, mock.Anything).
		Return(domain.ErrAdd)

	handler := AlbumHandler{
		albumUseCase: mockAlbumUseCase,
	}

	router := chi.NewRouter()

	payload, err := json.Marshal(mockAlbum)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/album", bytes.NewBuffer(payload))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	router.HandleFunc("/album", handler.Add)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	mockAlbumUseCase.AssertExpectations(t)

	// err decoding json

	mockAlbumUseCase.
		On("Add", mock.Anything, mock.Anything).
		Return(domain.ErrAdd)

	mockAlbum2 := []byte(`{"id:""1"}`)

	req, err = http.NewRequest(http.MethodPost, "/album", bytes.NewBuffer(mockAlbum2))
	assert.NoError(t, err)

	rec = httptest.NewRecorder()

	router.HandleFunc("/album", handler.Add)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	mockAlbumUseCase.AssertExpectations(t)

	// validation errors

	mockAlbumUseCase.
		On("Add", mock.Anything, mock.Anything).
		Return(errors.New("the length field is required"))

	mockAlbum3 := []byte(`{"name":"Nova Era"}`)

	req, err = http.NewRequest(http.MethodPost, "/album", bytes.NewBuffer(mockAlbum3))
	assert.NoError(t, err)

	rec = httptest.NewRecorder()

	router.HandleFunc("/album", handler.Add)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	mockAlbumUseCase.AssertExpectations(t)
}

func TestUpdate(t *testing.T) {
	now := time.Now()
	newUUID := uuid.New()
	mockAlbumUseCase := new(mocks.AlbumUseCase)

	mockAlbum := &domain.Album{Name: "St. Anger", Length: 75, UpdatedAt: now}

	mockAlbumUseCase.
		On("Update", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	handler := AlbumHandler{
		albumUseCase: mockAlbumUseCase,
	}

	router := chi.NewRouter()

	payload, err := json.Marshal(mockAlbum)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPut, "/album/"+newUUID.String(), bytes.NewBuffer(payload))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	router.HandleFunc("/album/{uuid}", handler.Update)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	mockAlbumUseCase.AssertExpectations(t)
}

func TestUpdateFail(t *testing.T) {
	now := time.Now()
	newUUID := uuid.New()
	mockAlbumUseCase := new(mocks.AlbumUseCase)

	mockAlbum := &domain.Album{Name: "St. Anger", Length: 75, UpdatedAt: now}

	mockAlbumUseCase.
		On("Update", mock.Anything, mock.Anything, mock.Anything).
		Return(domain.ErrUpdate)

	handler := AlbumHandler{
		albumUseCase: mockAlbumUseCase,
	}

	router := chi.NewRouter()

	payload, err := json.Marshal(mockAlbum)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPut, "/album/"+newUUID.String(), bytes.NewBuffer(payload))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	router.HandleFunc("/album/{uuid}", handler.Update)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	mockAlbumUseCase.AssertExpectations(t)

	// err uuid parsing

	mockAlbumUseCase.
		On("Update", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, domain.ErrUUIDParse)

	req, err = http.NewRequest(http.MethodPut, "/album/{uuid}", bytes.NewBuffer(payload))
	assert.NoError(t, err)

	rec = httptest.NewRecorder()

	router.HandleFunc("/album/{uuid}", handler.Update)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	mockAlbumUseCase.AssertExpectations(t)

	// err decoding json

	mockAlbumUseCase.
		On("Update", mock.Anything, mock.Anything, mock.Anything).
		Return(domain.ErrUpdate)

	mockAlbum2 := []byte(`{"id:""1"}`)

	req, err = http.NewRequest(http.MethodPut, "/album/"+newUUID.String(), bytes.NewBuffer(mockAlbum2))
	assert.NoError(t, err)

	rec = httptest.NewRecorder()

	router.HandleFunc("/album/{uuid}", handler.Update)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	mockAlbumUseCase.AssertExpectations(t)

	// validation errors

	mockAlbumUseCase.
		On("Update", mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("the length field is required"))

	mockAlbum3 := []byte(`{"name":"Nova Era"}`)

	req, err = http.NewRequest(http.MethodPut, "/album/"+newUUID.String(), bytes.NewBuffer(mockAlbum3))
	assert.NoError(t, err)

	rec = httptest.NewRecorder()

	router.HandleFunc("/album/{uuid}", handler.Update)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	mockAlbumUseCase.AssertExpectations(t)
}

func TestUpdateFailResouce(t *testing.T) {
	now := time.Now()
	newUUID := uuid.New()
	mockAlbumUseCase := new(mocks.AlbumUseCase)

	mockAlbum := &domain.Album{Name: "St. Anger", Length: 75, UpdatedAt: now}

	mockAlbumUseCase.
		On("Update", mock.Anything, mock.Anything, mock.Anything).
		Return(domain.ErrResourceNotFound)

	handler := AlbumHandler{
		albumUseCase: mockAlbumUseCase,
	}

	router := chi.NewRouter()

	payload, err := json.Marshal(mockAlbum)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPut, "/album/"+newUUID.String(), bytes.NewBuffer(payload))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	router.HandleFunc("/album/{uuid}", handler.Update)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)

	mockAlbumUseCase.AssertExpectations(t)
}

func TestDelete(t *testing.T) {
	newUUID := uuid.New()
	mockAlbumUseCase := new(mocks.AlbumUseCase)

	mockAlbumUseCase.
		On("Delete", mock.Anything, mock.Anything).
		Return(nil)

	handler := AlbumHandler{
		albumUseCase: mockAlbumUseCase,
	}

	router := chi.NewRouter()

	req, err := http.NewRequest(http.MethodDelete, "/album/"+newUUID.String(), nil)
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	router.HandleFunc("/album/{uuid}", handler.Delete)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	mockAlbumUseCase.AssertExpectations(t)
}

func TestDeleteFail(t *testing.T) {
	newUUID := uuid.New()
	mockAlbumUseCase := new(mocks.AlbumUseCase)

	mockAlbumUseCase.
		On("Delete", mock.Anything, mock.Anything).
		Return(domain.ErrDelete)

	handler := AlbumHandler{
		albumUseCase: mockAlbumUseCase,
	}

	router := chi.NewRouter()

	req, err := http.NewRequest(http.MethodDelete, "/album/"+newUUID.String(), nil)
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	router.HandleFunc("/album/{uuid}", handler.Delete)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	mockAlbumUseCase.AssertExpectations(t)

	// err uuid parsing

	mockAlbumUseCase.
		On("Delete", mock.Anything, mock.Anything).
		Return(nil, domain.ErrUUIDParse)

	req, err = http.NewRequest(http.MethodDelete, "/album/{uuid}", nil)
	assert.NoError(t, err)

	rec = httptest.NewRecorder()

	router.HandleFunc("/album/{uuid}", handler.Delete)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	mockAlbumUseCase.AssertExpectations(t)

	// err decoding json

	mockAlbumUseCase.
		On("Delete", mock.Anything, mock.Anything).
		Return(domain.ErrDelete)

	mockAlbum2 := []byte(`{"id":"1"}`)

	req, err = http.NewRequest(http.MethodDelete, "/album/"+newUUID.String(), bytes.NewBuffer(mockAlbum2))
	assert.NoError(t, err)

	rec = httptest.NewRecorder()

	router.HandleFunc("/album/{uuid}", handler.Delete)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	mockAlbumUseCase.AssertExpectations(t)
}

func TestDeleteFailResource(t *testing.T) {
	newUUID := uuid.New()
	mockAlbumUseCase := new(mocks.AlbumUseCase)

	mockAlbumUseCase.
		On("Delete", mock.Anything, mock.Anything).
		Return(domain.ErrResourceNotFound)

	handler := AlbumHandler{
		albumUseCase: mockAlbumUseCase,
	}

	router := chi.NewRouter()

	req, err := http.NewRequest(http.MethodDelete, "/album/"+newUUID.String(), nil)
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	router.HandleFunc("/album/{uuid}", handler.Delete)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)

	mockAlbumUseCase.AssertExpectations(t)
}
