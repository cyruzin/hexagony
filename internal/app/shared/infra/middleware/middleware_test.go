package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestLoggerMiddleware(t *testing.T) {
	router := chi.NewRouter()

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	router.Use(LoggerMiddleware)

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("Hexagony v1.0")); err != nil {
			t.Error(err)
		}
	})

	router.ServeHTTP(rec, req)
}
