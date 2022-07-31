package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Albums struct {
	UUID      uuid.UUID `db:"uuid" json:"id"`
	Name      string    `db:"name" json:"name"`
	Length    int       `db:"length" json:"length"`
	CreatedAt time.Time `db:"created_at" json:"created_at" `
	UpdatedAt time.Time `db:"updated_at" json:"updated_at" `
}

type AlbumsRepository interface {
	FindAll(context.Context) ([]*Albums, error)
	FindByID(context.Context, uuid.UUID) (*Albums, error)
	Add(context.Context, *Albums) error
	Update(context.Context, uuid.UUID, *Albums) error
	Delete(context.Context, uuid.UUID) error
}

type AlbumsUseCase interface {
	FindAll(ctx context.Context) ([]*Albums, error)
	FindByID(ctx context.Context, uuid uuid.UUID) (*Albums, error)
	Add(ctx context.Context, album *Albums) error
	Update(ctx context.Context, uuid uuid.UUID, album *Albums) error
	Delete(ctx context.Context, uuid uuid.UUID) error
}
