package inputs

import (
	"time"

	"github.com/google/uuid"
)

// SavePlatformRequest represents a save platform request input
type SavePlatformRequest struct {
	Identifier uuid.UUID `json:"identifier"`
	Name       string    `json:"name"`
	Handle     string    `json:"handle"`
	BaseURL    string    `json:"base_url"`
	CreatedOn  time.Time `json:"created_on"`
}
