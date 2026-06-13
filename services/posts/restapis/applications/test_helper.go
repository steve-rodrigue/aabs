package applications

import (
	communities_application "github.com/steve-rodrigue/aabs/services/posts/restapis/applications/communities"
	platforms_application "github.com/steve-rodrigue/aabs/services/posts/restapis/applications/platforms"
	posts_application "github.com/steve-rodrigue/aabs/services/posts/restapis/applications/posts"
	users_application "github.com/steve-rodrigue/aabs/services/posts/restapis/applications/users"
)

func NewMockEntitiesApplication() *MockEntitiesApplication {
	return &MockEntitiesApplication{
		PostsValue:       posts_application.NewMockPostsApplication(),
		UsersValue:       users_application.NewMockUsersApplication(),
		CommunitiesValue: communities_application.NewMockCommunitiesApplication(),
		PlatformsValue:   platforms_application.NewMockPlatformsApplication(),
	}
}

type MockEntitiesApplication struct {
	PostsCalls int
	PostsValue posts_application.Application

	UsersCalls int
	UsersValue users_application.Application

	CommunitiesCalls int
	CommunitiesValue communities_application.Application

	PlatformsCalls int
	PlatformsValue platforms_application.Application
}

func (application *MockEntitiesApplication) Posts() posts_application.Application {
	application.PostsCalls++

	return application.PostsValue
}

func (application *MockEntitiesApplication) Users() users_application.Application {
	application.UsersCalls++

	return application.UsersValue
}

func (application *MockEntitiesApplication) Communities() communities_application.Application {
	application.CommunitiesCalls++

	return application.CommunitiesValue
}

func (application *MockEntitiesApplication) Platforms() platforms_application.Application {
	application.PlatformsCalls++

	return application.PlatformsValue
}
