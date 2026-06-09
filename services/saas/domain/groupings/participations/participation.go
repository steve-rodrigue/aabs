package participations

import (
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
)

type participation struct {
	identifier uuid.UUID

	participant participatables.Participatable
	target      participatables.Participatable

	postCount      int
	totalPostCount int
	percentage     float64

	detectedOn time.Time
}

func (participation *participation) Identifier() uuid.UUID {
	return participation.identifier
}

func (participation *participation) Participant() participatables.Participatable {
	return participation.participant
}

func (participation *participation) Target() participatables.Participatable {
	return participation.target
}

func (participation *participation) PostCount() int {
	return participation.postCount
}

func (participation *participation) TotalPostCount() int {
	return participation.totalPostCount
}

func (participation *participation) Percentage() float64 {
	return participation.percentage
}

func (participation *participation) DetectedOn() time.Time {
	return participation.detectedOn
}
