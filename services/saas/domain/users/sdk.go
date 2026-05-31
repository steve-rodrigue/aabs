package users

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user
type User interface {
	Identifier() uuid.UUID
	ID() string
	Handle() string
	CreatedOn() time.Time
}

// Repository represents a user repository
type Repository interface {
	Save(user User) error
	FindByID(id uuid.UUID) (User, error)
	FindByExternalID(id string) (User, error)
	FindByHandle(handle string) (User, error)
}
