package scores

import (
	domain_scores "github.com/steve-rodrigue/aabs/services/saas/domain/scores"
	"github.com/steve-rodrigue/aabs/services/saas/domain/scores/scorables"
)

type applicationFixture struct {
	application Application

	repository         *domain_scores.MockScoreRepository
	scorableRepository *scorables.MockScorableRepository

	trustCalculator *domain_scores.MockCalculator
	spamCalculator  *domain_scores.MockCalculator
}

func newApplicationFixture() *applicationFixture {
	repository := domain_scores.NewMockScoreRepository()
	scorableRepository := scorables.NewMockScorableRepository()

	trustCalculator := domain_scores.NewMockCalculator(domain_scores.TrustType)
	spamCalculator := domain_scores.NewMockCalculator(domain_scores.SpamType)

	application := New(
		repository,
		scorableRepository,
		[]domain_scores.Calculator{
			trustCalculator,
			spamCalculator,
		},
	)

	return &applicationFixture{
		application:        application,
		repository:         repository,
		scorableRepository: scorableRepository,
		trustCalculator:    trustCalculator,
		spamCalculator:     spamCalculator,
	}
}
