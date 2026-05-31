package scorables

import "github.com/google/uuid"

type Kind string

const (
	UserKind         Kind = "user"
	PostKind         Kind = "post"
	CampaignKind     Kind = "campaign"
	CommunityKind    Kind = "community"
	TopicKind        Kind = "topic"
	NarrativeKind    Kind = "narrative"
	RelationshipKind Kind = "relationship"
	ClusterKind      Kind = "cluster"
)

// Scorable represents an entity that can receive scores.
type Scorable interface {
	Identifier() uuid.UUID
	ScoreKind() Kind
}
