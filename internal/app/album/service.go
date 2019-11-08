package albums

import (
	"context"

	"github.com/google/uuid"
)

// Service provides albums operations.
type Service interface {
	FindAll(ctx context.Context) ([]*Album, error)
	FindByID(ctx context.Context, uuid uuid.UUID) (*Album, error)
	Add(ctx context.Context, album *Album) error
	Update(ctx context.Context, uuid uuid.UUID, album *Album) error
	Delete(ctx context.Context, uuid uuid.UUID) error
}
