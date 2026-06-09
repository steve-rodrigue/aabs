package clusterables

import "github.com/google/uuid"

type adapter struct{}

func (adapter *adapter) ToDomain(
	input ClusterableInput,
) (Clusterable, error) {
	if input.Identifier == uuid.Nil {
		return nil, ErrInvalidClusterableIdentifier
	}

	if !isValidKind(input.ClusterKind) {
		return nil, ErrInvalidClusterableKind
	}

	return &clusterable{
		identifier:  input.Identifier,
		clusterKind: input.ClusterKind,
	}, nil
}

func isValidKind(
	kind Kind,
) bool {
	switch kind {
	case PostKind,
		UserKind,
		CommunityKind,
		PlatformKind,
		CampaignKind,
		TopicKind,
		NarrativeKind:
		return true

	default:
		return false
	}
}
