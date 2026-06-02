package searches

import (
	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/embeddings"
	domain_searches "github.com/steve-rodrigue/aabs/services/saas/domain/searches"
)

type ResultKind string

const (
	PostKind         ResultKind = "post"
	CampaignKind     ResultKind = "campaign"
	TopicKind        ResultKind = "topic"
	NarrativeKind    ResultKind = "narrative"
	UserKind         ResultKind = "user"
	CommunityKind    ResultKind = "community"
	RelationshipKind ResultKind = "relationship"
)

// New creates a new search application
func New(
	embedder embeddings.Embedder,
	searchRepository domain_searches.Repository,
	searchableRepository domain_searches.SearchableRepository,
) Application {
	return createApplication(
		embedder,
		searchRepository,
		searchableRepository,
	)
}

// Result represents a search result
type Result interface {
	Identifier() uuid.UUID
	Kind() ResultKind
	HasTitle() bool
	Title() string
	Text() string
	Score() float64
}

// Application represents a search application
type Application interface {
	Index(searchable domain_searches.Searchable) error
	Search(query string, limit int) ([]Result, error)
}
