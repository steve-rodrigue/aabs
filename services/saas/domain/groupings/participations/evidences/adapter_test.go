package evidences

import (
	"errors"
	"math"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/participatables"
	domain_posts "github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations"
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
	participation := participations.NewMockParticipation()
	participant := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.UserKind,
	)
	target := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.TopicKind,
	)
	post := domain_posts.NewMockPost("hello")
	detectedOn := time.Now()

	result, err := adapter.ToDomain(
		EvidenceInput{
			Identifier:    id,
			Participation: participation,
			Participant:   participant,
			Target:        target,
			Post:          post,
			Score:         0.75,
			DetectedOn:    detectedOn,
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	if result.Identifier() != id {
		t.Fatalf("expected identifier %s, got %s", id, result.Identifier())
	}

	if result.Participation() != participation {
		t.Fatalf("expected participation")
	}

	if result.Participant() != participant {
		t.Fatalf("expected participant")
	}

	if result.Target() != target {
		t.Fatalf("expected target")
	}

	if result.Post() != post {
		t.Fatalf("expected post")
	}

	if result.Score() != 0.75 {
		t.Fatalf("expected score 0.75, got %f", result.Score())
	}

	if !result.DetectedOn().Equal(detectedOn.UTC()) {
		t.Fatalf("expected UTC detected on")
	}
}

func TestAdapterToDomainAcceptsZeroScore(t *testing.T) {
	assertAdapterAcceptsScore(t, 0)
}

func TestAdapterToDomainAcceptsMaximumScore(t *testing.T) {
	assertAdapterAcceptsScore(t, 1)
}

func TestAdapterToDomainReturnsInvalidIdentifierError(t *testing.T) {
	assertAdapterError(
		t,
		func(input *EvidenceInput) {
			input.Identifier = uuid.Nil
		},
		ErrInvalidEvidenceIdentifier,
	)
}

func TestAdapterToDomainReturnsInvalidParticipationError(t *testing.T) {
	assertAdapterError(
		t,
		func(input *EvidenceInput) {
			input.Participation = nil
		},
		ErrInvalidEvidenceParticipation,
	)
}

func TestAdapterToDomainReturnsInvalidParticipantError(t *testing.T) {
	assertAdapterError(
		t,
		func(input *EvidenceInput) {
			input.Participant = nil
		},
		ErrInvalidEvidenceParticipant,
	)
}

func TestAdapterToDomainReturnsInvalidTargetError(t *testing.T) {
	assertAdapterError(
		t,
		func(input *EvidenceInput) {
			input.Target = nil
		},
		ErrInvalidEvidenceTarget,
	)
}

func TestAdapterToDomainReturnsInvalidPostError(t *testing.T) {
	assertAdapterError(
		t,
		func(input *EvidenceInput) {
			input.Post = nil
		},
		ErrInvalidEvidencePost,
	)
}

func TestAdapterToDomainReturnsInvalidScoreErrorWhenNegative(t *testing.T) {
	assertAdapterError(
		t,
		func(input *EvidenceInput) {
			input.Score = -0.1
		},
		ErrInvalidEvidenceScore,
	)
}

func TestAdapterToDomainReturnsInvalidScoreErrorWhenGreaterThanOne(t *testing.T) {
	assertAdapterError(
		t,
		func(input *EvidenceInput) {
			input.Score = 1.1
		},
		ErrInvalidEvidenceScore,
	)
}

func TestAdapterToDomainReturnsInvalidScoreErrorWhenNaN(t *testing.T) {
	assertAdapterError(
		t,
		func(input *EvidenceInput) {
			input.Score = math.NaN()
		},
		ErrInvalidEvidenceScore,
	)
}

func TestAdapterToDomainReturnsInvalidScoreErrorWhenInf(t *testing.T) {
	assertAdapterError(
		t,
		func(input *EvidenceInput) {
			input.Score = math.Inf(1)
		},
		ErrInvalidEvidenceScore,
	)
}

func TestAdapterToDomainReturnsInvalidDetectedOnError(t *testing.T) {
	assertAdapterError(
		t,
		func(input *EvidenceInput) {
			input.DetectedOn = time.Time{}
		},
		ErrInvalidEvidenceDetectedOn,
	)
}

func assertAdapterAcceptsScore(
	t *testing.T,
	score float64,
) {
	t.Helper()

	adapter := NewAdapter()

	result, err := adapter.ToDomain(
		validEvidenceInput(func(input *EvidenceInput) {
			input.Score = score
		}),
	)

	if err != nil {
		t.Fatal(err)
	}

	if result.Score() != score {
		t.Fatalf("expected score %f, got %f", score, result.Score())
	}
}

func assertAdapterError(
	t *testing.T,
	mutate func(input *EvidenceInput),
	expected error,
) {
	t.Helper()

	adapter := NewAdapter()

	_, err := adapter.ToDomain(
		validEvidenceInput(mutate),
	)

	if !errors.Is(err, expected) {
		t.Fatalf("expected error %v, got %v", expected, err)
	}
}

func validEvidenceInput(
	mutate func(input *EvidenceInput),
) EvidenceInput {
	input := EvidenceInput{
		Identifier:    uuid.New(),
		Participation: participations.NewMockParticipation(),
		Participant: participatables.NewMockParticipatable(
			uuid.New(),
			participatables.UserKind,
		),
		Target: participatables.NewMockParticipatable(
			uuid.New(),
			participatables.TopicKind,
		),
		Post:       domain_posts.NewMockPost("hello"),
		Score:      0.75,
		DetectedOn: time.Now(),
	}

	if mutate != nil {
		mutate(&input)
	}

	return input
}
