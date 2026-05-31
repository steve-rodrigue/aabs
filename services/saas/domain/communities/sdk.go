package communities

import (
	"time"

	"github.com/google/uuid"
)

// Community represents a community
type Community interface {
	Identifier() uuid.UUID
	Handle() string
	Title() string
	Text() string
	CreatedOn() time.Time
}

// Repository represents a community repository
type Repository interface {
	Save(community Community) error
	FindByID(id uuid.UUID) (Community, error)
	FindByHandle(handle string) (Community, error)
}
