package searches

import "github.com/google/uuid"

type Kind string

const (
	PostKind         Kind = "post"
	CampaignKind     Kind = "campaign"
	TopicKind        Kind = "topic"
	NarrativeKind    Kind = "narrative"
	UserKind         Kind = "user"
	CommunityKind    Kind = "community"
	RelationshipKind Kind = "relationship"
)

// MatchInput represents a match input
type MatchInput struct {
	Target     uuid.UUID
	Kind       Kind
	Similarity float64
}

// Adapter represents a match adapter
type Adapter interface {
	ToDomain(input MatchInput) (Match, error)
}

// Match represents a match
type Match interface {
	Target() uuid.UUID
	Kind() Kind
	Similarity() float64
}

// Repository represents a search repository
type Repository interface {
	Store(target uuid.UUID, kind Kind, vector []float32) error
	Search(vector []float32, limit int) ([]Match, error)
}
