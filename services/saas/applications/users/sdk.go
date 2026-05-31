package users

import (
	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/platforms"
	domain_users "github.com/steve-rodrigue/aabs/services/saas/domain/users"
)

// Application represents the users application
type Application interface {
	Save(user domain_users.User) error

	FindByID(id uuid.UUID) (domain_users.User, error)
	FindByExternalID(platform platforms.Platform, externalID string) (domain_users.User, error)

	FindAll() ([]domain_users.User, error)
}
