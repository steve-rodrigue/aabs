package communities

import (
	domain_communities "github.com/steve-rodrigue/aabs/services/saas/domain/communities"
)

type applicationFixture struct {
	application Application
	repository  *domain_communities.MockCommunityRepository
}

func newApplicationFixture() *applicationFixture {
	repository := domain_communities.NewMockCommunityRepository()

	application := New(repository)

	return &applicationFixture{
		application: application,
		repository:  repository,
	}
}
