package albums

import (
	"time"

	"github.com/google/uuid"
)

// Album is used for album type.
type Album struct {
	UUID      uuid.UUID     `db:"uuid" json:"id"`
	Name      string        `json:"name"`
	Length    time.Duration `json:"length"`
	CreatedAt time.Time     `db:"created_at" json:"created_at"`
	UpdatedAt time.Time     `db:"updated_at" json:"updated_at"`
}
