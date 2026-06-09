package threads

import (
	"testing"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/users"
)

func TestThread(t *testing.T) {
	id := uuid.New()
	creator := users.NewMockUser("@user", "User")

	thread := &thread{
		identifier: id,
		creator:    creator,
		title:      "Thread title",
		text:       "Thread text",
	}

	if thread.Identifier() != id {
		t.Fatalf("expected identifier %s, got %s", id, thread.Identifier())
	}

	if thread.Creator() != creator {
		t.Fatalf("expected creator")
	}

	if thread.Title() != "Thread title" {
		t.Fatalf("expected title %q, got %q", "Thread title", thread.Title())
	}

	if thread.Text() != "Thread text" {
		t.Fatalf("expected text %q, got %q", "Thread text", thread.Text())
	}
}
