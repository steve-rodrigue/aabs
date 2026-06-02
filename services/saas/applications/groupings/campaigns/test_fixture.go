package campaigns

import (
	app_participations "github.com/steve-rodrigue/aabs/services/saas/applications/groupings/participations"
	app_posts "github.com/steve-rodrigue/aabs/services/saas/applications/posts"
	domain_campaigns "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/campaigns"
)

type applicationFixture struct {
	application Application

	repository     *domain_campaigns.MockCampaignRepository
	posts          *app_posts.MockPostsApplication
	participations *app_participations.MockParticipationsApplication
	classifier     *domain_campaigns.MockCampaignClassifier
}

func newApplicationFixture() *applicationFixture {
	repository := domain_campaigns.NewMockCampaignRepository()
	posts := app_posts.NewMockPostsApplication()
	participations := app_participations.NewMockParticipationsApplication()
	classifier := domain_campaigns.NewMockCampaignClassifier()

	application := New(
		repository,
		posts,
		participations,
		classifier,
		25,
	)

	return &applicationFixture{
		application:    application,
		repository:     repository,
		posts:          posts,
		participations: participations,
		classifier:     classifier,
	}
}
