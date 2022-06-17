package controller

import (
	"encoding/json"
	"hexagony/lib/clog"
	"hexagony/lib/crypto"
	"hexagony/lib/rest"
	"hexagony/lib/validation"
	cmiddleware "hexagony/shared/infra/middleware"
	"hexagony/users/domain"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UserHandler struct {
	userUseCase domain.UserUseCase
}

func NewUserHandler(c *chi.Mux, as domain.UserUseCase) {
	handler := UserHandler{userUseCase: as}

	c.Route("/user", func(r chi.Router) {
		r.Use(cmiddleware.AuthMiddleware)

		r.Get("/", handler.FindAll)
		r.Get("/{uuid}", handler.FindByID)
		r.Post("/", handler.Add)
		r.Put("/{uuid}", handler.Update)
		r.Delete("/{uuid}", handler.Delete)
	})
}

// FindAll godoc
// @Summary      List of users
// @Description  lists all users
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Insert your access token"  default(Bearer <Add access token here>)
// @Success      200            {object}  []domain.User
// @Failure      422            {object}  rest.Message
// @Router       /user [get]
func (u *UserHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	users, err := u.userUseCase.FindAll(r.Context())
	if err != nil {
		clog.Error(err, domain.ErrFindAll.Error())
		rest.DecodeError(w, r, domain.ErrFindAll, http.StatusUnprocessableEntity)
		return
	}

	rest.JSON(w, http.StatusOK, &users)
}

// FindByID godoc
// @Summary      List an user
// @Description  lists an user by uuid
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Insert your access token"  default(Bearer <Add access token here>)
// @Param        uuid           path      string  true  "user uuid"
// @Success      200            {object}  domain.User
// @Failure      422            {object}  rest.Message
// @Router       /user/{uuid} [get]
func (u *UserHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		clog.Error(err, domain.ErrUUIDParse.Error())
		rest.DecodeError(w, r, domain.ErrUUIDParse, http.StatusUnprocessableEntity)
		return
	}

	user, err := u.userUseCase.FindByID(r.Context(), uuid)
	if err != nil {
		clog.Error(err, domain.ErrFindByID.Error())
		rest.DecodeError(w, r, domain.ErrFindByID, http.StatusUnprocessableEntity)
		return
	}

	rest.JSON(w, http.StatusOK, user)
}

// Add godoc
// @Summary      Add an user
// @Description  add a new user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string       true  "Insert your access token"  default(Bearer <Add access token here>)
// @Param        payload        body      domain.User  true  "add a new user"
// @Success      201            {object}  rest.Message
// @Failure      422            {object}  rest.Message
// @Router       /user [post]
func (u *UserHandler) Add(w http.ResponseWriter, r *http.Request) {
	var user domain.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		clog.Error(err, domain.ErrAdd.Error())
		rest.DecodeError(w, r, domain.ErrAdd, http.StatusUnprocessableEntity)
		return
	}

	validation := validation.New()

	if err := validation.Bind(r.Context(), user); err != nil {
		validation.DecodeError(w, err)
		return
	}

	user.UUID = uuid.New()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	bcrypt := crypto.New()

	hashPass, err := bcrypt.HashPassword(user.Password, 10)
	if err != nil {
		clog.Error(err, domain.ErrAdd.Error())
		rest.DecodeError(w, r, domain.ErrAdd, http.StatusUnprocessableEntity)
		return
	}

	user.Password = hashPass

	err = u.userUseCase.Add(r.Context(), &user)
	if err != nil {
		clog.Error(err, domain.ErrAdd.Error())
		rest.DecodeError(w, r, domain.ErrAdd, http.StatusUnprocessableEntity)
		return
	}

	rest.JSON(w, http.StatusCreated, &rest.Message{Message: "Created"})
}

// Update godoc
// @Summary      Update an user
// @Description  update an user by uuid
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string             true  "Insert your access token"  default(Bearer <Add access token here>)
// @Param        uuid           path      string             true  "user uuid"
// @Param        payload        body      domain.UserUpdate  true  "update an user by uuid"
// @Success      200            {object}  rest.Message
// @Failure      422            {object}  rest.Message
// @Router       /user/{uuid} [put]
func (u *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		clog.Error(err, domain.ErrUUIDParse.Error())
		rest.DecodeError(w, r, domain.ErrUUIDParse, http.StatusUnprocessableEntity)
		return
	}

	var user domain.UserUpdate

	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		clog.Error(err, domain.ErrUpdate.Error())
		rest.DecodeError(w, r, domain.ErrUpdate, http.StatusUnprocessableEntity)
		return
	}

	validation := validation.New()

	if err := validation.Bind(r.Context(), user); err != nil {
		clog.Error(err, domain.ErrUpdate.Error())
		validation.DecodeError(w, err)
		return
	}

	user.UpdatedAt = time.Now()

	bcrypt := crypto.New()

	hashPass, err := bcrypt.HashPassword(user.Password, 10)
	if err != nil {
		clog.Error(err, domain.ErrAdd.Error())
		rest.DecodeError(w, r, domain.ErrAdd, http.StatusUnprocessableEntity)
		return
	}

	user.Password = hashPass

	err = u.userUseCase.Update(r.Context(), uuid, &user)
	if err != nil {
		clog.Error(err, domain.ErrUpdate.Error())
		rest.DecodeError(w, r, domain.ErrUpdate, http.StatusUnprocessableEntity)
		return
	}

	rest.JSON(w, http.StatusOK, &rest.Message{Message: "Updated"})
}

// Update godoc
// @Summary      Delete an user
// @Description  delete an user by uuid
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Insert your access token"  default(Bearer <Add access token here>)
// @Param        uuid           path      string  true  "user uuid"
// @Success      200            {object}  rest.Message
// @Failure      422            {object}  rest.Message
// @Router       /user/{uuid} [delete]
func (u *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		clog.Error(err, domain.ErrDelete.Error())
		rest.DecodeError(w, r, domain.ErrDelete, http.StatusUnprocessableEntity)
		return
	}

	err = u.userUseCase.Delete(r.Context(), uuid)
	if err != nil {
		clog.Error(err, domain.ErrDelete.Error())
		rest.DecodeError(w, r, domain.ErrDelete, http.StatusUnprocessableEntity)
		return
	}

	rest.JSON(w, http.StatusOK, &rest.Message{Message: "Deleted"})
}
