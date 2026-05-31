package narratives

import (
	domain_narratives "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/narratives"
	domain_topics "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/topics"
)

// Builder builds narratives from topics
type Builder interface {
	BuildFromTopics(topics []domain_topics.Topic) ([]domain_narratives.Narrative, error)
}
