package application

import (
	"github.com/steve-rodrigue/aabs/services/saas/applications/groupings"
	"github.com/steve-rodrigue/aabs/services/saas/applications/pipelines"
	"github.com/steve-rodrigue/aabs/services/saas/applications/posts"
	"github.com/steve-rodrigue/aabs/services/saas/applications/relationships"
	"github.com/steve-rodrigue/aabs/services/saas/applications/scores"
	"github.com/steve-rodrigue/aabs/services/saas/applications/searches"
)

// Application represents the root application
type Application interface {
	Groupings() groupings.Application
	Pipeline() pipelines.Application
	Posts() posts.Application
	Relationships() relationships.Application
	Scores() scores.Application
	Searches() searches.Application
}
