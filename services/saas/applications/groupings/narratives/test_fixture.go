package narratives

import (
	app_participations "github.com/steve-rodrigue/aabs/services/saas/applications/groupings/participations"
	domain_narratives "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/narratives"
)

type applicationFixture struct {
	application Application

	repository     *domain_narratives.MockNarrativeRepository
	participations *app_participations.MockParticipationsApplication
}

func newApplicationFixture() *applicationFixture {
	repository := domain_narratives.NewMockNarrativeRepository()
	participations := app_participations.NewMockParticipationsApplication()

	application := New(
		repository,
		participations,
	)

	return &applicationFixture{
		application:    application,
		repository:     repository,
		participations: participations,
	}
}
