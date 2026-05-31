package participations

import (
	"github.com/google/uuid"
	domain_participations "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
)

// Application represents the participation application
type Application interface {
	FindByID(id uuid.UUID) (domain_participations.Participation, error)
	FindByParticipant(participant participatables.Participatable) ([]domain_participations.Participation, error)
	FindByTarget(target participatables.Participatable) ([]domain_participations.Participation, error)
	FindBetween(participant participatables.Participatable, target participatables.Participatable) (domain_participations.Participation, error)
	RebuildParticipations() error
}
