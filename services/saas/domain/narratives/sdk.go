package narratives

import (
	"time"

	"github.com/google/uuid"
)

// Narrative represents a narrative
type Narrative interface {
	Identifier() uuid.UUID
	Name() string
	Description() string
	CreatedOn() time.Time
}

// Repository represents a narrative repository
type Repository interface {
	Save(narrative Narrative) error
	FindByID(id uuid.UUID) (Narrative, error)
}
