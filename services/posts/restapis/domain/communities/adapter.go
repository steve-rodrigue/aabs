package communities

import (
	"strings"

	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/users"
)

type adapter struct{}

func (adapter *adapter) ToDomain(
	input CommunityInput,
) (Community, error) {
	if input.Identifier == uuid.Nil {
		return nil, ErrInvalidCommunityIdentifier
	}

	if input.Platform == nil {
		return nil, ErrInvalidCommunityPlatform
	}

	input.Handle = strings.TrimSpace(input.Handle)
	if input.Handle == "" {
		return nil, ErrInvalidCommunityHandle
	}

	input.Title = strings.TrimSpace(input.Title)
	if input.Title == "" {
		return nil, ErrInvalidCommunityTitle
	}

	input.Text = strings.TrimSpace(input.Text)
	if input.Text == "" {
		return nil, ErrInvalidCommunityText
	}

	if input.CreatedOn.IsZero() {
		return nil, ErrInvalidCommunityCreatedOn
	}

	for _, moderator := range input.Moderators {
		if moderator == nil {
			return nil, ErrInvalidCommunityModerator
		}
	}

	moderators := make([]users.User, len(input.Moderators))
	copy(moderators, input.Moderators)

	return &community{
		identifier: input.Identifier,
		platform:   input.Platform,
		handle:     input.Handle,
		title:      input.Title,
		text:       input.Text,
		createdOn:  input.CreatedOn.UTC(),
		moderators: moderators,
	}, nil
}
