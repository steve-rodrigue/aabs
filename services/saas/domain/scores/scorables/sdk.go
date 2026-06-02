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

// ScorableInput represents a scorable input
type ScorableInput struct {
	Identifier uuid.UUID
	ScoreKind  Kind
}

// Adapter represents a scorable adapter
type Adapter interface {
	ToDomain(input ScorableInput) (Scorable, error)
}

// Scorable represents an entity that can receive scores.
type Scorable interface {
	Identifier() uuid.UUID
	ScoreKind() Kind
}

// Repository represents a scorable repository
type Repository interface {
	Find(index int, amount int) ([]Scorable, error)
	FindAfter(cursor uuid.UUID, amount int) ([]Scorable, error)
	Count() (int64, error)

	FindByID(id uuid.UUID) (Scorable, error)
}
