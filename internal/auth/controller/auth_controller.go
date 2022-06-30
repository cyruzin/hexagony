package controller

import (
	"encoding/json"
	"hexagony/internal/auth/domain"
	"hexagony/pkg/clog"
	"hexagony/pkg/rest"
	"hexagony/pkg/validation"
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
// @Failure      422      {object}  rest.Message
// @Failure      400      {object}  rest.Message
// @Failure      500      {object}  rest.Message
// @Router       /auth [post]
func (a *AuthHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
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

	res, err := a.authUseCase.Authenticate(r.Context(), user.Email, user.Password)
	if err != nil {
		clog.Error(err, err.Error())
		rest.DecodeError(w, r, domain.ErrAuth, http.StatusUnprocessableEntity)
		return
	}

	rest.JSON(w, http.StatusOK, &res)
}
