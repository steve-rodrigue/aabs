package communities

import (
	"github.com/google/uuid"

	domain_communities "github.com/steve-rodrigue/aabs/services/saas/domain/communities"
	"github.com/steve-rodrigue/aabs/services/saas/domain/platforms"
)

// Application represents the communities application
type Application interface {
	Save(community domain_communities.Community) error

	FindByID(id uuid.UUID) (domain_communities.Community, error)

	FindByHandle(platform platforms.Platform, handle string) (domain_communities.Community, error)
	FindAll() ([]domain_communities.Community, error)
	FindByPlatform(platform platforms.Platform) ([]domain_communities.Community, error)
}
