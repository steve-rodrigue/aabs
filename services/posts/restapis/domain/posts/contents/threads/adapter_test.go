package threads

import (
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/users"
)

func TestNewAdapter(t *testing.T) {
	adapter := NewAdapter()

	if adapter == nil {
		t.Fatalf("expected adapter")
	}
}

func TestAdapterToDomain(t *testing.T) {
	adapter := NewAdapter()

	id := uuid.New()
	creator := users.NewMockUser("@user", "User")

	result, err := adapter.ToDomain(
		ThreadInput{
			Identifier: id,
			Creator:    creator,
			Title:      " Thread title ",
			Text:       " Thread text ",
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	if result.Identifier() != id {
		t.Fatalf("expected identifier %s, got %s", id, result.Identifier())
	}

	if result.Creator() != creator {
		t.Fatalf("expected creator")
	}

	if result.Title() != "Thread title" {
		t.Fatalf("expected trimmed title %q, got %q", "Thread title", result.Title())
	}

	if result.Text() != "Thread text" {
		t.Fatalf("expected trimmed text %q, got %q", "Thread text", result.Text())
	}
}

func TestAdapterToDomainReturnsInvalidIdentifierError(t *testing.T) {
	adapter := NewAdapter()

	_, err := adapter.ToDomain(
		validThreadInput(func(input *ThreadInput) {
			input.Identifier = uuid.Nil
		}),
	)

	if !errors.Is(err, ErrInvalidThreadIdentifier) {
		t.Fatalf("expected invalid identifier error, got %v", err)
	}
}

func TestAdapterToDomainReturnsInvalidCreatorError(t *testing.T) {
	adapter := NewAdapter()

	_, err := adapter.ToDomain(
		validThreadInput(func(input *ThreadInput) {
			input.Creator = nil
		}),
	)

	if !errors.Is(err, ErrInvalidThreadCreator) {
		t.Fatalf("expected invalid creator error, got %v", err)
	}
}

func TestAdapterToDomainReturnsInvalidTitleError(t *testing.T) {
	adapter := NewAdapter()

	_, err := adapter.ToDomain(
		validThreadInput(func(input *ThreadInput) {
			input.Title = "   "
		}),
	)

	if !errors.Is(err, ErrInvalidThreadTitle) {
		t.Fatalf("expected invalid title error, got %v", err)
	}
}

func TestAdapterToDomainReturnsInvalidTextError(t *testing.T) {
	adapter := NewAdapter()

	_, err := adapter.ToDomain(
		validThreadInput(func(input *ThreadInput) {
			input.Text = "   "
		}),
	)

	if !errors.Is(err, ErrInvalidThreadText) {
		t.Fatalf("expected invalid text error, got %v", err)
	}
}

func validThreadInput(
	mutate func(input *ThreadInput),
) ThreadInput {
	input := ThreadInput{
		Identifier: uuid.New(),
		Creator:    users.NewMockUser("@user", "User"),
		Title:      "Thread title",
		Text:       "Thread text",
	}

	if mutate != nil {
		mutate(&input)
	}

	return input
}
