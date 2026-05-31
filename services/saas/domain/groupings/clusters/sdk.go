package clusters

import (
	"time"

	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters/clusterables"
)

// Cluster represents a cluster
type Cluster interface {
	Identifier() uuid.UUID

	Target() clusterables.Clusterable

	MemberIDs() []uuid.UUID
	MemberKind() clusterables.Kind

	ConfidenceScore() float64
	Centroid() []float32

	CreatedOn() time.Time
}

// Repository represents a cluster repository
type Repository interface {
	Save(cluster Cluster) error
	FindByID(id uuid.UUID) (Cluster, error)
	FindByTarget(target uuid.UUID) ([]Cluster, error)
	FindByMember(member uuid.UUID) ([]Cluster, error)
}

// Detector represents a cluster detector
type Detector interface {
	Detect(
		target clusterables.Clusterable,
		members []clusterables.Clusterable,
	) ([]Cluster, error)
}
