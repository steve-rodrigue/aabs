package participations

import (
	"math"

	"github.com/google/uuid"
)

type adapter struct{}

func (adapter *adapter) ToDomain(
	input ParticipationInput,
) (Participation, error) {
	if input.Identifier == uuid.Nil {
		return nil, ErrInvalidParticipationIdentifier
	}

	if input.Participant == nil {
		return nil, ErrInvalidParticipationParticipant
	}

	if input.Target == nil {
		return nil, ErrInvalidParticipationTarget
	}

	if input.PostCount < 0 {
		return nil, ErrInvalidParticipationPostCount
	}

	if input.TotalPostCount < 0 ||
		input.PostCount > input.TotalPostCount {
		return nil, ErrInvalidParticipationTotalPostCount
	}

	if math.IsNaN(input.Percentage) ||
		math.IsInf(input.Percentage, 0) ||
		input.Percentage < 0 ||
		input.Percentage > 1 {
		return nil, ErrInvalidParticipationPercentage
	}

	if input.DetectedOn.IsZero() {
		return nil, ErrInvalidParticipationDetectedOn
	}

	return &participation{
		identifier:     input.Identifier,
		participant:    input.Participant,
		target:         input.Target,
		postCount:      input.PostCount,
		totalPostCount: input.TotalPostCount,
		percentage:     input.Percentage,
		detectedOn:     input.DetectedOn.UTC(),
	}, nil
}
