package relatables

import "github.com/google/uuid"

type Kind string

const (
	CampaignKind  Kind = "campaign"
	TopicKind     Kind = "topic"
	UserKind      Kind = "user"
	PostKind      Kind = "post"
	NarrativeKind Kind = "narrative"
)

type Relatable interface {
	Identifier() uuid.UUID
	RelationshipKind() Kind
}

type Repository interface {
	FindAll() ([]Relatable, error)
}
