package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func TestLoggerMiddleware(t *testing.T) {
	router := chi.NewRouter()

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	router.Use(LoggerMiddleware)

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hexagony v1.0"))
	})

	router.ServeHTTP(rec, req)
}
