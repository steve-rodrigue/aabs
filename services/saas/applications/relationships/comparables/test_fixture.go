package comparables

import domain_comparables "github.com/steve-rodrigue/aabs/services/saas/domain/relationships/comparables"

type applicationFixture struct {
	application Application
	comparator  *domain_comparables.MockComparator
}

func newApplicationFixture() *applicationFixture {
	comparator := domain_comparables.NewMockComparator()

	application := New(
		comparator,
	)

	return &applicationFixture{
		application: application,
		comparator:  comparator,
	}
}
