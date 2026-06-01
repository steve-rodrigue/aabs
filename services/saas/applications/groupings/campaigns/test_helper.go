package campaigns

import (
	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/communities"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/campaigns"
	"github.com/steve-rodrigue/aabs/services/saas/domain/platforms"
	"github.com/steve-rodrigue/aabs/services/saas/domain/users"
)

type MockCampaignsApplication struct {
	RebuildCampaignsCalls int
	RebuildCampaignsErr   error
}

func (application *MockCampaignsApplication) FindByID(id uuid.UUID) (campaigns.Campaign, error) {
	return nil, nil
}

func (application *MockCampaignsApplication) FindAll() ([]campaigns.Campaign, error) {
	return nil, nil
}

func (application *MockCampaignsApplication) FindCampaignsByUser(user users.User) ([]campaigns.Campaign, error) {
	return nil, nil
}

func (application *MockCampaignsApplication) FindCampaignsByCommunity(
	community communities.Community,
) ([]campaigns.Campaign, error) {
	return nil, nil
}

func (application *MockCampaignsApplication) FindCampaignsByPlatform(
	platform platforms.Platform,
) ([]campaigns.Campaign, error) {
	return nil, nil
}

func (application *MockCampaignsApplication) RebuildCampaigns() error {
	application.RebuildCampaignsCalls++

	return application.RebuildCampaignsErr
}
