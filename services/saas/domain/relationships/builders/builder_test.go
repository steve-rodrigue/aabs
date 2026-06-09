package builders

import (
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships"
	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships/comparables"
	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships/relatables"
)

func TestNewBuilder(t *testing.T) {
	builder := NewBuilder(
		comparables.NewMockComparator(),
	)

	if builder == nil {
		t.Fatalf("expected builder")
	}
}

func TestBuilderBuild(t *testing.T) {
	comparator := comparables.NewMockComparator()

	relationship := relationships.NewMockRelationship()
	comparator.CompareValue = relationship

	builder := NewBuilder(comparator)

	source := comparables.NewMockComparable(
		uuid.New(),
		relatables.PostKind,
		[]float32{1, 0},
	)

	target := comparables.NewMockComparable(
		uuid.New(),
		relatables.PostKind,
		[]float32{0, 1},
	)

	result, err := builder.Build(
		source,
		[]relatables.Relatable{
			target,
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 relationship, got %d", len(result))
	}

	if result[0] != relationship {
		t.Fatalf("expected relationship result")
	}

	if comparator.CompareCalls != 1 {
		t.Fatalf("expected 1 compare call, got %d", comparator.CompareCalls)
	}

	if comparator.LastSource != source {
		t.Fatalf("expected source to be passed")
	}

	if comparator.LastTarget != target {
		t.Fatalf("expected target to be passed")
	}
}

func TestBuilderBuildMultipleTargets(t *testing.T) {
	comparator := comparables.NewMockComparator()
	comparator.CompareValue = relationships.NewMockRelationship()

	builder := NewBuilder(comparator)

	source := comparables.NewMockComparable(
		uuid.New(),
		relatables.PostKind,
		[]float32{1, 0},
	)

	first := comparables.NewMockComparable(
		uuid.New(),
		relatables.PostKind,
		[]float32{0, 1},
	)

	second := comparables.NewMockComparable(
		uuid.New(),
		relatables.UserKind,
		[]float32{1, 1},
	)

	result, err := builder.Build(
		source,
		[]relatables.Relatable{
			first,
			second,
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 relationships, got %d", len(result))
	}

	if comparator.CompareCalls != 2 {
		t.Fatalf("expected 2 compare calls, got %d", comparator.CompareCalls)
	}

	if comparator.LastTarget != second {
		t.Fatalf("expected last target to be second")
	}
}

func TestBuilderBuildReturnsEmptyWhenTargetsAreEmpty(t *testing.T) {
	comparator := comparables.NewMockComparator()
	builder := NewBuilder(comparator)

	source := comparables.NewMockComparable(
		uuid.New(),
		relatables.PostKind,
		[]float32{1, 0},
	)

	result, err := builder.Build(
		source,
		[]relatables.Relatable{},
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 0 {
		t.Fatalf("expected empty relationships, got %d", len(result))
	}

	if comparator.CompareCalls != 0 {
		t.Fatalf("expected comparator not to be called")
	}
}

func TestBuilderBuildReturnsInvalidSourceError(t *testing.T) {
	builder := NewBuilder(
		comparables.NewMockComparator(),
	)

	source := relatables.NewMockRelatable(
		uuid.New(),
		relatables.PostKind,
	)

	target := comparables.NewMockComparable(
		uuid.New(),
		relatables.PostKind,
		[]float32{1, 0},
	)

	_, err := builder.Build(
		source,
		[]relatables.Relatable{
			target,
		},
	)

	if !errors.Is(err, ErrInvalidRelationshipBuilderSource) {
		t.Fatalf("expected invalid source error, got %v", err)
	}
}

func TestBuilderBuildReturnsInvalidTargetError(t *testing.T) {
	builder := NewBuilder(
		comparables.NewMockComparator(),
	)

	source := comparables.NewMockComparable(
		uuid.New(),
		relatables.PostKind,
		[]float32{1, 0},
	)

	target := relatables.NewMockRelatable(
		uuid.New(),
		relatables.PostKind,
	)

	_, err := builder.Build(
		source,
		[]relatables.Relatable{
			target,
		},
	)

	if !errors.Is(err, ErrInvalidRelationshipBuilderTarget) {
		t.Fatalf("expected invalid target error, got %v", err)
	}
}

func TestBuilderBuildReturnsComparatorError(t *testing.T) {
	comparator := comparables.NewMockComparator()
	comparator.CompareErr = errTest

	builder := NewBuilder(comparator)

	source := comparables.NewMockComparable(
		uuid.New(),
		relatables.PostKind,
		[]float32{1, 0},
	)

	target := comparables.NewMockComparable(
		uuid.New(),
		relatables.PostKind,
		[]float32{0, 1},
	)

	_, err := builder.Build(
		source,
		[]relatables.Relatable{
			target,
		},
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected comparator error, got %v", err)
	}
}

var errTest = errors.New("test error")
