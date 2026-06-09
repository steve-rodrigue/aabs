package posts

import (
	"context"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/communities"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/platforms"
	domain_posts "github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/users"
)

// New creates a new posts application
func New(
	repository domain_posts.Repository,
) Application {
	return createApplication(repository)
}

// Application represents the posts application
type Application interface {
	Save(ctx context.Context, post domain_posts.Post) error

	FindByID(ctx context.Context, id uuid.UUID) (domain_posts.Post, error)

	Find(ctx context.Context, index int, amount int) ([]domain_posts.Post, error)
	FindAfter(ctx context.Context, cursor uuid.UUID, amount int) ([]domain_posts.Post, error)

	Count(ctx context.Context) (int64, error)

	FindByUser(ctx context.Context, user users.User) ([]domain_posts.Post, error)
	FindByCommunity(ctx context.Context, community communities.Community) ([]domain_posts.Post, error)
	FindByPlatform(ctx context.Context, platform platforms.Platform) ([]domain_posts.Post, error)
}
