package campaigns

import (
	"time"

	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
	"github.com/steve-rodrigue/aabs/services/saas/domain/posts"
)

// CampaignInput represents a campaign input
type CampaignInput struct {
	Identifier        uuid.UUID
	ParticipationKind participatables.Kind
	Name              string
	Description       string
	Cluster           clusters.Cluster
	PostCount         int
	Confidence        float64
	CreatedOn         time.Time
}

// Adapter represents the campaign adapter
type Adapter interface {
	ToDomain(input CampaignInput) (Campaign, error)
}

// Campaign represents a campaign
type Campaign interface {
	Identifier() uuid.UUID
	ParticipationKind() participatables.Kind
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

	Find(index int, amount int) ([]Campaign, error)
	FindAfter(cursor uuid.UUID, amount int) ([]Campaign, error)

	Count() (int64, error)
}

// Classifier represents a campaign classifier
type Classifier interface {
	Classify(post posts.Post) (Campaign, float64, error)
}
