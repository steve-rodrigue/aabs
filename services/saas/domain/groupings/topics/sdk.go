package topics

import (
	"time"

	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
	"github.com/steve-rodrigue/aabs/services/saas/domain/posts"
)

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
	Save(topic Topic) error

	FindByID(id uuid.UUID) (Topic, error)
	FindByName(name string) (Topic, error)

	Find(index int, amount int) ([]Topic, error)
	FindAfter(cursor uuid.UUID, amount int) ([]Topic, error)
	Count() (int64, error)

	FindChildren(parent uuid.UUID) ([]Topic, error)
	FindRoots() ([]Topic, error)
}

// Builder represents a topic builder
type Builder interface {
	Build(posts []posts.Post) ([]Topic, error)
}
