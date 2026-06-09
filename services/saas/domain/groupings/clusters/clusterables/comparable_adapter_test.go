package clusterables

import (
	"errors"
	"testing"

	"github.com/google/uuid"
)

var errTest = errors.New("test error")

func TestNewComparableAdapter(t *testing.T) {
	adapter := NewComparableAdapter(
		NewMockClusterableAdapter(),
	)

	if adapter == nil {
		t.Fatalf("expected adapter")
	}
}

func TestComparableAdapterToDomain(t *testing.T) {
	clusterableAdapter := NewMockClusterableAdapter()
	adapter := NewComparableAdapter(clusterableAdapter)

	id := uuid.New()

	result, err := adapter.ToDomain(
		ComparableInput{
			Clusterable: ClusterableInput{
				Identifier:  id,
				ClusterKind: PostKind,
			},
			Vector: []float32{
				0.1,
				0.2,
				0.3,
			},
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	if clusterableAdapter.ToDomainCalls != 1 {
		t.Fatalf(
			"expected 1 clusterable adapter call, got %d",
			clusterableAdapter.ToDomainCalls,
		)
	}

	if result.Identifier() != id {
		t.Fatalf("expected identifier %s, got %s", id, result.Identifier())
	}

	if result.ClusterKind() != PostKind {
		t.Fatalf("expected kind %s, got %s", PostKind, result.ClusterKind())
	}

	vector := result.Vector()

	if len(vector) != 3 {
		t.Fatalf("expected vector length 3, got %d", len(vector))
	}

	if vector[0] != 0.1 ||
		vector[1] != 0.2 ||
		vector[2] != 0.3 {
		t.Fatalf("unexpected vector %+v", vector)
	}
}

func TestComparableAdapterToDomainReturnsClusterableAdapterError(t *testing.T) {
	clusterableAdapter := NewMockClusterableAdapter()
	clusterableAdapter.ToDomainErr = errTest

	adapter := NewComparableAdapter(clusterableAdapter)

	_, err := adapter.ToDomain(
		validComparableInput(nil),
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected clusterable adapter error, got %v", err)
	}
}

func TestComparableAdapterToDomainReturnsInvalidVectorErrorWhenVectorIsNil(t *testing.T) {
	adapter := NewComparableAdapter(
		NewMockClusterableAdapter(),
	)

	_, err := adapter.ToDomain(
		validComparableInput(func(input *ComparableInput) {
			input.Vector = nil
		}),
	)

	if !errors.Is(err, ErrInvalidComparableVector) {
		t.Fatalf("expected invalid comparable vector error, got %v", err)
	}
}

func TestComparableAdapterToDomainReturnsInvalidVectorErrorWhenVectorIsEmpty(t *testing.T) {
	adapter := NewComparableAdapter(
		NewMockClusterableAdapter(),
	)

	_, err := adapter.ToDomain(
		validComparableInput(func(input *ComparableInput) {
			input.Vector = []float32{}
		}),
	)

	if !errors.Is(err, ErrInvalidComparableVector) {
		t.Fatalf("expected invalid comparable vector error, got %v", err)
	}
}

func TestComparableAdapterToDomainCopiesVector(t *testing.T) {
	adapter := NewComparableAdapter(
		NewMockClusterableAdapter(),
	)

	vector := []float32{
		0.1,
		0.2,
	}

	result, err := adapter.ToDomain(
		validComparableInput(func(input *ComparableInput) {
			input.Vector = vector
		}),
	)

	if err != nil {
		t.Fatal(err)
	}

	vector[0] = 99

	resultVector := result.Vector()

	if resultVector[0] == 99 {
		t.Fatalf("expected vector copy")
	}
}

func validComparableInput(
	mutate func(input *ComparableInput),
) ComparableInput {
	input := ComparableInput{
		Clusterable: ClusterableInput{
			Identifier:  uuid.New(),
			ClusterKind: PostKind,
		},
		Vector: []float32{
			0.1,
			0.2,
		},
	}

	if mutate != nil {
		mutate(&input)
	}

	return input
}
