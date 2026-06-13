package dtos

import (
	"time"

	"github.com/google/uuid"

	domain_platforms "github.com/steve-rodrigue/aabs/services/posts/restapis/domain/platforms"
)

type PlatformResponse struct {
	Identifier uuid.UUID `json:"identifier"`
	Name       string    `json:"name"`
	Handle     string    `json:"handle"`
	BaseURL    string    `json:"base_url"`
	CreatedOn  time.Time `json:"created_on"`
}

func PlatformDTO(
	platform domain_platforms.Platform,
) PlatformResponse {
	return PlatformResponse{
		Identifier: platform.Identifier(),
		Name:       platform.Name(),
		Handle:     platform.Handle(),
		BaseURL:    platform.BaseURL(),
		CreatedOn:  platform.CreatedOn(),
	}
}
