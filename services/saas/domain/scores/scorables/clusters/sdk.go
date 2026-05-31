package clusters

import "github.com/google/uuid"

type Kind string

const (
	UserKind      Kind = "user"
	CampaignKind  Kind = "campaign"
	CommunityKind Kind = "community"
	TopicKind     Kind = "topic"
)

// Cluster represents any scorable cluster.
type Cluster interface {
	Identifier() uuid.UUID
	ClusterKind() Kind
}
