package narratives

import (
	"time"

	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters"
)

// Narrative represents a narrative
type Narrative interface {
	Identifier() uuid.UUID
	Cluster() clusters.Cluster
	Name() string
	Description() string
	CreatedOn() time.Time
}

// Repository represents a narrative repository
type Repository interface {
	Save(narrative Narrative) error
	FindByID(id uuid.UUID) (Narrative, error)
}
