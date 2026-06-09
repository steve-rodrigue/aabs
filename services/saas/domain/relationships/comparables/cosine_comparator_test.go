package comparables

import (
	"errors"
	"math"
	"testing"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships"
	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships/relatables"
)

func TestNewCosineComparator(t *testing.T) {
	comparator := NewCosineComparator(
		relationships.NewMockRelationshipAdapter(),
	)

	if comparator == nil {
		t.Fatalf("expected comparator")
	}
}

func TestCosineComparatorCompare(t *testing.T) {
	adapter := relationships.NewMockRelationshipAdapter()
	comparator := NewCosineComparator(adapter)

	source := NewMockComparable(
		uuid.New(),
		relatables.PostKind,
		[]float32{1, 0},
	)

	target := NewMockComparable(
		uuid.New(),
		relatables.PostKind,
		[]float32{0, 1},
	)

	result, err := comparator.Compare(source, target)

	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Fatalf("expected relationship")
	}

	if adapter.ToDomainCalls != 1 {
		t.Fatalf("expected 1 adapter call, got %d", adapter.ToDomainCalls)
	}

	if adapter.LastInput.Source != source {
		t.Fatalf("expected source to be passed")
	}

	if adapter.LastInput.Target != target {
		t.Fatalf("expected target to be passed")
	}

	if adapter.LastInput.Similarity != 0 {
		t.Fatalf("expected similarity 0, got %f", adapter.LastInput.Similarity)
	}
}

func TestCosineComparatorCompareSameDirection(t *testing.T) {
	adapter := relationships.NewMockRelationshipAdapter()
	comparator := NewCosineComparator(adapter)

	source := NewMockComparable(
		uuid.New(),
		relatables.PostKind,
		[]float32{1, 1},
	)

	target := NewMockComparable(
		uuid.New(),
		relatables.PostKind,
		[]float32{1, 1},
	)

	_, err := comparator.Compare(source, target)

	if err != nil {
		t.Fatal(err)
	}

	assertFloatEqual(t, adapter.LastInput.Similarity, 1)
}

func TestCosineComparatorCompareOppositeDirection(t *testing.T) {
	adapter := relationships.NewMockRelationshipAdapter()
	comparator := NewCosineComparator(adapter)

	source := NewMockComparable(
		uuid.New(),
		relatables.PostKind,
		[]float32{1, 0},
	)

	target := NewMockComparable(
		uuid.New(),
		relatables.PostKind,
		[]float32{-1, 0},
	)

	_, err := comparator.Compare(source, target)

	if err != nil {
		t.Fatal(err)
	}

	assertFloatEqual(t, adapter.LastInput.Similarity, -1)
}

func TestCosineComparatorCompareReturnsZeroWhenSourceVectorMagnitudeIsZero(t *testing.T) {
	adapter := relationships.NewMockRelationshipAdapter()
	comparator := NewCosineComparator(adapter)

	source := NewMockComparable(
		uuid.New(),
		relatables.PostKind,
		[]float32{0, 0},
	)

	target := NewMockComparable(
		uuid.New(),
		relatables.PostKind,
		[]float32{1, 1},
	)

	_, err := comparator.Compare(source, target)

	if err != nil {
		t.Fatal(err)
	}

	if adapter.LastInput.Similarity != 0 {
		t.Fatalf("expected similarity 0, got %f", adapter.LastInput.Similarity)
	}
}

func TestCosineComparatorCompareReturnsZeroWhenTargetVectorMagnitudeIsZero(t *testing.T) {
	adapter := relationships.NewMockRelationshipAdapter()
	comparator := NewCosineComparator(adapter)

	source := NewMockComparable(
		uuid.New(),
		relatables.PostKind,
		[]float32{1, 1},
	)

	target := NewMockComparable(
		uuid.New(),
		relatables.PostKind,
		[]float32{0, 0},
	)

	_, err := comparator.Compare(source, target)

	if err != nil {
		t.Fatal(err)
	}

	if adapter.LastInput.Similarity != 0 {
		t.Fatalf("expected similarity 0, got %f", adapter.LastInput.Similarity)
	}
}

func TestCosineComparatorCompareReturnsInvalidSourceComparableError(t *testing.T) {
	comparator := NewCosineComparator(
		relationships.NewMockRelationshipAdapter(),
	)

	_, err := comparator.Compare(
		nil,
		NewMockComparable(
			uuid.New(),
			relatables.PostKind,
			[]float32{1, 1},
		),
	)

	if !errors.Is(err, ErrInvalidSourceComparable) {
		t.Fatalf("expected invalid source comparable error, got %v", err)
	}
}

func TestCosineComparatorCompareReturnsInvalidTargetComparableError(t *testing.T) {
	comparator := NewCosineComparator(
		relationships.NewMockRelationshipAdapter(),
	)

	_, err := comparator.Compare(
		NewMockComparable(
			uuid.New(),
			relatables.PostKind,
			[]float32{1, 1},
		),
		nil,
	)

	if !errors.Is(err, ErrInvalidTargetComparable) {
		t.Fatalf("expected invalid target comparable error, got %v", err)
	}
}

func TestCosineComparatorCompareReturnsVectorSizeMismatchError(t *testing.T) {
	comparator := NewCosineComparator(
		relationships.NewMockRelationshipAdapter(),
	)

	_, err := comparator.Compare(
		NewMockComparable(
			uuid.New(),
			relatables.PostKind,
			[]float32{1, 1},
		),
		NewMockComparable(
			uuid.New(),
			relatables.PostKind,
			[]float32{1, 1, 1},
		),
	)

	if !errors.Is(err, ErrVectorSizeMismatch) {
		t.Fatalf("expected vector size mismatch error, got %v", err)
	}
}

func TestCosineComparatorCompareReturnsAdapterError(t *testing.T) {
	adapter := relationships.NewMockRelationshipAdapter()
	adapter.ToDomainErr = errTest

	comparator := NewCosineComparator(adapter)

	_, err := comparator.Compare(
		NewMockComparable(
			uuid.New(),
			relatables.PostKind,
			[]float32{1, 1},
		),
		NewMockComparable(
			uuid.New(),
			relatables.PostKind,
			[]float32{1, 1},
		),
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected adapter error, got %v", err)
	}
}

var errTest = errors.New("test error")

func assertFloatEqual(
	t *testing.T,
	result float64,
	expected float64,
) {
	t.Helper()

	const tolerance = 0.0000001

	if math.Abs(result-expected) > tolerance {
		t.Fatalf("expected %f, got %f", expected, result)
	}
}
