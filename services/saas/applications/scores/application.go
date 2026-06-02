package scores

import (
	"github.com/google/uuid"

	domain_scores "github.com/steve-rodrigue/aabs/services/saas/domain/scores"
	"github.com/steve-rodrigue/aabs/services/saas/domain/scores/scorables"
)

type application struct {
	repository         domain_scores.Repository
	scorableRepository scorables.Repository
	calculators        []domain_scores.Calculator
}

func createApplication(
	repository domain_scores.Repository,
	scorableRepository scorables.Repository,
	calculators []domain_scores.Calculator,
) Application {
	return &application{
		repository:         repository,
		scorableRepository: scorableRepository,
		calculators:        calculators,
	}
}

// Calculate calculates all scores for a target
func (app *application) Calculate(
	target scorables.Scorable,
) ([]domain_scores.Score, error) {
	out := make([]domain_scores.Score, 0, len(app.calculators))

	for _, calculator := range app.calculators {
		score, err := calculator.Calculate(target)
		if err != nil {
			return nil, err
		}

		if err := app.repository.Save(score); err != nil {
			return nil, err
		}

		out = append(out, score)
	}

	return out, nil
}

// LatestScore finds the latest score for a target
func (app *application) LatestScore(
	id uuid.UUID,
	scoreType domain_scores.Type,
) (domain_scores.Score, error) {
	target, err := app.scorableRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	return app.repository.FindLatestByTarget(
		target,
		scoreType,
	)
}

// ScoreHistory finds score history for a target
func (app *application) ScoreHistory(
	id uuid.UUID,
	scoreType domain_scores.Type,
) ([]domain_scores.Score, error) {
	target, err := app.scorableRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	return app.repository.FindHistoryByTarget(
		target,
		scoreType,
	)
}

// RecalculateScores recalculates scores for all scorable entities
func (app *application) RecalculateScores() error {
	targets, err := app.scorableRepository.FindAll()
	if err != nil {
		return err
	}

	for _, target := range targets {
		if _, err := app.Calculate(target); err != nil {
			return err
		}
	}

	return nil
}
