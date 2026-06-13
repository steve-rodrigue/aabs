package platforms

import (
	"net/http"
	"strings"

	application_platforms "github.com/steve-rodrigue/aabs/services/posts/restapis/applications/platforms"
)

// New creates a new platform client application
func New(
	baseURL string,
	client *http.Client,
) application_platforms.Application {
	if client == nil {
		client = http.DefaultClient
	}

	return &application{
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  client,
	}
}
