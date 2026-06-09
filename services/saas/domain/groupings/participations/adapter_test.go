package participations

import (
	"errors"
	"math"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
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
		participatables.TopicKind,
	)

	detectedOn := time.Now()

	result, err := adapter.ToDomain(
		ParticipationInput{
			Identifier:     id,
			Participant:    participant,
			Target:         target,
			PostCount:      4,
			TotalPostCount: 10,
			Percentage:     0.4,
			DetectedOn:     detectedOn,
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

	if result.PostCount() != 4 {
		t.Fatalf("expected post count 4, got %d", result.PostCount())
	}

	if result.TotalPostCount() != 10 {
		t.Fatalf(
			"expected total post count 10, got %d",
			result.TotalPostCount(),
		)
	}

	if result.Percentage() != 0.4 {
		t.Fatalf("expected percentage 0.4, got %f", result.Percentage())
	}

	if !result.DetectedOn().Equal(detectedOn.UTC()) {
		t.Fatalf("expected UTC detected on")
	}
}

func TestAdapterToDomainAcceptsZeroPostCount(t *testing.T) {
	result, err := NewAdapter().ToDomain(
		validParticipationInput(func(input *ParticipationInput) {
			input.PostCount = 0
			input.TotalPostCount = 0
			input.Percentage = 0
		}),
	)

	if err != nil {
		t.Fatal(err)
	}

	if result.PostCount() != 0 {
		t.Fatalf("expected post count 0")
	}
}

func TestAdapterToDomainAcceptsZeroPercentage(t *testing.T) {
	assertAdapterAcceptsPercentage(t, 0)
}

func TestAdapterToDomainAcceptsMaximumPercentage(t *testing.T) {
	assertAdapterAcceptsPercentage(t, 1)
}

func TestAdapterToDomainReturnsInvalidIdentifierError(t *testing.T) {
	assertAdapterError(
		t,
		func(input *ParticipationInput) {
			input.Identifier = uuid.Nil
		},
		ErrInvalidParticipationIdentifier,
	)
}

func TestAdapterToDomainReturnsInvalidParticipantError(t *testing.T) {
	assertAdapterError(
		t,
		func(input *ParticipationInput) {
			input.Participant = nil
		},
		ErrInvalidParticipationParticipant,
	)
}

func TestAdapterToDomainReturnsInvalidTargetError(t *testing.T) {
	assertAdapterError(
		t,
		func(input *ParticipationInput) {
			input.Target = nil
		},
		ErrInvalidParticipationTarget,
	)
}

func TestAdapterToDomainReturnsInvalidPostCountError(t *testing.T) {
	assertAdapterError(
		t,
		func(input *ParticipationInput) {
			input.PostCount = -1
		},
		ErrInvalidParticipationPostCount,
	)
}

func TestAdapterToDomainReturnsInvalidTotalPostCountErrorWhenNegative(t *testing.T) {
	assertAdapterError(
		t,
		func(input *ParticipationInput) {
			input.TotalPostCount = -1
		},
		ErrInvalidParticipationTotalPostCount,
	)
}

func TestAdapterToDomainReturnsInvalidTotalPostCountErrorWhenPostCountIsGreater(t *testing.T) {
	assertAdapterError(
		t,
		func(input *ParticipationInput) {
			input.PostCount = 11
			input.TotalPostCount = 10
		},
		ErrInvalidParticipationTotalPostCount,
	)
}

func TestAdapterToDomainReturnsInvalidPercentageErrorWhenNegative(t *testing.T) {
	assertAdapterError(
		t,
		func(input *ParticipationInput) {
			input.Percentage = -0.1
		},
		ErrInvalidParticipationPercentage,
	)
}

func TestAdapterToDomainReturnsInvalidPercentageErrorWhenGreaterThanOne(t *testing.T) {
	assertAdapterError(
		t,
		func(input *ParticipationInput) {
			input.Percentage = 1.1
		},
		ErrInvalidParticipationPercentage,
	)
}

func TestAdapterToDomainReturnsInvalidPercentageErrorWhenNaN(t *testing.T) {
	assertAdapterError(
		t,
		func(input *ParticipationInput) {
			input.Percentage = math.NaN()
		},
		ErrInvalidParticipationPercentage,
	)
}

func TestAdapterToDomainReturnsInvalidPercentageErrorWhenInf(t *testing.T) {
	assertAdapterError(
		t,
		func(input *ParticipationInput) {
			input.Percentage = math.Inf(1)
		},
		ErrInvalidParticipationPercentage,
	)
}

func TestAdapterToDomainReturnsInvalidDetectedOnError(t *testing.T) {
	assertAdapterError(
		t,
		func(input *ParticipationInput) {
			input.DetectedOn = time.Time{}
		},
		ErrInvalidParticipationDetectedOn,
	)
}

func assertAdapterAcceptsPercentage(
	t *testing.T,
	percentage float64,
) {
	t.Helper()

	result, err := NewAdapter().ToDomain(
		validParticipationInput(func(input *ParticipationInput) {
			input.Percentage = percentage
		}),
	)

	if err != nil {
		t.Fatal(err)
	}

	if result.Percentage() != percentage {
		t.Fatalf(
			"expected percentage %f, got %f",
			percentage,
			result.Percentage(),
		)
	}
}

func assertAdapterError(
	t *testing.T,
	mutate func(input *ParticipationInput),
	expected error,
) {
	t.Helper()

	_, err := NewAdapter().ToDomain(
		validParticipationInput(mutate),
	)

	if !errors.Is(err, expected) {
		t.Fatalf("expected error %v, got %v", expected, err)
	}
}

func validParticipationInput(
	mutate func(input *ParticipationInput),
) ParticipationInput {
	input := ParticipationInput{
		Identifier: uuid.New(),
		Participant: participatables.NewMockParticipatable(
			uuid.New(),
			participatables.UserKind,
		),
		Target: participatables.NewMockParticipatable(
			uuid.New(),
			participatables.TopicKind,
		),
		PostCount:      4,
		TotalPostCount: 10,
		Percentage:     0.4,
		DetectedOn:     time.Now(),
	}

	if mutate != nil {
		mutate(&input)
	}

	return input
}
