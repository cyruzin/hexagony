package controller

import (
	"encoding/json"
	"errors"
	"hexagony/app/domain"
	"hexagony/libs/clog"
	"hexagony/libs/rest"
	"hexagony/libs/validation"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type AlbumsController struct {
	AlbumsUseCase domain.AlbumsUseCase
}

type albumRequest struct {
	Name   string `json:"name" validate:"required"`
	Length int    `json:"length" validate:"required"`
}

// FindAll godoc
// @Summary      List of albums
// @Description  lists all albums
// @Tags         album
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Insert your access token"  default(Bearer <Add access token here>)
// @Success      200            {object}  []domain.Album
// @Failure      500            {object}  rest.Message
// @Router       /album [get]
func (a *AlbumsController) FindAll(w http.ResponseWriter, r *http.Request) {
	albums, err := a.AlbumsUseCase.FindAll(r.Context())
	if err != nil {
		clog.Error(err, domain.ErrAlbumsFindAll.Error())
		rest.DecodeError(w, r, domain.ErrAlbumsFindAll, http.StatusInternalServerError)
		return
	}

	rest.JSON(w, http.StatusOK, &albums)
}

// FindByID godoc
// @Summary      List an album
// @Description  lists an album by uuid
// @Tags         album
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Insert your access token"  default(Bearer <Add access token here>)
// @Param        uuid           path      string  true  "album uuid"
// @Success      200            {object}  domain.Album
// @Failure      404            {object}  rest.Message
// @Failure      500            {object}  rest.Message
// @Router       /album/{uuid} [get]
func (a *AlbumsController) FindByID(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		clog.Error(err, domain.ErrAlbumsUUIDParse.Error())
		rest.DecodeError(w, r, domain.ErrAlbumsFindByID, http.StatusInternalServerError)
		return
	}

	album, err := a.AlbumsUseCase.FindByID(r.Context(), uuid)

	exists := errors.Is(err, domain.ErrResourceNotFound)
	if exists {
		clog.Error(err, domain.ErrResourceNotFound.Error())
		rest.DecodeError(w, r, domain.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	if err != nil {
		clog.Error(err, err.Error())
		rest.DecodeError(w, r, err, http.StatusInternalServerError)
		return
	}

	rest.JSON(w, http.StatusOK, album)
}

// Add godoc
// @Summary      Add an album
// @Description  add a new album
// @Tags         album
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string        true  "Insert your access token"  default(Bearer <Add access token here>)
// @Param        payload        body      albumRequest  true  "add a new album"
// @Success      201            {object}  rest.Message
// @Failure      400            {object}  rest.Message
// @Failure      500            {object}  rest.Message
// @Router       /album [post]
func (a *AlbumsController) Add(w http.ResponseWriter, r *http.Request) {
	var payload albumRequest

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		clog.Error(err, domain.ErrAlbumsAdd.Error())
		rest.DecodeError(w, r, domain.ErrAlbumsAdd, http.StatusInternalServerError)
		return
	}

	validation := validation.New()

	if err := validation.BindStruct(r.Context(), payload); err != nil {
		validation.DecodeError(w, err)
		return
	}

	album := domain.Albums{
		UUID:      uuid.New(),
		Name:      payload.Name,
		Length:    payload.Length,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = a.AlbumsUseCase.Add(r.Context(), &album)
	if err != nil {
		clog.Error(err, err.Error())
		rest.DecodeError(w, r, err, http.StatusInternalServerError)
		return
	}

	rest.JSON(w, http.StatusCreated, &rest.Message{Message: "Created"})
}

// Update godoc
// @Summary      Update an album
// @Description  update an album by uuid
// @Tags         album
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string        true  "Insert your access token"  default(Bearer <Add access token here>)
// @Param        uuid           path      string        true  "album uuid"
// @Param        payload        body      albumRequest  true  "update an album by uuid"
// @Success      200            {object}  rest.Message
// @Failure      404            {object}  rest.Message
// @Failure      500            {object}  rest.Message
// @Router       /album/{uuid} [put]
func (a *AlbumsController) Update(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		clog.Error(err, domain.ErrAlbumsUUIDParse.Error())
		rest.DecodeError(w, r, domain.ErrAlbumsUpdate, http.StatusInternalServerError)
		return
	}

	var payload albumRequest

	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		clog.Error(err, err.Error())
		rest.DecodeError(w, r, err, http.StatusInternalServerError)
		return
	}

	validation := validation.New()

	if err := validation.BindStruct(r.Context(), payload); err != nil {
		clog.Error(err, domain.ErrAlbumsUpdate.Error())
		validation.DecodeError(w, err)
		return
	}

	album := domain.Albums{
		Name:      payload.Name,
		Length:    payload.Length,
		UpdatedAt: time.Now(),
	}

	err = a.AlbumsUseCase.Update(r.Context(), uuid, &album)

	exists := errors.Is(err, domain.ErrResourceNotFound)
	if exists {
		clog.Error(err, domain.ErrResourceNotFound.Error())
		rest.DecodeError(w, r, domain.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	if err != nil {
		clog.Error(err, err.Error())
		rest.DecodeError(w, r, err, http.StatusInternalServerError)
		return
	}

	rest.JSON(w, http.StatusOK, &rest.Message{Message: "Updated"})
}

// Update godoc
// @Summary      Delete an album
// @Description  delete an album by uuid
// @Tags         album
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Insert your access token"  default(Bearer <Add access token here>)
// @Param        uuid           path      string  true  "album uuid"
// @Success      200            {object}  rest.Message
// @Failure      404            {object}  rest.Message
// @Failure      500            {object}  rest.Message
// @Router       /album/{uuid} [delete]
func (a *AlbumsController) Delete(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		clog.Error(err, domain.ErrAlbumsUUIDParse.Error())
		rest.DecodeError(w, r, domain.ErrAlbumsDelete, http.StatusInternalServerError)
		return
	}

	err = a.AlbumsUseCase.Delete(r.Context(), uuid)

	exists := errors.Is(err, domain.ErrResourceNotFound)
	if exists {
		clog.Error(err, domain.ErrResourceNotFound.Error())
		rest.DecodeError(w, r, domain.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	if err != nil {
		clog.Error(err, err.Error())
		rest.DecodeError(w, r, err, http.StatusInternalServerError)
		return
	}

	rest.JSON(w, http.StatusOK, &rest.Message{Message: "Deleted"})
}
