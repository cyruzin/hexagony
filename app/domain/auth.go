package domain

import (
	"context"
)

// Auth represent the auth's model.
type Auth struct {
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

// AuthToken represent the token payload.
type AuthToken struct {
	Token string `json:"token,omitempty"`
}

// AuthRepository represent the auth's repository contract.
type AuthRepository interface {
	Authenticate(ctx context.Context, email string) (*Users, error)
}

// AuthUsecase represent the auth's usecases.
type AuthUseCase interface {
	Authenticate(ctx context.Context, email, password string) (*AuthToken, error)
}
