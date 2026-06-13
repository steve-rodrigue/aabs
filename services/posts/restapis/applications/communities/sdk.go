package communities

import (
	"context"

	"github.com/google/uuid"

	domain_communities "github.com/steve-rodrigue/aabs/services/posts/restapis/domain/communities"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/platforms"
)

// New creates a new communities application
func New(
	repository domain_communities.Repository,
) Application {
	return createApplication(repository)
}

// Application represents the communities application
type Application interface {
	Save(ctx context.Context, community domain_communities.Community) error

	FindByID(ctx context.Context, id uuid.UUID) (domain_communities.Community, error)
	FindByHandle(
		ctx context.Context,
		platform platforms.Platform,
		handle string,
	) (domain_communities.Community, error)

	Find(ctx context.Context, index int, amount int) ([]domain_communities.Community, error)
	FindAfter(
		ctx context.Context,
		cursor uuid.UUID,
		amount int,
	) ([]domain_communities.Community, error)

	FindByPlatform(
		ctx context.Context,
		platform platforms.Platform,
	) ([]domain_communities.Community, error)

	Count(ctx context.Context) (int64, error)
}
