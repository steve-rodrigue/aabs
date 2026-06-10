package participations

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/participatables"
)

func TestParticipation(t *testing.T) {
	id := uuid.New()

	participant := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.UserKind,
	)

	target := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.TopicKind,
	)

	detectedOn := time.Now().UTC()

	participation := &participation{
		identifier:     id,
		participant:    participant,
		target:         target,
		postCount:      4,
		totalPostCount: 10,
		percentage:     0.4,
		detectedOn:     detectedOn,
	}

	if participation.Identifier() != id {
		t.Fatalf("expected identifier %s, got %s", id, participation.Identifier())
	}

	if participation.Participant() != participant {
		t.Fatalf("expected participant")
	}

	if participation.Target() != target {
		t.Fatalf("expected target")
	}

	if participation.PostCount() != 4 {
		t.Fatalf("expected post count 4, got %d", participation.PostCount())
	}

	if participation.TotalPostCount() != 10 {
		t.Fatalf(
			"expected total post count 10, got %d",
			participation.TotalPostCount(),
		)
	}

	if participation.Percentage() != 0.4 {
		t.Fatalf("expected percentage 0.4, got %f", participation.Percentage())
	}

	if !participation.DetectedOn().Equal(detectedOn) {
		t.Fatalf(
			"expected detected on %s, got %s",
			detectedOn,
			participation.DetectedOn(),
		)
	}
}
