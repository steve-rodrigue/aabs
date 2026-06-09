package topics

import (
	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/communities"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/users"
	domain_topics "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/topics"
)

func NewMockTopicsApplication() *MockTopicsApplication {
	return &MockTopicsApplication{}
}

type MockTopicsApplication struct {
	FindByIDCalls int
	FindByIDErr   error
	FindByIDValue domain_topics.Topic

	FindCalls int
	FindErr   error
	FindValue []domain_topics.Topic

	FindAfterCalls int
	FindAfterErr   error
	FindAfterValue []domain_topics.Topic

	FindTopicsByUserCalls int
	FindTopicsByUserErr   error
	FindTopicsByUserValue []domain_topics.Topic

	FindTopicsByCommunityCalls int
	FindTopicsByCommunityErr   error
	FindTopicsByCommunityValue []domain_topics.Topic

	CountCalls int
	CountErr   error
	CountValue int64

	RebuildTopicsCalls int
	RebuildTopicsErr   error
}

func (application *MockTopicsApplication) FindByID(
	id uuid.UUID,
) (domain_topics.Topic, error) {
	application.FindByIDCalls++

	return application.FindByIDValue, application.FindByIDErr
}

func (application *MockTopicsApplication) Find(
	index int,
	amount int,
) ([]domain_topics.Topic, error) {
	application.FindCalls++

	return application.FindValue, application.FindErr
}

func (application *MockTopicsApplication) FindAfter(
	cursor uuid.UUID,
	amount int,
) ([]domain_topics.Topic, error) {
	application.FindAfterCalls++

	return application.FindAfterValue, application.FindAfterErr
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

	return application.FindTopicsByCommunityValue,
		application.FindTopicsByCommunityErr
}

func (application *MockTopicsApplication) Count() (int64, error) {
	application.CountCalls++

	return application.CountValue, application.CountErr
}

func (application *MockTopicsApplication) RebuildTopics() error {
	application.RebuildTopicsCalls++

	return application.RebuildTopicsErr
}
