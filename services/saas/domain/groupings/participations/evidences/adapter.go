package evidences

import (
	"math"

	"github.com/google/uuid"
)

type adapter struct{}

func (adapter *adapter) ToDomain(
	input EvidenceInput,
) (Evidence, error) {
	if input.Identifier == uuid.Nil {
		return nil, ErrInvalidEvidenceIdentifier
	}

	if input.Participation == nil {
		return nil, ErrInvalidEvidenceParticipation
	}

	if input.Participant == nil {
		return nil, ErrInvalidEvidenceParticipant
	}

	if input.Target == nil {
		return nil, ErrInvalidEvidenceTarget
	}

	if input.Post == nil {
		return nil, ErrInvalidEvidencePost
	}

	if math.IsNaN(input.Score) ||
		math.IsInf(input.Score, 0) ||
		input.Score < 0 ||
		input.Score > 1 {
		return nil, ErrInvalidEvidenceScore
	}

	if input.DetectedOn.IsZero() {
		return nil, ErrInvalidEvidenceDetectedOn
	}

	return &evidence{
		identifier:    input.Identifier,
		participation: input.Participation,
		participant:   input.Participant,
		target:        input.Target,
		post:          input.Post,
		score:         input.Score,
		detectedOn:    input.DetectedOn.UTC(),
	}, nil
}
