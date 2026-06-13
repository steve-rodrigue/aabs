package applications

import (
	communities_application "github.com/steve-rodrigue/aabs/services/posts/restapis/applications/communities"
	platforms_application "github.com/steve-rodrigue/aabs/services/posts/restapis/applications/platforms"
	posts_application "github.com/steve-rodrigue/aabs/services/posts/restapis/applications/posts"
	users_application "github.com/steve-rodrigue/aabs/services/posts/restapis/applications/users"
)

type application struct {
	posts       posts_application.Application
	users       users_application.Application
	communities communities_application.Application
	platforms   platforms_application.Application
}

func createApplication(
	posts posts_application.Application,
	users users_application.Application,
	communities communities_application.Application,
	platforms platforms_application.Application,
) Application {
	return &application{
		posts:       posts,
		users:       users,
		communities: communities,
		platforms:   platforms,
	}
}

func (application *application) Posts() posts_application.Application {
	return application.posts
}

func (application *application) Users() users_application.Application {
	return application.users
}

func (application *application) Communities() communities_application.Application {
	return application.communities
}

func (application *application) Platforms() platforms_application.Application {
	return application.platforms
}
