package searches

import (
	"testing"

	"github.com/google/uuid"
)

func TestResult(t *testing.T) {
	id := uuid.New()

	result := &result{
		identifier: id,
		kind:       PostKind,
		title:      "Post",
		text:       "hello world",
		score:      0.95,
	}

	if result.Identifier() != id {
		t.Fatalf("expected identifier %s, got %s", id, result.Identifier())
	}

	if result.Kind() != PostKind {
		t.Fatalf("expected kind %s, got %s", PostKind, result.Kind())
	}

	if !result.HasTitle() {
		t.Fatalf("expected title to be present")
	}

	if result.Title() != "Post" {
		t.Fatalf("expected title %q, got %q", "Post", result.Title())
	}

	if result.Text() != "hello world" {
		t.Fatalf("expected text %q, got %q", "hello world", result.Text())
	}

	if result.Score() != 0.95 {
		t.Fatalf("expected score %.2f, got %.2f", 0.95, result.Score())
	}
}

func TestResultWithoutTitle(t *testing.T) {
	result := &result{
		identifier: uuid.New(),
		kind:       PostKind,
		title:      "",
		text:       "reply text",
		score:      0.80,
	}

	if result.HasTitle() {
		t.Fatalf("expected title to be absent")
	}

	if result.Title() != "" {
		t.Fatalf("expected empty title")
	}

	if result.Text() != "reply text" {
		t.Fatalf("expected text %q, got %q", "reply text", result.Text())
	}
}
