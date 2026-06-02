package scores

import (
	"github.com/google/uuid"

	domain_scores "github.com/steve-rodrigue/aabs/services/saas/domain/scores"
	"github.com/steve-rodrigue/aabs/services/saas/domain/scores/scorables"
)

func NewMockScoresApplication() *MockScoresApplication {
	return &MockScoresApplication{}
}

type MockScoresApplication struct {
	CalculateCalls int
	CalculateErr   error
	CalculateValue []domain_scores.Score

	LatestScoreCalls int
	LatestScoreErr   error
	LatestScoreValue domain_scores.Score

	ScoreHistoryCalls int
	ScoreHistoryErr   error
	ScoreHistoryValue []domain_scores.Score

	RecalculateScoresCalls int
	RecalculateScoresErr   error
}

func (application *MockScoresApplication) Calculate(
	target scorables.Scorable,
) ([]domain_scores.Score, error) {
	application.CalculateCalls++

	return application.CalculateValue, application.CalculateErr
}

func (application *MockScoresApplication) LatestScore(
	id uuid.UUID,
	scoreType domain_scores.Type,
) (domain_scores.Score, error) {
	application.LatestScoreCalls++

	return application.LatestScoreValue, application.LatestScoreErr
}

func (application *MockScoresApplication) ScoreHistory(
	id uuid.UUID,
	scoreType domain_scores.Type,
) ([]domain_scores.Score, error) {
	application.ScoreHistoryCalls++

	return application.ScoreHistoryValue, application.ScoreHistoryErr
}

func (application *MockScoresApplication) RecalculateScores() error {
	application.RecalculateScoresCalls++

	return application.RecalculateScoresErr
}
