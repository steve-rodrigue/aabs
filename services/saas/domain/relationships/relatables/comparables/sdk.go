package comparables

import (
	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships"
	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships/relatables"
)

// Comparable represents a comparable relatable
type Comparable interface {
	relatables.Relatable
	Vector() []float32
}

// Comparator represents a comparator
type Comparator interface {
	Compare(source Comparable, target Comparable) (relationships.Relationship, error)
}
