package applications

import (
	communities_application "github.com/steve-rodrigue/aabs/services/posts/restapis/applications/communities"
	platforms_application "github.com/steve-rodrigue/aabs/services/posts/restapis/applications/platforms"
	posts_application "github.com/steve-rodrigue/aabs/services/posts/restapis/applications/posts"
	users_application "github.com/steve-rodrigue/aabs/services/posts/restapis/applications/users"
)

type applicationFixture struct {
	application Application

	posts       *posts_application.MockPostsApplication
	users       *users_application.MockUsersApplication
	communities *communities_application.MockCommunitiesApplication
	platforms   *platforms_application.MockPlatformsApplication
}

func newApplicationFixture() *applicationFixture {
	posts := posts_application.NewMockPostsApplication()
	users := users_application.NewMockUsersApplication()
	communities := communities_application.NewMockCommunitiesApplication()
	platforms := platforms_application.NewMockPlatformsApplication()

	return &applicationFixture{
		application: New(
			posts,
			users,
			communities,
			platforms,
		),

		posts:       posts,
		users:       users,
		communities: communities,
		platforms:   platforms,
	}
}
