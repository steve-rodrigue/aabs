package narratives

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
)

var (
	ErrInvalidNarrativeIdentifier        = errors.New("invalid narrative identifier")
	ErrInvalidNarrativeParticipationKind = errors.New("invalid narrative participation kind")
	ErrInvalidNarrativeCluster           = errors.New("invalid narrative cluster")
	ErrInvalidNarrativeName              = errors.New("invalid narrative name")
	ErrInvalidNarrativeDescription       = errors.New("invalid narrative description")
	ErrInvalidNarrativeCreatedOn         = errors.New("invalid narrative created on")
)

// NewAdapter creates a new narrative adapter
func NewAdapter() Adapter {
	return &adapter{}
}

// NarrativeInput represents a narrative input
type NarrativeInput struct {
	Identifier uuid.UUID

	ParticipationKind participatables.Kind

	Cluster clusters.Cluster

	Name        string
	Description string

	CreatedOn time.Time
}

// Adapter represents a narrative adapter
type Adapter interface {
	ToDomain(input NarrativeInput) (Narrative, error)
}

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
	Save(
		ctx context.Context,
		narrative Narrative,
	) error

	FindByID(
		ctx context.Context,
		id uuid.UUID,
	) (Narrative, error)

	FindByName(
		ctx context.Context,
		name string,
	) (Narrative, error)

	Find(
		ctx context.Context,
		index int,
		amount int,
	) ([]Narrative, error)

	FindAfter(
		ctx context.Context,
		cursor uuid.UUID,
		amount int,
	) ([]Narrative, error)

	Count(
		ctx context.Context,
	) (int64, error)
}
