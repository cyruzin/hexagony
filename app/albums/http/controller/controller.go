package controller

import (
	"encoding/json"
	"hexagony/app/albums/domain"
	cmiddleware "hexagony/app/shared/http/middleware"
	"hexagony/lib/clog"
	"hexagony/lib/rest"
	"hexagony/lib/validation"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type AlbumHandler struct {
	albumUseCase domain.AlbumUseCase
}

func NewAlbumHandler(c *chi.Mux, as domain.AlbumUseCase) {
	handler := AlbumHandler{albumUseCase: as}

	c.Route("/album", func(r chi.Router) {
		r.Use(cmiddleware.AuthMiddleware)

		r.Get("/", handler.FindAll)
		r.Get("/{uuid}", handler.FindByID)
		r.Post("/", handler.Add)
		r.Put("/{uuid}", handler.Update)
		r.Delete("/{uuid}", handler.Delete)
	})
}

type albumRequest struct {
	Name   string        `json:"name" validate:"required"`
	Length time.Duration `json:"length" validate:"required" swaggertype:"primitive,integer"`
}

// FindAll godoc
// @Summary      List of albums
// @Description  lists all albums
// @Tags         album
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Insert your access token"  default(Bearer <Add access token here>)
// @Success      200            {object}  []domain.Album
// @Failure      422            {object}  rest.Message
// @Router       /album [get]
func (a *AlbumHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	albums, err := a.albumUseCase.FindAll(r.Context())
	if err != nil {
		clog.Error(err, domain.ErrFindAll.Error())
		rest.DecodeError(w, r, domain.ErrFindAll, http.StatusUnprocessableEntity)
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
// @Failure      422            {object}  rest.Message
// @Router       /album/{uuid} [get]
func (a *AlbumHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		clog.Error(err, domain.ErrUUIDParse.Error())
		rest.DecodeError(w, r, domain.ErrUUIDParse, http.StatusUnprocessableEntity)
		return
	}

	album, err := a.albumUseCase.FindByID(r.Context(), uuid)
	if err != nil {
		clog.Error(err, domain.ErrFindByID.Error())
		rest.DecodeError(w, r, domain.ErrFindByID, http.StatusUnprocessableEntity)
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
// @Failure      422            {object}  rest.Message
// @Router       /album [post]
func (a *AlbumHandler) Add(w http.ResponseWriter, r *http.Request) {
	var payload albumRequest

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		clog.Error(err, domain.ErrAdd.Error())
		rest.DecodeError(w, r, domain.ErrAdd, http.StatusUnprocessableEntity)
		return
	}

	validation := validation.New()

	if err := validation.Bind(r.Context(), payload); err != nil {
		validation.DecodeError(w, err)
		return
	}

	album := domain.Album{
		UUID:      uuid.New(),
		Name:      payload.Name,
		Length:    payload.Length,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = a.albumUseCase.Add(r.Context(), &album)
	if err != nil {
		clog.Error(err, domain.ErrAdd.Error())
		rest.DecodeError(w, r, domain.ErrAdd, http.StatusUnprocessableEntity)
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
// @Failure      422            {object}  rest.Message
// @Router       /album/{uuid} [put]
func (a *AlbumHandler) Update(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		clog.Error(err, domain.ErrUUIDParse.Error())
		rest.DecodeError(w, r, domain.ErrUUIDParse, http.StatusUnprocessableEntity)
		return
	}

	var payload albumRequest

	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		clog.Error(err, domain.ErrUpdate.Error())
		rest.DecodeError(w, r, domain.ErrUpdate, http.StatusUnprocessableEntity)
		return
	}

	validation := validation.New()

	if err := validation.Bind(r.Context(), payload); err != nil {
		clog.Error(err, domain.ErrUpdate.Error())
		validation.DecodeError(w, err)
		return
	}

	album := domain.Album{
		Name:      payload.Name,
		Length:    payload.Length,
		UpdatedAt: time.Now(),
	}

	err = a.albumUseCase.Update(r.Context(), uuid, &album)
	if err != nil {
		clog.Error(err, domain.ErrUpdate.Error())
		rest.DecodeError(w, r, domain.ErrUpdate, http.StatusUnprocessableEntity)
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
// @Failure      422            {object}  rest.Message
// @Router       /album/{uuid} [delete]
func (a *AlbumHandler) Delete(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		clog.Error(err, domain.ErrDelete.Error())
		rest.DecodeError(w, r, domain.ErrDelete, http.StatusUnprocessableEntity)
		return
	}

	err = a.albumUseCase.Delete(r.Context(), uuid)
	if err != nil {
		clog.Error(err, domain.ErrDelete.Error())
		rest.DecodeError(w, r, domain.ErrDelete, http.StatusUnprocessableEntity)
		return
	}

	rest.JSON(w, http.StatusOK, &rest.Message{Message: "Deleted"})
}
