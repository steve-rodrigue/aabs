package builders

import (
	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/relatables"
	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships"
	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships/comparables"
)

type builder struct {
	comparator comparables.Comparator
}

func (builder *builder) Build(
	source relatables.Relatable,
	targets []relatables.Relatable,
) ([]relationships.Relationship, error) {
	sourceComparable, ok := source.(comparables.Comparable)
	if !ok {
		return nil, ErrInvalidRelationshipBuilderSource
	}

	out := make([]relationships.Relationship, 0, len(targets))

	for _, target := range targets {
		targetComparable, ok := target.(comparables.Comparable)
		if !ok {
			return nil, ErrInvalidRelationshipBuilderTarget
		}

		relationship, err := builder.comparator.Compare(
			sourceComparable,
			targetComparable,
		)
		if err != nil {
			return nil, err
		}

		out = append(out, relationship)
	}

	return out, nil
}
