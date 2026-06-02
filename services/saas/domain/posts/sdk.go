package posts

import (
	"time"

	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/communities"
	"github.com/steve-rodrigue/aabs/services/saas/domain/platforms"
	"github.com/steve-rodrigue/aabs/services/saas/domain/posts/contents"
	"github.com/steve-rodrigue/aabs/services/saas/domain/users"
)

// PostInput represents a post input
type PostInput struct {
	Identifier   uuid.UUID
	CommunityIDs []uuid.UUID
	Creator      users.User
	Content      contents.ContentInput
	CreatedOn    time.Time
}

// Adapter represents a post adapter
type Adapter interface {
	ToDomain(input PostInput) (Post, error)
}

// Post represents a post
type Post interface {
	Identifier() uuid.UUID
	CommunityIDs() []uuid.UUID
	Creator() users.User
	Content() contents.Content
	CreatedOn() time.Time
}

// Repository represents a post repository
type Repository interface {
	Save(post Post) error

	FindByID(id uuid.UUID) (Post, error)

	Find(index int, amount int) ([]Post, error)
	FindAfter(cursor uuid.UUID, amount int) ([]Post, error)
	Count() (int64, error)

	FindByUser(user users.User) ([]Post, error)
	FindByCommunity(community communities.Community) ([]Post, error)
	FindByPlatform(platform platforms.Platform) ([]Post, error)
}
