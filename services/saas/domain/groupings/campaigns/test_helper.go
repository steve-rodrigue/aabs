package campaigns

import (
	"time"

	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters"
)

func NewMockCampaign(
	name string,
	description string,
) Campaign {
	return &MockCampaign{
		id:          uuid.New(),
		name:        name,
		description: description,
	}
}

type MockCampaign struct {
	id          uuid.UUID
	name        string
	description string
	cluster     clusters.Cluster
}

func (campaign *MockCampaign) Identifier() uuid.UUID {
	return campaign.id
}

func (campaign *MockCampaign) Name() string {
	return campaign.name
}

func (campaign *MockCampaign) Description() string {
	return campaign.description
}

func (campaign *MockCampaign) Cluster() clusters.Cluster {
	return campaign.cluster
}

func (campaign *MockCampaign) PostCount() int {
	return 0
}

func (campaign *MockCampaign) Confidence() float64 {
	return 0
}

func (campaign *MockCampaign) CreatedOn() time.Time {
	return time.Time{}
}

type MockCampaignRepository struct {
	SaveCalls int
	SaveErr   error

	Items map[uuid.UUID]Campaign

	FindByIDCalls int
	FindByIDErr   error

	FindByNameCalls int
	FindByNameErr   error
}

func (repository *MockCampaignRepository) Save(campaign Campaign) error {
	repository.SaveCalls++

	return repository.SaveErr
}

func (repository *MockCampaignRepository) FindByID(id uuid.UUID) (Campaign, error) {
	repository.FindByIDCalls++

	if repository.FindByIDErr != nil {
		return nil, repository.FindByIDErr
	}

	return repository.Items[id], nil
}

func (repository *MockCampaignRepository) FindByName(name string) (Campaign, error) {
	repository.FindByNameCalls++

	if repository.FindByNameErr != nil {
		return nil, repository.FindByNameErr
	}

	for _, campaign := range repository.Items {
		if campaign.Name() == name {
			return campaign, nil
		}
	}

	return nil, nil
}
