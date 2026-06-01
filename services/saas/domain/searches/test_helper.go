package searches

import "github.com/google/uuid"

func NewMockSearchRepository() *MockSearchRepository {
	return &MockSearchRepository{
		StoreCalls: 0,
		StoreErr:   nil,

		LastStoredTarget: uuid.UUID{},
		LastStoredKind:   "",
		LastStoredVector: nil,

		SearchCalls: 0,
		SearchErr:   nil,
		Matches:     make([]Match, 0),

		LastSearchVector: nil,
		LastSearchLimit:  0,
	}
}

type MockSearchRepository struct {
	StoreCalls int
	StoreErr   error

	LastStoredTarget uuid.UUID
	LastStoredKind   Kind
	LastStoredVector []float32

	SearchCalls int
	SearchErr   error
	Matches     []Match

	LastSearchVector []float32
	LastSearchLimit  int
}

func (repository *MockSearchRepository) Store(
	target uuid.UUID,
	kind Kind,
	vector []float32,
) error {
	repository.StoreCalls++
	repository.LastStoredTarget = target
	repository.LastStoredKind = kind
	repository.LastStoredVector = vector

	return repository.StoreErr
}

func (repository *MockSearchRepository) Search(
	vector []float32,
	limit int,
) ([]Match, error) {
	repository.SearchCalls++
	repository.LastSearchVector = vector
	repository.LastSearchLimit = limit

	return repository.Matches, repository.SearchErr
}

func NewMockMatch(
	target uuid.UUID,
	kind Kind,
	similarity float64,
) Match {
	return &MockMatch{
		target:     target,
		kind:       kind,
		similarity: similarity,
	}
}

type MockMatch struct {
	target     uuid.UUID
	kind       Kind
	similarity float64
}

func (match *MockMatch) Target() uuid.UUID {
	return match.target
}

func (match *MockMatch) Kind() Kind {
	return match.kind
}

func (match *MockMatch) Similarity() float64 {
	return match.similarity
}
