package communities

import (
	"net/http"
	"strings"

	"github.com/steve-rodrigue/aabs/services/posts/restapis/applications/communities"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/applications/platforms"
)

// New creates a new community client application
func New(
	baseURL string,
	client *http.Client,
	platforms platforms.Application,
) communities.Application {
	if client == nil {
		client = http.DefaultClient
	}

	return &application{
		baseURL:   strings.TrimRight(baseURL, "/"),
		client:    client,
		platforms: platforms,
	}
}
