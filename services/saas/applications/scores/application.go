package scores

import (
	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/scorables"
	domain_scores "github.com/steve-rodrigue/aabs/services/saas/domain/scores"
)

type application struct {
	repository           domain_scores.Repository
	scorableRepository   scorables.Repository
	calculators          []domain_scores.Calculator
	recalculateBatchSize int
}

func createApplication(
	repository domain_scores.Repository,
	scorableRepository scorables.Repository,
	calculators []domain_scores.Calculator,
	recalculateBatchSize int,
) Application {
	return &application{
		repository:           repository,
		scorableRepository:   scorableRepository,
		calculators:          calculators,
		recalculateBatchSize: recalculateBatchSize,
	}
}

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

func (app *application) RecalculateScores() error {
	cursor := uuid.Nil

	for {
		targets, err := app.scorableRepository.FindAfter(
			cursor,
			app.recalculateBatchSize,
		)
		if err != nil {
			return err
		}

		if len(targets) == 0 {
			return nil
		}

		for _, target := range targets {
			if _, err := app.Calculate(target); err != nil {
				return err
			}
		}

		cursor = targets[len(targets)-1].Identifier()
	}
}
