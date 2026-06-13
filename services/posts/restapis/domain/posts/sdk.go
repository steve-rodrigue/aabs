package posts

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts/contents"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/users"
)

var (
	ErrInvalidPostIdentifier  = errors.New("invalid post identifier")
	ErrInvalidPostCreator     = errors.New("invalid post creator")
	ErrInvalidPostContent     = errors.New("invalid post content")
	ErrInvalidPostCommunityID = errors.New("invalid post community id")
	ErrInvalidPostCreatedOn   = errors.New("invalid post created on")
)

// NewService creates a post service that executes sub-services in order
func NewService(
	services ...Service,
) Service {
	return &service{
		services: services,
	}
}

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

// Criteria represents filters used to search posts
type Criteria struct {
	UserIDs      []uuid.UUID
	CommunityIDs []uuid.UUID
	PlatformIDs  []uuid.UUID
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

// Service represents a post service
type Service interface {
	Save(ctx context.Context, post Post) error
}

// Repository represents a post repository
type Repository interface {
	FindByID(ctx context.Context, id uuid.UUID) (Post, error)

	Find(ctx context.Context, index int, amount int) ([]Post, error)
	FindAfter(ctx context.Context, cursor uuid.UUID, amount int) ([]Post, error)

	FindByCriteria(
		ctx context.Context,
		criteria Criteria,
		index int,
		amount int,
	) ([]Post, error)

	FindByCriteriaAfter(
		ctx context.Context,
		criteria Criteria,
		cursor uuid.UUID,
		amount int,
	) ([]Post, error)

	Count(ctx context.Context) (int64, error)

	CountByCriteria(
		ctx context.Context,
		criteria Criteria,
	) (int64, error)
}
