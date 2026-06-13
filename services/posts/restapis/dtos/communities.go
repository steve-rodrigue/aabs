package dtos

import (
	"time"

	"github.com/google/uuid"

	domain_communities "github.com/steve-rodrigue/aabs/services/posts/restapis/domain/communities"
)

type CommunityResponse struct {
	Identifier uuid.UUID        `json:"identifier"`
	Platform   PlatformResponse `json:"platform"`
	Handle     string           `json:"handle"`
	Title      string           `json:"title"`
	Text       string           `json:"text"`
	CreatedOn  time.Time        `json:"created_on"`
	Moderators []UserResponse   `json:"moderators"`
}

func CommunityDTO(
	community domain_communities.Community,
) CommunityResponse {
	moderators := make(
		[]UserResponse,
		0,
		len(community.Moderators()),
	)

	for _, moderator := range community.Moderators() {
		moderators = append(
			moderators,
			UserDTO(moderator),
		)
	}

	return CommunityResponse{
		Identifier: community.Identifier(),
		Platform:   PlatformDTO(community.Platform()),
		Handle:     community.Handle(),
		Title:      community.Title(),
		Text:       community.Text(),
		CreatedOn:  community.CreatedOn(),
		Moderators: moderators,
	}
}
