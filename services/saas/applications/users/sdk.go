package users

import (
	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/platforms"
	domain_users "github.com/steve-rodrigue/aabs/services/saas/domain/users"
)

// New creates a new users application
func New(
	repository domain_users.Repository,
) Application {
	return createApplication(repository)
}

// Application represents the users application
type Application interface {
	Save(user domain_users.User) error

	FindByID(id uuid.UUID) (domain_users.User, error)
	FindByExternalID(platform platforms.Platform, externalID string) (domain_users.User, error)

	Find(index int, amount int) ([]domain_users.User, error)
	FindAfter(cursor uuid.UUID, amount int) ([]domain_users.User, error)

	Count() (int64, error)
}
