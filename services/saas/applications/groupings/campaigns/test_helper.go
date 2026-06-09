package campaigns

import (
	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/communities"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/platforms"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/users"
	domain_campaigns "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/campaigns"
)

func NewMockCampaignsApplication() *MockCampaignsApplication {
	return &MockCampaignsApplication{}
}

type MockCampaignsApplication struct {
	FindByIDCalls int
	FindByIDErr   error
	FindByIDValue domain_campaigns.Campaign

	FindCalls int
	FindErr   error
	FindValue []domain_campaigns.Campaign

	FindAfterCalls int
	FindAfterErr   error
	FindAfterValue []domain_campaigns.Campaign

	FindCampaignsByUserCalls int
	FindCampaignsByUserErr   error
	FindCampaignsByUserValue []domain_campaigns.Campaign

	FindCampaignsByCommunityCalls int
	FindCampaignsByCommunityErr   error
	FindCampaignsByCommunityValue []domain_campaigns.Campaign

	FindCampaignsByPlatformCalls int
	FindCampaignsByPlatformErr   error
	FindCampaignsByPlatformValue []domain_campaigns.Campaign

	CountCalls int
	CountErr   error
	CountValue int64

	RebuildCampaignsCalls int
	RebuildCampaignsErr   error
}

func (application *MockCampaignsApplication) FindByID(
	id uuid.UUID,
) (domain_campaigns.Campaign, error) {
	application.FindByIDCalls++

	return application.FindByIDValue, application.FindByIDErr
}

func (application *MockCampaignsApplication) Find(
	index int,
	amount int,
) ([]domain_campaigns.Campaign, error) {
	application.FindCalls++

	return application.FindValue, application.FindErr
}

func (application *MockCampaignsApplication) FindAfter(
	cursor uuid.UUID,
	amount int,
) ([]domain_campaigns.Campaign, error) {
	application.FindAfterCalls++

	return application.FindAfterValue, application.FindAfterErr
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

func (application *MockCampaignsApplication) Count() (int64, error) {
	application.CountCalls++

	return application.CountValue, application.CountErr
}

func (application *MockCampaignsApplication) RebuildCampaigns() error {
	application.RebuildCampaignsCalls++

	return application.RebuildCampaignsErr
}
