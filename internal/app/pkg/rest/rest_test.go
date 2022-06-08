package rest

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestDecodeError(t *testing.T) {
	router := chi.NewRouter()

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		DecodeError(w, r, errors.New("testing decode error func..."), http.StatusBadRequest)
	})

	router.ServeHTTP(rec, req)
}

func TestEncodeJSON(t *testing.T) {
	router := chi.NewRouter()

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	albums := []byte(`{"uuid": "b6548fd6-23e8-41f0-af74-6a45f7c930fa", "name":"St. Anger", "length": 75}`)

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		EncodeJSON(w, http.StatusOK, albums)
	})

	router.ServeHTTP(rec, req)
}
