package clusterables

import (
	"testing"

	"github.com/google/uuid"
)

func TestComparable(t *testing.T) {
	id := uuid.New()

	clusterable := NewMockClusterable(
		PostKind,
	)

	clusterable.ID = id

	vector := []float32{
		0.1,
		0.2,
		0.3,
	}

	comparable := &comparable{
		clusterable: clusterable,
		vector:      vector,
	}

	if comparable.Identifier() != id {
		t.Fatalf(
			"expected identifier %s, got %s",
			id,
			comparable.Identifier(),
		)
	}

	if comparable.ClusterKind() != PostKind {
		t.Fatalf(
			"expected kind %s, got %s",
			PostKind,
			comparable.ClusterKind(),
		)
	}

	result := comparable.Vector()

	if len(result) != len(vector) {
		t.Fatalf(
			"expected vector length %d, got %d",
			len(vector),
			len(result),
		)
	}

	for index := range vector {
		if result[index] != vector[index] {
			t.Fatalf(
				"expected vector[%d] %f, got %f",
				index,
				vector[index],
				result[index],
			)
		}
	}
}

func TestComparableVectorReturnsCopy(t *testing.T) {
	clusterable := NewMockClusterable(
		PostKind,
	)

	comparable := &comparable{
		clusterable: clusterable,
		vector: []float32{
			0.1,
			0.2,
		},
	}

	vector := comparable.Vector()
	vector[0] = 99

	again := comparable.Vector()

	if again[0] == 99 {
		t.Fatalf("expected vector copy")
	}
}
