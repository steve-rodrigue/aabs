package users

import (
	domain_users "github.com/steve-rodrigue/aabs/services/saas/domain/entities/users"
)

type applicationFixture struct {
	application Application
	repository  *domain_users.MockUserRepository
}

func newApplicationFixture() *applicationFixture {
	repository := domain_users.NewMockUserRepository()

	application := New(repository)

	return &applicationFixture{
		application: application,
		repository:  repository,
	}
}
