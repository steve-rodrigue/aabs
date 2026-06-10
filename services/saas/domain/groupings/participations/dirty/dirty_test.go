package dirty

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/participatables"
)

func TestDirty(t *testing.T) {
	id := uuid.New()

	participant := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.UserKind,
	)

	target := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.CampaignKind,
	)

	markedOn := time.Now().UTC()

	dirty := &dirty{
		identifier:  id,
		participant: participant,
		target:      target,
		markedOn:    markedOn,
	}

	if dirty.Identifier() != id {
		t.Fatalf("expected identifier %s, got %s", id, dirty.Identifier())
	}

	if dirty.Participant() != participant {
		t.Fatalf("expected participant")
	}

	if dirty.Target() != target {
		t.Fatalf("expected target")
	}

	if !dirty.MarkedOn().Equal(markedOn) {
		t.Fatalf("expected marked on %s, got %s", markedOn, dirty.MarkedOn())
	}
}
