package scores

import (
	"github.com/google/uuid"
	domain_scores "github.com/steve-rodrigue/aabs/services/saas/domain/scores"
	"github.com/steve-rodrigue/aabs/services/saas/domain/scores/scorables"
)

type MockScoresApplication struct {
	RecalculateScoresCalls int
	RecalculateScoresErr   error
}

func (application *MockScoresApplication) Calculate(target scorables.Scorable) ([]domain_scores.Score, error) {
	return nil, nil
}

func (application *MockScoresApplication) LatestScore(id uuid.UUID) (domain_scores.Score, error) {
	return nil, nil
}

func (application *MockScoresApplication) ScoreHistory(id uuid.UUID) ([]domain_scores.Score, error) {
	return nil, nil
}

func (application *MockScoresApplication) RecalculateScores() error {
	application.RecalculateScoresCalls++

	return application.RecalculateScoresErr
}
