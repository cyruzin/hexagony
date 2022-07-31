package controller

import (
	"encoding/json"
	"errors"
	"hexagony/app/domain"
	"hexagony/libs/clog"
	"hexagony/libs/crypto"
	"hexagony/libs/rest"
	"hexagony/libs/validation"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UsersController struct {
	UsersUseCase domain.UsersUseCase
}

type createUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required,gte=8"`
}

type updateUserRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required"`
}

// FindAll godoc
// @Summary      List of users
// @Description  lists all users
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Insert your access token"  default(Bearer <Add access token here>)
// @Success      200            {object}  []domain.User
// @Failure      500            {object}  rest.Message
// @Router       /user [get]
func (u *UsersController) FindAll(w http.ResponseWriter, r *http.Request) {
	users, err := u.UsersUseCase.FindAll(r.Context())
	if err != nil {
		clog.Error(err, err.Error())
		rest.DecodeError(w, r, err, http.StatusInternalServerError)
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
// @Failure      404            {object}  rest.Message
// @Failure      500            {object}  rest.Message
// @Router       /user/{uuid} [get]
func (u *UsersController) FindByID(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		clog.Error(err, domain.ErrUsersUUIDParse.Error())
		rest.DecodeError(w, r, domain.ErrUsersFindByID, http.StatusInternalServerError)
		return
	}

	user, err := u.UsersUseCase.FindByID(r.Context(), uuid)

	exists := errors.Is(err, domain.ErrResourceNotFound)
	if exists {
		clog.Error(err, err.Error())
		rest.DecodeError(w, r, err, http.StatusNotFound)
		return
	}

	if err != nil {
		clog.Error(err, err.Error())
		rest.DecodeError(w, r, err, http.StatusInternalServerError)
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
// @Param        Authorization  header    string             true  "Insert your access token"  default(Bearer <Add access token here>)
// @Param        payload        body      createUserRequest  true  "add a new user"
// @Success      201            {object}  rest.Message
// @Failure      400            {object}  rest.Message
// @Failure      409            {object}  rest.Message
// @Failure      422            {object}  rest.Message
// @Failure      500            {object}  rest.Message
// @Router       /user [post]
func (u *UsersController) Add(w http.ResponseWriter, r *http.Request) {
	var payload createUserRequest

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		clog.Error(err, err.Error())
		rest.DecodeError(w, r, domain.ErrUsersAdd, http.StatusInternalServerError)
		return
	}

	validation := validation.New()

	if err := validation.BindStruct(r.Context(), payload); err != nil {
		validation.DecodeError(w, err)
		return
	}

	bcrypt := crypto.New()

	hashPass, err := bcrypt.HashPassword(payload.Password, 10)
	if err != nil {
		clog.Error(err, domain.ErrUsersHashPassword.Error())
		rest.DecodeError(w, r, domain.ErrUsersAdd, http.StatusUnprocessableEntity)
		return
	}

	user := domain.Users{
		UUID:      uuid.New(),
		Name:      payload.Name,
		Email:     payload.Email,
		Password:  hashPass,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = u.UsersUseCase.Add(r.Context(), &user)

	isDuplicate := errors.Is(err, domain.ErrUsersDuplicateEmail)

	if isDuplicate {
		clog.Error(err, domain.ErrUsersDuplicateEmail.Error())
		rest.DecodeError(w, r, domain.ErrUsersDuplicateEmail, http.StatusConflict)
		return
	}

	if err != nil {
		clog.Error(err, err.Error())
		rest.DecodeError(w, r, err, http.StatusInternalServerError)
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
// @Param        payload        body      updateUserRequest  true  "update an user by uuid"
// @Success      200            {object}  rest.Message
// @Failure      400            {object}  rest.Message
// @Failure      404            {object}  rest.Message
// @Failure      500            {object}  rest.Message
// @Router       /user/{uuid} [put]
func (u *UsersController) Update(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		clog.Error(err, domain.ErrUsersUUIDParse.Error())
		rest.DecodeError(w, r, domain.ErrUsersUpdate, http.StatusInternalServerError)
		return
	}

	var payload updateUserRequest

	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		clog.Error(err, err.Error())
		rest.DecodeError(w, r, err, http.StatusInternalServerError)
		return
	}

	validation := validation.New()

	if err := validation.BindStruct(r.Context(), payload); err != nil {
		clog.Error(err, domain.ErrUsersUpdate.Error())
		validation.DecodeError(w, err)
		return
	}

	user := domain.Users{
		Name:      payload.Name,
		Email:     payload.Email,
		UpdatedAt: time.Now(),
	}

	err = u.UsersUseCase.Update(r.Context(), uuid, &user)

	exists := errors.Is(err, domain.ErrResourceNotFound)
	if exists {
		clog.Error(err, err.Error())
		rest.DecodeError(w, r, err, http.StatusNotFound)
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
// @Summary      Delete an user
// @Description  delete an user by uuid
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Insert your access token"  default(Bearer <Add access token here>)
// @Param        uuid           path      string  true  "user uuid"
// @Success      200            {object}  rest.Message
// @Failure      404            {object}  rest.Message
// @Failure      500            {object}  rest.Message
// @Router       /user/{uuid} [delete]
func (u *UsersController) Delete(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		clog.Error(err, err.Error())
		rest.DecodeError(w, r, err, http.StatusInternalServerError)
		return
	}

	err = u.UsersUseCase.Delete(r.Context(), uuid)

	exists := errors.Is(err, domain.ErrResourceNotFound)
	if exists {
		clog.Error(err, err.Error())
		rest.DecodeError(w, r, err, http.StatusNotFound)
		return
	}

	if err != nil {
		clog.Error(err, err.Error())
		rest.DecodeError(w, r, err, http.StatusInternalServerError)
		return
	}

	rest.JSON(w, http.StatusOK, &rest.Message{Message: "Deleted"})
}
