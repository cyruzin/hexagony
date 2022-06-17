package domain

import (
	"context"
	"hexagony/users/domain"

	"github.com/google/uuid"
)

// Auth represent the auth's model.
type Auth struct {
	UUID     uuid.UUID `db:"uuid" json:"id,omitempty"`
	Name     string    `json:"name,omitempty"`
	Email    string    `json:"email" validate:"required,email"`
	Password string    `json:"password,omitempty" validate:"required,gte=8"`
	Token    string    `json:"token,omitempty"`
}

// AuthToken represent the token payload.
type AuthToken struct {
	Token string `json:"token,omitempty"`
}

// AuthUsecase represent the auth's usecases.
type AuthUsecase interface {
	Authenticate(ctx context.Context, email, password string) (*AuthToken, error)
}

// AuthRepository represent the auth's repository contract.
type AuthRepository interface {
	Authenticate(ctx context.Context, email string, password string) (*domain.User, error)
}
