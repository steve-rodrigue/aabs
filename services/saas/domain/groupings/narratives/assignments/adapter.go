package assignments

import (
	"math"

	"github.com/google/uuid"
)

type adapter struct{}

func (adapter *adapter) ToDomain(
	input AssignmentInput,
) (Assignment, error) {
	if input.Identifier == uuid.Nil {
		return nil, ErrInvalidAssignmentIdentifier
	}

	if input.Narrative == nil {
		return nil, ErrInvalidAssignmentNarrative
	}

	if input.Campaign == nil {
		return nil, ErrInvalidAssignmentCampaign
	}

	if math.IsNaN(input.Confidence) ||
		math.IsInf(input.Confidence, 0) ||
		input.Confidence < 0 ||
		input.Confidence > 1 {
		return nil, ErrInvalidAssignmentConfidence
	}

	if input.AssignedOn.IsZero() {
		return nil, ErrInvalidAssignmentAssignedOn
	}

	return &assignment{
		identifier: input.Identifier,
		narrative:  input.Narrative,
		campaign:   input.Campaign,
		confidence: input.Confidence,
		assignedOn: input.AssignedOn.UTC(),
	}, nil
}
