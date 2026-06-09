package clusters

import (
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters/clusterables"
)

type cluster struct {
	identifier uuid.UUID

	target clusterables.Clusterable

	memberIDs  []uuid.UUID
	memberKind clusterables.Kind

	confidenceScore float64
	centroid        []float32

	createdOn time.Time
}

func (cluster *cluster) Identifier() uuid.UUID {
	return cluster.identifier
}

func (cluster *cluster) Target() clusterables.Clusterable {
	return cluster.target
}

func (cluster *cluster) MemberIDs() []uuid.UUID {
	out := make([]uuid.UUID, len(cluster.memberIDs))
	copy(out, cluster.memberIDs)

	return out
}

func (cluster *cluster) MemberKind() clusterables.Kind {
	return cluster.memberKind
}

func (cluster *cluster) ConfidenceScore() float64 {
	return cluster.confidenceScore
}

func (cluster *cluster) Centroid() []float32 {
	out := make([]float32, len(cluster.centroid))
	copy(out, cluster.centroid)

	return out
}

func (cluster *cluster) CreatedOn() time.Time {
	return cluster.createdOn
}
