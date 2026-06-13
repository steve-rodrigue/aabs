package inputs

import (
	"time"

	"github.com/google/uuid"
)

// SaveCommunityRequest represents a save community request input
type SaveCommunityRequest struct {
	Identifier   uuid.UUID   `json:"identifier"`
	PlatformID   uuid.UUID   `json:"platform_id"`
	Handle       string      `json:"handle"`
	Title        string      `json:"title"`
	Text         string      `json:"text"`
	CreatedOn    time.Time   `json:"created_on"`
	ModeratorIDs []uuid.UUID `json:"moderator_ids"`
}
