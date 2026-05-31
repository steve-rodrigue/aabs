package campaigns

import (
	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/communities"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/campaigns"
	"github.com/steve-rodrigue/aabs/services/saas/domain/platforms"
	"github.com/steve-rodrigue/aabs/services/saas/domain/users"
)

// Application represents the campaign application
type Application interface {
	FindByID(id uuid.UUID) (campaigns.Campaign, error)
	FindAll() ([]campaigns.Campaign, error)
	FindCampaignsByUser(user users.User) ([]campaigns.Campaign, error)
	FindCampaignsByCommunity(community communities.Community) ([]campaigns.Campaign, error)
	FindCampaignsByPlatform(platform platforms.Platform) ([]campaigns.Campaign, error)
	RebuildCampaigns() error
}
