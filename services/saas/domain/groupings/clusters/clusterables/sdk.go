package clusterables

import "github.com/google/uuid"

type Kind string

const (
	PostKind      Kind = "post"
	UserKind      Kind = "user"
	CommunityKind Kind = "community"
	PlatformKind  Kind = "platform"
	CampaignKind  Kind = "campaign"
	TopicKind     Kind = "topic"
	NarrativeKind Kind = "narrative"
)

// Clusterable represents a clusterable entity
type Clusterable interface {
	Identifier() uuid.UUID
	ClusterKind() Kind
}
