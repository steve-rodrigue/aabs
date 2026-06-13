package dtos

import (
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/users"
)

type UserResponse struct {
	Identifier  uuid.UUID        `json:"identifier"`
	Platform    PlatformResponse `json:"platform"`
	ExternalID  string           `json:"external_id"`
	Handle      string           `json:"handle"`
	DisplayName string           `json:"display_name"`
	ProfileURL  string           `json:"profile_url"`
	CreatedOn   time.Time        `json:"created_on"`
}

func UserDTO(
	user users.User,
) UserResponse {
	return UserResponse{
		Identifier:  user.Identifier(),
		Platform:    PlatformDTO(user.Platform()),
		ExternalID:  user.ExternalID(),
		Handle:      user.Handle(),
		DisplayName: user.DisplayName(),
		ProfileURL:  user.ProfileURL(),
		CreatedOn:   user.CreatedOn(),
	}
}
