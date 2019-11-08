package rest

import (
	"hexagony/internal/app/albums"
	"net/http"
)

// AlbumHandler interface for the album handlers.
type AlbumHandler interface {
	Index(http.ResponseWriter, *http.Request)
	Show(http.ResponseWriter, *http.Request)
	Store(http.ResponseWriter, *http.Request)
	Update(http.ResponseWriter, *http.Request)
	Delete(http.ResponseWriter, *http.Request)
}

type handler struct {
	albumService albums.Service
}

// NewHandler will instantiate the handlers.
func NewHandler(albumService albums.Service) AlbumHandler {
	return &handler{albumService}
}

// Index is responsable for find the latest albums.
func (h *handler) Index(w http.ResponseWriter, r *http.Request) {

}

// Show is responsable for find an album by ID.
func (h *handler) Show(w http.ResponseWriter, r *http.Request) {

}

// Store is responsable for add new albums.
func (h *handler) Store(w http.ResponseWriter, r *http.Request) {

}

// Update is responsable for update an album by ID.
func (h *handler) Update(w http.ResponseWriter, r *http.Request) {

}

// Delete is responsable for delete an album by ID.
func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {

}
