package campaigns

import (
	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/campaigns/clusters"
	"github.com/steve-rodrigue/aabs/services/saas/domain/posts"
)

// Campaign represents a campaign
type Campaign interface {
	Identifier() uuid.UUID
	Name() string
	Cluster() clusters.Cluster
	PostCount() int
	Confidence() float64
}

// Repository represents a campaign repository
type Repository interface {
	Save(campaign Campaign) error
	FindByID(id uuid.UUID) (Campaign, error)
	FindByName(name string) (Campaign, error)
}

// Relationship a relationship
type Relationship interface {
	Source() Campaign
	Target() Campaign
	Similarity() float64
}

// GraphRepository represents a graph repository of the campaign
type GraphRepository interface {
	SaveCampaign(campaign Campaign) error
	SaveRelationship(relationship Relationship) error
}

// Match represents a match
type Match interface {
	Post() uuid.UUID
	Similarity() float64
}

// SearchRepository represents a search repository
type SearchRepository interface {
	Store(
		post uuid.UUID,
		vector []float32,
	) error

	Search(
		vector []float32,
		limit int,
	) ([]Match, error)
}

// Classifier represents a campaign classifier
type Classifier interface {
	Classify(
		post posts.Post,
	) (Campaign, float64, error)
}
