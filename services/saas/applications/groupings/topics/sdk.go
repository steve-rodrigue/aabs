package topics

import (
	"github.com/google/uuid"

	app_participations "github.com/steve-rodrigue/aabs/services/saas/applications/groupings/participations"
	app_posts "github.com/steve-rodrigue/aabs/services/saas/applications/posts"
	"github.com/steve-rodrigue/aabs/services/saas/domain/communities"
	domain_topics "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/topics"
	"github.com/steve-rodrigue/aabs/services/saas/domain/users"
)

// New creates a new topic application
func New(
	repository domain_topics.Repository,
	posts app_posts.Application,
	participations app_participations.Application,
	builder domain_topics.Builder,
) Application {
	return createApplication(repository, posts, participations, builder)
}

// Application represents a topic application
type Application interface {
	FindByID(id uuid.UUID) (domain_topics.Topic, error)

	Find(index int, amount int) ([]domain_topics.Topic, error)
	FindAfter(cursor uuid.UUID, amount int) ([]domain_topics.Topic, error)

	FindTopicsByUser(user users.User) ([]domain_topics.Topic, error)
	FindTopicsByCommunity(community communities.Community) ([]domain_topics.Topic, error)

	Count() (int64, error)

	RebuildTopics() error
}
