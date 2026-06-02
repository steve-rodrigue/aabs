package topics

import (
	app_participations "github.com/steve-rodrigue/aabs/services/saas/applications/groupings/participations"
	app_posts "github.com/steve-rodrigue/aabs/services/saas/applications/posts"
	domain_topics "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/topics"
)

type applicationFixture struct {
	application Application

	repository     *domain_topics.MockTopicRepository
	posts          *app_posts.MockPostsApplication
	participations *app_participations.MockParticipationsApplication
	builder        *domain_topics.MockTopicBuilder
}

func newApplicationFixture() *applicationFixture {
	repository := domain_topics.NewMockTopicRepository()
	posts := app_posts.NewMockPostsApplication()
	participations := app_participations.NewMockParticipationsApplication()
	builder := domain_topics.NewMockTopicBuilder()

	application := New(
		repository,
		posts,
		participations,
		builder,
	)

	return &applicationFixture{
		application:    application,
		repository:     repository,
		posts:          posts,
		participations: participations,
		builder:        builder,
	}
}
