package rest

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// Router contains all routes for albums service.
func Router(h AlbumHandler) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(middleware.Logger)

	router.Route("/album", func(r chi.Router) {
		router.Get("/", h.Index)
		router.Get("/{id}", h.Show)
		router.Post("/", h.Store)
		router.Put("/{id}", h.Update)
		router.Delete("/{id}", h.Delete)
	})

	return router
}
