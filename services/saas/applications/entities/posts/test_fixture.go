package posts

import (
	domain_posts "github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts"
)

type applicationFixture struct {
	application Application
	repository  *domain_posts.MockPostRepository
}

func newApplicationFixture() *applicationFixture {
	repository := domain_posts.NewMockPostRepository()

	application := New(repository)

	return &applicationFixture{
		application: application,
		repository:  repository,
	}
}
