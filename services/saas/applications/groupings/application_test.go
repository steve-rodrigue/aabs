package groupings

import (
	"testing"
)

func TestApplications(t *testing.T) {
	fixture := newApplicationFixture()

	if fixture.application.Campaigns() != fixture.campaigns {
		t.Fatalf("expected campaigns application")
	}

	if fixture.application.Topics() != fixture.topics {
		t.Fatalf("expected topics application")
	}

	if fixture.application.Narratives() != fixture.narratives {
		t.Fatalf("expected narratives application")
	}

	if fixture.application.Participations() != fixture.participations {
		t.Fatalf("expected participations application")
	}

	if fixture.application.Clusters() != fixture.clusters {
		t.Fatalf("expected clusters application")
	}
}
