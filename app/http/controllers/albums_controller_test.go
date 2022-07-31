package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"hexagony/app/domain"
	"hexagony/app/domain/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAlbumsFindAll(t *testing.T) {
	now := time.Now()
	mockAlbumUseCase := new(mocks.AlbumUseCase)
	mockAlbumList := make([]*domain.Albums, 0)

	mockAlbum := domain.Albums{UUID: uuid.New(), Name: "St. Anger", Length: 75, CreatedAt: now, UpdatedAt: now}
	mockAlbumList = append(mockAlbumList, &mockAlbum)

	mockAlbumUseCase.
		On("FindAll", mock.Anything).
		Return(mockAlbumList, nil)

	handler := AlbumsController{
		AlbumsUseCase: mockAlbumUseCase,
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

func TestAlbumsFindAllFail(t *testing.T) {
	mockAlbumUseCase := new(mocks.AlbumUseCase)

	mockAlbumUseCase.
		On("FindAll", mock.Anything).
		Return(nil, domain.ErrAlbumsFindAll)

	handler := AlbumsController{
		AlbumsUseCase: mockAlbumUseCase,
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

func TestAlbumsFetchByID(t *testing.T) {
	now := time.Now()
	newUUID := uuid.New()
	mockAlbumUseCase := new(mocks.AlbumUseCase)

	mockAlbum := &domain.Albums{UUID: newUUID, Name: "St. Anger", Length: 75, CreatedAt: now, UpdatedAt: now}

	mockAlbumUseCase.
		On("FindByID", mock.Anything, mock.Anything).
		Return(mockAlbum, nil)

	handler := AlbumsController{
		AlbumsUseCase: mockAlbumUseCase,
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

func TestAlbumsFetchByIDFail(t *testing.T) {
	newUUID := uuid.New()
	mockAlbumUseCase := new(mocks.AlbumUseCase)

	mockAlbumUseCase.
		On("FindByID", mock.Anything, mock.Anything).
		Return(nil, domain.ErrAlbumsFindByID)

	handler := AlbumsController{
		AlbumsUseCase: mockAlbumUseCase,
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
		Return(nil, domain.ErrAlbumsUUIDParse)

	req, err = http.NewRequest(http.MethodGet, "/album/{uuid}", nil)
	assert.NoError(t, err)

	rec = httptest.NewRecorder()

	router.HandleFunc("/album/{uuid}", handler.FindByID)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	mockAlbumUseCase.AssertExpectations(t)
}

func TestAlbumsFetchByIDFailResource(t *testing.T) {
	newUUID := uuid.New()
	mockAlbumUseCase := new(mocks.AlbumUseCase)

	mockAlbumUseCase.
		On("FindByID", mock.Anything, mock.Anything).
		Return(nil, domain.ErrResourceNotFound)

	handler := AlbumsController{
		AlbumsUseCase: mockAlbumUseCase,
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

func TestAlbumsAdd(t *testing.T) {
	now := time.Now()
	newUUID := uuid.New()
	mockAlbumUseCase := new(mocks.AlbumUseCase)

	mockAlbum := &domain.Albums{UUID: newUUID, Name: "St. Anger", Length: 75, CreatedAt: now, UpdatedAt: now}

	mockAlbumUseCase.
		On("Add", mock.Anything, mock.Anything).
		Return(nil)

	handler := AlbumsController{
		AlbumsUseCase: mockAlbumUseCase,
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

func TestAlbumsAddFail(t *testing.T) {
	now := time.Now()
	newUUID := uuid.New()
	mockAlbumUseCase := new(mocks.AlbumUseCase)

	mockAlbum := &domain.Albums{UUID: newUUID, Name: "St. Anger", Length: 75, CreatedAt: now, UpdatedAt: now}

	mockAlbumUseCase.
		On("Add", mock.Anything, mock.Anything).
		Return(domain.ErrAlbumsAdd)

	handler := AlbumsController{
		AlbumsUseCase: mockAlbumUseCase,
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
		Return(domain.ErrAlbumsAdd)

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

func TestAlbumsUpdate(t *testing.T) {
	now := time.Now()
	newUUID := uuid.New()
	mockAlbumUseCase := new(mocks.AlbumUseCase)

	mockAlbum := &domain.Albums{Name: "St. Anger", Length: 75, UpdatedAt: now}

	mockAlbumUseCase.
		On("Update", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	handler := AlbumsController{
		AlbumsUseCase: mockAlbumUseCase,
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

func TestAlbumsUpdateFail(t *testing.T) {
	now := time.Now()
	newUUID := uuid.New()
	mockAlbumUseCase := new(mocks.AlbumUseCase)

	mockAlbum := &domain.Albums{Name: "St. Anger", Length: 75, UpdatedAt: now}

	mockAlbumUseCase.
		On("Update", mock.Anything, mock.Anything, mock.Anything).
		Return(domain.ErrAlbumsUpdate)

	handler := AlbumsController{
		AlbumsUseCase: mockAlbumUseCase,
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
		Return(nil, domain.ErrAlbumsUUIDParse)

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
		Return(domain.ErrAlbumsUpdate)

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

func TestAlbumsUpdateFailResouce(t *testing.T) {
	now := time.Now()
	newUUID := uuid.New()
	mockAlbumUseCase := new(mocks.AlbumUseCase)

	mockAlbum := &domain.Albums{Name: "St. Anger", Length: 75, UpdatedAt: now}

	mockAlbumUseCase.
		On("Update", mock.Anything, mock.Anything, mock.Anything).
		Return(domain.ErrResourceNotFound)

	handler := AlbumsController{
		AlbumsUseCase: mockAlbumUseCase,
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

func TestAlbumsDelete(t *testing.T) {
	newUUID := uuid.New()
	mockAlbumUseCase := new(mocks.AlbumUseCase)

	mockAlbumUseCase.
		On("Delete", mock.Anything, mock.Anything).
		Return(nil)

	handler := AlbumsController{
		AlbumsUseCase: mockAlbumUseCase,
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

func TestAlbumsDeleteFail(t *testing.T) {
	newUUID := uuid.New()
	mockAlbumUseCase := new(mocks.AlbumUseCase)

	mockAlbumUseCase.
		On("Delete", mock.Anything, mock.Anything).
		Return(domain.ErrAlbumsDelete)

	handler := AlbumsController{
		AlbumsUseCase: mockAlbumUseCase,
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
		Return(nil, domain.ErrAlbumsUUIDParse)

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
		Return(domain.ErrAlbumsDelete)

	mockAlbum2 := []byte(`{"id":"1"}`)

	req, err = http.NewRequest(http.MethodDelete, "/album/"+newUUID.String(), bytes.NewBuffer(mockAlbum2))
	assert.NoError(t, err)

	rec = httptest.NewRecorder()

	router.HandleFunc("/album/{uuid}", handler.Delete)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	mockAlbumUseCase.AssertExpectations(t)
}

func TestAlbumsDeleteFailResource(t *testing.T) {
	newUUID := uuid.New()
	mockAlbumUseCase := new(mocks.AlbumUseCase)

	mockAlbumUseCase.
		On("Delete", mock.Anything, mock.Anything).
		Return(domain.ErrResourceNotFound)

	handler := AlbumsController{
		AlbumsUseCase: mockAlbumUseCase,
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
