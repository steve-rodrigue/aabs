package platforms

import (
	"github.com/google/uuid"

	domain_platforms "github.com/steve-rodrigue/aabs/services/saas/domain/platforms"
)

// New creates a new platforms application
func New(
	repository domain_platforms.Repository,
) Application {
	return createApplication(repository)
}

// Application represents the platforms application
type Application interface {
	Save(platform domain_platforms.Platform) error

	FindByID(id uuid.UUID) (domain_platforms.Platform, error)
	FindByHandle(handle string) (domain_platforms.Platform, error)

	FindAll() ([]domain_platforms.Platform, error)
}
