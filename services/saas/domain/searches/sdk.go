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

// SearchableInput represents a searchable input
type SearchableInput struct {
	Identifier uuid.UUID
	Kind       Kind
	Title      string
	Text       string
}

// Searchable represents an entity that can be indexed and displayed in search results.
type Searchable interface {
	Identifier() uuid.UUID
	SearchKind() Kind
	SearchTitle() string
	SearchText() string
}

// SearchableAdapter represents a searchable adapter
type SearchableAdapter interface {
	ToDomain(input SearchableInput) (Searchable, error)
}

// SearchableRepository represents a searchable repository
type SearchableRepository interface {
	FindByID(id uuid.UUID) (Searchable, error)

	Find(index int, amount int) ([]Searchable, error)
	FindAfter(cursor uuid.UUID, amount int) ([]Searchable, error)

	Count() (int64, error)
}

// Repository represents a search repository
type Repository interface {
	Store(target uuid.UUID, kind Kind, vector []float32) error
	Search(vector []float32, limit int) ([]Match, error)
}
