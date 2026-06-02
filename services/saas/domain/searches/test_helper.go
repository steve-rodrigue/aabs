package searches

import "github.com/google/uuid"

func NewMockSearchRepository() *MockSearchRepository {
	return &MockSearchRepository{
		Matches: make([]Match, 0),
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

func NewMockSearchable(
	id uuid.UUID,
	kind Kind,
	title string,
	text string,
) Searchable {
	return &MockSearchable{
		id:    id,
		kind:  kind,
		title: title,
		text:  text,
	}
}

type MockSearchable struct {
	id    uuid.UUID
	kind  Kind
	title string
	text  string
}

func (searchable *MockSearchable) Identifier() uuid.UUID {
	return searchable.id
}

func (searchable *MockSearchable) SearchKind() Kind {
	return searchable.kind
}

func (searchable *MockSearchable) SearchTitle() string {
	return searchable.title
}

func (searchable *MockSearchable) SearchText() string {
	return searchable.text
}

func NewMockSearchableRepository() *MockSearchableRepository {
	return &MockSearchableRepository{
		Items: map[uuid.UUID]Searchable{},
	}
}

type MockSearchableRepository struct {
	Items map[uuid.UUID]Searchable

	FindByIDCalls int
	FindByIDErr   error

	FindCalls int
	FindErr   error
	FindValue []Searchable

	FindAfterCalls int
	FindAfterErr   error
	FindAfterValue []Searchable

	CountCalls int
	CountErr   error
	CountValue int64

	LastIndex  int
	LastAmount int

	LastCursor uuid.UUID
}

func (repository *MockSearchableRepository) FindByID(
	id uuid.UUID,
) (Searchable, error) {
	repository.FindByIDCalls++

	if repository.FindByIDErr != nil {
		return nil, repository.FindByIDErr
	}

	return repository.Items[id], nil
}

func (repository *MockSearchableRepository) Find(
	index int,
	amount int,
) ([]Searchable, error) {
	repository.FindCalls++

	repository.LastIndex = index
	repository.LastAmount = amount

	if repository.FindErr != nil {
		return nil, repository.FindErr
	}

	if repository.FindValue != nil {
		return repository.FindValue, nil
	}

	out := make([]Searchable, 0, len(repository.Items))

	for _, item := range repository.Items {
		out = append(out, item)
	}

	return out, nil
}

func (repository *MockSearchableRepository) FindAfter(
	cursor uuid.UUID,
	amount int,
) ([]Searchable, error) {
	repository.FindAfterCalls++

	repository.LastCursor = cursor
	repository.LastAmount = amount

	if repository.FindAfterErr != nil {
		return nil, repository.FindAfterErr
	}

	return repository.FindAfterValue, nil
}

func (repository *MockSearchableRepository) Count() (int64, error) {
	repository.CountCalls++

	if repository.CountErr != nil {
		return 0, repository.CountErr
	}

	return repository.CountValue, nil
}
