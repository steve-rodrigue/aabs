package pipelines

import (
	"github.com/steve-rodrigue/aabs/services/saas/applications/groupings"
	"github.com/steve-rodrigue/aabs/services/saas/applications/groupings/campaigns"
	"github.com/steve-rodrigue/aabs/services/saas/applications/groupings/clusters"
	"github.com/steve-rodrigue/aabs/services/saas/applications/groupings/narratives"
	"github.com/steve-rodrigue/aabs/services/saas/applications/groupings/participations"
	"github.com/steve-rodrigue/aabs/services/saas/applications/groupings/topics"
	app_posts "github.com/steve-rodrigue/aabs/services/saas/applications/posts"
	"github.com/steve-rodrigue/aabs/services/saas/applications/relationships"
	"github.com/steve-rodrigue/aabs/services/saas/applications/scores"
	app_searches "github.com/steve-rodrigue/aabs/services/saas/applications/searches"
)

type applicationFixture struct {
	application Application

	posts    *app_posts.MockPostsApplication
	searches *app_searches.MockSearchApplication

	groupings      *groupings.MockGroupingsApplication
	clusters       *clusters.MockClustersApplication
	campaigns      *campaigns.MockCampaignsApplication
	topics         *topics.MockTopicsApplication
	narratives     *narratives.MockNarrativesApplication
	participations *participations.MockParticipationsApplication

	relationships *relationships.MockRelationshipsApplication
	scores        *scores.MockScoresApplication
}

func newApplicationFixture() *applicationFixture {
	posts := app_posts.NewMockPostsApplication()
	searches := app_searches.NewMockSearchApplication()

	groupings := groupings.NewMockGroupingsApplication()

	clusters := groupings.Clusters().(*clusters.MockClustersApplication)
	campaigns := groupings.Campaigns().(*campaigns.MockCampaignsApplication)
	topics := groupings.Topics().(*topics.MockTopicsApplication)
	narratives := groupings.Narratives().(*narratives.MockNarrativesApplication)
	participations := groupings.Participations().(*participations.MockParticipationsApplication)

	relationships := relationships.NewMockRelationshipsApplication()
	scores := scores.NewMockScoresApplication()

	application := New(
		posts,
		searches,
		groupings,
		relationships,
		scores,
	)

	return &applicationFixture{
		application: application,

		posts:    posts,
		searches: searches,

		groupings:      groupings,
		clusters:       clusters,
		campaigns:      campaigns,
		topics:         topics,
		narratives:     narratives,
		participations: participations,

		relationships: relationships,
		scores:        scores,
	}
}
