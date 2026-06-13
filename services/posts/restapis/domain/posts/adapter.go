package posts

import (
	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts/contents"
)

type adapter struct {
	contents contents.Adapter
}

func (adapter *adapter) ToDomain(
	input PostInput,
) (Post, error) {
	if input.Identifier == uuid.Nil {
		return nil, ErrInvalidPostIdentifier
	}

	if input.Creator == nil {
		return nil, ErrInvalidPostCreator
	}

	if input.CreatedOn.IsZero() {
		return nil, ErrInvalidPostCreatedOn
	}

	for _, communityID := range input.CommunityIDs {
		if communityID == uuid.Nil {
			return nil, ErrInvalidPostCommunityID
		}
	}

	content, err := adapter.contents.ToDomain(input.Content)
	if err != nil {
		return nil, err
	}

	if content == nil {
		return nil, ErrInvalidPostContent
	}

	communityIDs := make([]uuid.UUID, len(input.CommunityIDs))
	copy(communityIDs, input.CommunityIDs)

	return &post{
		identifier:   input.Identifier,
		communityIDs: communityIDs,
		creator:      input.Creator,
		content:      content,
		createdOn:    input.CreatedOn.UTC(),
	}, nil
}
