package rest

import (
	"hexagony/internal/app/albums"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

const (
	errFindAll   = "Couldn't list the albums"
	errFindByID  = "Couldn't get the album"
	errAdd       = "Couldn't insert the album"
	errUpdate    = "Couldn't update the album"
	errDelete    = "Couldn't delete the album"
	errUUIDParse = "Couldn't parse the UUID"
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

// Index finds the latest albums.
func (h *handler) Index(w http.ResponseWriter, r *http.Request) {
	albums, err := h.albumService.FindAll(r.Context())
	if err != nil {
		InvalidRequest(w, r, err, errFindAll, http.StatusUnprocessableEntity)
		return
	}

	ToJSON(w, http.StatusOK, &albums)
}

// Show finds an album by ID.
func (h *handler) Show(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		InvalidRequest(w, r, err, errUUIDParse, http.StatusUnprocessableEntity)
		return
	}

	album, err := h.albumService.FindByID(r.Context(), uuid)
	if err != nil {
		InvalidRequest(w, r, err, errFindByID, http.StatusUnprocessableEntity)
		return
	}

	ToJSON(w, http.StatusOK, album)
}

// Store adds new albums.
func (h *handler) Store(w http.ResponseWriter, r *http.Request) {
	var album albums.Album

	err := json.NewDecoder(r.Body).Decode(&album)
	if err != nil {
		InvalidRequest(w, r, err, errAdd, http.StatusUnprocessableEntity)
		return
	}

	album.UUID = uuid.New()
	album.CreatedAt = time.Now()
	album.UpdatedAt = time.Now()

	err = h.albumService.Add(r.Context(), &album)
	if err != nil {
		InvalidRequest(w, r, err, errAdd, http.StatusUnprocessableEntity)
		return
	}

	ToJSON(w, http.StatusOK, &APIMessage{Message: "Created"})
}

// Update updates an album by ID.
func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		InvalidRequest(w, r, err, errUUIDParse, http.StatusUnprocessableEntity)
		return
	}

	var album albums.Album

	err = json.NewDecoder(r.Body).Decode(&album)
	if err != nil {
		InvalidRequest(w, r, err, errUpdate, http.StatusUnprocessableEntity)
		return
	}

	album.UpdatedAt = time.Now()

	err = h.albumService.Update(r.Context(), uuid, &album)
	if err != nil {
		InvalidRequest(w, r, err, errUpdate, http.StatusUnprocessableEntity)
		return
	}

	ToJSON(w, http.StatusOK, &APIMessage{Message: "Updated"})
}

// Delete deletes an album by ID.
func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		InvalidRequest(w, r, err, errUUIDParse, http.StatusUnprocessableEntity)
		return
	}

	err = h.albumService.Delete(r.Context(), uuid)
	if err != nil {
		InvalidRequest(w, r, err, errDelete, http.StatusUnprocessableEntity)
		return
	}

	ToJSON(w, http.StatusOK, &APIMessage{Message: "Deleted"})
}
