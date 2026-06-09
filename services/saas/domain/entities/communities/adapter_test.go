package communities

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/platforms"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/users"
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
	platform := platforms.NewMockPlatform("Reddit", "reddit")
	createdOn := time.Now()
	moderator := users.NewMockUser("@mod", "Moderator")

	result, err := adapter.ToDomain(
		CommunityInput{
			Identifier: id,
			Platform:   platform,
			Handle:     " aabs ",
			Title:      " AABS ",
			Text:       " Anti-AI Bot Spam community ",
			CreatedOn:  createdOn,
			Moderators: []users.User{
				moderator,
			},
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	if result.Identifier() != id {
		t.Fatalf("expected identifier %s, got %s", id, result.Identifier())
	}

	if result.Platform() != platform {
		t.Fatalf("expected platform")
	}

	if result.Handle() != "aabs" {
		t.Fatalf("expected trimmed handle %q, got %q", "aabs", result.Handle())
	}

	if result.Title() != "AABS" {
		t.Fatalf("expected trimmed title %q, got %q", "AABS", result.Title())
	}

	if result.Text() != "Anti-AI Bot Spam community" {
		t.Fatalf(
			"expected trimmed text %q, got %q",
			"Anti-AI Bot Spam community",
			result.Text(),
		)
	}

	if !result.CreatedOn().Equal(createdOn.UTC()) {
		t.Fatalf(
			"expected created on %s, got %s",
			createdOn.UTC(),
			result.CreatedOn(),
		)
	}

	if !result.HasModerators() {
		t.Fatalf("expected moderators")
	}

	if len(result.Moderators()) != 1 || result.Moderators()[0] != moderator {
		t.Fatalf("expected moderator")
	}
}

func TestAdapterToDomainAllowsNoModerators(t *testing.T) {
	adapter := NewAdapter()

	result, err := adapter.ToDomain(validCommunityInput(nil))

	if err != nil {
		t.Fatal(err)
	}

	if result.HasModerators() {
		t.Fatalf("expected no moderators")
	}

	if len(result.Moderators()) != 0 {
		t.Fatalf("expected empty moderators")
	}
}

func TestAdapterToDomainReturnsInvalidIdentifierError(t *testing.T) {
	adapter := NewAdapter()

	_, err := adapter.ToDomain(
		validCommunityInput(func(input *CommunityInput) {
			input.Identifier = uuid.Nil
		}),
	)

	if !errors.Is(err, ErrInvalidCommunityIdentifier) {
		t.Fatalf("expected invalid identifier error, got %v", err)
	}
}

func TestAdapterToDomainReturnsInvalidPlatformError(t *testing.T) {
	adapter := NewAdapter()

	_, err := adapter.ToDomain(
		validCommunityInput(func(input *CommunityInput) {
			input.Platform = nil
		}),
	)

	if !errors.Is(err, ErrInvalidCommunityPlatform) {
		t.Fatalf("expected invalid platform error, got %v", err)
	}
}

func TestAdapterToDomainReturnsInvalidHandleError(t *testing.T) {
	adapter := NewAdapter()

	_, err := adapter.ToDomain(
		validCommunityInput(func(input *CommunityInput) {
			input.Handle = "   "
		}),
	)

	if !errors.Is(err, ErrInvalidCommunityHandle) {
		t.Fatalf("expected invalid handle error, got %v", err)
	}
}

func TestAdapterToDomainReturnsInvalidTitleError(t *testing.T) {
	adapter := NewAdapter()

	_, err := adapter.ToDomain(
		validCommunityInput(func(input *CommunityInput) {
			input.Title = "   "
		}),
	)

	if !errors.Is(err, ErrInvalidCommunityTitle) {
		t.Fatalf("expected invalid title error, got %v", err)
	}
}

func TestAdapterToDomainReturnsInvalidTextError(t *testing.T) {
	adapter := NewAdapter()

	_, err := adapter.ToDomain(
		validCommunityInput(func(input *CommunityInput) {
			input.Text = "   "
		}),
	)

	if !errors.Is(err, ErrInvalidCommunityText) {
		t.Fatalf("expected invalid text error, got %v", err)
	}
}

func TestAdapterToDomainReturnsInvalidCreatedOnError(t *testing.T) {
	adapter := NewAdapter()

	_, err := adapter.ToDomain(
		validCommunityInput(func(input *CommunityInput) {
			input.CreatedOn = time.Time{}
		}),
	)

	if !errors.Is(err, ErrInvalidCommunityCreatedOn) {
		t.Fatalf("expected invalid created on error, got %v", err)
	}
}

func TestAdapterToDomainReturnsInvalidModeratorError(t *testing.T) {
	adapter := NewAdapter()

	_, err := adapter.ToDomain(
		validCommunityInput(func(input *CommunityInput) {
			input.Moderators = []users.User{
				nil,
			}
		}),
	)

	if !errors.Is(err, ErrInvalidCommunityModerator) {
		t.Fatalf("expected invalid moderator error, got %v", err)
	}
}

func TestAdapterToDomainCopiesModerators(t *testing.T) {
	adapter := NewAdapter()

	moderator := users.NewMockUser("@mod", "Moderator")
	moderators := []users.User{
		moderator,
	}

	result, err := adapter.ToDomain(
		validCommunityInput(func(input *CommunityInput) {
			input.Moderators = moderators
		}),
	)

	if err != nil {
		t.Fatal(err)
	}

	moderators[0] = users.NewMockUser("@other", "Other")

	if result.Moderators()[0] != moderator {
		t.Fatalf("expected moderators to be copied")
	}
}

func validCommunityInput(
	mutate func(input *CommunityInput),
) CommunityInput {
	input := CommunityInput{
		Identifier: uuid.New(),
		Platform:   platforms.NewMockPlatform("Reddit", "reddit"),
		Handle:     "aabs",
		Title:      "AABS",
		Text:       "Anti-AI Bot Spam community",
		CreatedOn:  time.Now().UTC(),
		Moderators: []users.User{},
	}

	if mutate != nil {
		mutate(&input)
	}

	return input
}
