package posts

import (
	domain_posts "github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts"
)

type applicationFixture struct {
	application Application
	repository  *domain_posts.MockPostRepository
	service     *domain_posts.MockPostService
}

func newApplicationFixture() *applicationFixture {
	repository := domain_posts.NewMockPostRepository()
	service := domain_posts.NewMockPostService()

	application := New(
		repository,
		service,
	)

	return &applicationFixture{
		application: application,
		repository:  repository,
		service:     service,
	}
}
