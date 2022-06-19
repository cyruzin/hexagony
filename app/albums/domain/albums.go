package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Album struct {
	UUID      uuid.UUID     `db:"uuid" json:"id"`
	Name      string        `db:"name" json:"name"`
	Length    time.Duration `db:"length" json:"length"`
	CreatedAt time.Time     `db:"created_at" json:"created_at" `
	UpdatedAt time.Time     `db:"updated_at" json:"updated_at" `
}

type AlbumRepository interface {
	FindAll(context.Context) ([]*Album, error)
	FindByID(context.Context, uuid.UUID) (*Album, error)
	Add(context.Context, *Album) error
	Update(context.Context, uuid.UUID, *Album) error
	Delete(context.Context, uuid.UUID) error
}

type AlbumUseCase interface {
	FindAll(ctx context.Context) ([]*Album, error)
	FindByID(ctx context.Context, uuid uuid.UUID) (*Album, error)
	Add(ctx context.Context, album *Album) error
	Update(ctx context.Context, uuid uuid.UUID, album *Album) error
	Delete(ctx context.Context, uuid uuid.UUID) error
}
