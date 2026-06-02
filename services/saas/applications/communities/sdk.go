package communities

import (
	"github.com/google/uuid"

	domain_communities "github.com/steve-rodrigue/aabs/services/saas/domain/communities"
	"github.com/steve-rodrigue/aabs/services/saas/domain/platforms"
)

// New creates a new communities application
func New(
	repository domain_communities.Repository,
) Application {
	return createApplication(repository)
}

// Application represents the communities application
type Application interface {
	Save(community domain_communities.Community) error

	FindByID(id uuid.UUID) (domain_communities.Community, error)
	FindByHandle(platform platforms.Platform, handle string) (domain_communities.Community, error)

	Find(index int, amount int) ([]domain_communities.Community, error)
	FindAfter(cursor uuid.UUID, amount int) ([]domain_communities.Community, error)

	FindByPlatform(platform platforms.Platform) ([]domain_communities.Community, error)

	Count() (int64, error)
}
