package campaigns

import (
	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/communities"
	domain_campaigns "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/campaigns"
	"github.com/steve-rodrigue/aabs/services/saas/domain/platforms"
	"github.com/steve-rodrigue/aabs/services/saas/domain/users"
)

func NewMockCampaignsApplication() *MockCampaignsApplication {
	return &MockCampaignsApplication{}
}

type MockCampaignsApplication struct {
	FindByIDCalls int
	FindByIDErr   error
	FindByIDValue domain_campaigns.Campaign

	FindAllCalls int
	FindAllErr   error
	FindAllValue []domain_campaigns.Campaign

	FindCampaignsByUserCalls int
	FindCampaignsByUserErr   error
	FindCampaignsByUserValue []domain_campaigns.Campaign

	FindCampaignsByCommunityCalls int
	FindCampaignsByCommunityErr   error
	FindCampaignsByCommunityValue []domain_campaigns.Campaign

	FindCampaignsByPlatformCalls int
	FindCampaignsByPlatformErr   error
	FindCampaignsByPlatformValue []domain_campaigns.Campaign

	RebuildCampaignsCalls int
	RebuildCampaignsErr   error
}

func (application *MockCampaignsApplication) FindByID(
	id uuid.UUID,
) (domain_campaigns.Campaign, error) {
	application.FindByIDCalls++

	return application.FindByIDValue, application.FindByIDErr
}

func (application *MockCampaignsApplication) FindAll() (
	[]domain_campaigns.Campaign,
	error,
) {
	application.FindAllCalls++

	return application.FindAllValue, application.FindAllErr
}

func (application *MockCampaignsApplication) FindCampaignsByUser(
	user users.User,
) ([]domain_campaigns.Campaign, error) {
	application.FindCampaignsByUserCalls++

	return application.FindCampaignsByUserValue, application.FindCampaignsByUserErr
}

func (application *MockCampaignsApplication) FindCampaignsByCommunity(
	community communities.Community,
) ([]domain_campaigns.Campaign, error) {
	application.FindCampaignsByCommunityCalls++

	return application.FindCampaignsByCommunityValue, application.FindCampaignsByCommunityErr
}

func (application *MockCampaignsApplication) FindCampaignsByPlatform(
	platform platforms.Platform,
) ([]domain_campaigns.Campaign, error) {
	application.FindCampaignsByPlatformCalls++

	return application.FindCampaignsByPlatformValue, application.FindCampaignsByPlatformErr
}

func (application *MockCampaignsApplication) RebuildCampaigns() error {
	application.RebuildCampaignsCalls++

	return application.RebuildCampaignsErr
}
