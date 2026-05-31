package assignments

import (
	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/campaigns"
	"github.com/steve-rodrigue/aabs/services/saas/domain/narratives"
)

// Assignment represents an assignment
type Assignment interface {
	Identifier() uuid.UUID
	Narrative() narratives.Narrative
	Campaign() campaigns.Campaign
	Confidence() float64
}
