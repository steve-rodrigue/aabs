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

type application struct {
	pipeline pipelines.Application

	posts       posts.Application
	users       users.Application
	communities communities.Application
	platforms   platforms.Application

	groupings     groupings.Application
	relationships relationships.Application
	scores        scores.Application
	searches      searches.Application
}

func createApplication(
	pipeline pipelines.Application,
	posts posts.Application,
	users users.Application,
	communities communities.Application,
	platforms platforms.Application,
	groupings groupings.Application,
	relationships relationships.Application,
	scores scores.Application,
	searches searches.Application,
) Application {
	return &application{
		pipeline: pipeline,

		posts:       posts,
		users:       users,
		communities: communities,
		platforms:   platforms,

		groupings:     groupings,
		relationships: relationships,
		scores:        scores,
		searches:      searches,
	}
}

func (app *application) Pipeline() pipelines.Application {
	return app.pipeline
}

func (app *application) Posts() posts.Application {
	return app.posts
}

func (app *application) Users() users.Application {
	return app.users
}

func (app *application) Communities() communities.Application {
	return app.communities
}

func (app *application) Platforms() platforms.Application {
	return app.platforms
}

func (app *application) Groupings() groupings.Application {
	return app.groupings
}

func (app *application) Relationships() relationships.Application {
	return app.relationships
}

func (app *application) Scores() scores.Application {
	return app.scores
}

func (app *application) Searches() searches.Application {
	return app.searches
}
