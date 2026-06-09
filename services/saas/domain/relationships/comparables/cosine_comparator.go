package comparables

import (
	"math"

	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships"
)

type cosineComparator struct {
	adapter relationships.Adapter
}

func (comparator *cosineComparator) Compare(
	source Comparable,
	target Comparable,
) (relationships.Relationship, error) {
	if source == nil {
		return nil, ErrInvalidSourceComparable
	}

	if target == nil {
		return nil, ErrInvalidTargetComparable
	}

	sourceVector := source.Vector()
	targetVector := target.Vector()

	if len(sourceVector) != len(targetVector) {
		return nil, ErrVectorSizeMismatch
	}

	similarity := cosineSimilarity(
		sourceVector,
		targetVector,
	)

	return comparator.adapter.ToDomain(
		relationships.RelationshipInput{
			Source:     source,
			Target:     target,
			Similarity: similarity,
		},
	)
}

func cosineSimilarity(
	source []float32,
	target []float32,
) float64 {
	var dot float64
	var sourceMagnitude float64
	var targetMagnitude float64

	for index := range source {
		sourceValue := float64(source[index])
		targetValue := float64(target[index])

		dot += sourceValue * targetValue

		sourceMagnitude += sourceValue * sourceValue
		targetMagnitude += targetValue * targetValue
	}

	if sourceMagnitude == 0 ||
		targetMagnitude == 0 {
		return 0
	}

	return dot /
		(math.Sqrt(sourceMagnitude) *
			math.Sqrt(targetMagnitude))
}
