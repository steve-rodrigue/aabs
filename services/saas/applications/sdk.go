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

// New creates a new application
func New(
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
	return createApplication(
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
}

// Application represents the root application
type Application interface {
	Pipeline() pipelines.Application

	Posts() posts.Application
	Users() users.Application
	Communities() communities.Application
	Platforms() platforms.Application

	Groupings() groupings.Application
	Relationships() relationships.Application
	Scores() scores.Application
	Searches() searches.Application
}
