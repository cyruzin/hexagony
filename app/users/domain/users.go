package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type User struct {
	UUID      uuid.UUID `db:"uuid" json:"id"`
	Name      string    `db:"name" json:"name"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password" json:"password"`
	CreatedAt time.Time `db:"created_at" json:"created_at" `
	UpdatedAt time.Time `db:"updated_at" json:"updated_at" `
}

type UserRepository interface {
	FindAll(context.Context) ([]*User, error)
	FindByID(context.Context, uuid.UUID) (*User, error)
	Add(context.Context, *User) error
	Update(context.Context, uuid.UUID, *User) error
	Delete(context.Context, uuid.UUID) error
}

type UserUseCase interface {
	FindAll(ctx context.Context) ([]*User, error)
	FindByID(ctx context.Context, uuid uuid.UUID) (*User, error)
	Add(ctx context.Context, user *User) error
	Update(ctx context.Context, uuid uuid.UUID, user *User) error
	Delete(ctx context.Context, uuid uuid.UUID) error
}
