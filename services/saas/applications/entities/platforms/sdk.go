package platforms

import (
	"context"

	"github.com/google/uuid"

	domain_platforms "github.com/steve-rodrigue/aabs/services/saas/domain/entities/platforms"
)

// New creates a new platforms application
func New(
	repository domain_platforms.Repository,
) Application {
	return createApplication(repository)
}

// Application represents the platforms application
type Application interface {
	Save(ctx context.Context, platform domain_platforms.Platform) error

	FindByID(ctx context.Context, id uuid.UUID) (domain_platforms.Platform, error)
	FindByHandle(ctx context.Context, handle string) (domain_platforms.Platform, error)

	Find(ctx context.Context, index int, amount int) ([]domain_platforms.Platform, error)
	FindAfter(ctx context.Context, cursor uuid.UUID, amount int) ([]domain_platforms.Platform, error)

	Count(ctx context.Context) (int64, error)
}
