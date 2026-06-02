package scores

import (
	"github.com/google/uuid"
	domain_scores "github.com/steve-rodrigue/aabs/services/saas/domain/scores"
	"github.com/steve-rodrigue/aabs/services/saas/domain/scores/scorables"
)

// New creates a new scores application
func New(
	repository domain_scores.Repository,
	scorableRepository scorables.Repository,
	calculators []domain_scores.Calculator,
	recalculateBatchSize int,
) Application {
	return createApplication(
		repository,
		scorableRepository,
		calculators,
		recalculateBatchSize,
	)

}

// Application represents the scores application
type Application interface {
	Calculate(target scorables.Scorable) ([]domain_scores.Score, error)
	LatestScore(id uuid.UUID, scoreType domain_scores.Type) (domain_scores.Score, error)
	ScoreHistory(id uuid.UUID, scoreType domain_scores.Type) ([]domain_scores.Score, error)
	RecalculateScores() error
}
