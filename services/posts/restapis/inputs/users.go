package inputs

import (
	"time"

	"github.com/google/uuid"
)

// SaveUserRequest represents a save platform request input
type SaveUserRequest struct {
	Identifier  uuid.UUID `json:"identifier"`
	PlatformID  uuid.UUID `json:"platform_id"`
	ExternalID  string    `json:"external_id"`
	Handle      string    `json:"handle"`
	DisplayName string    `json:"display_name"`
	ProfileURL  string    `json:"profile_url"`
	CreatedOn   time.Time `json:"created_on"`
}
