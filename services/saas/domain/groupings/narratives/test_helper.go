package narratives

import (
	"context"
	"sort"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
)

func NewMockNarrative(
	name string,
	description string,
) *MockNarrative {
	return &MockNarrative{
		ID:                     uuid.New(),
		ParticipationKindValue: participatables.NarrativeKind,
		NameValue:              name,
		DescriptionValue:       description,
		CreatedOnValue:         time.Now().UTC(),
	}
}

func NewMockNarrativeWithID(
	id uuid.UUID,
	name string,
	description string,
) *MockNarrative {
	return &MockNarrative{
		ID:                     id,
		ParticipationKindValue: participatables.NarrativeKind,
		NameValue:              name,
		DescriptionValue:       description,
		CreatedOnValue:         time.Now().UTC(),
	}
}

func NewMockNarrativeRepository() *MockNarrativeRepository {
	return &MockNarrativeRepository{
		Items: map[uuid.UUID]Narrative{},
	}
}

func NewMockNarrativeAdapter() *MockNarrativeAdapter {
	return &MockNarrativeAdapter{}
}

type MockNarrative struct {
	ID uuid.UUID

	ParticipationKindValue participatables.Kind

	ClusterValue clusters.Cluster

	NameValue        string
	DescriptionValue string

	CreatedOnValue time.Time
}

func (narrative *MockNarrative) SetCluster(
	cluster clusters.Cluster,
) {
	narrative.ClusterValue = cluster
}

func (narrative *MockNarrative) Identifier() uuid.UUID {
	return narrative.ID
}

func (narrative *MockNarrative) ParticipationKind() participatables.Kind {
	return narrative.ParticipationKindValue
}

func (narrative *MockNarrative) Cluster() clusters.Cluster {
	return narrative.ClusterValue
}

func (narrative *MockNarrative) Name() string {
	return narrative.NameValue
}

func (narrative *MockNarrative) Description() string {
	return narrative.DescriptionValue
}

func (narrative *MockNarrative) CreatedOn() time.Time {
	return narrative.CreatedOnValue
}

type MockNarrativeAdapter struct {
	ToDomainCalls int
	ToDomainErr   error
	ToDomainValue Narrative

	LastInput NarrativeInput
}

func (adapter *MockNarrativeAdapter) ToDomain(
	input NarrativeInput,
) (Narrative, error) {
	adapter.ToDomainCalls++
	adapter.LastInput = input

	if adapter.ToDomainErr != nil {
		return nil, adapter.ToDomainErr
	}

	if adapter.ToDomainValue != nil {
		return adapter.ToDomainValue, nil
	}

	return &MockNarrative{
		ID:                     input.Identifier,
		ParticipationKindValue: input.ParticipationKind,
		ClusterValue:           input.Cluster,
		NameValue:              input.Name,
		DescriptionValue:       input.Description,
		CreatedOnValue:         input.CreatedOn,
	}, nil
}

type MockNarrativeRepository struct {
	SaveCalls int
	SaveErr   error

	Items map[uuid.UUID]Narrative

	FindByIDCalls int
	FindByIDErr   error
	FindByIDValue Narrative

	FindCalls int
	FindErr   error
	FindValue []Narrative

	FindByNameCalls int
	FindByNameErr   error
	FindByNameValue Narrative

	FindAfterCalls int
	FindAfterErr   error
	FindAfterValue []Narrative

	CountCalls int
	CountErr   error
	CountValue int64

	LastContext context.Context
	LastSaved   Narrative
	LastID      uuid.UUID
	LastName    string
	LastIndex   int
	LastAmount  int
	LastCursor  uuid.UUID
}

func (repository *MockNarrativeRepository) Save(
	ctx context.Context,
	narrative Narrative,
) error {
	repository.SaveCalls++
	repository.LastContext = ctx
	repository.LastSaved = narrative

	if repository.Items != nil && narrative != nil {
		repository.Items[narrative.Identifier()] = narrative
	}

	return repository.SaveErr
}

func (repository *MockNarrativeRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (Narrative, error) {
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

func (repository *MockNarrativeRepository) Find(
	ctx context.Context,
	index int,
	amount int,
) ([]Narrative, error) {
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

	items := repository.sortedNarratives()

	if index >= len(items) {
		return []Narrative{}, nil
	}

	end := index + amount
	if end > len(items) {
		end = len(items)
	}

	return items[index:end], nil
}

func (repository *MockNarrativeRepository) FindByName(
	ctx context.Context,
	name string,
) (Narrative, error) {
	repository.FindByNameCalls++
	repository.LastContext = ctx
	repository.LastName = name

	if repository.FindByNameErr != nil {
		return nil, repository.FindByNameErr
	}

	if repository.FindByNameValue != nil {
		return repository.FindByNameValue, nil
	}

	for _, narrative := range repository.Items {
		if narrative.Name() == name {
			return narrative, nil
		}
	}

	return nil, nil
}

func (repository *MockNarrativeRepository) FindAfter(
	ctx context.Context,
	cursor uuid.UUID,
	amount int,
) ([]Narrative, error) {
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

	items := repository.sortedNarratives()

	start := 0

	if cursor != uuid.Nil {
		for i, narrative := range items {
			if narrative.Identifier() == cursor {
				start = i + 1
				break
			}
		}
	}

	if start >= len(items) {
		return []Narrative{}, nil
	}

	end := start + amount
	if end > len(items) {
		end = len(items)
	}

	return items[start:end], nil
}

func (repository *MockNarrativeRepository) Count(
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

func (repository *MockNarrativeRepository) sortedNarratives() []Narrative {
	out := make([]Narrative, 0, len(repository.Items))

	for _, narrative := range repository.Items {
		out = append(out, narrative)
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].Identifier().String() <
			out[j].Identifier().String()
	})

	return out
}
