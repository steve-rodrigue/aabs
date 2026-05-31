package posts

import (
	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/communities"
	"github.com/steve-rodrigue/aabs/services/saas/domain/platforms"
	domain_posts "github.com/steve-rodrigue/aabs/services/saas/domain/posts"
	"github.com/steve-rodrigue/aabs/services/saas/domain/users"
)

// Application represents the posts application
type Application interface {
	Save(post domain_posts.Post) error

	FindByID(id uuid.UUID) (domain_posts.Post, error)
	FindAll() ([]domain_posts.Post, error)

	FindByUser(user users.User) ([]domain_posts.Post, error)
	FindByCommunity(community communities.Community) ([]domain_posts.Post, error)
	FindByPlatform(platform platforms.Platform) ([]domain_posts.Post, error)
}
