package comparables

import (
	domain_relationships "github.com/steve-rodrigue/aabs/services/saas/domain/relationships"
	domain_comparables "github.com/steve-rodrigue/aabs/services/saas/domain/relationships/comparables"
)

func NewMockComparablesApplication() *MockComparablesApplication {
	return &MockComparablesApplication{}
}

type MockComparablesApplication struct {
	CompareCalls int
	CompareErr   error

	CompareValue domain_relationships.Relationship

	LastSource domain_comparables.Comparable
	LastTarget domain_comparables.Comparable
}

func (application *MockComparablesApplication) Compare(
	source domain_comparables.Comparable,
	target domain_comparables.Comparable,
) (domain_relationships.Relationship, error) {
	application.CompareCalls++

	application.LastSource = source
	application.LastTarget = target

	return application.CompareValue, application.CompareErr
}
