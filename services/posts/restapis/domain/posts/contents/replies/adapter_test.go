package replies

import (
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts/contents/threads"
)

func TestNewAdapter(t *testing.T) {
	adapter := NewAdapter()

	if adapter == nil {
		t.Fatalf("expected adapter")
	}
}

func TestAdapterToDomainWithThreadTarget(t *testing.T) {
	adapter := NewAdapter()

	id := uuid.New()
	thread := newMockThread()

	result, err := adapter.ToDomain(
		ReplyInput{
			Identifier: id,
			Target: TargetInput{
				Thread: thread,
			},
			Text: " Reply text ",
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	if result.Identifier() != id {
		t.Fatalf("expected identifier %s, got %s", id, result.Identifier())
	}

	if !result.Target().IsThread() {
		t.Fatalf("expected thread target")
	}

	if result.Target().Thread() != thread {
		t.Fatalf("expected thread")
	}

	if result.Target().IsReply() {
		t.Fatalf("expected target not to be reply")
	}

	if result.Text() != "Reply text" {
		t.Fatalf("expected trimmed text %q, got %q", "Reply text", result.Text())
	}
}

func TestAdapterToDomainWithReplyTarget(t *testing.T) {
	adapter := NewAdapter()

	parent := NewMockReplyWithThreadTarget(
		"Parent reply",
		newMockThread(),
	)

	result, err := adapter.ToDomain(
		ReplyInput{
			Identifier: uuid.New(),
			Target: TargetInput{
				Reply: parent,
			},
			Text: " Child reply ",
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	if !result.Target().IsReply() {
		t.Fatalf("expected reply target")
	}

	if result.Target().Reply() != parent {
		t.Fatalf("expected parent reply")
	}

	if result.Target().IsThread() {
		t.Fatalf("expected target not to be thread")
	}

	if result.Text() != "Child reply" {
		t.Fatalf("expected trimmed text")
	}
}

func TestAdapterToDomainReturnsInvalidIdentifierError(t *testing.T) {
	adapter := NewAdapter()

	_, err := adapter.ToDomain(
		validReplyInput(func(input *ReplyInput) {
			input.Identifier = uuid.Nil
		}),
	)

	if !errors.Is(err, ErrInvalidReplyIdentifier) {
		t.Fatalf("expected invalid identifier error, got %v", err)
	}
}

func TestAdapterToDomainReturnsInvalidTargetErrorWhenEmpty(t *testing.T) {
	adapter := NewAdapter()

	_, err := adapter.ToDomain(
		validReplyInput(func(input *ReplyInput) {
			input.Target = TargetInput{}
		}),
	)

	if !errors.Is(err, ErrInvalidReplyTarget) {
		t.Fatalf("expected invalid target error, got %v", err)
	}
}

func TestAdapterToDomainReturnsInvalidTargetErrorWhenReplyAndThreadAreSet(t *testing.T) {
	adapter := NewAdapter()

	_, err := adapter.ToDomain(
		validReplyInput(func(input *ReplyInput) {
			input.Target = TargetInput{
				Reply:  NewMockReply("Parent reply"),
				Thread: newMockThread(),
			}
		}),
	)

	if !errors.Is(err, ErrInvalidReplyTarget) {
		t.Fatalf("expected invalid target error, got %v", err)
	}
}

func TestAdapterToDomainReturnsInvalidTextError(t *testing.T) {
	adapter := NewAdapter()

	_, err := adapter.ToDomain(
		validReplyInput(func(input *ReplyInput) {
			input.Text = "   "
		}),
	)

	if !errors.Is(err, ErrInvalidReplyText) {
		t.Fatalf("expected invalid text error, got %v", err)
	}
}

func validReplyInput(
	mutate func(input *ReplyInput),
) ReplyInput {
	input := ReplyInput{
		Identifier: uuid.New(),
		Target: TargetInput{
			Thread: newMockThread(),
		},
		Text: "Reply text",
	}

	if mutate != nil {
		mutate(&input)
	}

	return input
}

func newMockThread() threads.Thread {
	return threads.NewMockThread("Thread title", "Thread text")
}
