package clusterables

import (
	"context"

	"github.com/google/uuid"
)

func NewMockClusterable(
	kind Kind,
) *MockClusterable {
	return &MockClusterable{
		ID:   uuid.New(),
		Kind: kind,
	}
}

func NewMockClusterableWithID(
	id uuid.UUID,
	kind Kind,
) *MockClusterable {
	return &MockClusterable{
		ID:   id,
		Kind: kind,
	}
}

type MockClusterable struct {
	ID   uuid.UUID
	Kind Kind
}

func (clusterable *MockClusterable) Identifier() uuid.UUID {
	return clusterable.ID
}

func (clusterable *MockClusterable) ClusterKind() Kind {
	return clusterable.Kind
}

func NewMockComparable(
	kind Kind,
	vector []float32,
) *MockComparable {
	return &MockComparable{
		ID:          uuid.New(),
		Kind:        kind,
		VectorValue: copyMockVector(vector),
	}
}

func NewMockComparableWithID(
	id uuid.UUID,
	kind Kind,
	vector []float32,
) *MockComparable {
	return &MockComparable{
		ID:          id,
		Kind:        kind,
		VectorValue: copyMockVector(vector),
	}
}

type MockComparable struct {
	ID          uuid.UUID
	Kind        Kind
	VectorValue []float32
}

func (comparable *MockComparable) Identifier() uuid.UUID {
	return comparable.ID
}

func (comparable *MockComparable) ClusterKind() Kind {
	return comparable.Kind
}

func (comparable *MockComparable) Vector() []float32 {
	return copyMockVector(comparable.VectorValue)
}

func NewMockClusterableAdapter() *MockClusterableAdapter {
	return &MockClusterableAdapter{}
}

type MockClusterableAdapter struct {
	ToDomainCalls int
	ToDomainErr   error
	ToDomainValue Clusterable

	ToDomainReturnsNil bool

	LastInput ClusterableInput
}

func (adapter *MockClusterableAdapter) ToDomain(
	input ClusterableInput,
) (Clusterable, error) {
	adapter.ToDomainCalls++
	adapter.LastInput = input

	if adapter.ToDomainErr != nil {
		return nil, adapter.ToDomainErr
	}

	if adapter.ToDomainReturnsNil {
		return nil, nil
	}

	if adapter.ToDomainValue != nil {
		return adapter.ToDomainValue, nil
	}

	return &MockClusterable{
		ID:   input.Identifier,
		Kind: input.ClusterKind,
	}, nil
}

func NewMockComparableAdapter() *MockComparableAdapter {
	return &MockComparableAdapter{}
}

type MockComparableAdapter struct {
	ToDomainCalls int
	ToDomainErr   error
	ToDomainValue Comparable

	LastInput ComparableInput
}

func (adapter *MockComparableAdapter) ToDomain(
	input ComparableInput,
) (Comparable, error) {
	adapter.ToDomainCalls++
	adapter.LastInput = input

	if adapter.ToDomainErr != nil {
		return nil, adapter.ToDomainErr
	}

	if adapter.ToDomainValue != nil {
		return adapter.ToDomainValue, nil
	}

	return &MockComparable{
		ID:          input.Clusterable.Identifier,
		Kind:        input.Clusterable.ClusterKind,
		VectorValue: copyMockVector(input.Vector),
	}, nil
}

func NewMockClusterableRepository() *MockClusterableRepository {
	return &MockClusterableRepository{}
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

	LastContext context.Context
	LastKind    Kind
	LastIndex   int
	LastCursor  uuid.UUID
	LastAmount  int

	FailOnCall int
}

func (repository *MockClusterableRepository) FindByKind(
	ctx context.Context,
	kind Kind,
	index int,
	amount int,
) ([]Clusterable, error) {
	repository.FindByKindCalls++
	repository.LastContext = ctx
	repository.LastKind = kind
	repository.LastIndex = index
	repository.LastAmount = amount

	return repository.FindByKindValue,
		repository.FindByKindErr
}

func (repository *MockClusterableRepository) FindByKindAfter(
	ctx context.Context,
	kind Kind,
	cursor uuid.UUID,
	amount int,
) ([]Clusterable, error) {
	repository.FindByKindAfterCalls++
	repository.LastContext = ctx
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
	ctx context.Context,
	kind Kind,
) (int64, error) {
	repository.CountByKindCalls++
	repository.LastContext = ctx
	repository.LastKind = kind

	return repository.CountByKindValue,
		repository.CountByKindErr
}

func NewMockCandidateRepository() *MockCandidateRepository {
	return &MockCandidateRepository{}
}

type MockCandidateRepository struct {
	FindCandidatesCalls int
	FindCandidatesErr   error
	FindCandidatesValue []Clusterable

	LastContext context.Context
	LastTarget  Clusterable
	LastKind    Kind
	LastAmount  int
}

func (repository *MockCandidateRepository) FindCandidates(
	ctx context.Context,
	target Clusterable,
	kind Kind,
	amount int,
) ([]Clusterable, error) {
	repository.FindCandidatesCalls++
	repository.LastContext = ctx
	repository.LastTarget = target
	repository.LastKind = kind
	repository.LastAmount = amount

	return repository.FindCandidatesValue,
		repository.FindCandidatesErr
}

func NewMockComparableRepository() *MockComparableRepository {
	return &MockComparableRepository{
		Items: map[uuid.UUID]Comparable{},
	}
}

type MockComparableRepository struct {
	Items map[uuid.UUID]Comparable

	FindByIDCalls int
	FindByIDErr   error
	FindByIDValue Comparable

	FindByKindCalls int
	FindByKindErr   error
	FindByKindValue []Comparable

	FindNearestCalls int
	FindNearestErr   error
	FindNearestValue []Comparable

	LastContext context.Context
	LastID      uuid.UUID
	LastKind    Kind
	LastIndex   int
	LastAmount  int
	LastTarget  Comparable
}

func (repository *MockComparableRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (Comparable, error) {
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

func (repository *MockComparableRepository) FindByKind(
	ctx context.Context,
	kind Kind,
	index int,
	amount int,
) ([]Comparable, error) {
	repository.FindByKindCalls++
	repository.LastContext = ctx
	repository.LastKind = kind
	repository.LastIndex = index
	repository.LastAmount = amount

	if repository.FindByKindErr != nil {
		return nil, repository.FindByKindErr
	}

	if repository.FindByKindValue != nil {
		return repository.FindByKindValue, nil
	}

	if repository.Items == nil {
		return []Comparable{}, nil
	}

	out := []Comparable{}

	for _, item := range repository.Items {
		if item.ClusterKind() == kind {
			out = append(out, item)
		}
	}

	if index >= len(out) {
		return []Comparable{}, nil
	}

	end := index + amount
	if end > len(out) {
		end = len(out)
	}

	return out[index:end], nil
}

func (repository *MockComparableRepository) FindNearest(
	ctx context.Context,
	target Comparable,
	kind Kind,
	amount int,
) ([]Comparable, error) {
	repository.FindNearestCalls++
	repository.LastContext = ctx
	repository.LastTarget = target
	repository.LastKind = kind
	repository.LastAmount = amount

	if repository.FindNearestErr != nil {
		return nil, repository.FindNearestErr
	}

	if repository.FindNearestValue != nil {
		return repository.FindNearestValue, nil
	}

	return []Comparable{}, nil
}

func copyMockVector(
	vector []float32,
) []float32 {
	out := make([]float32, len(vector))
	copy(out, vector)

	return out
}
