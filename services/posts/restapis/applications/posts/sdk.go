package posts

import (
	"context"

	"github.com/google/uuid"

	domain_posts "github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts"
)

// New creates a new posts application
func New(
	repository domain_posts.Repository,
	service domain_posts.Service,
) Application {
	return createApplication(repository, service)
}

// Application represents the posts application
type Application interface {
	Save(ctx context.Context, post domain_posts.Post) error

	FindByID(ctx context.Context, id uuid.UUID) (domain_posts.Post, error)

	Find(ctx context.Context, index int, amount int) ([]domain_posts.Post, error)
	FindAfter(ctx context.Context, cursor uuid.UUID, amount int) ([]domain_posts.Post, error)

	FindByCriteria(
		ctx context.Context,
		criteria domain_posts.Criteria,
		index int,
		amount int,
	) ([]domain_posts.Post, error)

	FindByCriteriaAfter(
		ctx context.Context,
		criteria domain_posts.Criteria,
		cursor uuid.UUID,
		amount int,
	) ([]domain_posts.Post, error)

	Count(ctx context.Context) (int64, error)

	CountByCriteria(
		ctx context.Context,
		criteria domain_posts.Criteria,
	) (int64, error)
}
