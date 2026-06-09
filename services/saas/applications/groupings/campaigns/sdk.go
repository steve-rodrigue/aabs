package campaigns

import (
	"github.com/google/uuid"

	app_participations "github.com/steve-rodrigue/aabs/services/saas/applications/groupings/participations"
	app_posts "github.com/steve-rodrigue/aabs/services/saas/applications/posts"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/communities"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/platforms"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/users"
	domain_campaigns "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/campaigns"
)

func New(
	repository domain_campaigns.Repository,
	posts app_posts.Application,
	participations app_participations.Application,
	classifier domain_campaigns.Classifier,
	rebuildBatchSize int,
) Application {
	return createApplication(
		repository,
		posts,
		participations,
		classifier,
		rebuildBatchSize,
	)
}

type Application interface {
	FindByID(id uuid.UUID) (domain_campaigns.Campaign, error)

	Find(index int, amount int) ([]domain_campaigns.Campaign, error)
	FindAfter(cursor uuid.UUID, amount int) ([]domain_campaigns.Campaign, error)

	FindCampaignsByUser(user users.User) ([]domain_campaigns.Campaign, error)
	FindCampaignsByCommunity(community communities.Community) ([]domain_campaigns.Campaign, error)
	FindCampaignsByPlatform(platform platforms.Platform) ([]domain_campaigns.Campaign, error)

	Count() (int64, error)

	RebuildCampaigns() error
}
