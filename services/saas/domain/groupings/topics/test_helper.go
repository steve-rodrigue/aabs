package topics

import (
	"context"
	"sort"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/participatables"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters"
)

func NewMockTopicAdapter() *MockTopicAdapter {
	return &MockTopicAdapter{}
}

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

func NewMockTopicRepository() *MockTopicRepository {
	return &MockTopicRepository{
		Items: map[uuid.UUID]Topic{},
	}
}

func NewMockTopicBuilder() *MockTopicBuilder {
	return &MockTopicBuilder{}
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

func (topic *MockTopic) ParticipationKind() participatables.Kind {
	return participatables.TopicKind
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

	FindCalls int
	FindErr   error
	FindValue []Topic

	FindAfterCalls int
	FindAfterErr   error
	FindAfterValue []Topic

	CountCalls int
	CountErr   error
	CountValue int64

	FindChildrenCalls int
	FindChildrenErr   error
	FindChildrenValue []Topic

	FindRootsCalls int
	FindRootsErr   error
	FindRootsValue []Topic

	LastContext context.Context
	LastTopic   Topic
	LastID      uuid.UUID
	LastName    string
	LastIndex   int
	LastAmount  int
	LastCursor  uuid.UUID
	LastParent  uuid.UUID
}

func (repository *MockTopicRepository) Save(
	ctx context.Context,
	topic Topic,
) error {
	repository.SaveCalls++
	repository.LastContext = ctx
	repository.LastTopic = topic

	if repository.Items != nil && topic != nil {
		repository.Items[topic.Identifier()] = topic
	}

	return repository.SaveErr
}

func (repository *MockTopicRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (Topic, error) {
	repository.FindByIDCalls++
	repository.LastContext = ctx
	repository.LastID = id

	if repository.FindByIDErr != nil {
		return nil, repository.FindByIDErr
	}

	return repository.Items[id], nil
}

func (repository *MockTopicRepository) FindByName(
	ctx context.Context,
	name string,
) (Topic, error) {
	repository.FindByNameCalls++
	repository.LastContext = ctx
	repository.LastName = name

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

func (repository *MockTopicRepository) Find(
	ctx context.Context,
	index int,
	amount int,
) ([]Topic, error) {
	repository.FindCalls++
	repository.LastContext = ctx
	repository.LastIndex = index
	repository.LastAmount = amount

	if repository.FindErr != nil {
		return nil, repository.FindErr
	}

	if repository.FindValue != nil {
		return repository.FindValue, nil
	}

	topics := repository.sortedTopics()

	if index >= len(topics) {
		return []Topic{}, nil
	}

	end := index + amount
	if end > len(topics) {
		end = len(topics)
	}

	return topics[index:end], nil
}

func (repository *MockTopicRepository) FindAfter(
	ctx context.Context,
	cursor uuid.UUID,
	amount int,
) ([]Topic, error) {
	repository.FindAfterCalls++
	repository.LastContext = ctx
	repository.LastCursor = cursor
	repository.LastAmount = amount

	if repository.FindAfterErr != nil {
		return nil, repository.FindAfterErr
	}

	if repository.FindAfterValue != nil {
		return repository.FindAfterValue, nil
	}

	topics := repository.sortedTopics()

	start := 0

	if cursor != uuid.Nil {
		for index, topic := range topics {
			if topic.Identifier() == cursor {
				start = index + 1
				break
			}
		}
	}

	if start >= len(topics) {
		return []Topic{}, nil
	}

	end := start + amount
	if end > len(topics) {
		end = len(topics)
	}

	return topics[start:end], nil
}

func (repository *MockTopicRepository) Count(
	ctx context.Context,
) (int64, error) {
	repository.CountCalls++
	repository.LastContext = ctx

	if repository.CountErr != nil {
		return 0, repository.CountErr
	}

	if repository.CountValue != 0 {
		return repository.CountValue, nil
	}

	return int64(len(repository.Items)), nil
}

func (repository *MockTopicRepository) FindChildren(
	ctx context.Context,
	parent uuid.UUID,
) ([]Topic, error) {
	repository.FindChildrenCalls++
	repository.LastContext = ctx
	repository.LastParent = parent

	if repository.FindChildrenErr != nil {
		return nil, repository.FindChildrenErr
	}

	if repository.FindChildrenValue != nil {
		return repository.FindChildrenValue, nil
	}

	out := []Topic{}

	for _, topic := range repository.Items {
		if !topic.HasParent() {
			continue
		}

		if topic.Parent().Identifier() == parent {
			out = append(out, topic)
		}
	}

	return out, nil
}

func (repository *MockTopicRepository) FindRoots(
	ctx context.Context,
) ([]Topic, error) {
	repository.FindRootsCalls++
	repository.LastContext = ctx

	if repository.FindRootsErr != nil {
		return nil, repository.FindRootsErr
	}

	if repository.FindRootsValue != nil {
		return repository.FindRootsValue, nil
	}

	out := []Topic{}

	for _, topic := range repository.Items {
		if !topic.HasParent() {
			out = append(out, topic)
		}
	}

	return out, nil
}

func (repository *MockTopicRepository) sortedTopics() []Topic {
	out := make([]Topic, 0, len(repository.Items))

	for _, topic := range repository.Items {
		out = append(out, topic)
	}

	sort.Slice(out, func(left int, right int) bool {
		return out[left].Identifier().String() <
			out[right].Identifier().String()
	})

	return out
}

type MockTopicBuilder struct {
	BuildCalls int
	BuildErr   error

	LastContext context.Context
	LastPosts   []posts.Post

	BuildValue []Topic
}

func (builder *MockTopicBuilder) Build(
	ctx context.Context,
	posts []posts.Post,
) ([]Topic, error) {
	builder.BuildCalls++
	builder.LastContext = ctx
	builder.LastPosts = posts

	return builder.BuildValue,
		builder.BuildErr
}

type MockTopicAdapter struct {
	ToDomainCalls int
	ToDomainErr   error
	ToDomainValue Topic

	LastInput TopicInput
}

func (adapter *MockTopicAdapter) ToDomain(
	input TopicInput,
) (Topic, error) {
	adapter.ToDomainCalls++
	adapter.LastInput = input

	if adapter.ToDomainErr != nil {
		return nil, adapter.ToDomainErr
	}

	if adapter.ToDomainValue != nil {
		return adapter.ToDomainValue, nil
	}

	return &MockTopic{
		id:          input.Identifier,
		cluster:     input.Cluster,
		name:        input.Name,
		description: input.Description,
		parent:      input.Parent,
	}, nil
}
