package clients

import (
	"net/http"
	"strings"

	"github.com/steve-rodrigue/aabs/services/posts/restapis/applications"

	client_communities "github.com/steve-rodrigue/aabs/services/posts/restapis/clients/communities"
	client_platforms "github.com/steve-rodrigue/aabs/services/posts/restapis/clients/platforms"
	client_posts "github.com/steve-rodrigue/aabs/services/posts/restapis/clients/posts"
	client_users "github.com/steve-rodrigue/aabs/services/posts/restapis/clients/users"
)

func New(
	baseURL string,
	client *http.Client,
) applications.Application {
	if client == nil {
		client = http.DefaultClient
	}

	baseURL = strings.TrimRight(baseURL, "/")

	platforms := client_platforms.New(baseURL, client)
	users := client_users.New(baseURL, client, platforms)
	communities := client_communities.New(baseURL, client, platforms)
	posts := client_posts.New(baseURL, client)

	return applications.New(
		posts,
		users,
		communities,
		platforms,
	)
}
