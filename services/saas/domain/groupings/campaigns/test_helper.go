package campaigns

import (
	"sort"
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
	FindByNameValue Campaign

	FindCalls int
	FindErr   error
	FindValue []Campaign

	FindAfterCalls int
	FindAfterErr   error
	FindAfterValue []Campaign

	CountCalls int
	CountErr   error
	CountValue int64
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

	if repository.FindByNameValue != nil {
		return repository.FindByNameValue, nil
	}

	for _, campaign := range repository.Items {
		if campaign.Name() == name {
			return campaign, nil
		}
	}

	return nil, nil
}

func (repository *MockCampaignRepository) Find(
	index int,
	amount int,
) ([]Campaign, error) {
	repository.FindCalls++

	if repository.FindErr != nil {
		return nil, repository.FindErr
	}

	if repository.FindValue != nil {
		return repository.FindValue, nil
	}

	items := repository.sortedCampaigns()

	if index >= len(items) {
		return []Campaign{}, nil
	}

	end := index + amount
	if end > len(items) {
		end = len(items)
	}

	return items[index:end], nil
}

func (repository *MockCampaignRepository) FindAfter(
	cursor uuid.UUID,
	amount int,
) ([]Campaign, error) {
	repository.FindAfterCalls++

	if repository.FindAfterErr != nil {
		return nil, repository.FindAfterErr
	}

	if repository.FindAfterValue != nil {
		return repository.FindAfterValue, nil
	}

	items := repository.sortedCampaigns()

	start := 0

	if cursor != uuid.Nil {
		for index, campaign := range items {
			if campaign.Identifier() == cursor {
				start = index + 1
				break
			}
		}
	}

	if start >= len(items) {
		return []Campaign{}, nil
	}

	end := start + amount
	if end > len(items) {
		end = len(items)
	}

	return items[start:end], nil
}

func (repository *MockCampaignRepository) Count() (int64, error) {
	repository.CountCalls++

	if repository.CountErr != nil {
		return 0, repository.CountErr
	}

	if repository.CountValue != 0 {
		return repository.CountValue, nil
	}

	return int64(len(repository.Items)), nil
}

func (repository *MockCampaignRepository) sortedCampaigns() []Campaign {
	out := make([]Campaign, 0, len(repository.Items))

	for _, campaign := range repository.Items {
		out = append(out, campaign)
	}

	sort.Slice(out, func(left int, right int) bool {
		return out[left].Identifier().String() <
			out[right].Identifier().String()
	})

	return out
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
