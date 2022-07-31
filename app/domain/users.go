package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Users struct {
	UUID      uuid.UUID `db:"uuid" json:"id"`
	Name      string    `db:"name" json:"name"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password" json:"password"`
	CreatedAt time.Time `db:"created_at" json:"created_at" `
	UpdatedAt time.Time `db:"updated_at" json:"updated_at" `
}

type UsersList struct {
	UUID      uuid.UUID `db:"uuid" json:"id"`
	Name      string    `db:"name" json:"name"`
	Email     string    `db:"email" json:"email"`
	CreatedAt time.Time `db:"created_at" json:"created_at" `
	UpdatedAt time.Time `db:"updated_at" json:"updated_at" `
}

type UsersRepository interface {
	FindAll(context.Context) ([]*UsersList, error)
	FindByID(context.Context, uuid.UUID) (*UsersList, error)
	Add(context.Context, *Users) error
	Update(context.Context, uuid.UUID, *Users) error
	Delete(context.Context, uuid.UUID) error
}

type UsersUseCase interface {
	FindAll(ctx context.Context) ([]*UsersList, error)
	FindByID(ctx context.Context, uuid uuid.UUID) (*UsersList, error)
	Add(ctx context.Context, user *Users) error
	Update(ctx context.Context, uuid uuid.UUID, user *Users) error
	Delete(ctx context.Context, uuid uuid.UUID) error
}
