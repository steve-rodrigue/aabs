package contents

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts/contents/replies"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts/contents/threads"
)

func TestContentWithReply(t *testing.T) {
	id := uuid.New()
	createdAt := time.Now().UTC()
	reply := replies.NewMockReply("Reply text")

	content := &content{
		identifier: id,
		reply:      reply,
		createdAt:  createdAt,
	}

	if content.Identifier() != id {
		t.Fatalf("expected identifier %s, got %s", id, content.Identifier())
	}

	if !content.IsReply() {
		t.Fatalf("expected content to be reply")
	}

	if content.Reply() != reply {
		t.Fatalf("expected reply")
	}

	if content.IsThread() {
		t.Fatalf("expected content not to be thread")
	}

	if content.Thread() != nil {
		t.Fatalf("expected nil thread")
	}

	if content.Text() != "Reply text" {
		t.Fatalf("expected reply text %q, got %q", "Reply text", content.Text())
	}

	if !content.CreatedAt().Equal(createdAt) {
		t.Fatalf("expected created at %s, got %s", createdAt, content.CreatedAt())
	}
}

func TestContentWithThread(t *testing.T) {
	id := uuid.New()
	createdAt := time.Now().UTC()
	thread := threads.NewMockThread("Thread title", "Thread text")

	content := &content{
		identifier: id,
		thread:     thread,
		createdAt:  createdAt,
	}

	if content.Identifier() != id {
		t.Fatalf("expected identifier %s, got %s", id, content.Identifier())
	}

	if content.IsReply() {
		t.Fatalf("expected content not to be reply")
	}

	if content.Reply() != nil {
		t.Fatalf("expected nil reply")
	}

	if !content.IsThread() {
		t.Fatalf("expected content to be thread")
	}

	if content.Thread() != thread {
		t.Fatalf("expected thread")
	}

	if content.Text() != "Thread text" {
		t.Fatalf("expected thread text %q, got %q", "Thread text", content.Text())
	}

	if !content.CreatedAt().Equal(createdAt) {
		t.Fatalf("expected created at %s, got %s", createdAt, content.CreatedAt())
	}
}

func TestContentTextReturnsEmptyWhenNoSubContent(t *testing.T) {
	content := &content{}

	if content.Text() != "" {
		t.Fatalf("expected empty text, got %q", content.Text())
	}
}
