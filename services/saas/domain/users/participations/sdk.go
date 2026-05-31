package participations

import (
	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/campaigns"
)

// Participation represents how much a user participates in a campaign
type Participation interface {
	User() uuid.UUID
	Campaign() campaigns.Campaign
	PostCount() int
	TotalUserPostCount() int
	Percentage() float64
}

// Repository stores user participation in campaigns
type Repository interface {
	Save(participation Participation) error
	FindByUser(user uuid.UUID) ([]Participation, error)
	FindByCampaign(campaign uuid.UUID) ([]Participation, error)
	FindByUserAndCampaign(user uuid.UUID, campaign uuid.UUID) (Participation, error)
}

// Calculator calculates how much each user contributed to a campaign
type Calculator interface {
	Calculate(campaign campaigns.Campaign) ([]Participation, error)
}
