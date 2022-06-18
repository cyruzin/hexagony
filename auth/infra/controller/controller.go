package controller

import (
	"encoding/json"
	"hexagony/auth/domain"
	"hexagony/lib/clog"
	"hexagony/lib/rest"
	"hexagony/lib/validation"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type AuthHandler struct {
	authUseCase domain.AuthUseCase
}

func NewAuthHandler(c *chi.Mux, auc domain.AuthUseCase) {
	handler := AuthHandler{authUseCase: auc}

	c.Post("/auth", handler.Authenticate)
}

// Auth godoc
// @Summary      Authenticate a user
// @Description  authenticate a user and returns a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        payload  body      domain.User  true  "authenticates the user"
// @Success      200      {object}  domain.User
// @Failure      422      {object}  rest.Message
// @Router       /auth [post]
func (a *AuthHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	var user domain.Auth

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		clog.Error(err, domain.ErrAuth.Error())
		rest.DecodeError(w, r, domain.ErrAuth, http.StatusUnprocessableEntity)
		return
	}

	validation := validation.New()

	if err := validation.Bind(r.Context(), user); err != nil {
		validation.DecodeError(w, err)
		return
	}

	res, err := a.authUseCase.Authenticate(r.Context(), user.Email, user.Password)
	if err != nil {
		clog.Error(err, err.Error())
		rest.DecodeError(w, r, domain.ErrAuth, http.StatusUnprocessableEntity)
		return
	}

	rest.JSON(w, http.StatusOK, &res)
}
