package participations

import (
	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/campaigns"
	"github.com/steve-rodrigue/aabs/services/saas/domain/communities"
)

// Participation represents a campaign participation
type Participation interface {
	Identifier() uuid.UUID
	Community() communities.Community
	Campaign() campaigns.Campaign
	PostCount() int
	TotalCommunityPostCount() int
	Percentage() float64
}

// Repository represents a participation repository
type Repository interface {
	Save(participation Participation) error
	FindByCommunity(community uuid.UUID) ([]Participation, error)
	FindByCampaign(campaign uuid.UUID) ([]Participation, error)
}

// Calculator represents a participation calculator
type Calculator interface {
	CalculateCommunityParticipation(campaign campaigns.Campaign) ([]Participation, error)
}
