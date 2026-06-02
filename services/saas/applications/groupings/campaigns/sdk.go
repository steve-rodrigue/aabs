package campaigns

import (
	"github.com/google/uuid"

	app_participations "github.com/steve-rodrigue/aabs/services/saas/applications/groupings/participations"
	app_posts "github.com/steve-rodrigue/aabs/services/saas/applications/posts"
	"github.com/steve-rodrigue/aabs/services/saas/domain/communities"
	domain_campaigns "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/campaigns"
	"github.com/steve-rodrigue/aabs/services/saas/domain/platforms"
	"github.com/steve-rodrigue/aabs/services/saas/domain/users"
)

// New creates a new campaign application
func New(
	repository domain_campaigns.Repository,
	posts app_posts.Application,
	participations app_participations.Application,
	classifier domain_campaigns.Classifier,
) Application {
	return createApplication(
		repository,
		posts,
		participations,
		classifier,
	)
}

// Application represents the campaign application
type Application interface {
	FindByID(id uuid.UUID) (domain_campaigns.Campaign, error)
	FindAll() ([]domain_campaigns.Campaign, error)
	FindCampaignsByUser(user users.User) ([]domain_campaigns.Campaign, error)
	FindCampaignsByCommunity(community communities.Community) ([]domain_campaigns.Campaign, error)
	FindCampaignsByPlatform(platform platforms.Platform) ([]domain_campaigns.Campaign, error)
	RebuildCampaigns() error
}
