package participations

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
)

var (
	ErrInvalidParticipationIdentifier     = errors.New("invalid participation identifier")
	ErrInvalidParticipationParticipant    = errors.New("invalid participation participant")
	ErrInvalidParticipationTarget         = errors.New("invalid participation target")
	ErrInvalidParticipationPostCount      = errors.New("invalid participation post count")
	ErrInvalidParticipationTotalPostCount = errors.New("invalid participation total post count")
	ErrInvalidParticipationPercentage     = errors.New("invalid participation percentage")
	ErrInvalidParticipationDetectedOn     = errors.New("invalid participation detected on")
)

// NewAdapter creates a new adapter
func NewAdapter() Adapter {
	return &adapter{}
}

// ParticipationInput represents a participation input
type ParticipationInput struct {
	Identifier uuid.UUID

	Participant participatables.Participatable
	Target      participatables.Participatable

	PostCount      int
	TotalPostCount int
	Percentage     float64

	DetectedOn time.Time
}

// Adapter represents a participation adapter
type Adapter interface {
	ToDomain(input ParticipationInput) (Participation, error)
}

// Participation represents how much one entity participates in another entity
type Participation interface {
	Identifier() uuid.UUID
	Participant() participatables.Participatable
	Target() participatables.Participatable
	PostCount() int
	TotalPostCount() int
	Percentage() float64
	DetectedOn() time.Time
}

// Repository represents a participation repository
type Repository interface {
	Save(ctx context.Context, participation Participation) error

	FindByID(ctx context.Context, id uuid.UUID) (Participation, error)

	FindByParticipant(
		ctx context.Context,
		participant participatables.Participatable,
	) ([]Participation, error)

	FindByTarget(
		ctx context.Context,
		target participatables.Participatable,
	) ([]Participation, error)

	FindBetween(
		ctx context.Context,
		participant participatables.Participatable,
		target participatables.Participatable,
	) (Participation, error)
}

// Calculator represents a participation calculator
type Calculator interface {
	Calculate(
		ctx context.Context,
		participant participatables.Participatable,
		target participatables.Participatable,
	) (Participation, error)
}
