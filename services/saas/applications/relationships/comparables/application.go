package comparables

import (
	domain_relationships "github.com/steve-rodrigue/aabs/services/saas/domain/relationships"
	domain_comparables "github.com/steve-rodrigue/aabs/services/saas/domain/relationships/comparables"
)

type application struct {
	comparator domain_comparables.Comparator
}

func createApplication(
	comparator domain_comparables.Comparator,
) Application {
	return &application{
		comparator: comparator,
	}
}

func (app *application) Compare(
	source domain_comparables.Comparable,
	target domain_comparables.Comparable,
) (domain_relationships.Relationship, error) {
	return app.comparator.Compare(
		source,
		target,
	)
}
