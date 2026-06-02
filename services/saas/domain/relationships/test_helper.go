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

func NewMockRelationshipRepository() *MockRelationshipRepository {
	return &MockRelationshipRepository{
		Items: map[uuid.UUID]Relationship{},
	}
}

func NewMockRelationshipBuilder() *MockRelationshipBuilder {
	return &MockRelationshipBuilder{}
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

	FindAllCalls int
	FindAllErr   error
	FindAllValue []Relationship

	FindBySourceIDCalls int
	FindBySourceIDErr   error
	FindBySourceIDValue []Relationship

	FindByTargetIDCalls int
	FindByTargetIDErr   error
	FindByTargetIDValue []Relationship

	FindBySourceCalls int
	FindBySourceErr   error
	FindBySourceValue []Relationship

	FindByTargetCalls int
	FindByTargetErr   error
	FindByTargetValue []Relationship

	FindBetweenCalls int
	FindBetweenErr   error
	FindBetweenValue Relationship
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

	if repository.Items == nil {
		return nil, nil
	}

	return repository.Items[id], nil
}

func (repository *MockRelationshipRepository) FindAll() ([]Relationship, error) {
	repository.FindAllCalls++

	if repository.FindAllErr != nil {
		return nil, repository.FindAllErr
	}

	if repository.FindAllValue != nil {
		return repository.FindAllValue, nil
	}

	out := make([]Relationship, 0, len(repository.Items))

	for _, relationship := range repository.Items {
		out = append(out, relationship)
	}

	return out, nil
}

func (repository *MockRelationshipRepository) FindBySourceID(
	source uuid.UUID,
) ([]Relationship, error) {
	repository.FindBySourceIDCalls++

	if repository.FindBySourceIDErr != nil {
		return nil, repository.FindBySourceIDErr
	}

	return repository.FindBySourceIDValue, nil
}

func (repository *MockRelationshipRepository) FindByTargetID(
	target uuid.UUID,
) ([]Relationship, error) {
	repository.FindByTargetIDCalls++

	if repository.FindByTargetIDErr != nil {
		return nil, repository.FindByTargetIDErr
	}

	return repository.FindByTargetIDValue, nil
}

func (repository *MockRelationshipRepository) FindBySource(
	source relatables.Relatable,
) ([]Relationship, error) {
	repository.FindBySourceCalls++

	if repository.FindBySourceErr != nil {
		return nil, repository.FindBySourceErr
	}

	return repository.FindBySourceValue, nil
}

func (repository *MockRelationshipRepository) FindByTarget(
	target relatables.Relatable,
) ([]Relationship, error) {
	repository.FindByTargetCalls++

	if repository.FindByTargetErr != nil {
		return nil, repository.FindByTargetErr
	}

	return repository.FindByTargetValue, nil
}

func (repository *MockRelationshipRepository) FindBetween(
	source relatables.Relatable,
	target relatables.Relatable,
) (Relationship, error) {
	repository.FindBetweenCalls++

	if repository.FindBetweenErr != nil {
		return nil, repository.FindBetweenErr
	}

	return repository.FindBetweenValue, nil
}

type MockRelationshipBuilder struct {
	BuildCalls  int
	BuildErr    error
	BuildValue  []Relationship
	LastSource  relatables.Relatable
	LastTargets []relatables.Relatable
}

func (builder *MockRelationshipBuilder) Build(

	source relatables.Relatable,
	targets []relatables.Relatable,

) ([]Relationship, error) {
	builder.BuildCalls++
	builder.LastSource = source
	builder.LastTargets = targets
	return builder.BuildValue, builder.BuildErr

}
