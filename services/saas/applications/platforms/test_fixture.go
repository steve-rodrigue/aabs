package platforms

import (
	domain_platforms "github.com/steve-rodrigue/aabs/services/saas/domain/platforms"
)

type applicationFixture struct {
	application Application
	repository  *domain_platforms.MockPlatformRepository
}

func newApplicationFixture() *applicationFixture {
	repository := domain_platforms.NewMockPlatformRepository()

	application := New(repository)

	return &applicationFixture{
		application: application,
		repository:  repository,
	}
}
