package users

import (
	"net/http"
	"strings"

	"github.com/steve-rodrigue/aabs/services/posts/restapis/applications/platforms"
	application_users "github.com/steve-rodrigue/aabs/services/posts/restapis/applications/users"
)

// New creates a new client user application
func New(
	baseURL string,
	client *http.Client,
	platforms platforms.Application,
) application_users.Application {
	if client == nil {
		client = http.DefaultClient
	}

	return &application{
		baseURL:   strings.TrimRight(baseURL, "/"),
		client:    client,
		platforms: platforms,
	}
}
