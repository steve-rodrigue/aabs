package evidences

import (
	domain_evidences "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/evidences"
)

type applicationFixture struct {
	application Application
	repository  *domain_evidences.MockEvidenceRepository
}

func newApplicationFixture() *applicationFixture {
	repository := domain_evidences.NewMockEvidenceRepository()

	application := New(repository)

	return &applicationFixture{
		application: application,
		repository:  repository,
	}
}
