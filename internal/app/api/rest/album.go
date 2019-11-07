package rest

import (
	"hexagony/internal/app/albums"
	"net/http"
)

type albumHandler struct {
	albumService albums.Service
}

// FindAllAlbums is responsable for find the latest albums.
func (h *albumHandler) FindAllAlbums(w http.ResponseWriter, r *http.Request) {

}

// FindAlbumByID is responsable for find an album by ID.
func (h *albumHandler) FindAlbumByID(w http.ResponseWriter, r *http.Request) {

}

// AddAlbum is responsable for add new albums.
func (h *albumHandler) AddAlbum(w http.ResponseWriter, r *http.Request) {

}

// UpdateAlbum is responsable for update an album by ID.
func (h *albumHandler) UpdateAlbum(w http.ResponseWriter, r *http.Request) {

}

// DeleteAlbum is responsable for delete an album by ID.
func (h *albumHandler) DeleteAlbum(w http.ResponseWriter, r *http.Request) {

}
