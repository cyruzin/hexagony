package controller

import (
	"encoding/json"
	"hexagony/internal/app/domain"
	"hexagony/internal/app/pkg/rest"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type AlbumHandler struct {
	albumUseCase domain.AlbumUseCase
}

func NewAlbumHandler(c *chi.Mux, as domain.AlbumUseCase) {
	handler := AlbumHandler{albumUseCase: as}

	c.Route("/album", func(r chi.Router) {
		r.Get("/", handler.FindAll)
		r.Get("/{uuid}", handler.FindByID)
		r.Post("/", handler.Add)
		r.Put("/{uuid}", handler.Update)
		r.Delete("/{uuid}", handler.Delete)
	})
}

func (a *AlbumHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	albums, err := a.albumUseCase.FindAll(r.Context())
	if err != nil {
		rest.DecodeError(w, r, domain.ErrFindAll, http.StatusUnprocessableEntity)
		return
	}

	rest.EncodeJSON(w, http.StatusOK, &albums)
}

func (a *AlbumHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		rest.DecodeError(w, r, domain.ErrUUIDParse, http.StatusUnprocessableEntity)
		return
	}

	album, err := a.albumUseCase.FindByID(r.Context(), uuid)
	if err != nil {
		rest.DecodeError(w, r, domain.ErrFindByID, http.StatusUnprocessableEntity)
		return
	}

	rest.EncodeJSON(w, http.StatusOK, album)
}

func (a *AlbumHandler) Add(w http.ResponseWriter, r *http.Request) {
	var album domain.Album

	err := json.NewDecoder(r.Body).Decode(&album)
	if err != nil {
		rest.DecodeError(w, r, domain.ErrAdd, http.StatusUnprocessableEntity)
		return
	}

	album.UUID = uuid.New()
	album.CreatedAt = time.Now()
	album.UpdatedAt = time.Now()

	err = a.albumUseCase.Add(r.Context(), &album)
	if err != nil {
		rest.DecodeError(w, r, domain.ErrAdd, http.StatusUnprocessableEntity)
		return
	}

	rest.EncodeJSON(w, http.StatusOK, &rest.APIMessage{Message: "Created"})
}

func (a *AlbumHandler) Update(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		rest.DecodeError(w, r, domain.ErrUUIDParse, http.StatusUnprocessableEntity)
		return
	}

	var album domain.Album

	err = json.NewDecoder(r.Body).Decode(&album)
	if err != nil {
		rest.DecodeError(w, r, domain.ErrUpdate, http.StatusUnprocessableEntity)
		return
	}

	album.UpdatedAt = time.Now()

	err = a.albumUseCase.Update(r.Context(), uuid, &album)
	if err != nil {
		rest.DecodeError(w, r, domain.ErrUpdate, http.StatusUnprocessableEntity)
		return
	}

	rest.EncodeJSON(w, http.StatusOK, &rest.APIMessage{Message: "Updated"})
}

func (a *AlbumHandler) Delete(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		rest.DecodeError(w, r, domain.ErrDelete, http.StatusUnprocessableEntity)
		return
	}

	err = a.albumUseCase.Delete(r.Context(), uuid)
	if err != nil {
		rest.DecodeError(w, r, domain.ErrDelete, http.StatusUnprocessableEntity)
		return
	}

	rest.EncodeJSON(w, http.StatusOK, &rest.APIMessage{Message: "Deleted"})
}
