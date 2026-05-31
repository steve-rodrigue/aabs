package campaigns

import (
	"time"

	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/campaigns/clusters"
	"github.com/steve-rodrigue/aabs/services/saas/domain/posts"
)

// Campaign represents a campaign
type Campaign interface {
	Identifier() uuid.UUID
	Name() string
	Description() string
	Cluster() clusters.Cluster
	PostCount() int
	Confidence() float64
	CreatedOn() time.Time
}

// Repository represents a campaign repository
type Repository interface {
	Save(campaign Campaign) error
	FindByID(id uuid.UUID) (Campaign, error)
	FindByName(name string) (Campaign, error)
}

// Classifier represents a campaign classifier
type Classifier interface {
	Classify(post posts.Post) (Campaign, float64, error)
}
