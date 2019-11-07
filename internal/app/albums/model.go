package albums

import (
	"time"

	"github.com/google/uuid"
)

// Album is used for album type.
type Album struct {
	ID        uuid.UUID     `json:"id"`
	Name      string        `json:"string"`
	Length    time.Duration `json:"length"`
	CreatedAt int64         `db:"created_at" json:"created_at"`
	UpdatedAt int64         `db:"updated_at" json:"updated_at"`
}
