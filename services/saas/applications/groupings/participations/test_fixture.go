package participations

import (
	app_evidences "github.com/steve-rodrigue/aabs/services/saas/applications/groupings/participations/evidences"
	domain_participations "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations"
	domain_evidences "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/evidences"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
)

type applicationFixture struct {
	application Application

	repository               *domain_participations.MockParticipationRepository
	calculator               *domain_participations.MockParticipationCalculator
	participatableRepository *participatables.MockParticipatableRepository
	evidenceRepository       *domain_evidences.MockEvidenceRepository
	evidenceCalculator       *domain_evidences.MockEvidenceCalculator
	evidenceApplication      *app_evidences.MockEvidencesApplication
}

func newApplicationFixture() *applicationFixture {
	repository := domain_participations.NewMockParticipationRepository()
	calculator := domain_participations.NewMockParticipationCalculator()
	participatableRepository := participatables.NewMockParticipatableRepository()
	evidenceRepository := domain_evidences.NewMockEvidenceRepository()
	evidenceCalculator := domain_evidences.NewMockEvidenceCalculator()
	evidenceApplication := app_evidences.NewMockEvidencesApplication()

	application := New(
		repository,
		calculator,
		participatableRepository,
		evidenceRepository,
		evidenceCalculator,
		evidenceApplication,
	)

	return &applicationFixture{
		application:              application,
		repository:               repository,
		calculator:               calculator,
		participatableRepository: participatableRepository,
		evidenceRepository:       evidenceRepository,
		evidenceCalculator:       evidenceCalculator,
		evidenceApplication:      evidenceApplication,
	}
}
