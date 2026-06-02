package participatables

import "github.com/google/uuid"

type Kind string

const (
	UserKind      Kind = "user"
	CommunityKind Kind = "community"
	PlatformKind  Kind = "platform"

	CampaignKind  Kind = "campaign"
	TopicKind     Kind = "topic"
	NarrativeKind Kind = "narrative"

	ClusterKind Kind = "cluster"
)

// Participatable represents an entity that can participate in another entity.
type Participatable interface {
	Identifier() uuid.UUID
	ParticipationKind() Kind
}

// Repository represents a participatable repository
type Repository interface {
	FindAllParticipants() ([]Participatable, error)
	FindAllTargets() ([]Participatable, error)
}
