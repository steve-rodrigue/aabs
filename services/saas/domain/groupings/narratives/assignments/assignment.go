package assignments

import (
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/campaigns"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/narratives"
)

type assignment struct {
	identifier uuid.UUID

	narrative narratives.Narrative
	campaign  campaigns.Campaign

	confidence float64
	assignedOn time.Time
}

func (assignment *assignment) Identifier() uuid.UUID {
	return assignment.identifier
}

func (assignment *assignment) Narrative() narratives.Narrative {
	return assignment.narrative
}

func (assignment *assignment) Campaign() campaigns.Campaign {
	return assignment.campaign
}

func (assignment *assignment) Confidence() float64 {
	return assignment.confidence
}

func (assignment *assignment) AssignedOn() time.Time {
	return assignment.assignedOn
}
