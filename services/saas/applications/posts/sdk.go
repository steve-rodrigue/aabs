package posts

import (
	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/communities"
	"github.com/steve-rodrigue/aabs/services/saas/domain/platforms"
	domain_posts "github.com/steve-rodrigue/aabs/services/saas/domain/posts"
	"github.com/steve-rodrigue/aabs/services/saas/domain/users"
)

// New creates a new posts application
func New(
	repository domain_posts.Repository,
) Application {
	return createApplication(repository)
}

// Application represents the posts application
type Application interface {
	Save(post domain_posts.Post) error

	FindByID(id uuid.UUID) (domain_posts.Post, error)

	Find(index int, amount int) ([]domain_posts.Post, error)
	FindAfter(cursor uuid.UUID, amount int) ([]domain_posts.Post, error)

	Count() (int64, error)

	FindByUser(user users.User) ([]domain_posts.Post, error)
	FindByCommunity(community communities.Community) ([]domain_posts.Post, error)
	FindByPlatform(platform platforms.Platform) ([]domain_posts.Post, error)
}
