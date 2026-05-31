package scores

import (
	domain_scores "github.com/steve-rodrigue/aabs/services/saas/domain/scores"
	"github.com/steve-rodrigue/aabs/services/saas/domain/scores/scorables"
)

// Calculator calculates all required scores for a target
type Calculator interface {
	Calculate(target scorables.Scorable) ([]domain_scores.Score, error)
}
