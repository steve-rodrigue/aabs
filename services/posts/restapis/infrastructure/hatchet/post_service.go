package hatchet

import (
	"context"
	"time"

	"github.com/google/uuid"
	hatchet "github.com/hatchet-dev/hatchet/sdks/go"

	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts"
)

type postService struct {
	client *hatchet.Client
}

type PostSavedPayload struct {
	PostID       uuid.UUID   `json:"post_id"`
	CommunityIDs []uuid.UUID `json:"community_ids"`
	CreatorID    uuid.UUID   `json:"creator_id"`
	Text         string      `json:"text"`
	CreatedOn    time.Time   `json:"created_on"`
}

func (service *postService) Save(
	ctx context.Context,
	post posts.Post,
) error {
	return service.client.Events().Push(
		ctx,
		PostSavedEventName,
		PostSavedPayload{
			PostID:       post.Identifier(),
			CommunityIDs: post.CommunityIDs(),
			CreatorID:    post.Creator().Identifier(),
			Text:         post.Content().Text(),
			CreatedOn:    post.CreatedOn().UTC(),
		},
	)
}
