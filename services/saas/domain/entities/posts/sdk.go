package posts

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/communities"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/platforms"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts/contents"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/users"
)

var (
	ErrInvalidPostIdentifier  = errors.New("invalid post identifier")
	ErrInvalidPostCreator     = errors.New("invalid post creator")
	ErrInvalidPostContent     = errors.New("invalid post content")
	ErrInvalidPostCommunityID = errors.New("invalid post community id")
	ErrInvalidPostCreatedOn   = errors.New("invalid post created on")
)

// NewAdapter creates a new post adapter
func NewAdapter(
	contents contents.Adapter,
) Adapter {
	return &adapter{
		contents: contents,
	}
}

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
	Save(ctx context.Context, post Post) error

	FindByID(ctx context.Context, id uuid.UUID) (Post, error)

	Find(ctx context.Context, index int, amount int) ([]Post, error)
	FindAfter(ctx context.Context, cursor uuid.UUID, amount int) ([]Post, error)
	Count(ctx context.Context) (int64, error)

	FindByUser(ctx context.Context, user users.User) ([]Post, error)
	FindByCommunity(ctx context.Context, community communities.Community) ([]Post, error)
	FindByPlatform(ctx context.Context, platform platforms.Platform) ([]Post, error)
}
