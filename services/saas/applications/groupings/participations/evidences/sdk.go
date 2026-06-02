package evidences

import (
	"github.com/google/uuid"

	domain_evidences "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/evidences"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
)

// New creates a new participation evidence application
func New(
	repository domain_evidences.Repository,
) Application {
	return createApplication(repository)
}

// Application represents the participation evidence application
type Application interface {
	FindByID(id uuid.UUID) (domain_evidences.Evidence, error)
	FindByParticipation(participation uuid.UUID) ([]domain_evidences.Evidence, error)
	FindByPost(post uuid.UUID) ([]domain_evidences.Evidence, error)
	FindByParticipant(participant participatables.Participatable) ([]domain_evidences.Evidence, error)
	FindByTarget(target participatables.Participatable) ([]domain_evidences.Evidence, error)
}
