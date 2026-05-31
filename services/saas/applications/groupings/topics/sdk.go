package topics

import (
	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/communities"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/topics"
	"github.com/steve-rodrigue/aabs/services/saas/domain/users"
)

// Application represents the topics application
type Application interface {
	FindByID(id uuid.UUID) (topics.Topic, error)
	FindAll() ([]topics.Topic, error)
	FindTopicsByUser(user users.User) ([]topics.Topic, error)
	FindTopicsByCommunity(community communities.Community) ([]topics.Topic, error)
	RebuildTopics() error
}
