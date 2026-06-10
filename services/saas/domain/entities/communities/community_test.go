package communities

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/participatables"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/platforms"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/users"
)

func TestCommunity(t *testing.T) {
	id := uuid.New()
	platform := platforms.NewMockPlatform("Reddit", "reddit")
	createdOn := time.Now().UTC()
	moderator := users.NewMockUser("@mod", "Moderator")

	community := &community{
		identifier: id,
		platform:   platform,
		handle:     "aabs",
		title:      "AABS",
		text:       "Anti-AI Bot Spam community",
		createdOn:  createdOn,
		moderators: []users.User{
			moderator,
		},
	}

	if community.Identifier() != id {
		t.Fatalf("expected identifier %s, got %s", id, community.Identifier())
	}

	if community.ParticipationKind() != participatables.CommunityKind {
		t.Fatalf(
			"expected participation kind %s, got %s",
			participatables.CommunityKind,
			community.ParticipationKind(),
		)
	}

	if community.Platform() != platform {
		t.Fatalf("expected platform")
	}

	if community.Handle() != "aabs" {
		t.Fatalf("expected handle %q, got %q", "aabs", community.Handle())
	}

	if community.Title() != "AABS" {
		t.Fatalf("expected title %q, got %q", "AABS", community.Title())
	}

	if community.Text() != "Anti-AI Bot Spam community" {
		t.Fatalf(
			"expected text %q, got %q",
			"Anti-AI Bot Spam community",
			community.Text(),
		)
	}

	if !community.CreatedOn().Equal(createdOn) {
		t.Fatalf(
			"expected created on %s, got %s",
			createdOn,
			community.CreatedOn(),
		)
	}

	if !community.HasModerators() {
		t.Fatalf("expected community to have moderators")
	}

	moderators := community.Moderators()

	if len(moderators) != 1 {
		t.Fatalf("expected 1 moderator, got %d", len(moderators))
	}

	if moderators[0] != moderator {
		t.Fatalf("expected moderator")
	}
}

func TestCommunityHasModeratorsReturnsFalseWhenEmpty(t *testing.T) {
	community := &community{}

	if community.HasModerators() {
		t.Fatalf("expected community not to have moderators")
	}
}

func TestCommunityModeratorsReturnsCopy(t *testing.T) {
	moderator := users.NewMockUser("@mod", "Moderator")

	community := &community{
		moderators: []users.User{
			moderator,
		},
	}

	moderators := community.Moderators()
	moderators[0] = users.NewMockUser("@other", "Other")

	if community.Moderators()[0] != moderator {
		t.Fatalf("expected moderators copy not to mutate original")
	}
}
