package platforms

import (
	"net/url"
	"strings"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
)

type adapter struct{}

func (adapter *adapter) ToDomain(
	input PlatformInput,
) (Platform, error) {
	if input.Identifier == uuid.Nil {
		return nil, ErrInvalidPlatformIdentifier
	}

	if input.ParticipationKind != participatables.PlatformKind {
		return nil, ErrInvalidPlatformParticipationKind
	}

	input.Name = strings.TrimSpace(input.Name)
	if input.Name == "" {
		return nil, ErrInvalidPlatformName
	}

	input.Handle = strings.TrimSpace(input.Handle)
	if input.Handle == "" {
		return nil, ErrInvalidPlatformHandle
	}

	input.BaseURL = strings.TrimSpace(input.BaseURL)
	if input.BaseURL == "" {
		return nil, ErrInvalidPlatformBaseURL
	}

	parsedURL, err := url.ParseRequestURI(input.BaseURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return nil, ErrInvalidPlatformBaseURL
	}

	if input.CreatedOn.IsZero() {
		return nil, ErrInvalidPlatformCreatedOn
	}

	return &platform{
		identifier:        input.Identifier,
		participationKind: input.ParticipationKind,
		name:              input.Name,
		handle:            input.Handle,
		baseURL:           input.BaseURL,
		createdOn:         input.CreatedOn.UTC(),
	}, nil
}
