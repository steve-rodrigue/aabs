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
