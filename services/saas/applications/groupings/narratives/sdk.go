package narratives

import (
	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/communities"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/narratives"
	"github.com/steve-rodrigue/aabs/services/saas/domain/users"
)

// Application represents the narratives application
type Application interface {
	FindByID(id uuid.UUID) (narratives.Narrative, error)
	FindAll() ([]narratives.Narrative, error)
	FindNarrativesByUser(user users.User) ([]narratives.Narrative, error)
	FindNarrativesByCommunity(community communities.Community) ([]narratives.Narrative, error)
	RebuildNarratives() error
}
