package clusters

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/clusterables"
)

var (
	// adapter
	ErrInvalidClusterIdentifier      = errors.New("invalid cluster identifier")
	ErrInvalidClusterTarget          = errors.New("invalid cluster target")
	ErrInvalidClusterMemberID        = errors.New("invalid cluster member id")
	ErrInvalidClusterMemberKind      = errors.New("invalid cluster member kind")
	ErrInvalidClusterConfidenceScore = errors.New("invalid cluster confidence score")
	ErrInvalidClusterCentroid        = errors.New("invalid cluster centroid")
	ErrInvalidClusterCreatedOn       = errors.New("invalid cluster created on")

	// detector
	ErrInvalidClusterDetectorTarget     = errors.New("invalid cluster detector target")
	ErrInvalidClusterDetectorMember     = errors.New("invalid cluster detector member")
	ErrInvalidClusterDetectorComparable = errors.New("invalid cluster detector comparable")
	ErrInvalidClusterDetectorVectorSize = errors.New("invalid cluster detector vector size")
	ErrInvalidClusterDetectorMemberKind = errors.New("invalid cluster detector member kind")
)

// NewAdapter creates a new cluster adapter
func NewAdapter(
	clusterables clusterables.Adapter,
) Adapter {
	return &adapter{
		clusterables: clusterables,
	}
}

// NewDetector creates a new cluster detector
func NewDetector(
	adapter Adapter,
	comparables clusterables.ComparableRepository,
) Detector {
	return &detector{
		adapter:     adapter,
		comparables: comparables,
	}
}

// ClusterInput represents a cluster input
type ClusterInput struct {
	Identifier uuid.UUID

	Target clusterables.ClusterableInput

	MemberIDs  []uuid.UUID
	MemberKind clusterables.Kind

	ConfidenceScore float64
	Centroid        []float32

	CreatedOn time.Time
}

// Adapter represents a cluster adapter
type Adapter interface {
	ToDomain(input ClusterInput) (Cluster, error)
}

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
	Save(ctx context.Context, cluster Cluster) error

	FindByID(ctx context.Context, id uuid.UUID) (Cluster, error)
	FindByTarget(ctx context.Context, target uuid.UUID) ([]Cluster, error)
	FindByMember(ctx context.Context, member uuid.UUID) ([]Cluster, error)

	Find(ctx context.Context, index int, amount int) ([]Cluster, error)
	FindAfter(ctx context.Context, cursor uuid.UUID, amount int) ([]Cluster, error)
	Count(ctx context.Context) (int64, error)
}

// Detector represents a cluster detector
type Detector interface {
	Detect(
		ctx context.Context,
		target clusterables.Clusterable,
		members []clusterables.Clusterable,
	) ([]Cluster, error)
}
