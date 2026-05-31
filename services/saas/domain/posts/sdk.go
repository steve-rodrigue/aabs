package posts

import (
	"time"

	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/communities"
	"github.com/steve-rodrigue/aabs/services/saas/domain/posts/contents"
	"github.com/steve-rodrigue/aabs/services/saas/domain/users"
)

// Post represents a post
type Post interface {
	Identifier() uuid.UUID
	Community() communities.Community
	Creator() users.User
	Content() contents.Content
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
