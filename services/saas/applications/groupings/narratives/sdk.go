package narratives

import (
	"github.com/google/uuid"

	app_participations "github.com/steve-rodrigue/aabs/services/saas/applications/groupings/participations"
	"github.com/steve-rodrigue/aabs/services/saas/domain/communities"
	domain_narratives "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/narratives"
	"github.com/steve-rodrigue/aabs/services/saas/domain/users"
)

func New(
	repository domain_narratives.Repository,
	participations app_participations.Application,
) Application {
	return createApplication(repository, participations)
}

type Application interface {
	FindByID(id uuid.UUID) (domain_narratives.Narrative, error)

	Find(index int, amount int) ([]domain_narratives.Narrative, error)
	FindAfter(cursor uuid.UUID, amount int) ([]domain_narratives.Narrative, error)

	FindNarrativesByUser(user users.User) ([]domain_narratives.Narrative, error)
	FindNarrativesByCommunity(community communities.Community) ([]domain_narratives.Narrative, error)

	Count() (int64, error)

	RebuildNarratives() error
}
