package topics

import (
	domain_campaigns "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/campaigns"
	domain_topics "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/topics"
)

// Classifier assigns campaigns to topics
type Classifier interface {
	Classify(campaign domain_campaigns.Campaign) (domain_topics.Topic, float64, error)
}

// Builder builds or updates topics from campaigns.
type Builder interface {
	BuildFromCampaigns(campaigns []domain_campaigns.Campaign) ([]domain_topics.Topic, error)
}
