package relationships

import (
	"errors"
	"math"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/relatables"
)

func TestNewAdapter(t *testing.T) {
	adapter := NewAdapter()

	if adapter == nil {
		t.Fatalf("expected adapter")
	}
}

func TestAdapterToDomain(t *testing.T) {
	adapter := NewAdapter()

	id := uuid.New()
	source := relatables.NewMockRelatable(uuid.New(), relatables.PostKind)
	target := relatables.NewMockRelatable(uuid.New(), relatables.UserKind)
	createdOn := time.Now()

	result, err := adapter.ToDomain(
		RelationshipInput{
			Identifier: id,
			Source:     source,
			Target:     target,
			Similarity: 0.75,
			CreatedOn:  createdOn,
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	if result.Identifier() != id {
		t.Fatalf("expected identifier %s, got %s", id, result.Identifier())
	}

	if result.Source() != source {
		t.Fatalf("expected source")
	}

	if result.Target() != target {
		t.Fatalf("expected target")
	}

	if result.Similarity() != 0.75 {
		t.Fatalf("expected similarity %f, got %f", 0.75, result.Similarity())
	}

	if !result.CreatedOn().Equal(createdOn.UTC()) {
		t.Fatalf(
			"expected created on %s, got %s",
			createdOn.UTC(),
			result.CreatedOn(),
		)
	}
}

func TestAdapterToDomainAcceptsMinimumSimilarity(t *testing.T) {
	assertAdapterAcceptsSimilarity(t, -1)
}

func TestAdapterToDomainAcceptsZeroSimilarity(t *testing.T) {
	assertAdapterAcceptsSimilarity(t, 0)
}

func TestAdapterToDomainAcceptsMaximumSimilarity(t *testing.T) {
	assertAdapterAcceptsSimilarity(t, 1)
}

func TestAdapterToDomainReturnsInvalidIdentifierError(t *testing.T) {
	adapter := NewAdapter()

	_, err := adapter.ToDomain(
		validRelationshipInput(func(input *RelationshipInput) {
			input.Identifier = uuid.Nil
		}),
	)

	if !errors.Is(err, ErrInvalidRelationshipIdentifier) {
		t.Fatalf("expected invalid relationship identifier error, got %v", err)
	}
}

func TestAdapterToDomainReturnsInvalidSourceError(t *testing.T) {
	adapter := NewAdapter()

	_, err := adapter.ToDomain(
		validRelationshipInput(func(input *RelationshipInput) {
			input.Source = nil
		}),
	)

	if !errors.Is(err, ErrInvalidRelationshipSource) {
		t.Fatalf("expected invalid relationship source error, got %v", err)
	}
}

func TestAdapterToDomainReturnsInvalidTargetError(t *testing.T) {
	adapter := NewAdapter()

	_, err := adapter.ToDomain(
		validRelationshipInput(func(input *RelationshipInput) {
			input.Target = nil
		}),
	)

	if !errors.Is(err, ErrInvalidRelationshipTarget) {
		t.Fatalf("expected invalid relationship target error, got %v", err)
	}
}

func TestAdapterToDomainReturnsInvalidSimilarityErrorWhenLessThanMinusOne(t *testing.T) {
	adapter := NewAdapter()

	_, err := adapter.ToDomain(
		validRelationshipInput(func(input *RelationshipInput) {
			input.Similarity = -1.000001
		}),
	)

	if !errors.Is(err, ErrInvalidRelationshipSimilarity) {
		t.Fatalf("expected invalid relationship similarity error, got %v", err)
	}
}

func TestAdapterToDomainReturnsInvalidSimilarityErrorWhenGreaterThanOne(t *testing.T) {
	adapter := NewAdapter()

	_, err := adapter.ToDomain(
		validRelationshipInput(func(input *RelationshipInput) {
			input.Similarity = 1.000001
		}),
	)

	if !errors.Is(err, ErrInvalidRelationshipSimilarity) {
		t.Fatalf("expected invalid relationship similarity error, got %v", err)
	}
}

func TestAdapterToDomainReturnsInvalidSimilarityErrorWhenNaN(t *testing.T) {
	adapter := NewAdapter()

	_, err := adapter.ToDomain(
		validRelationshipInput(func(input *RelationshipInput) {
			input.Similarity = math.NaN()
		}),
	)

	if !errors.Is(err, ErrInvalidRelationshipSimilarity) {
		t.Fatalf("expected invalid relationship similarity error, got %v", err)
	}
}

func TestAdapterToDomainReturnsInvalidSimilarityErrorWhenPositiveInf(t *testing.T) {
	adapter := NewAdapter()

	_, err := adapter.ToDomain(
		validRelationshipInput(func(input *RelationshipInput) {
			input.Similarity = math.Inf(1)
		}),
	)

	if !errors.Is(err, ErrInvalidRelationshipSimilarity) {
		t.Fatalf("expected invalid relationship similarity error, got %v", err)
	}
}

func TestAdapterToDomainReturnsInvalidSimilarityErrorWhenNegativeInf(t *testing.T) {
	adapter := NewAdapter()

	_, err := adapter.ToDomain(
		validRelationshipInput(func(input *RelationshipInput) {
			input.Similarity = math.Inf(-1)
		}),
	)

	if !errors.Is(err, ErrInvalidRelationshipSimilarity) {
		t.Fatalf("expected invalid relationship similarity error, got %v", err)
	}
}

func TestAdapterToDomainReturnsInvalidCreatedOnError(t *testing.T) {
	adapter := NewAdapter()

	_, err := adapter.ToDomain(
		validRelationshipInput(func(input *RelationshipInput) {
			input.CreatedOn = time.Time{}
		}),
	)

	if !errors.Is(err, ErrInvalidRelationshipCreatedOn) {
		t.Fatalf("expected invalid relationship created on error, got %v", err)
	}
}

func assertAdapterAcceptsSimilarity(
	t *testing.T,
	similarity float64,
) {
	t.Helper()

	adapter := NewAdapter()

	result, err := adapter.ToDomain(
		validRelationshipInput(func(input *RelationshipInput) {
			input.Similarity = similarity
		}),
	)

	if err != nil {
		t.Fatal(err)
	}

	if result.Similarity() != similarity {
		t.Fatalf(
			"expected similarity %f, got %f",
			similarity,
			result.Similarity(),
		)
	}
}

func validRelationshipInput(
	mutate func(input *RelationshipInput),
) RelationshipInput {
	input := RelationshipInput{
		Identifier: uuid.New(),
		Source: relatables.NewMockRelatable(
			uuid.New(),
			relatables.PostKind,
		),
		Target: relatables.NewMockRelatable(
			uuid.New(),
			relatables.UserKind,
		),
		Similarity: 0.75,
		CreatedOn:  time.Now().UTC(),
	}

	if mutate != nil {
		mutate(&input)
	}

	return input
}
