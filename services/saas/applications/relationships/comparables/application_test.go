package comparables

import (
	"errors"
	"testing"

	"github.com/google/uuid"

	domain_relationships "github.com/steve-rodrigue/aabs/services/saas/domain/relationships"
	domain_comparables "github.com/steve-rodrigue/aabs/services/saas/domain/relationships/comparables"
	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships/relatables"
)

var errTest = errors.New("test error")

func TestCompare(t *testing.T) {
	fixture := newApplicationFixture()

	source := domain_comparables.NewMockComparable(
		uuid.New(),
		relatables.PostKind,
		[]float32{1, 2, 3},
	)

	target := domain_comparables.NewMockComparable(
		uuid.New(),
		relatables.TopicKind,
		[]float32{4, 5, 6},
	)

	relationship := domain_relationships.NewMockRelationship()

	fixture.comparator.CompareValue = relationship

	result, err := fixture.application.Compare(
		source,
		target,
	)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.comparator.CompareCalls != 1 {
		t.Fatalf("expected 1 compare call")
	}

	if fixture.comparator.LastSource != source {
		t.Fatalf("expected source passed to comparator")
	}

	if fixture.comparator.LastTarget != target {
		t.Fatalf("expected target passed to comparator")
	}

	if result != relationship {
		t.Fatalf("expected relationship result")
	}
}

func TestCompareReturnsError(t *testing.T) {
	fixture := newApplicationFixture()

	source := domain_comparables.NewMockComparable(
		uuid.New(),
		relatables.PostKind,
		nil,
	)

	target := domain_comparables.NewMockComparable(
		uuid.New(),
		relatables.TopicKind,
		nil,
	)

	fixture.comparator.CompareErr = errTest

	_, err := fixture.application.Compare(
		source,
		target,
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected compare error, got %v", err)
	}
}
