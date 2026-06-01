package topics

import (
	"time"

	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters"
)

func NewMockTopic(
	name string,
	description string,
) Topic {
	return &MockTopic{
		id:          uuid.New(),
		name:        name,
		description: description,
	}
}

type MockTopic struct {
	id          uuid.UUID
	cluster     clusters.Cluster
	name        string
	description string
	parent      Topic
}

func (topic *MockTopic) Identifier() uuid.UUID {
	return topic.id
}

func (topic *MockTopic) Cluster() clusters.Cluster {
	return topic.cluster
}

func (topic *MockTopic) Name() string {
	return topic.name
}

func (topic *MockTopic) Description() string {
	return topic.description
}

func (topic *MockTopic) CreatedOn() time.Time {
	return time.Time{}
}

func (topic *MockTopic) HasParent() bool {
	return topic.parent != nil
}

func (topic *MockTopic) Parent() Topic {
	return topic.parent
}

type MockTopicRepository struct {
	SaveCalls int
	SaveErr   error

	Items map[uuid.UUID]Topic

	FindByIDCalls int
	FindByIDErr   error

	FindByNameCalls int
	FindByNameErr   error

	FindChildrenCalls int
	FindChildrenErr   error

	FindRootsCalls int
	FindRootsErr   error
}

func (repository *MockTopicRepository) Save(topic Topic) error {
	repository.SaveCalls++

	return repository.SaveErr
}

func (repository *MockTopicRepository) FindByID(id uuid.UUID) (Topic, error) {
	repository.FindByIDCalls++

	if repository.FindByIDErr != nil {
		return nil, repository.FindByIDErr
	}

	return repository.Items[id], nil
}

func (repository *MockTopicRepository) FindByName(name string) (Topic, error) {
	repository.FindByNameCalls++

	if repository.FindByNameErr != nil {
		return nil, repository.FindByNameErr
	}

	for _, topic := range repository.Items {
		if topic.Name() == name {
			return topic, nil
		}
	}

	return nil, nil
}

func (repository *MockTopicRepository) FindChildren(parent uuid.UUID) ([]Topic, error) {
	repository.FindChildrenCalls++

	if repository.FindChildrenErr != nil {
		return nil, repository.FindChildrenErr
	}

	return []Topic{}, nil
}

func (repository *MockTopicRepository) FindRoots() ([]Topic, error) {
	repository.FindRootsCalls++

	if repository.FindRootsErr != nil {
		return nil, repository.FindRootsErr
	}

	return []Topic{}, nil
}
