package contents

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts/contents/replies"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts/contents/threads"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/users"
)

func TestNewAdapter(t *testing.T) {
	adapter := NewAdapter(
		replies.NewMockReplyAdapter(),
		threads.NewMockThreadAdapter(),
	)

	if adapter == nil {
		t.Fatalf("expected adapter")
	}
}

func TestAdapterToDomainWithReply(t *testing.T) {
	replyAdapter := replies.NewMockReplyAdapter()
	threadAdapter := threads.NewMockThreadAdapter()
	adapter := NewAdapter(replyAdapter, threadAdapter)

	id := uuid.New()
	createdAt := time.Now()

	result, err := adapter.ToDomain(
		ContentInput{
			Identifier: id,
			Reply: &replies.ReplyInput{
				Identifier: uuid.New(),
				Target: replies.TargetInput{
					Thread: threads.NewMockThread("Thread title", "Thread text"),
				},
				Text: "Reply text",
			},
			CreatedAt: createdAt,
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	if result.Identifier() != id {
		t.Fatalf("expected identifier %s, got %s", id, result.Identifier())
	}

	if !result.IsReply() {
		t.Fatalf("expected content to be reply")
	}

	if result.Reply() == nil {
		t.Fatalf("expected reply")
	}

	if result.IsThread() {
		t.Fatalf("expected content not to be thread")
	}

	if result.Thread() != nil {
		t.Fatalf("expected nil thread")
	}

	if result.Text() != "Reply text" {
		t.Fatalf("expected reply text %q, got %q", "Reply text", result.Text())
	}

	if !result.CreatedAt().Equal(createdAt.UTC()) {
		t.Fatalf(
			"expected created at %s, got %s",
			createdAt.UTC(),
			result.CreatedAt(),
		)
	}

	if replyAdapter.ToDomainCalls != 1 {
		t.Fatalf("expected 1 reply adapter call, got %d", replyAdapter.ToDomainCalls)
	}

	if threadAdapter.ToDomainCalls != 0 {
		t.Fatalf("expected thread adapter not to be called")
	}
}

func TestAdapterToDomainWithThread(t *testing.T) {
	replyAdapter := replies.NewMockReplyAdapter()
	threadAdapter := threads.NewMockThreadAdapter()
	adapter := NewAdapter(replyAdapter, threadAdapter)

	id := uuid.New()
	createdAt := time.Now()

	result, err := adapter.ToDomain(
		ContentInput{
			Identifier: id,
			Thread: &threads.ThreadInput{
				Identifier: uuid.New(),
				Creator:    users.NewMockUser("@user", "User"),
				Title:      "Thread title",
				Text:       "Thread text",
			},
			CreatedAt: createdAt,
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	if result.Identifier() != id {
		t.Fatalf("expected identifier %s, got %s", id, result.Identifier())
	}

	if result.IsReply() {
		t.Fatalf("expected content not to be reply")
	}

	if result.Reply() != nil {
		t.Fatalf("expected nil reply")
	}

	if !result.IsThread() {
		t.Fatalf("expected content to be thread")
	}

	if result.Thread() == nil {
		t.Fatalf("expected thread")
	}

	if result.Text() != "Thread text" {
		t.Fatalf("expected thread text %q, got %q", "Thread text", result.Text())
	}

	if !result.CreatedAt().Equal(createdAt.UTC()) {
		t.Fatalf(
			"expected created at %s, got %s",
			createdAt.UTC(),
			result.CreatedAt(),
		)
	}

	if replyAdapter.ToDomainCalls != 0 {
		t.Fatalf("expected reply adapter not to be called")
	}

	if threadAdapter.ToDomainCalls != 1 {
		t.Fatalf("expected 1 thread adapter call, got %d", threadAdapter.ToDomainCalls)
	}
}

func TestAdapterToDomainReturnsInvalidIdentifierError(t *testing.T) {
	adapter := newTestAdapter()

	_, err := adapter.ToDomain(
		validContentInput(func(input *ContentInput) {
			input.Identifier = uuid.Nil
		}),
	)

	if !errors.Is(err, ErrInvalidContentIdentifier) {
		t.Fatalf("expected invalid identifier error, got %v", err)
	}
}

func TestAdapterToDomainReturnsInvalidCreatedAtError(t *testing.T) {
	adapter := newTestAdapter()

	_, err := adapter.ToDomain(
		validContentInput(func(input *ContentInput) {
			input.CreatedAt = time.Time{}
		}),
	)

	if !errors.Is(err, ErrInvalidContentCreatedAt) {
		t.Fatalf("expected invalid created at error, got %v", err)
	}
}

func TestAdapterToDomainReturnsInvalidTargetErrorWhenEmpty(t *testing.T) {
	adapter := newTestAdapter()

	_, err := adapter.ToDomain(
		validContentInput(func(input *ContentInput) {
			input.Reply = nil
			input.Thread = nil
		}),
	)

	if !errors.Is(err, ErrInvalidContentTarget) {
		t.Fatalf("expected invalid content target error, got %v", err)
	}
}

func TestAdapterToDomainReturnsInvalidTargetErrorWhenReplyAndThreadAreSet(t *testing.T) {
	adapter := newTestAdapter()

	_, err := adapter.ToDomain(
		validContentInput(func(input *ContentInput) {
			input.Thread = &threads.ThreadInput{
				Identifier: uuid.New(),
				Creator:    users.NewMockUser("@user", "User"),
				Title:      "Thread title",
				Text:       "Thread text",
			}
		}),
	)

	if !errors.Is(err, ErrInvalidContentTarget) {
		t.Fatalf("expected invalid content target error, got %v", err)
	}
}

func TestAdapterToDomainReturnsReplyAdapterError(t *testing.T) {
	replyAdapter := replies.NewMockReplyAdapter()
	threadAdapter := threads.NewMockThreadAdapter()
	adapter := NewAdapter(replyAdapter, threadAdapter)

	replyAdapter.ToDomainErr = errTest

	_, err := adapter.ToDomain(validContentInput(nil))

	if !errors.Is(err, errTest) {
		t.Fatalf("expected reply adapter error, got %v", err)
	}
}

func TestAdapterToDomainReturnsThreadAdapterError(t *testing.T) {
	replyAdapter := replies.NewMockReplyAdapter()
	threadAdapter := threads.NewMockThreadAdapter()
	adapter := NewAdapter(replyAdapter, threadAdapter)

	threadAdapter.ToDomainErr = errTest

	_, err := adapter.ToDomain(
		ContentInput{
			Identifier: uuid.New(),
			Thread: &threads.ThreadInput{
				Identifier: uuid.New(),
				Creator:    users.NewMockUser("@user", "User"),
				Title:      "Thread title",
				Text:       "Thread text",
			},
			CreatedAt: time.Now().UTC(),
		},
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected thread adapter error, got %v", err)
	}
}

var errTest = errors.New("test error")

func newTestAdapter() Adapter {
	return NewAdapter(
		replies.NewMockReplyAdapter(),
		threads.NewMockThreadAdapter(),
	)
}

func validContentInput(
	mutate func(input *ContentInput),
) ContentInput {
	input := ContentInput{
		Identifier: uuid.New(),
		Reply: &replies.ReplyInput{
			Identifier: uuid.New(),
			Target: replies.TargetInput{
				Thread: threads.NewMockThread("Thread title", "Thread text"),
			},
			Text: "Reply text",
		},
		CreatedAt: time.Now().UTC(),
	}

	if mutate != nil {
		mutate(&input)
	}

	return input
}
