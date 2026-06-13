package users

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/platforms"
)

func TestUser(t *testing.T) {
	id := uuid.New()
	createdOn := time.Now().UTC()
	platform := platforms.NewMockPlatform("Reddit", "reddit")

	user := &user{
		identifier:  id,
		platform:    platform,
		externalID:  "123",
		handle:      "steve-rodrigue",
		displayName: "Steve Rodrigue",
		profileURL:  "https://reddit.com/u/steve-rodrigue",
		createdOn:   createdOn,
	}

	if user.Identifier() != id {
		t.Fatalf("expected identifier %s, got %s", id, user.Identifier())
	}

	if user.Platform() != platform {
		t.Fatalf("expected platform")
	}

	if user.ExternalID() != "123" {
		t.Fatalf("expected external id %q, got %q", "123", user.ExternalID())
	}

	if user.Handle() != "steve-rodrigue" {
		t.Fatalf("expected handle %q, got %q", "steve", user.Handle())
	}

	if user.DisplayName() != "Steve Rodrigue" {
		t.Fatalf(
			"expected display name %q, got %q",
			"Steve Rodrigue",
			user.DisplayName(),
		)
	}

	if user.ProfileURL() != "https://reddit.com/u/steve-rodrigue" {
		t.Fatalf(
			"expected profile url %q, got %q",
			"https://reddit.com/u/steve",
			user.ProfileURL(),
		)
	}

	if !user.CreatedOn().Equal(createdOn) {
		t.Fatalf(
			"expected created on %s, got %s",
			createdOn,
			user.CreatedOn(),
		)
	}
}
