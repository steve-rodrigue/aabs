package campaigns

import (
	"context"
	"sort"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters/clusterables"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
)

func NewMockCampaign(
	name string,
	description string,
) Campaign {
	return &MockCampaign{
		ID:                     uuid.New(),
		ParticipationKindValue: participatables.CampaignKind,
		NameValue:              name,
		DescriptionValue:       description,
		CreatedOnValue:         time.Now().UTC(),
	}
}

func NewMockCampaignWithID(
	id uuid.UUID,
	name string,
	description string,
) Campaign {
	return &MockCampaign{
		ID:                     id,
		ParticipationKindValue: participatables.CampaignKind,
		NameValue:              name,
		DescriptionValue:       description,
		CreatedOnValue:         time.Now().UTC(),
	}
}

func NewMockCampaignRepository() *MockCampaignRepository {
	return &MockCampaignRepository{
		Items: map[uuid.UUID]Campaign{},
	}
}

func NewMockCampaignAdapter() *MockCampaignAdapter {
	return &MockCampaignAdapter{}
}

func NewMockCampaignClassifier() *MockCampaignClassifier {
	return &MockCampaignClassifier{}
}

func NewMockCampaignDetector() *MockCampaignDetector {
	return &MockCampaignDetector{}
}

type MockCampaign struct {
	ID uuid.UUID

	ParticipationKindValue participatables.Kind

	NameValue        string
	DescriptionValue string

	ClusterValue clusters.Cluster

	PostCountValue  int
	ConfidenceValue float64

	CreatedOnValue time.Time
}

func (campaign *MockCampaign) Identifier() uuid.UUID {
	return campaign.ID
}

func (campaign *MockCampaign) ParticipationKind() participatables.Kind {
	return campaign.ParticipationKindValue
}

func (campaign *MockCampaign) Name() string {
	return campaign.NameValue
}

func (campaign *MockCampaign) Description() string {
	return campaign.DescriptionValue
}

func (campaign *MockCampaign) Cluster() clusters.Cluster {
	return campaign.ClusterValue
}

func (campaign *MockCampaign) PostCount() int {
	return campaign.PostCountValue
}

func (campaign *MockCampaign) Confidence() float64 {
	return campaign.ConfidenceValue
}

func (campaign *MockCampaign) CreatedOn() time.Time {
	return campaign.CreatedOnValue
}

type MockCampaignAdapter struct {
	ToDomainCalls int
	ToDomainErr   error
	ToDomainValue Campaign

	LastInput CampaignInput
}

func (adapter *MockCampaignAdapter) ToDomain(
	input CampaignInput,
) (Campaign, error) {
	adapter.ToDomainCalls++
	adapter.LastInput = input

	if adapter.ToDomainErr != nil {
		return nil, adapter.ToDomainErr
	}

	if adapter.ToDomainValue != nil {
		return adapter.ToDomainValue, nil
	}

	return &MockCampaign{
		ID:                     input.Identifier,
		ParticipationKindValue: participatables.CampaignKind,
		NameValue:              input.Name,
		DescriptionValue:       input.Description,
		ClusterValue:           input.Cluster,
		PostCountValue:         input.PostCount,
		ConfidenceValue:        input.Confidence,
		CreatedOnValue:         input.CreatedOn,
	}, nil
}

type MockCampaignRepository struct {
	SaveCalls int
	SaveErr   error

	Items map[uuid.UUID]Campaign

	FindByIDCalls int
	FindByIDErr   error
	FindByIDValue Campaign

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

	LastContext context.Context
	LastSaved   Campaign
	LastID      uuid.UUID
	LastName    string
	LastIndex   int
	LastAmount  int
	LastCursor  uuid.UUID
}

func (repository *MockCampaignRepository) Save(
	ctx context.Context,
	campaign Campaign,
) error {
	repository.SaveCalls++
	repository.LastContext = ctx
	repository.LastSaved = campaign

	if repository.Items != nil && campaign != nil {
		repository.Items[campaign.Identifier()] = campaign
	}

	return repository.SaveErr
}

func (repository *MockCampaignRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (Campaign, error) {
	repository.FindByIDCalls++
	repository.LastContext = ctx
	repository.LastID = id

	if repository.FindByIDErr != nil {
		return nil, repository.FindByIDErr
	}

	if repository.FindByIDValue != nil {
		return repository.FindByIDValue, nil
	}

	if repository.Items == nil {
		return nil, nil
	}

	return repository.Items[id], nil
}

func (repository *MockCampaignRepository) FindByName(
	ctx context.Context,
	name string,
) (Campaign, error) {
	repository.FindByNameCalls++
	repository.LastContext = ctx
	repository.LastName = name

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
	ctx context.Context,
	index int,
	amount int,
) ([]Campaign, error) {
	repository.FindCalls++
	repository.LastContext = ctx
	repository.LastIndex = index
	repository.LastAmount = amount

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
	ctx context.Context,
	cursor uuid.UUID,
	amount int,
) ([]Campaign, error) {
	repository.FindAfterCalls++
	repository.LastContext = ctx
	repository.LastCursor = cursor
	repository.LastAmount = amount

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

func (repository *MockCampaignRepository) Count(
	ctx context.Context,
) (int64, error) {
	repository.CountCalls++
	repository.LastContext = ctx

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

	ClassifyValue      Campaign
	ClassifyConfidence float64

	LastContext context.Context
	LastPost    posts.Post
}

func (classifier *MockCampaignClassifier) Classify(
	ctx context.Context,
	post posts.Post,
) (Campaign, float64, error) {
	classifier.ClassifyCalls++
	classifier.LastContext = ctx
	classifier.LastPost = post

	return classifier.ClassifyValue,
		classifier.ClassifyConfidence,
		classifier.ClassifyErr
}

type MockCampaignDetector struct {
	DetectCalls int
	DetectErr   error
	DetectValue []Campaign

	LastContext    context.Context
	LastCandidates []clusterables.Clusterable
}

func (detector *MockCampaignDetector) Detect(
	ctx context.Context,
	candidates []clusterables.Clusterable,
) ([]Campaign, error) {
	detector.DetectCalls++
	detector.LastContext = ctx
	detector.LastCandidates = candidates

	return detector.DetectValue,
		detector.DetectErr
}
