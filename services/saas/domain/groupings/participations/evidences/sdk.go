package evidences

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters/clusterables"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
)

var (
	// adapter
	ErrInvalidEvidenceIdentifier    = errors.New("invalid evidence identifier")
	ErrInvalidEvidenceParticipation = errors.New("invalid evidence participation")
	ErrInvalidEvidenceParticipant   = errors.New("invalid evidence participant")
	ErrInvalidEvidenceTarget        = errors.New("invalid evidence target")
	ErrInvalidEvidencePost          = errors.New("invalid evidence post")
	ErrInvalidEvidenceScore         = errors.New("invalid evidence score")
	ErrInvalidEvidenceDetectedOn    = errors.New("invalid evidence detected on")

	// calculator
	ErrInvalidEvidenceCalculatorParticipation = errors.New("invalid evidence calculator participation")
	ErrInvalidEvidenceCalculatorParticipant   = errors.New("invalid evidence calculator participant")
	ErrInvalidEvidenceCalculatorTarget        = errors.New("invalid evidence calculator target")
	ErrInvalidEvidenceCalculatorComparable    = errors.New("invalid evidence calculator comparable")
	ErrInvalidEvidenceCalculatorVector        = errors.New("invalid evidence calculator vector")
)

// NewAdapter creates a new evidence adapter
func NewAdapter() Adapter {
	return &adapter{}
}

// NewCalculator creates a new participation evidence calculator
func NewCalculator(
	adapter Adapter,
	posts posts.Repository,
	comparables clusterables.ComparableRepository,
	threshold float64,
) Calculator {
	return &calculator{
		adapter:     adapter,
		posts:       posts,
		comparables: comparables,
		threshold:   threshold,
	}
}

// EvidenceInput represents a participation evidence input
type EvidenceInput struct {
	Identifier uuid.UUID

	Participation participations.Participation

	Participant participatables.Participatable
	Target      participatables.Participatable

	Post  posts.Post
	Score float64

	DetectedOn time.Time
}

// Adapter represents a participation evidence adapter
type Adapter interface {
	ToDomain(input EvidenceInput) (Evidence, error)
}

// Evidence represents a post that contributed to a participation score
type Evidence interface {
	Identifier() uuid.UUID

	Participation() participations.Participation

	Participant() participatables.Participatable
	Target() participatables.Participatable

	Post() posts.Post
	Score() float64

	DetectedOn() time.Time
}

// Repository represents a participation evidence repository
type Repository interface {
	Save(ctx context.Context, evidence Evidence) error

	FindByID(ctx context.Context, id uuid.UUID) (Evidence, error)

	FindByParticipation(
		ctx context.Context,
		participation uuid.UUID,
	) ([]Evidence, error)

	FindByPost(
		ctx context.Context,
		post uuid.UUID,
	) ([]Evidence, error)

	FindByParticipant(
		ctx context.Context,
		participant participatables.Participatable,
	) ([]Evidence, error)

	FindByTarget(
		ctx context.Context,
		target participatables.Participatable,
	) ([]Evidence, error)
}

// Calculator represents a participation evidence calculator
type Calculator interface {
	Calculate(
		ctx context.Context,
		participation participations.Participation,
	) ([]Evidence, error)
}
