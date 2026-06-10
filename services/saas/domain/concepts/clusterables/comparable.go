package clusterables

import (
	"github.com/google/uuid"
)

type comparable struct {
	clusterable Clusterable
	vector      []float32
}

func (comparable *comparable) Identifier() uuid.UUID {
	return comparable.clusterable.Identifier()
}

func (comparable *comparable) ClusterKind() Kind {
	return comparable.clusterable.ClusterKind()
}

func (comparable *comparable) Vector() []float32 {
	out := make([]float32, len(comparable.vector))
	copy(out, comparable.vector)

	return out
}
