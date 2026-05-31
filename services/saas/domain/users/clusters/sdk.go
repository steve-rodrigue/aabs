package clusters

import "github.com/google/uuid"

// Cluster represents a user post's cluster
type Cluster interface {
	Identifier() uuid.UUID
	MemberPostIDs() []uuid.UUID
	ConfidenceScore() float64
	Centroid() []float32
}
