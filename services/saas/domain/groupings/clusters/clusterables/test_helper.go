package clusterables

import (
	"github.com/google/uuid"
)

func NewMockClusterable(
	kind Kind,
) *MockClusterable {
	return &MockClusterable{
		id:   uuid.New(),
		kind: kind,
	}
}

type MockClusterable struct {
	id   uuid.UUID
	kind Kind
}

func (clusterable *MockClusterable) Identifier() uuid.UUID {
	return clusterable.id
}

func (clusterable *MockClusterable) ClusterKind() Kind {
	return clusterable.kind
}

type MockClusterableRepository struct {
	FindByKindCalls int
	FindByKindErr   error
	FindByKindValue []Clusterable

	FindByKindAfterCalls int
	FindByKindAfterErr   error
	FindByKindAfterValue []Clusterable

	CountByKindCalls int
	CountByKindErr   error
	CountByKindValue int64

	LastKind   Kind
	LastIndex  int
	LastCursor uuid.UUID
	LastAmount int

	FailOnCall int
}

func NewMockClusterableRepository() *MockClusterableRepository {
	return &MockClusterableRepository{}
}

func (repository *MockClusterableRepository) FindByKind(
	kind Kind,
	index int,
	amount int,
) ([]Clusterable, error) {
	repository.FindByKindCalls++

	repository.LastKind = kind
	repository.LastIndex = index
	repository.LastAmount = amount

	return repository.FindByKindValue,
		repository.FindByKindErr
}

func (repository *MockClusterableRepository) FindByKindAfter(
	kind Kind,
	cursor uuid.UUID,
	amount int,
) ([]Clusterable, error) {
	repository.FindByKindAfterCalls++

	repository.LastKind = kind
	repository.LastCursor = cursor
	repository.LastAmount = amount

	if repository.FailOnCall > 0 {
		if repository.FindByKindAfterCalls == repository.FailOnCall {
			return nil, repository.FindByKindAfterErr
		}
	} else if repository.FindByKindAfterErr != nil {
		return nil, repository.FindByKindAfterErr
	}

	if repository.FindByKindAfterValue != nil {
		if repository.FindByKindAfterCalls == 1 {
			return repository.FindByKindAfterValue, nil
		}

		return []Clusterable{}, nil
	}

	return []Clusterable{}, nil
}

func (repository *MockClusterableRepository) CountByKind(
	kind Kind,
) (int64, error) {
	repository.CountByKindCalls++

	repository.LastKind = kind

	return repository.CountByKindValue,
		repository.CountByKindErr
}

type MockCandidateRepository struct {
	FindCandidatesCalls int
	FindCandidatesErr   error
	FindCandidatesValue []Clusterable

	LastTarget Clusterable
	LastKind   Kind
	LastAmount int
}

func NewMockCandidateRepository() *MockCandidateRepository {
	return &MockCandidateRepository{}
}

func (repository *MockCandidateRepository) FindCandidates(
	target Clusterable,
	kind Kind,
	amount int,
) ([]Clusterable, error) {
	repository.FindCandidatesCalls++

	repository.LastTarget = target
	repository.LastKind = kind
	repository.LastAmount = amount

	return repository.FindCandidatesValue,
		repository.FindCandidatesErr
}
