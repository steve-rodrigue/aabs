package builders

import (
	"errors"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/relatables"
	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships"
	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships/comparables"
)

var (
	ErrInvalidRelationshipBuilderSource = errors.New("invalid relationship builder source")
	ErrInvalidRelationshipBuilderTarget = errors.New("invalid relationship builder target")
)

// NewBuilder creates a new relationship builder
func NewBuilder(
	comparator comparables.Comparator,
) Builder {
	return &builder{
		comparator: comparator,
	}
}

// Builder represents a relationship builder
type Builder interface {
	Build(source relatables.Relatable, targets []relatables.Relatable) ([]relationships.Relationship, error)
}
