package evidences

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/participatables"
	domain_posts "github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations"
)

func TestEvidence(t *testing.T) {
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
	detectedOn := time.Now().UTC()

	evidence := &evidence{
		identifier:    id,
		participation: participation,
		participant:   participant,
		target:        target,
		post:          post,
		score:         0.85,
		detectedOn:    detectedOn,
	}

	if evidence.Identifier() != id {
		t.Fatalf("expected identifier %s, got %s", id, evidence.Identifier())
	}

	if evidence.Participation() != participation {
		t.Fatalf("expected participation")
	}

	if evidence.Participant() != participant {
		t.Fatalf("expected participant")
	}

	if evidence.Target() != target {
		t.Fatalf("expected target")
	}

	if evidence.Post() != post {
		t.Fatalf("expected post")
	}

	if evidence.Score() != 0.85 {
		t.Fatalf("expected score 0.85, got %f", evidence.Score())
	}

	if !evidence.DetectedOn().Equal(detectedOn) {
		t.Fatalf(
			"expected detected on %s, got %s",
			detectedOn,
			evidence.DetectedOn(),
		)
	}
}
