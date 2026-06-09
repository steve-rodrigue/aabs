package relatables

import (
	"context"

	"github.com/google/uuid"
)

func NewMockRelatable(
	id uuid.UUID,
	kind Kind,
) Relatable {
	return &MockRelatable{
		id:   id,
		kind: kind,
	}
}

type MockRelatable struct {
	id   uuid.UUID
	kind Kind
}

func (relatable *MockRelatable) Identifier() uuid.UUID {
	return relatable.id
}

func (relatable *MockRelatable) RelationshipKind() Kind {
	return relatable.kind
}

func NewMockRelatableAdapter() *MockRelatableAdapter {
	return &MockRelatableAdapter{}
}

type MockRelatableAdapter struct {
	ToDomainCalls int
	ToDomainErr   error
	ToDomainValue Relatable

	LastInput RelatableInput
}

func (adapter *MockRelatableAdapter) ToDomain(
	input RelatableInput,
) (Relatable, error) {
	adapter.ToDomainCalls++
	adapter.LastInput = input

	if adapter.ToDomainErr != nil {
		return nil, adapter.ToDomainErr
	}

	if adapter.ToDomainValue != nil {
		return adapter.ToDomainValue, nil
	}

	return &MockRelatable{
		id:   input.Identifier,
		kind: input.RelationshipKind,
	}, nil
}

func NewMockRelatableRepository() *MockRelatableRepository {
	return &MockRelatableRepository{
		Items: map[uuid.UUID]Relatable{},
	}
}

type MockRelatableRepository struct {
	SaveCalls int
	SaveErr   error

	DeleteCalls int
	DeleteErr   error

	DeleteByIDCalls int
	DeleteByIDErr   error

	FindCalls int
	FindErr   error
	FindValue []Relatable

	FindAfterCalls int
	FindAfterErr   error
	FindAfterValue []Relatable

	CountCalls int
	CountErr   error
	CountValue int64

	FindByKindCalls int
	FindByKindErr   error
	FindByKindValue []Relatable

	CountByKindCalls int
	CountByKindErr   error
	CountByKindValue int64

	Items map[uuid.UUID]Relatable

	LastContext   context.Context
	LastRelatable Relatable
	LastID        uuid.UUID
	LastKind      Kind
	LastIndex     int
	LastAmount    int
	LastCursor    uuid.UUID
}

func (repository *MockRelatableRepository) Save(
	ctx context.Context,
	relatable Relatable,
) error {
	repository.SaveCalls++
	repository.LastContext = ctx
	repository.LastRelatable = relatable

	if repository.SaveErr != nil {
		return repository.SaveErr
	}

	if repository.Items != nil && relatable != nil {
		repository.Items[relatable.Identifier()] = relatable
	}

	return nil
}

func (repository *MockRelatableRepository) Delete(
	ctx context.Context,
	relatable Relatable,
) error {
	repository.DeleteCalls++
	repository.LastContext = ctx
	repository.LastRelatable = relatable

	if repository.DeleteErr != nil {
		return repository.DeleteErr
	}

	if repository.Items != nil && relatable != nil {
		delete(repository.Items, relatable.Identifier())
	}

	return nil
}

func (repository *MockRelatableRepository) DeleteByID(
	ctx context.Context,
	id uuid.UUID,
) error {
	repository.DeleteByIDCalls++
	repository.LastContext = ctx
	repository.LastID = id

	if repository.DeleteByIDErr != nil {
		return repository.DeleteByIDErr
	}

	if repository.Items != nil {
		delete(repository.Items, id)
	}

	return nil
}

func (repository *MockRelatableRepository) Find(
	ctx context.Context,
	index int,
	amount int,
) ([]Relatable, error) {
	repository.FindCalls++
	repository.LastContext = ctx
	repository.LastIndex = index
	repository.LastAmount = amount

	return repository.FindValue, repository.FindErr
}

func (repository *MockRelatableRepository) FindAfter(
	ctx context.Context,
	cursor uuid.UUID,
	amount int,
) ([]Relatable, error) {
	repository.FindAfterCalls++
	repository.LastContext = ctx
	repository.LastCursor = cursor
	repository.LastAmount = amount

	if repository.FindAfterErr != nil {
		return nil, repository.FindAfterErr
	}

	if repository.FindAfterValue != nil {
		if repository.FindAfterCalls == 1 {
			return repository.FindAfterValue, nil
		}

		return []Relatable{}, nil
	}

	return []Relatable{}, nil
}

func (repository *MockRelatableRepository) Count(
	ctx context.Context,
) (
	int64,
	error,
) {
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

func (repository *MockRelatableRepository) FindByKind(
	ctx context.Context,
	kind Kind,
	index int,
	amount int,
) ([]Relatable, error) {
	repository.FindByKindCalls++
	repository.LastContext = ctx
	repository.LastKind = kind
	repository.LastIndex = index
	repository.LastAmount = amount

	return repository.FindByKindValue, repository.FindByKindErr
}

func (repository *MockRelatableRepository) CountByKind(
	ctx context.Context,
	kind Kind,
) (
	int64,
	error,
) {
	repository.CountByKindCalls++
	repository.LastContext = ctx
	repository.LastKind = kind

	return repository.CountByKindValue, repository.CountByKindErr
}

func NewMockCandidateRepository() *MockCandidateRepository {
	return &MockCandidateRepository{}
}

type MockCandidateRepository struct {
	FindCandidatesCalls int
	FindCandidatesErr   error
	FindCandidatesValue []Relatable

	LastContext context.Context
	LastSource  Relatable
	LastAmount  int
}

func (repository *MockCandidateRepository) FindCandidates(
	ctx context.Context,
	source Relatable,
	amount int,
) ([]Relatable, error) {
	repository.FindCandidatesCalls++
	repository.LastContext = ctx
	repository.LastSource = source
	repository.LastAmount = amount

	return repository.FindCandidatesValue,
		repository.FindCandidatesErr
}
