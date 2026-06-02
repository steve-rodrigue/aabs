package comparables

import (
	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships"
	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships/relatables"
)

func NewMockComparable(
	id uuid.UUID,
	kind relatables.Kind,
	vector []float32,
) Comparable {
	return &MockComparable{
		id:     id,
		kind:   kind,
		vector: vector,
	}
}

type MockComparable struct {
	id     uuid.UUID
	kind   relatables.Kind
	vector []float32
}

func (comparable *MockComparable) Identifier() uuid.UUID {
	return comparable.id
}

func (comparable *MockComparable) RelationshipKind() relatables.Kind {
	return comparable.kind
}

func (comparable *MockComparable) Vector() []float32 {
	return comparable.vector
}

func NewMockComparator() *MockComparator {
	return &MockComparator{}
}

type MockComparator struct {
	CompareCalls int
	CompareErr   error

	LastSource Comparable
	LastTarget Comparable

	CompareValue relationships.Relationship
}

func (comparator *MockComparator) Compare(
	source Comparable,
	target Comparable,
) (relationships.Relationship, error) {
	comparator.CompareCalls++

	comparator.LastSource = source
	comparator.LastTarget = target

	if comparator.CompareErr != nil {
		return nil, comparator.CompareErr
	}

	return comparator.CompareValue, nil
}
