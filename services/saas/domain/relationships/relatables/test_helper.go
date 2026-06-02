package relatables

import "github.com/google/uuid"

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

func NewMockRelatableRepository() *MockRelatableRepository {
	return &MockRelatableRepository{}
}

func NewMockCandidateRepository() *MockCandidateRepository {
	return &MockCandidateRepository{}
}

type MockRelatableRepository struct {
	FindCalls int
	FindErr   error
	FindValue []Relatable

	FindAfterCalls int
	FindAfterErr   error
	FindAfterValue []Relatable

	CountCalls int
	CountErr   error
	CountValue int64

	LastIndex  int
	LastAmount int

	LastCursor uuid.UUID
}

func (repository *MockRelatableRepository) Find(
	index int,
	amount int,
) ([]Relatable, error) {
	repository.FindCalls++

	repository.LastIndex = index
	repository.LastAmount = amount

	return repository.FindValue,
		repository.FindErr
}

func (repository *MockRelatableRepository) FindAfter(
	cursor uuid.UUID,
	amount int,
) ([]Relatable, error) {
	repository.FindAfterCalls++

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

func (repository *MockRelatableRepository) Count() (
	int64,
	error,
) {
	repository.CountCalls++

	return repository.CountValue,
		repository.CountErr
}

type MockCandidateRepository struct {
	FindCandidatesCalls int
	FindCandidatesErr   error
	FindCandidatesValue []Relatable

	LastSource Relatable
	LastAmount int
}

func (repository *MockCandidateRepository) FindCandidates(
	source Relatable,
	amount int,
) ([]Relatable, error) {
	repository.FindCandidatesCalls++

	repository.LastSource = source
	repository.LastAmount = amount

	return repository.FindCandidatesValue,
		repository.FindCandidatesErr
}
