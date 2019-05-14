package albums

import (
	"time"

	"github.com/google/uuid"
)

// Album is used for album type.
type Album struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"string"`
	Length    time.Time `json:"length"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
