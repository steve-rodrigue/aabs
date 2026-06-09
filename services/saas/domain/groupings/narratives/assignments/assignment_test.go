package assignments

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/campaigns"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/narratives"
)

func TestAssignment(t *testing.T) {
	id := uuid.New()

	narrative := narratives.NewMockNarrative(
		"Narrative A",
		"Description A",
	)

	campaign := campaigns.NewMockCampaign(
		"Campaign A",
		"Description A",
	)

	assignedOn := time.Now().UTC()

	assignment := &assignment{
		identifier: id,
		narrative:  narrative,
		campaign:   campaign,
		confidence: 0.8,
		assignedOn: assignedOn,
	}

	if assignment.Identifier() != id {
		t.Fatalf("expected identifier %s, got %s", id, assignment.Identifier())
	}

	if assignment.Narrative() != narrative {
		t.Fatalf("expected narrative")
	}

	if assignment.Campaign() != campaign {
		t.Fatalf("expected campaign")
	}

	if assignment.Confidence() != 0.8 {
		t.Fatalf("expected confidence 0.8, got %f", assignment.Confidence())
	}

	if !assignment.AssignedOn().Equal(assignedOn) {
		t.Fatalf(
			"expected assigned on %s, got %s",
			assignedOn,
			assignment.AssignedOn(),
		)
	}
}
