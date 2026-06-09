package participations

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
)

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
