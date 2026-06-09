package relationships

import (
	"context"
	"sort"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships/relatables"
)

func NewMockRelationship() Relationship {
	return &MockRelationship{
		id:        uuid.New(),
		createdOn: time.Now().UTC(),
	}
}

func NewMockRelationshipWithRelatables(
	source relatables.Relatable,
	target relatables.Relatable,
	similarity float64,
) Relationship {
	return &MockRelationship{
		id:         uuid.New(),
		source:     source,
		target:     target,
		similarity: similarity,
		createdOn:  time.Now().UTC(),
	}
}

type MockRelationship struct {
	id uuid.UUID

	source     relatables.Relatable
	target     relatables.Relatable
	similarity float64
	createdOn  time.Time
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
	return relationship.createdOn
}

func NewMockRelationshipAdapter() *MockRelationshipAdapter {
	return &MockRelationshipAdapter{}
}

type MockRelationshipAdapter struct {
	ToDomainCalls int
	ToDomainErr   error
	ToDomainValue Relationship

	LastInput RelationshipInput
}

func (adapter *MockRelationshipAdapter) ToDomain(
	input RelationshipInput,
) (Relationship, error) {
	adapter.ToDomainCalls++
	adapter.LastInput = input

	if adapter.ToDomainErr != nil {
		return nil, adapter.ToDomainErr
	}

	if adapter.ToDomainValue != nil {
		return adapter.ToDomainValue, nil
	}

	return &MockRelationship{
		id:         input.Identifier,
		source:     input.Source,
		target:     input.Target,
		similarity: input.Similarity,
		createdOn:  input.CreatedOn,
	}, nil
}

func NewMockRelationshipRepository() *MockRelationshipRepository {
	return &MockRelationshipRepository{
		Items: map[uuid.UUID]Relationship{},
	}
}

type MockRelationshipRepository struct {
	SaveCalls int
	SaveErr   error

	Items map[uuid.UUID]Relationship

	FindByIDCalls int
	FindByIDErr   error

	FindCalls int
	FindErr   error
	FindValue []Relationship

	FindAfterCalls int
	FindAfterErr   error
	FindAfterValue []Relationship

	CountCalls int
	CountErr   error
	CountValue int64

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

	LastContext context.Context

	LastSaved Relationship

	LastID       uuid.UUID
	LastSourceID uuid.UUID
	LastTargetID uuid.UUID

	LastIndex  int
	LastAmount int
	LastCursor uuid.UUID

	LastSource relatables.Relatable
	LastTarget relatables.Relatable
}

func (repository *MockRelationshipRepository) Save(
	ctx context.Context,
	relationship Relationship,
) error {
	repository.SaveCalls++
	repository.LastContext = ctx
	repository.LastSaved = relationship

	if repository.SaveErr != nil {
		return repository.SaveErr
	}

	if repository.Items != nil && relationship != nil {
		repository.Items[relationship.Identifier()] = relationship
	}

	return nil
}

func (repository *MockRelationshipRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (Relationship, error) {
	repository.FindByIDCalls++
	repository.LastContext = ctx
	repository.LastID = id

	if repository.FindByIDErr != nil {
		return nil, repository.FindByIDErr
	}

	if repository.Items == nil {
		return nil, nil
	}

	return repository.Items[id], nil
}

func (repository *MockRelationshipRepository) Find(
	ctx context.Context,
	index int,
	amount int,
) ([]Relationship, error) {
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

	items := repository.sortedRelationships()

	if index >= len(items) {
		return []Relationship{}, nil
	}

	end := index + amount
	if end > len(items) {
		end = len(items)
	}

	return items[index:end], nil
}

func (repository *MockRelationshipRepository) FindAfter(
	ctx context.Context,
	cursor uuid.UUID,
	amount int,
) ([]Relationship, error) {
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

	items := repository.sortedRelationships()

	start := 0

	if cursor != uuid.Nil {
		for index, relationship := range items {
			if relationship.Identifier() == cursor {
				start = index + 1
				break
			}
		}
	}

	if start >= len(items) {
		return []Relationship{}, nil
	}

	end := start + amount
	if end > len(items) {
		end = len(items)
	}

	return items[start:end], nil
}

func (repository *MockRelationshipRepository) Count(
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

func (repository *MockRelationshipRepository) FindBySourceID(
	ctx context.Context,
	source uuid.UUID,
) ([]Relationship, error) {
	repository.FindBySourceIDCalls++
	repository.LastContext = ctx
	repository.LastSourceID = source

	if repository.FindBySourceIDErr != nil {
		return nil, repository.FindBySourceIDErr
	}

	if repository.FindBySourceIDValue != nil {
		return repository.FindBySourceIDValue, nil
	}

	out := []Relationship{}

	for _, relationship := range repository.Items {
		if relationship.Source() != nil &&
			relationship.Source().Identifier() == source {
			out = append(out, relationship)
		}
	}

	return out, nil
}

func (repository *MockRelationshipRepository) FindByTargetID(
	ctx context.Context,
	target uuid.UUID,
) ([]Relationship, error) {
	repository.FindByTargetIDCalls++
	repository.LastContext = ctx
	repository.LastTargetID = target

	if repository.FindByTargetIDErr != nil {
		return nil, repository.FindByTargetIDErr
	}

	if repository.FindByTargetIDValue != nil {
		return repository.FindByTargetIDValue, nil
	}

	out := []Relationship{}

	for _, relationship := range repository.Items {
		if relationship.Target() != nil &&
			relationship.Target().Identifier() == target {
			out = append(out, relationship)
		}
	}

	return out, nil
}

func (repository *MockRelationshipRepository) FindBySource(
	ctx context.Context,
	source relatables.Relatable,
) ([]Relationship, error) {
	repository.FindBySourceCalls++
	repository.LastContext = ctx
	repository.LastSource = source

	if repository.FindBySourceErr != nil {
		return nil, repository.FindBySourceErr
	}

	if repository.FindBySourceValue != nil {
		return repository.FindBySourceValue, nil
	}

	return repository.FindBySourceID(ctx, source.Identifier())
}

func (repository *MockRelationshipRepository) FindByTarget(
	ctx context.Context,
	target relatables.Relatable,
) ([]Relationship, error) {
	repository.FindByTargetCalls++
	repository.LastContext = ctx
	repository.LastTarget = target

	if repository.FindByTargetErr != nil {
		return nil, repository.FindByTargetErr
	}

	if repository.FindByTargetValue != nil {
		return repository.FindByTargetValue, nil
	}

	return repository.FindByTargetID(ctx, target.Identifier())
}

func (repository *MockRelationshipRepository) FindBetween(
	ctx context.Context,
	source relatables.Relatable,
	target relatables.Relatable,
) (Relationship, error) {
	repository.FindBetweenCalls++
	repository.LastContext = ctx
	repository.LastSource = source
	repository.LastTarget = target

	if repository.FindBetweenErr != nil {
		return nil, repository.FindBetweenErr
	}

	if repository.FindBetweenValue != nil {
		return repository.FindBetweenValue, nil
	}

	for _, relationship := range repository.Items {
		if relationship.Source() == nil || relationship.Target() == nil {
			continue
		}

		if relationship.Source().Identifier() == source.Identifier() &&
			relationship.Target().Identifier() == target.Identifier() {
			return relationship, nil
		}
	}

	return nil, nil
}

func (repository *MockRelationshipRepository) sortedRelationships() []Relationship {
	out := make([]Relationship, 0, len(repository.Items))

	for _, relationship := range repository.Items {
		out = append(out, relationship)
	}

	sort.Slice(out, func(left int, right int) bool {
		return out[left].Identifier().String() <
			out[right].Identifier().String()
	})

	return out
}
