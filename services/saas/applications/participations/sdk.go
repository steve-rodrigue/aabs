package participations

import (
	domain_participations "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
)

// Calculator calculates participation between two domain objects
type Calculator interface {
	Calculate(participant participatables.Participatable, target participatables.Participatable) (domain_participations.Participation, error)
}
