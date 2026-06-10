package comparables

import (
	"errors"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/relatables"
	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships"
)

var (
	ErrInvalidSourceComparable = errors.New("invalid source comparable")
	ErrInvalidTargetComparable = errors.New("invalid target comparable")
	ErrVectorSizeMismatch      = errors.New("vector size mismatch")
)

// NewCosineComparator creates a new cosine comparator
func NewCosineComparator(
	adapter relationships.Adapter,
) Comparator {
	return &cosineComparator{
		adapter: adapter,
	}
}

// Comparable represents a comparable relatable
type Comparable interface {
	relatables.Relatable
	Vector() []float32
}

// Comparator represents a comparator
type Comparator interface {
	Compare(
		source Comparable,
		target Comparable,
	) (relationships.Relationship, error)
}
