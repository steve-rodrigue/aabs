package relationships

import (
	"math"
	"time"

	"github.com/google/uuid"
)

type adapter struct{}

func (adapter *adapter) ToDomain(
	input RelationshipInput,
) (Relationship, error) {
	if input.Identifier == uuid.Nil {
		return nil, ErrInvalidRelationshipIdentifier
	}

	if input.Source == nil {
		return nil, ErrInvalidRelationshipSource
	}

	if input.Target == nil {
		return nil, ErrInvalidRelationshipTarget
	}

	if math.IsNaN(input.Similarity) ||
		math.IsInf(input.Similarity, 0) ||
		input.Similarity < -1 ||
		input.Similarity > 1 {
		return nil, ErrInvalidRelationshipSimilarity
	}

	if input.CreatedOn.IsZero() {
		return nil, ErrInvalidRelationshipCreatedOn
	}

	return &relationship{
		identifier: input.Identifier,
		source:     input.Source,
		target:     input.Target,
		similarity: input.Similarity,
		createdOn:  input.CreatedOn.UTC(),
	}, nil
}

func nowUTC() time.Time {
	return time.Now().UTC()
}
