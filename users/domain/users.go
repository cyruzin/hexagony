package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type User struct {
	UUID      uuid.UUID `db:"uuid" json:"id"`
	Name      string    `json:"name" validate:"required"`
	Email     string    `json:"email" validate:"required"`
	Password  string    `json:"password" validate:"required,gte=8"`
	Token     string    `json:"token,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at" `
	UpdatedAt time.Time `db:"updated_at" json:"updated_at" `
}

type UserUpdate struct {
	UUID      uuid.UUID `db:"uuid" json:"id"`
	Name      string    `json:"name" validate:"required"`
	Email     string    `json:"email" validate:"required"`
	Password  string    `json:"password" validate:"gte=8"`
	CreatedAt time.Time `db:"created_at" json:"created_at" `
	UpdatedAt time.Time `db:"updated_at" json:"updated_at" `
}

type UserRepository interface {
	FindAll(context.Context) ([]*User, error)
	FindByID(context.Context, uuid.UUID) (*User, error)
	Add(context.Context, *User) error
	Update(context.Context, uuid.UUID, *UserUpdate) error
	Delete(context.Context, uuid.UUID) error
}

type UserUseCase interface {
	FindAll(ctx context.Context) ([]*User, error)
	FindByID(ctx context.Context, uuid uuid.UUID) (*User, error)
	Add(ctx context.Context, user *User) error
	Update(ctx context.Context, uuid uuid.UUID, user *UserUpdate) error
	Delete(ctx context.Context, uuid uuid.UUID) error
}
