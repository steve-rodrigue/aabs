package topics

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/participatables"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/namers"
)

var (
	// adapter
	ErrInvalidTopicIdentifier = errors.New("invalid topic identifier")
	ErrInvalidTopicCluster    = errors.New("invalid topic cluster")
	ErrInvalidTopicName       = errors.New("invalid topic name")
	ErrInvalidTopicCreatedOn  = errors.New("invalid topic created on")

	// builder
	ErrInvalidTopicBuilderPosts   = errors.New("invalid topic builder posts")
	ErrInvalidTopicBuilderCluster = errors.New("invalid topic builder cluster")
)

// NewAdapter creates a new topic adapter
func NewAdapter() Adapter {
	return &adapter{}
}

// NewBuilder creates a new topic builder
func NewBuilder(
	adapter Adapter,
	namer namers.Namer,
) Builder {
	return &builder{
		adapter: adapter,
		namer:   namer,
	}
}

// TopicInput represents a topic input
type TopicInput struct {
	Identifier  uuid.UUID
	Cluster     clusters.Cluster
	Name        string
	Description string
	Parent      Topic
	CreatedOn   time.Time
}

// Adapter represents a topic adapter
type Adapter interface {
	ToDomain(input TopicInput) (Topic, error)
}

// Topic represents a semantic subject that can have a parent topic
type Topic interface {
	Identifier() uuid.UUID
	ParticipationKind() participatables.Kind
	Cluster() clusters.Cluster
	Name() string
	Description() string
	CreatedOn() time.Time
	HasParent() bool
	Parent() Topic
}

// Repository represents a topic repository
type Repository interface {
	Save(ctx context.Context, topic Topic) error

	FindByID(ctx context.Context, id uuid.UUID) (Topic, error)
	FindByName(ctx context.Context, name string) (Topic, error)

	Find(ctx context.Context, index int, amount int) ([]Topic, error)
	FindAfter(ctx context.Context, cursor uuid.UUID, amount int) ([]Topic, error)
	Count(ctx context.Context) (int64, error)

	FindChildren(ctx context.Context, parent uuid.UUID) ([]Topic, error)
	FindRoots(ctx context.Context) ([]Topic, error)
}

// Builder represents a topic builder
type Builder interface {
	Build(ctx context.Context, posts []posts.Post) ([]Topic, error)
}
