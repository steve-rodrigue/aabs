package participations

import (
	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/applications/groupings/participations/evidences"
	domain_participations "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations"
	domain_evidences "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/evidences"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
)

// New creates a new participation application
func New(
	repository domain_participations.Repository,
	calculator domain_participations.Calculator,
	participatableRepository participatables.Repository,
	evidenceRepository domain_evidences.Repository,
	evidenceCalculator domain_evidences.Calculator,
	evidenceApplication evidences.Application,
) Application {
	return createApplication(
		repository,
		calculator,
		participatableRepository,
		evidenceRepository,
		evidenceCalculator,
		evidenceApplication,
	)
}

// Application represents the participation application
type Application interface {
	Evidences() evidences.Application

	FindByID(id uuid.UUID) (domain_participations.Participation, error)
	FindByParticipant(participant participatables.Participatable) ([]domain_participations.Participation, error)
	FindByTarget(target participatables.Participatable) ([]domain_participations.Participation, error)
	FindBetween(participant participatables.Participatable, target participatables.Participatable) (domain_participations.Participation, error)

	RebuildParticipations() error
}
