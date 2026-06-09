package users

import (
	"context"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/platforms"
	domain_users "github.com/steve-rodrigue/aabs/services/saas/domain/entities/users"
)

// New creates a new users application
func New(
	repository domain_users.Repository,
) Application {
	return createApplication(repository)
}

// Application represents the users application
type Application interface {
	Save(ctx context.Context, user domain_users.User) error

	FindByID(ctx context.Context, id uuid.UUID) (domain_users.User, error)
	FindByExternalID(
		ctx context.Context,
		platform platforms.Platform,
		externalID string,
	) (domain_users.User, error)
	FindByHandle(
		ctx context.Context,
		platform platforms.Platform,
		handle string,
	) (domain_users.User, error)

	Find(ctx context.Context, index int, amount int) ([]domain_users.User, error)
	FindAfter(ctx context.Context, cursor uuid.UUID, amount int) ([]domain_users.User, error)

	Count(ctx context.Context) (int64, error)
}
