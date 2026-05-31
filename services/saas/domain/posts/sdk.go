package posts

import (
	"time"

	"github.com/google/uuid"
)

// Post represents a post
type Post interface {
	Identifier() uuid.UUID
	UserID() string
	Text() string
	CreatedOn() time.Time
}

// Repository represents a post repository
type Repository interface {
	Save(post Post) error
	FindByID(id uuid.UUID) (Post, error)
}

// Processor represents a post processor
type Processor interface {
	Process(post Post) error
}
