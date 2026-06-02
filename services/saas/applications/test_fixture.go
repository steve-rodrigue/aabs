package application

import (
	"github.com/steve-rodrigue/aabs/services/saas/applications/communities"
	"github.com/steve-rodrigue/aabs/services/saas/applications/groupings"
	"github.com/steve-rodrigue/aabs/services/saas/applications/pipelines"
	"github.com/steve-rodrigue/aabs/services/saas/applications/platforms"
	"github.com/steve-rodrigue/aabs/services/saas/applications/posts"
	"github.com/steve-rodrigue/aabs/services/saas/applications/relationships"
	"github.com/steve-rodrigue/aabs/services/saas/applications/scores"
	"github.com/steve-rodrigue/aabs/services/saas/applications/searches"
	"github.com/steve-rodrigue/aabs/services/saas/applications/users"
)

type applicationFixture struct {
	application   Application
	pipeline      *pipelines.MockPipelineApplication
	posts         *posts.MockPostsApplication
	users         *users.MockUsersApplication
	communities   *communities.MockCommunitiesApplication
	platforms     *platforms.MockPlatformsApplication
	groupings     *groupings.MockGroupingsApplication
	relationships *relationships.MockRelationshipsApplication
	scores        *scores.MockScoresApplication
	searches      *searches.MockSearchApplication
}

func newApplicationFixture() *applicationFixture {

	pipeline := pipelines.NewMockPipelineApplication()
	posts := posts.NewMockPostsApplication()
	users := users.NewMockUsersApplication()
	communities := communities.NewMockCommunitiesApplication()
	platforms := platforms.NewMockPlatformsApplication()
	groupings := groupings.NewMockGroupingsApplication()
	relationships := relationships.NewMockRelationshipsApplication()
	scores := scores.NewMockScoresApplication()
	searches := searches.NewMockSearchApplication()
	application := New(
		pipeline,
		posts,
		users,
		communities,
		platforms,
		groupings,
		relationships,
		scores,
		searches,
	)
	return &applicationFixture{
		application:   application,
		pipeline:      pipeline,
		posts:         posts,
		users:         users,
		communities:   communities,
		platforms:     platforms,
		groupings:     groupings,
		relationships: relationships,
		scores:        scores,
		searches:      searches,
	}

}
