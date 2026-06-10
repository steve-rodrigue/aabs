package dirty

import (
	"github.com/google/uuid"
)

type adapter struct{}

func (adapter *adapter) ToDomain(
	input DirtyInput,
) (Dirty, error) {
	if input.Identifier == uuid.Nil {
		return nil, ErrInvalidDirtyIdentifier
	}

	if input.Participant == nil {
		return nil, ErrInvalidDirtyParticipant
	}

	if input.Target == nil {
		return nil, ErrInvalidDirtyTarget
	}

	if input.MarkedOn.IsZero() {
		return nil, ErrInvalidDirtyMarkedOn
	}

	return &dirty{
		identifier:  input.Identifier,
		participant: input.Participant,
		target:      input.Target,
		markedOn:    input.MarkedOn.UTC(),
	}, nil
}
