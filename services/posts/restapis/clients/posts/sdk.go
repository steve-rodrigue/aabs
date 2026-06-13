package posts

import (
	"net/http"
	"strings"

	"github.com/steve-rodrigue/aabs/services/posts/restapis/applications/posts"
)

// New creates a new post client application
func New(
	baseURL string,
	client *http.Client,
) posts.Application {
	if client == nil {
		client = http.DefaultClient
	}

	return &application{
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  client,
	}
}
