package assignments

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/campaigns"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters/clusterables"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/narratives"
)

var (
	// adapter
	ErrInvalidAssignmentIdentifier = errors.New("invalid assignment identifier")
	ErrInvalidAssignmentNarrative  = errors.New("invalid assignment narrative")
	ErrInvalidAssignmentCampaign   = errors.New("invalid assignment campaign")
	ErrInvalidAssignmentConfidence = errors.New("invalid assignment confidence")
	ErrInvalidAssignmentAssignedOn = errors.New("invalid assignment assigned on")

	// assigner
	ErrInvalidAssignmentAssignerCampaign   = errors.New("invalid assignment assigner campaign")
	ErrInvalidAssignmentAssignerNarrative  = errors.New("invalid assignment assigner narrative")
	ErrInvalidAssignmentAssignerComparable = errors.New("invalid assignment assigner comparable")
	ErrInvalidAssignmentAssignerVector     = errors.New("invalid assignment assigner vector")
)

// NewAdapter creates a new assignment adapter
func NewAdapter() Adapter {
	return &adapter{}
}

// NewAssigner creates a new assignment assigner
func NewAssigner(
	adapter Adapter,
	comparables clusterables.ComparableRepository,
	threshold float64,
) Assigner {
	return &assigner{
		adapter:     adapter,
		comparables: comparables,
		threshold:   threshold,
	}
}

// AssignmentInput represents an assignment input
type AssignmentInput struct {
	Identifier uuid.UUID

	Narrative narratives.Narrative
	Campaign  campaigns.Campaign

	Confidence float64
	AssignedOn time.Time
}

// Adapter represents an assignment adapter
type Adapter interface {
	ToDomain(input AssignmentInput) (Assignment, error)
}

// Assignment represents an assignment between a campaign and a narrative
type Assignment interface {
	Identifier() uuid.UUID

	Narrative() narratives.Narrative
	Campaign() campaigns.Campaign

	Confidence() float64
	AssignedOn() time.Time
}

// Repository represents an assignment repository
type Repository interface {
	Save(ctx context.Context, assignment Assignment) error

	FindByID(ctx context.Context, id uuid.UUID) (Assignment, error)

	FindByNarrative(
		ctx context.Context,
		narrative uuid.UUID,
	) ([]Assignment, error)

	FindByCampaign(
		ctx context.Context,
		campaign uuid.UUID,
	) ([]Assignment, error)

	FindBetween(
		ctx context.Context,
		narrative uuid.UUID,
		campaign uuid.UUID,
	) (Assignment, error)

	Find(ctx context.Context, index int, amount int) ([]Assignment, error)
	FindAfter(ctx context.Context, cursor uuid.UUID, amount int) ([]Assignment, error)
	Count(ctx context.Context) (int64, error)
}

// Assigner detects which narratives belong to a campaign
type Assigner interface {
	Assign(
		ctx context.Context,
		campaign campaigns.Campaign,
		narratives []narratives.Narrative,
	) ([]Assignment, error)
}
