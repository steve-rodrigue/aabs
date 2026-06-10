package dirty

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/participatables"
)

func TestNewAdapter(t *testing.T) {
	adapter := NewAdapter()

	if adapter == nil {
		t.Fatalf("expected adapter")
	}
}

func TestAdapterToDomain(t *testing.T) {
	adapter := NewAdapter()

	id := uuid.New()

	participant := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.UserKind,
	)

	target := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.CampaignKind,
	)

	markedOn := time.Now()

	result, err := adapter.ToDomain(
		DirtyInput{
			Identifier:  id,
			Participant: participant,
			Target:      target,
			MarkedOn:    markedOn,
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	if result.Identifier() != id {
		t.Fatalf("expected identifier %s, got %s", id, result.Identifier())
	}

	if result.Participant() != participant {
		t.Fatalf("expected participant")
	}

	if result.Target() != target {
		t.Fatalf("expected target")
	}

	if !result.MarkedOn().Equal(markedOn.UTC()) {
		t.Fatalf("expected UTC marked on")
	}
}

func TestAdapterToDomainReturnsInvalidIdentifierError(t *testing.T) {
	assertAdapterError(
		t,
		func(input *DirtyInput) {
			input.Identifier = uuid.Nil
		},
		ErrInvalidDirtyIdentifier,
	)
}

func TestAdapterToDomainReturnsInvalidParticipantError(t *testing.T) {
	assertAdapterError(
		t,
		func(input *DirtyInput) {
			input.Participant = nil
		},
		ErrInvalidDirtyParticipant,
	)
}

func TestAdapterToDomainReturnsInvalidTargetError(t *testing.T) {
	assertAdapterError(
		t,
		func(input *DirtyInput) {
			input.Target = nil
		},
		ErrInvalidDirtyTarget,
	)
}

func TestAdapterToDomainReturnsInvalidMarkedOnError(t *testing.T) {
	assertAdapterError(
		t,
		func(input *DirtyInput) {
			input.MarkedOn = time.Time{}
		},
		ErrInvalidDirtyMarkedOn,
	)
}

func assertAdapterError(
	t *testing.T,
	mutate func(input *DirtyInput),
	expected error,
) {
	t.Helper()

	_, err := NewAdapter().ToDomain(
		validDirtyInput(mutate),
	)

	if !errors.Is(err, expected) {
		t.Fatalf("expected error %v, got %v", expected, err)
	}
}

func validDirtyInput(
	mutate func(input *DirtyInput),
) DirtyInput {
	input := DirtyInput{
		Identifier: uuid.New(),
		Participant: participatables.NewMockParticipatable(
			uuid.New(),
			participatables.UserKind,
		),
		Target: participatables.NewMockParticipatable(
			uuid.New(),
			participatables.CampaignKind,
		),
		MarkedOn: time.Now(),
	}

	if mutate != nil {
		mutate(&input)
	}

	return input
}
