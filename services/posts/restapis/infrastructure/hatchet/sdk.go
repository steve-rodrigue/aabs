package hatchet

import (
	hatchet "github.com/hatchet-dev/hatchet/sdks/go"

	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts"
)

const PostSavedEventName = "post:saved"

// NewPostService creates a Hatchet post service.
func NewPostService(
	client *hatchet.Client,
) posts.Service {
	return &postService{
		client: client,
	}
}
