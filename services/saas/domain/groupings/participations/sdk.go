package participations

import (
	"time"

	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
)

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
	Save(participation Participation) error
	FindByID(id uuid.UUID) (Participation, error)
	FindByParticipant(participant participatables.Participatable) ([]Participation, error)
	FindByTarget(target participatables.Participatable) ([]Participation, error)
	FindBetween(
		participant participatables.Participatable,
		target participatables.Participatable,
	) (Participation, error)
}

// Calculator represents a participation calculator
type Calculator interface {
	Calculate(
		participant participatables.Participatable,
		target participatables.Participatable,
	) (Participation, error)
}
