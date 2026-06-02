package scorables

import "github.com/google/uuid"

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

	FindAllCalls int
	FindAllErr   error
	FindAllValue []Scorable

	FindByIDCalls int
	FindByIDErr   error
}

func (repository *MockScorableRepository) FindAll() ([]Scorable, error) {
	repository.FindAllCalls++

	if repository.FindAllErr != nil {
		return nil, repository.FindAllErr
	}

	if repository.FindAllValue != nil {
		return repository.FindAllValue, nil
	}

	out := make([]Scorable, 0, len(repository.Items))

	for _, item := range repository.Items {
		out = append(out, item)
	}

	return out, nil
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
