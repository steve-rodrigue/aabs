package topics

import (
	"time"

	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters"
)

// Topic represents a semantic subject that can a parent topic
type Topic interface {
	Identifier() uuid.UUID
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
	FindChildren(parent uuid.UUID) ([]Topic, error)
	FindRoots() ([]Topic, error)
}
