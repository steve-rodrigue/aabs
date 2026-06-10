package campaigns

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/clusterables"
	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/participatables"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters"
)

var (
	// adapter
	ErrInvalidCampaignIdentifier = errors.New("invalid campaign identifier")
	ErrInvalidCampaignName       = errors.New("invalid campaign name")
	ErrInvalidCampaignCluster    = errors.New("invalid campaign cluster")
	ErrInvalidCampaignPostCount  = errors.New("invalid campaign post count")
	ErrInvalidCampaignConfidence = errors.New("invalid campaign confidence")
	ErrInvalidCampaignCreatedOn  = errors.New("invalid campaign created on")

	// classifier
	ErrInvalidCampaignClassifierPost       = errors.New("invalid campaign classifier post")
	ErrInvalidCampaignClassifierText       = errors.New("invalid campaign classifier text")
	ErrInvalidCampaignClassifierVector     = errors.New("invalid campaign classifier vector")
	ErrInvalidCampaignClassifierComparable = errors.New("invalid campaign classifier comparable")
)

// NewAdapter creates a new campaign adapter
func NewAdapter() Adapter {
	return &adapter{}
}

// CampaignInput represents a campaign input
type CampaignInput struct {
	Identifier uuid.UUID

	Name        string
	Description string

	Cluster clusters.Cluster

	PostCount  int
	Confidence float64

	CreatedOn time.Time
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
	Save(ctx context.Context, campaign Campaign) error

	FindByID(ctx context.Context, id uuid.UUID) (Campaign, error)
	FindByName(ctx context.Context, name string) (Campaign, error)

	Find(ctx context.Context, index int, amount int) ([]Campaign, error)
	FindAfter(ctx context.Context, cursor uuid.UUID, amount int) ([]Campaign, error)

	Count(ctx context.Context) (int64, error)
}

// Classifier represents a campaign classifier
type Classifier interface {
	Classify(
		ctx context.Context,
		post posts.Post,
	) (Campaign, float64, error)
}

// Detector discovers and creates campaigns from clusters, posts, or participants
type Detector interface {
	Detect(
		ctx context.Context,
		candidates []clusterables.Clusterable,
	) ([]Campaign, error)
}
