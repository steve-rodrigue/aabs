package scorables

import (
	"sort"

	"github.com/google/uuid"
)

func NewMockScorable(
	id uuid.UUID,
	kind Kind,
) Scorable {
	return &MockScorable{
		id:   id,
		kind: kind,
	}
}

type MockScorable struct {
	id   uuid.UUID
	kind Kind
}

func (scorable *MockScorable) Identifier() uuid.UUID {
	return scorable.id
}

func (scorable *MockScorable) ScoreKind() Kind {
	return scorable.kind
}

func NewMockScorableRepository() *MockScorableRepository {
	return &MockScorableRepository{
		Items: map[uuid.UUID]Scorable{},
	}
}

type MockScorableRepository struct {
	Items map[uuid.UUID]Scorable

	FindCalls int
	FindErr   error
	FindValue []Scorable

	FindAfterCalls int
	FindAfterErr   error
	FindAfterValue []Scorable

	CountCalls int
	CountErr   error
	CountValue int64

	FindByIDCalls int
	FindByIDErr   error
}

func (repository *MockScorableRepository) Find(
	index int,
	amount int,
) ([]Scorable, error) {
	repository.FindCalls++

	if repository.FindErr != nil {
		return nil, repository.FindErr
	}

	if repository.FindValue != nil {
		return repository.FindValue, nil
	}

	items := repository.sortedScorables()

	if index >= len(items) {
		return []Scorable{}, nil
	}

	end := index + amount
	if end > len(items) {
		end = len(items)
	}

	return items[index:end], nil
}

func (repository *MockScorableRepository) FindAfter(
	cursor uuid.UUID,
	amount int,
) ([]Scorable, error) {
	repository.FindAfterCalls++

	if repository.FindAfterErr != nil {
		return nil, repository.FindAfterErr
	}

	if repository.FindAfterValue != nil {
		if repository.FindAfterCalls == 1 {
			return repository.FindAfterValue, nil
		}

		return []Scorable{}, nil
	}

	items := repository.sortedScorables()

	start := 0

	if cursor != uuid.Nil {
		for index, scorable := range items {
			if scorable.Identifier() == cursor {
				start = index + 1
				break
			}
		}
	}

	if start >= len(items) {
		return []Scorable{}, nil
	}

	end := start + amount
	if end > len(items) {
		end = len(items)
	}

	return items[start:end], nil
}

func (repository *MockScorableRepository) Count() (int64, error) {
	repository.CountCalls++

	if repository.CountErr != nil {
		return 0, repository.CountErr
	}

	if repository.CountValue != 0 {
		return repository.CountValue, nil
	}

	return int64(len(repository.Items)), nil
}

func (repository *MockScorableRepository) FindByID(id uuid.UUID) (Scorable, error) {
	repository.FindByIDCalls++

	if repository.FindByIDErr != nil {
		return nil, repository.FindByIDErr
	}

	if repository.Items == nil {
		return nil, nil
	}

	return repository.Items[id], nil
}

func (repository *MockScorableRepository) sortedScorables() []Scorable {
	out := make([]Scorable, 0, len(repository.Items))

	for _, item := range repository.Items {
		out = append(out, item)
	}

	sort.Slice(out, func(left int, right int) bool {
		return out[left].Identifier().String() <
			out[right].Identifier().String()
	})

	return out
}
