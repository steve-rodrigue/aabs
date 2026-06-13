package applications

import (
	"github.com/steve-rodrigue/aabs/services/posts/restapis/applications/communities"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/applications/platforms"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/applications/posts"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/applications/users"
)

// New creates a new entities application
func New(
	posts posts.Application,
	users users.Application,
	communities communities.Application,
	platforms platforms.Application,
) Application {
	return createApplication(
		posts,
		users,
		communities,
		platforms,
	)
}

// Application represents the entities application
type Application interface {
	Posts() posts.Application
	Users() users.Application
	Communities() communities.Application
	Platforms() platforms.Application
}
