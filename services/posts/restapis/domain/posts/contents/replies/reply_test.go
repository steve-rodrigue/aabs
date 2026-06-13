package replies

import (
	"testing"

	"github.com/google/uuid"
)

func TestReply(t *testing.T) {
	id := uuid.New()
	target := NewMockTargetWithThread(newMockThread())

	reply := &reply{
		identifier: id,
		target:     target,
		text:       "Reply text",
	}

	if reply.Identifier() != id {
		t.Fatalf("expected identifier %s, got %s", id, reply.Identifier())
	}

	if reply.Target() != target {
		t.Fatalf("expected target")
	}

	if reply.Text() != "Reply text" {
		t.Fatalf("expected text %q, got %q", "Reply text", reply.Text())
	}
}
