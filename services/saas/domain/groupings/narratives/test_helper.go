package narratives

import (
	"time"

	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters"
)

func NewMockNarrative(
	name string,
	description string,
) Narrative {
	return &MockNarrative{
		id:          uuid.New(),
		name:        name,
		description: description,
	}
}

type MockNarrative struct {
	id          uuid.UUID
	cluster     clusters.Cluster
	name        string
	description string
}

func (narrative *MockNarrative) Identifier() uuid.UUID {
	return narrative.id
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
}

func (repository *MockNarrativeRepository) Save(narrative Narrative) error {
	repository.SaveCalls++

	return repository.SaveErr
}

func (repository *MockNarrativeRepository) FindByID(id uuid.UUID) (Narrative, error) {
	repository.FindByIDCalls++

	if repository.FindByIDErr != nil {
		return nil, repository.FindByIDErr
	}

	return repository.Items[id], nil
}
