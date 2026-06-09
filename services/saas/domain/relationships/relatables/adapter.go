package relatables

import (
	"github.com/google/uuid"
)

type adapter struct{}

func (adapter *adapter) ToDomain(
	input RelatableInput,
) (Relatable, error) {
	if input.Identifier == uuid.Nil {
		return nil, ErrInvalidRelatableIdentifier
	}

	if !isValidKind(input.RelationshipKind) {
		return nil, ErrInvalidRelatableRelationshipKind
	}

	return &relatable{
		identifier:       input.Identifier,
		relationshipKind: input.RelationshipKind,
	}, nil
}

func isValidKind(kind Kind) bool {
	switch kind {
	case CampaignKind,
		TopicKind,
		UserKind,
		PostKind,
		NarrativeKind:
		return true
	default:
		return false
	}
}
