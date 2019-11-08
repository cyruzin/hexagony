package albums

import (
	"context"

	"github.com/google/uuid"
)

// Repository provides access to albums repository.
type Repository interface {
	FindAll(context.Context) ([]*Album, error)
	FindByID(context.Context, uuid.UUID) (*Album, error)
	Add(context.Context, *Album) error
	Update(context.Context, uuid.UUID, *Album) error
	Delete(context.Context, uuid.UUID) error
}
