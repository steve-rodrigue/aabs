package narratives

import (
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
		id:                uuid.New(),
		participationKind: participatables.NarrativeKind,
		name:              name,
		description:       description,
	}
}

func NewMockNarrativeRepository() *MockNarrativeRepository {
	return &MockNarrativeRepository{
		Items: map[uuid.UUID]Narrative{},
	}
}

type MockNarrative struct {
	id                uuid.UUID
	participationKind participatables.Kind
	cluster           clusters.Cluster
	name              string
	description       string
}

func (narrative *MockNarrative) SetCluster(
	cluster clusters.Cluster,
) {
	narrative.cluster = cluster
}

func (narrative *MockNarrative) Identifier() uuid.UUID {
	return narrative.id
}

func (narrative *MockNarrative) ParticipationKind() participatables.Kind {
	return participatables.NarrativeKind
}

func (narrative *MockNarrative) Cluster() clusters.Cluster {
	return narrative.cluster
}

func (narrative *MockNarrative) Name() string {
	return narrative.name
}

func (narrative *MockNarrative) Description() string {
	return narrative.description
}

func (narrative *MockNarrative) CreatedOn() time.Time {
	return time.Time{}
}

type MockNarrativeRepository struct {
	SaveCalls int
	SaveErr   error

	Items map[uuid.UUID]Narrative

	FindByIDCalls int
	FindByIDErr   error

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
}

func (repository *MockNarrativeRepository) Save(
	narrative Narrative,
) error {
	repository.SaveCalls++

	return repository.SaveErr
}

func (repository *MockNarrativeRepository) FindByID(
	id uuid.UUID,
) (Narrative, error) {
	repository.FindByIDCalls++

	if repository.FindByIDErr != nil {
		return nil, repository.FindByIDErr
	}

	return repository.Items[id], nil
}

func (repository *MockNarrativeRepository) Find(
	index int,
	amount int,
) ([]Narrative, error) {
	repository.FindCalls++

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
	name string,
) (Narrative, error) {
	repository.FindByNameCalls++

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
	cursor uuid.UUID,
	amount int,
) ([]Narrative, error) {
	repository.FindAfterCalls++

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

func (repository *MockNarrativeRepository) Count() (int64, error) {
	repository.CountCalls++

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
