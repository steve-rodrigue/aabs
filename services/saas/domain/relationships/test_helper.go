package relationships

import (
	"time"

	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships/relatables"
)

func NewMockRelationship() Relationship {
	return &MockRelationship{
		id: uuid.New(),
	}
}

type MockRelationship struct {
	id         uuid.UUID
	source     relatables.Relatable
	target     relatables.Relatable
	similarity float64
}

func (relationship *MockRelationship) Identifier() uuid.UUID {
	return relationship.id
}

func (relationship *MockRelationship) Source() relatables.Relatable {
	return relationship.source
}

func (relationship *MockRelationship) Target() relatables.Relatable {
	return relationship.target
}

func (relationship *MockRelationship) Similarity() float64 {
	return relationship.similarity
}

func (relationship *MockRelationship) CreatedOn() time.Time {
	return time.Time{}
}

type MockRelationshipRepository struct {
	SaveCalls int
	SaveErr   error

	Items map[uuid.UUID]Relationship

	FindByIDCalls int
	FindByIDErr   error

	FindBySourceCalls int
	FindBySourceErr   error

	FindByTargetCalls int
	FindByTargetErr   error

	FindBetweenCalls int
	FindBetweenErr   error
}

func (repository *MockRelationshipRepository) Save(relationship Relationship) error {
	repository.SaveCalls++

	return repository.SaveErr
}

func (repository *MockRelationshipRepository) FindByID(id uuid.UUID) (Relationship, error) {
	repository.FindByIDCalls++

	if repository.FindByIDErr != nil {
		return nil, repository.FindByIDErr
	}

	return repository.Items[id], nil
}

func (repository *MockRelationshipRepository) FindBySource(
	source relatables.Relatable,
) ([]Relationship, error) {
	repository.FindBySourceCalls++

	if repository.FindBySourceErr != nil {
		return nil, repository.FindBySourceErr
	}

	return []Relationship{}, nil
}

func (repository *MockRelationshipRepository) FindByTarget(
	target relatables.Relatable,
) ([]Relationship, error) {
	repository.FindByTargetCalls++

	if repository.FindByTargetErr != nil {
		return nil, repository.FindByTargetErr
	}

	return []Relationship{}, nil
}

func (repository *MockRelationshipRepository) FindBetween(
	source relatables.Relatable,
	target relatables.Relatable,
) (Relationship, error) {
	repository.FindBetweenCalls++

	if repository.FindBetweenErr != nil {
		return nil, repository.FindBetweenErr
	}

	return nil, nil
}
