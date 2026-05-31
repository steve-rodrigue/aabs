package campaigns

import (
	"github.com/google/uuid"
	domain_campaigns "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/campaigns"
	domain_posts "github.com/steve-rodrigue/aabs/services/saas/domain/posts"
)

// Builder builds or updates campaigns from posts.
type Builder interface {
	BuildFromPosts(posts []domain_posts.Post) ([]domain_campaigns.Campaign, error)
}

// Classifier assigns a post to an existing campaign.
type Classifier interface {
	Classify(post domain_posts.Post) (domain_campaigns.Campaign, float64, error)
}

// Reader reads campaign data for UI/API use.
type Reader interface {
	FindByID(id uuid.UUID) (domain_campaigns.Campaign, error)
	FindAll() ([]domain_campaigns.Campaign, error)
}
