package controller

import (
	"encoding/json"
	"errors"
	"hexagony/app/domain"
	"hexagony/libs/clog"
	"hexagony/libs/rest"
	"hexagony/libs/validation"
	"net/http"
)

type AuthController struct {
	AuthUseCase domain.AuthUseCase
}

type authRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password,omitempty" validate:"required,gte=8"`
}

// Auth godoc
// @Summary      Authenticate a user
// @Description  authenticate a user and returns a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        payload  body      authRequest  true  "authenticates the user"
// @Success      200      {object}  domain.AuthToken
// @Failure      400      {object}  rest.Message
// @Failure      401      {object}  rest.Message
// @Failure      500      {object}  rest.Message
// @Router       /auth [post]
func (a *AuthController) Authenticate(w http.ResponseWriter, r *http.Request) {
	var payload authRequest

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		clog.Error(err, domain.ErrAuth.Error())
		rest.DecodeError(w, r, domain.ErrAuth, http.StatusInternalServerError)
		return
	}

	validation := validation.New()

	if err := validation.BindStruct(r.Context(), payload); err != nil {
		validation.DecodeError(w, err)
		return
	}

	user := domain.Auth{
		Email:    payload.Email,
		Password: payload.Password,
	}

	res, err := a.AuthUseCase.Authenticate(r.Context(), user.Email, user.Password)

	userNotFound := errors.Is(err, domain.ErrAuthUserNotFound)
	if userNotFound {
		clog.Error(err, domain.ErrAuthUserNotFound.Error())
		rest.DecodeError(w, r, domain.ErrAuthUserNotFound, http.StatusUnauthorized)
		return
	}

	passwordMatch := errors.Is(err, domain.ErrAuthPassword)
	if passwordMatch {
		clog.Error(err, domain.ErrAuthPassword.Error())
		rest.DecodeError(w, r, domain.ErrAuthPassword, http.StatusUnauthorized)
		return
	}

	if err != nil {
		clog.Error(err, err.Error())
		rest.DecodeError(w, r, domain.ErrAuth, http.StatusInternalServerError)
		return
	}

	rest.JSON(w, http.StatusOK, &res)
}
