package users

import (
	"net/url"
	"strings"

	"github.com/google/uuid"
)

type adapter struct{}

func (adapter *adapter) ToDomain(
	input UserInput,
) (User, error) {
	if input.Identifier == uuid.Nil {
		return nil, ErrInvalidUserIdentifier
	}

	if input.Platform == nil {
		return nil, ErrInvalidUserPlatform
	}

	input.ExternalID = strings.TrimSpace(input.ExternalID)
	if input.ExternalID == "" {
		return nil, ErrInvalidUserExternalID
	}

	input.Handle = strings.TrimSpace(input.Handle)
	if input.Handle == "" {
		return nil, ErrInvalidUserHandle
	}

	input.DisplayName = strings.TrimSpace(input.DisplayName)
	if input.DisplayName == "" {
		return nil, ErrInvalidUserDisplayName
	}

	input.ProfileURL = strings.TrimSpace(input.ProfileURL)
	if input.ProfileURL == "" {
		return nil, ErrInvalidUserProfileURL
	}

	parsedURL, err := url.ParseRequestURI(input.ProfileURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return nil, ErrInvalidUserProfileURL
	}

	if input.CreatedOn.IsZero() {
		return nil, ErrInvalidUserCreatedOn
	}

	return &user{
		identifier:  input.Identifier,
		platform:    input.Platform,
		externalID:  input.ExternalID,
		handle:      input.Handle,
		displayName: input.DisplayName,
		profileURL:  input.ProfileURL,
		createdOn:   input.CreatedOn.UTC(),
	}, nil
}
