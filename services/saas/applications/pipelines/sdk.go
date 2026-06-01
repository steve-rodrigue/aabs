package pipelines

import (
	"github.com/steve-rodrigue/aabs/services/saas/applications/groupings"
	"github.com/steve-rodrigue/aabs/services/saas/applications/relationships"
	"github.com/steve-rodrigue/aabs/services/saas/applications/scores"
	"github.com/steve-rodrigue/aabs/services/saas/applications/searches"
	"github.com/steve-rodrigue/aabs/services/saas/domain/posts"
)

// New creates a new pipeline application
func New(
	postRepository posts.Repository,
	searches searches.Application,
	groupings groupings.Application,
	relationships relationships.Application,
	scores scores.Application,
) Application {
	return createApplication(
		postRepository,
		searches,
		groupings,
		relationships,
		scores,
	)
}

// Application represents the post pipeline application
type Application interface {
	ProcessPost(post posts.Post) error
	ProcessPosts(posts []posts.Post) error
	Rebuild() error
}
