package topics

import (
	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/communities"
	domain_topics "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/topics"
	"github.com/steve-rodrigue/aabs/services/saas/domain/users"
)

func NewMockTopicsApplication() *MockTopicsApplication {
	return &MockTopicsApplication{}
}

type MockTopicsApplication struct {
	FindByIDCalls int
	FindByIDErr   error
	FindByIDValue domain_topics.Topic

	FindAllCalls int
	FindAllErr   error
	FindAllValue []domain_topics.Topic

	FindTopicsByUserCalls int
	FindTopicsByUserErr   error
	FindTopicsByUserValue []domain_topics.Topic

	FindTopicsByCommunityCalls int
	FindTopicsByCommunityErr   error
	FindTopicsByCommunityValue []domain_topics.Topic

	RebuildTopicsCalls int
	RebuildTopicsErr   error
}

func (application *MockTopicsApplication) FindByID(
	id uuid.UUID,
) (domain_topics.Topic, error) {
	application.FindByIDCalls++

	return application.FindByIDValue, application.FindByIDErr
}

func (application *MockTopicsApplication) FindAll() (
	[]domain_topics.Topic,
	error,
) {
	application.FindAllCalls++

	return application.FindAllValue, application.FindAllErr
}

func (application *MockTopicsApplication) FindTopicsByUser(
	user users.User,
) ([]domain_topics.Topic, error) {
	application.FindTopicsByUserCalls++

	return application.FindTopicsByUserValue, application.FindTopicsByUserErr
}

func (application *MockTopicsApplication) FindTopicsByCommunity(
	community communities.Community,
) ([]domain_topics.Topic, error) {
	application.FindTopicsByCommunityCalls++

	return application.FindTopicsByCommunityValue, application.FindTopicsByCommunityErr
}

func (application *MockTopicsApplication) RebuildTopics() error {
	application.RebuildTopicsCalls++

	return application.RebuildTopicsErr
}
