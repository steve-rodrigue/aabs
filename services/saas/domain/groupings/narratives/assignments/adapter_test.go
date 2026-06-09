package assignments

import (
	"errors"
	"math"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/campaigns"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/narratives"
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

	narrative := narratives.NewMockNarrative(
		"Narrative A",
		"Description A",
	)

	campaign := campaigns.NewMockCampaign(
		"Campaign A",
		"Description A",
	)

	assignedOn := time.Now()

	result, err := adapter.ToDomain(
		AssignmentInput{
			Identifier: id,
			Narrative:  narrative,
			Campaign:   campaign,
			Confidence: 0.8,
			AssignedOn: assignedOn,
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	if result.Identifier() != id {
		t.Fatalf("expected identifier %s, got %s", id, result.Identifier())
	}

	if result.Narrative() != narrative {
		t.Fatalf("expected narrative")
	}

	if result.Campaign() != campaign {
		t.Fatalf("expected campaign")
	}

	if result.Confidence() != 0.8 {
		t.Fatalf("expected confidence 0.8, got %f", result.Confidence())
	}

	if !result.AssignedOn().Equal(assignedOn.UTC()) {
		t.Fatalf("expected UTC assigned on")
	}
}

func TestAdapterToDomainAcceptsZeroConfidence(t *testing.T) {
	assertAdapterAcceptsConfidence(t, 0)
}

func TestAdapterToDomainAcceptsMaximumConfidence(t *testing.T) {
	assertAdapterAcceptsConfidence(t, 1)
}

func TestAdapterToDomainReturnsInvalidIdentifierError(t *testing.T) {
	assertAdapterError(
		t,
		func(input *AssignmentInput) {
			input.Identifier = uuid.Nil
		},
		ErrInvalidAssignmentIdentifier,
	)
}

func TestAdapterToDomainReturnsInvalidNarrativeError(t *testing.T) {
	assertAdapterError(
		t,
		func(input *AssignmentInput) {
			input.Narrative = nil
		},
		ErrInvalidAssignmentNarrative,
	)
}

func TestAdapterToDomainReturnsInvalidCampaignError(t *testing.T) {
	assertAdapterError(
		t,
		func(input *AssignmentInput) {
			input.Campaign = nil
		},
		ErrInvalidAssignmentCampaign,
	)
}

func TestAdapterToDomainReturnsInvalidConfidenceErrorWhenNegative(t *testing.T) {
	assertAdapterError(
		t,
		func(input *AssignmentInput) {
			input.Confidence = -0.1
		},
		ErrInvalidAssignmentConfidence,
	)
}

func TestAdapterToDomainReturnsInvalidConfidenceErrorWhenGreaterThanOne(t *testing.T) {
	assertAdapterError(
		t,
		func(input *AssignmentInput) {
			input.Confidence = 1.1
		},
		ErrInvalidAssignmentConfidence,
	)
}

func TestAdapterToDomainReturnsInvalidConfidenceErrorWhenNaN(t *testing.T) {
	assertAdapterError(
		t,
		func(input *AssignmentInput) {
			input.Confidence = math.NaN()
		},
		ErrInvalidAssignmentConfidence,
	)
}

func TestAdapterToDomainReturnsInvalidConfidenceErrorWhenInf(t *testing.T) {
	assertAdapterError(
		t,
		func(input *AssignmentInput) {
			input.Confidence = math.Inf(1)
		},
		ErrInvalidAssignmentConfidence,
	)
}

func TestAdapterToDomainReturnsInvalidAssignedOnError(t *testing.T) {
	assertAdapterError(
		t,
		func(input *AssignmentInput) {
			input.AssignedOn = time.Time{}
		},
		ErrInvalidAssignmentAssignedOn,
	)
}

func assertAdapterAcceptsConfidence(
	t *testing.T,
	confidence float64,
) {
	t.Helper()

	result, err := NewAdapter().ToDomain(
		validAssignmentInput(func(input *AssignmentInput) {
			input.Confidence = confidence
		}),
	)

	if err != nil {
		t.Fatal(err)
	}

	if result.Confidence() != confidence {
		t.Fatalf(
			"expected confidence %f, got %f",
			confidence,
			result.Confidence(),
		)
	}
}

func assertAdapterError(
	t *testing.T,
	mutate func(input *AssignmentInput),
	expected error,
) {
	t.Helper()

	_, err := NewAdapter().ToDomain(
		validAssignmentInput(mutate),
	)

	if !errors.Is(err, expected) {
		t.Fatalf("expected error %v, got %v", expected, err)
	}
}

func validAssignmentInput(
	mutate func(input *AssignmentInput),
) AssignmentInput {
	input := AssignmentInput{
		Identifier: uuid.New(),
		Narrative: narratives.NewMockNarrative(
			"Narrative A",
			"Description A",
		),
		Campaign: campaigns.NewMockCampaign(
			"Campaign A",
			"Description A",
		),
		Confidence: 0.8,
		AssignedOn: time.Now(),
	}

	if mutate != nil {
		mutate(&input)
	}

	return input
}
