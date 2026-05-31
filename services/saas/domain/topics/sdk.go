package topics

import (
	"time"

	"github.com/google/uuid"
)

// Topic represents a semantic subject that can a parent topic
type Topic interface {
	Identifier() uuid.UUID
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
