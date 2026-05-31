package platforms

import (
	"time"

	"github.com/google/uuid"
)

// Platform represents a website or app where users publish content
type Platform interface {
	Identifier() uuid.UUID
	Name() string
	Handle() string
	BaseURL() string
	CreatedOn() time.Time
}

// Repository represents a platform repository
type Repository interface {
	Save(platform Platform) error
	FindByID(id uuid.UUID) (Platform, error)
	FindByHandle(handle string) (Platform, error)
	FindByName(name string) (Platform, error)
}
