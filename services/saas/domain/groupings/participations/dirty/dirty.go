package dirty

import (
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/participatables"
)

type dirty struct {
	identifier uuid.UUID

	participant participatables.Participatable
	target      participatables.Participatable

	markedOn time.Time
}

func (dirty *dirty) Identifier() uuid.UUID {
	return dirty.identifier
}

func (dirty *dirty) Participant() participatables.Participatable {
	return dirty.participant
}

func (dirty *dirty) Target() participatables.Participatable {
	return dirty.target
}

func (dirty *dirty) MarkedOn() time.Time {
	return dirty.markedOn
}
