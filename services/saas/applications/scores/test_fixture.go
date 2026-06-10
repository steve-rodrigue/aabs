package scores

import (
	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/scorables"
	domain_scores "github.com/steve-rodrigue/aabs/services/saas/domain/scores"
)

type applicationFixture struct {
	application Application

	repository  *domain_scores.MockScoreRepository
	scorables   *scorables.MockScorableRepository
	calculators []*domain_scores.MockScoreCalculator
}

func newApplicationFixture() *applicationFixture {
	repository := domain_scores.NewMockScoreRepository()
	scorableRepository := scorables.NewMockScorableRepository()

	trustCalculator := domain_scores.NewMockScoreCalculator(
		domain_scores.TrustType,
	)

	spamCalculator := domain_scores.NewMockScoreCalculator(
		domain_scores.SpamType,
	)

	application := New(
		repository,
		scorableRepository,
		[]domain_scores.Calculator{
			trustCalculator,
			spamCalculator,
		},
		25,
	)

	return &applicationFixture{
		application: application,
		repository:  repository,
		scorables:   scorableRepository,
		calculators: []*domain_scores.MockScoreCalculator{
			trustCalculator,
			spamCalculator,
		},
	}
}
