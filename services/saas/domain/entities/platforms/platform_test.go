package platforms

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/participatables"
)

func TestPlatform(t *testing.T) {
	id := uuid.New()
	createdOn := time.Now().UTC()

	platform := &platform{
		identifier:        id,
		participationKind: participatables.PlatformKind,
		name:              "Reddit",
		handle:            "reddit",
		baseURL:           "https://reddit.com",
		createdOn:         createdOn,
	}

	if platform.Identifier() != id {
		t.Fatalf("expected identifier %s, got %s", id, platform.Identifier())
	}

	if platform.ParticipationKind() != participatables.PlatformKind {
		t.Fatalf(
			"expected participation kind %s, got %s",
			participatables.PlatformKind,
			platform.ParticipationKind(),
		)
	}

	if platform.Name() != "Reddit" {
		t.Fatalf("expected name %q, got %q", "Reddit", platform.Name())
	}

	if platform.Handle() != "reddit" {
		t.Fatalf("expected handle %q, got %q", "reddit", platform.Handle())
	}

	if platform.BaseURL() != "https://reddit.com" {
		t.Fatalf(
			"expected base url %q, got %q",
			"https://reddit.com",
			platform.BaseURL(),
		)
	}

	if !platform.CreatedOn().Equal(createdOn) {
		t.Fatalf(
			"expected created on %s, got %s",
			createdOn,
			platform.CreatedOn(),
		)
	}
}
