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

func NewMockApplication() *MockApplication {
	return &MockApplication{
		PipelineValue:      pipelines.NewMockPipelineApplication(),
		PostsValue:         posts.NewMockPostsApplication(),
		UsersValue:         users.NewMockUsersApplication(),
		CommunitiesValue:   communities.NewMockCommunitiesApplication(),
		PlatformsValue:     platforms.NewMockPlatformsApplication(),
		GroupingsValue:     groupings.NewMockGroupingsApplication(),
		RelationshipsValue: relationships.NewMockRelationshipsApplication(),
		ScoresValue:        scores.NewMockScoresApplication(),
		SearchesValue:      searches.NewMockSearchApplication(),
	}
}

type MockApplication struct {
	PipelineCalls int
	PipelineValue pipelines.Application

	PostsCalls int
	PostsValue posts.Application

	UsersCalls int
	UsersValue users.Application

	CommunitiesCalls int
	CommunitiesValue communities.Application

	PlatformsCalls int
	PlatformsValue platforms.Application

	GroupingsCalls int
	GroupingsValue groupings.Application

	RelationshipsCalls int
	RelationshipsValue relationships.Application

	ScoresCalls int
	ScoresValue scores.Application

	SearchesCalls int
	SearchesValue searches.Application
}

func (application *MockApplication) Pipeline() pipelines.Application {
	application.PipelineCalls++

	return application.PipelineValue
}

func (application *MockApplication) Posts() posts.Application {
	application.PostsCalls++

	return application.PostsValue
}

func (application *MockApplication) Users() users.Application {
	application.UsersCalls++

	return application.UsersValue
}

func (application *MockApplication) Communities() communities.Application {
	application.CommunitiesCalls++

	return application.CommunitiesValue
}

func (application *MockApplication) Platforms() platforms.Application {
	application.PlatformsCalls++

	return application.PlatformsValue
}

func (application *MockApplication) Groupings() groupings.Application {
	application.GroupingsCalls++

	return application.GroupingsValue
}

func (application *MockApplication) Relationships() relationships.Application {
	application.RelationshipsCalls++

	return application.RelationshipsValue
}

func (application *MockApplication) Scores() scores.Application {
	application.ScoresCalls++

	return application.ScoresValue
}

func (application *MockApplication) Searches() searches.Application {
	application.SearchesCalls++

	return application.SearchesValue
}
