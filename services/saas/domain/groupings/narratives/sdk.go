package narratives

import (
	"time"

	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
)

// Narrative represents a narrative
type Narrative interface {
	Identifier() uuid.UUID
	ParticipationKind() participatables.Kind
	Cluster() clusters.Cluster
	Name() string
	Description() string
	CreatedOn() time.Time
}

// Repository represents a narrative repository
type Repository interface {
	Save(narrative Narrative) error

	FindByID(id uuid.UUID) (Narrative, error)
	FindByName(name string) (Narrative, error)

	Find(index int, amount int) ([]Narrative, error)
	FindAfter(cursor uuid.UUID, amount int) ([]Narrative, error)

	Count() (int64, error)
}
