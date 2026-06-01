package searches

import "github.com/google/uuid"

type MockSearchRepository struct {
	StoreCalls  int
	SearchCalls int
	StoreErr    error
	SearchErr   error
}

func (repository *MockSearchRepository) Store(target uuid.UUID, vector []float32) error {
	repository.StoreCalls++

	return repository.StoreErr
}

func (repository *MockSearchRepository) Search(vector []float32, limit int) ([]Match, error) {
	repository.SearchCalls++

	return nil, repository.SearchErr
}
