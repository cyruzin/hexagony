package domain

import (
	"context"
	"hexagony/app/users/domain"

	"github.com/google/uuid"
)

// Auth represent the auth's model.
type Auth struct {
	UUID     uuid.UUID `db:"uuid" json:"id,omitempty"`
	Name     string    `json:"name,omitempty"`
	Email    string    `json:"email" validate:"required,email"`
	Password string    `json:"password,omitempty" validate:"required,gte=8"`
}

// AuthToken represent the token payload.
type AuthToken struct {
	Token string `json:"token,omitempty"`
}

// AuthRepository represent the auth's repository contract.
type AuthRepository interface {
	Authenticate(ctx context.Context, email string) (*domain.User, error)
}

// AuthUsecase represent the auth's usecases.
type AuthUseCase interface {
	Authenticate(ctx context.Context, email, password string) (*AuthToken, error)
}
