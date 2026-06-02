package comparables

import (
	domain_relationships "github.com/steve-rodrigue/aabs/services/saas/domain/relationships"
	domain_comparables "github.com/steve-rodrigue/aabs/services/saas/domain/relationships/comparables"
)

// New creates a new comparables application
func New(
	comparator domain_comparables.Comparator,
) Application {
	return createApplication(
		comparator,
	)
}

// Application represents a comparables application
type Application interface {
	Compare(
		source domain_comparables.Comparable,
		target domain_comparables.Comparable,
	) (domain_relationships.Relationship, error)
}
