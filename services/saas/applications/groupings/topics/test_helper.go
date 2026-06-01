package topics

import (
	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/communities"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/topics"
	"github.com/steve-rodrigue/aabs/services/saas/domain/users"
)

type MockTopicsApplication struct {
	RebuildTopicsCalls int
	RebuildTopicsErr   error
}

func (application *MockTopicsApplication) FindByID(id uuid.UUID) (topics.Topic, error) {
	return nil, nil
}

func (application *MockTopicsApplication) FindAll() ([]topics.Topic, error) {
	return nil, nil
}

func (application *MockTopicsApplication) FindTopicsByUser(user users.User) ([]topics.Topic, error) {
	return nil, nil
}

func (application *MockTopicsApplication) FindTopicsByCommunity(community communities.Community) ([]topics.Topic, error) {
	return nil, nil
}

func (application *MockTopicsApplication) RebuildTopics() error {
	application.RebuildTopicsCalls++

	return application.RebuildTopicsErr
}
