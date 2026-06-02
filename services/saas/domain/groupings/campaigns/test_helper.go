package campaigns

import (
	"time"

	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
	"github.com/steve-rodrigue/aabs/services/saas/domain/posts"
)

func NewMockCampaign(
	name string,
	description string,
) Campaign {
	return &MockCampaign{
		id:                uuid.New(),
		participationKind: participatables.CampaignKind,
		name:              name,
		description:       description,
	}
}

func NewMockCampaignClassifier() *MockCampaignClassifier {
	return &MockCampaignClassifier{}
}

func NewMockCampaignRepository() *MockCampaignRepository {
	return &MockCampaignRepository{
		Items: map[uuid.UUID]Campaign{},
	}
}

type MockCampaign struct {
	id                uuid.UUID
	participationKind participatables.Kind

	name        string
	description string

	cluster    clusters.Cluster
	postCount  int
	confidence float64
}

func (campaign *MockCampaign) Identifier() uuid.UUID {
	return campaign.id
}

func (campaign *MockCampaign) ParticipationKind() participatables.Kind {
	return campaign.participationKind
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
	return campaign.postCount
}

func (campaign *MockCampaign) Confidence() float64 {
	return campaign.confidence
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

	FindAllCalls int
	FindAllErr   error
}

func (repository *MockCampaignRepository) Save(
	campaign Campaign,
) error {
	repository.SaveCalls++

	return repository.SaveErr
}

func (repository *MockCampaignRepository) FindByID(
	id uuid.UUID,
) (Campaign, error) {
	repository.FindByIDCalls++

	if repository.FindByIDErr != nil {
		return nil, repository.FindByIDErr
	}

	return repository.Items[id], nil
}

func (repository *MockCampaignRepository) FindByName(
	name string,
) (Campaign, error) {
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

func (repository *MockCampaignRepository) FindAll() (
	[]Campaign,
	error,
) {
	repository.FindAllCalls++

	if repository.FindAllErr != nil {
		return nil, repository.FindAllErr
	}

	out := make([]Campaign, 0, len(repository.Items))

	for _, campaign := range repository.Items {
		out = append(out, campaign)
	}

	return out, nil
}

type MockCampaignClassifier struct {
	ClassifyCalls int
	ClassifyErr   error
	ClassifyValue Campaign

	LastPost posts.Post
}

func (classifier *MockCampaignClassifier) Classify(
	post posts.Post,
) (Campaign, float64, error) {
	classifier.ClassifyCalls++
	classifier.LastPost = post

	return classifier.ClassifyValue, 1.0, classifier.ClassifyErr
}
