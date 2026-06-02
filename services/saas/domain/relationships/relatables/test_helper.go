package relatables

import "github.com/google/uuid"

func NewMockRelatable(id uuid.UUID, kind Kind) Relatable {
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
	return &MockRelatableRepository{
		Items: []Relatable{},
	}
}

type MockRelatableRepository struct {
	Items []Relatable

	FindAllCalls int
	FindAllErr   error
}

func (repository *MockRelatableRepository) FindAll() ([]Relatable, error) {
	repository.FindAllCalls++

	if repository.FindAllErr != nil {
		return nil, repository.FindAllErr
	}

	return repository.Items, nil
}
