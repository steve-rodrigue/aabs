package clusters

import (
	"math"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/clusterables"
)

type adapter struct {
	clusterables clusterables.Adapter
}

func (adapter *adapter) ToDomain(
	input ClusterInput,
) (Cluster, error) {
	if input.Identifier == uuid.Nil {
		return nil, ErrInvalidClusterIdentifier
	}

	target, err := adapter.clusterables.ToDomain(input.Target)
	if err != nil {
		return nil, err
	}

	if target == nil {
		return nil, ErrInvalidClusterTarget
	}

	if len(input.MemberIDs) == 0 {
		return nil, ErrInvalidClusterMemberID
	}

	for _, memberID := range input.MemberIDs {
		if memberID == uuid.Nil {
			return nil, ErrInvalidClusterMemberID
		}
	}

	if !isValidMemberKind(input.MemberKind) {
		return nil, ErrInvalidClusterMemberKind
	}

	if math.IsNaN(input.ConfidenceScore) ||
		math.IsInf(input.ConfidenceScore, 0) ||
		input.ConfidenceScore < 0 ||
		input.ConfidenceScore > 1 {
		return nil, ErrInvalidClusterConfidenceScore
	}

	if len(input.Centroid) == 0 {
		return nil, ErrInvalidClusterCentroid
	}

	if input.CreatedOn.IsZero() {
		return nil, ErrInvalidClusterCreatedOn
	}

	memberIDs := make([]uuid.UUID, len(input.MemberIDs))
	copy(memberIDs, input.MemberIDs)

	centroid := make([]float32, len(input.Centroid))
	copy(centroid, input.Centroid)

	return &cluster{
		identifier:      input.Identifier,
		target:          target,
		memberIDs:       memberIDs,
		memberKind:      input.MemberKind,
		confidenceScore: input.ConfidenceScore,
		centroid:        centroid,
		createdOn:       input.CreatedOn.UTC(),
	}, nil
}

func isValidMemberKind(
	kind clusterables.Kind,
) bool {
	switch kind {
	case clusterables.PostKind,
		clusterables.UserKind,
		clusterables.CommunityKind,
		clusterables.PlatformKind,
		clusterables.CampaignKind,
		clusterables.TopicKind,
		clusterables.NarrativeKind:
		return true

	default:
		return false
	}
}
