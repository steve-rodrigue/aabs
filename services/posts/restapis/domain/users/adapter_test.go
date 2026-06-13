package users

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/platforms"
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

	result, err := adapter.ToDomain(
		UserInput{
			Identifier:  id,
			Platform:    platform,
			ExternalID:  " 123 ",
			Handle:      " steve ",
			DisplayName: " Steve Rodrigue ",
			ProfileURL:  " https://reddit.com/u/steve-rodrigue ",
			CreatedOn:   createdOn,
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

	if result.ExternalID() != "123" {
		t.Fatalf("expected trimmed external id %q, got %q", "123", result.ExternalID())
	}

	if result.Handle() != "steve" {
		t.Fatalf("expected trimmed handle %q, got %q", "steve", result.Handle())
	}

	if result.DisplayName() != "Steve Rodrigue" {
		t.Fatalf(
			"expected trimmed display name %q, got %q",
			"Steve Rodrigue",
			result.DisplayName(),
		)
	}

	if result.ProfileURL() != "https://reddit.com/u/steve-rodrigue" {
		t.Fatalf(
			"expected trimmed profile url %q, got %q",
			"https://reddit.com/u/steve",
			result.ProfileURL(),
		)
	}

	if !result.CreatedOn().Equal(createdOn.UTC()) {
		t.Fatalf(
			"expected created on %s, got %s",
			createdOn.UTC(),
			result.CreatedOn(),
		)
	}
}

func TestAdapterToDomainReturnsInvalidIdentifierError(t *testing.T) {
	adapter := NewAdapter()

	_, err := adapter.ToDomain(
		validUserInput(func(input *UserInput) {
			input.Identifier = uuid.Nil
		}),
	)

	if !errors.Is(err, ErrInvalidUserIdentifier) {
		t.Fatalf("expected invalid identifier error, got %v", err)
	}
}

func TestAdapterToDomainReturnsInvalidPlatformError(t *testing.T) {
	adapter := NewAdapter()

	_, err := adapter.ToDomain(
		validUserInput(func(input *UserInput) {
			input.Platform = nil
		}),
	)

	if !errors.Is(err, ErrInvalidUserPlatform) {
		t.Fatalf("expected invalid platform error, got %v", err)
	}
}

func TestAdapterToDomainReturnsInvalidExternalIDError(t *testing.T) {
	adapter := NewAdapter()

	_, err := adapter.ToDomain(
		validUserInput(func(input *UserInput) {
			input.ExternalID = "   "
		}),
	)

	if !errors.Is(err, ErrInvalidUserExternalID) {
		t.Fatalf("expected invalid external id error, got %v", err)
	}
}

func TestAdapterToDomainReturnsInvalidHandleError(t *testing.T) {
	adapter := NewAdapter()

	_, err := adapter.ToDomain(
		validUserInput(func(input *UserInput) {
			input.Handle = "   "
		}),
	)

	if !errors.Is(err, ErrInvalidUserHandle) {
		t.Fatalf("expected invalid handle error, got %v", err)
	}
}

func TestAdapterToDomainReturnsInvalidDisplayNameError(t *testing.T) {
	adapter := NewAdapter()

	_, err := adapter.ToDomain(
		validUserInput(func(input *UserInput) {
			input.DisplayName = "   "
		}),
	)

	if !errors.Is(err, ErrInvalidUserDisplayName) {
		t.Fatalf("expected invalid display name error, got %v", err)
	}
}

func TestAdapterToDomainReturnsInvalidEmptyProfileURLError(t *testing.T) {
	adapter := NewAdapter()

	_, err := adapter.ToDomain(
		validUserInput(func(input *UserInput) {
			input.ProfileURL = "   "
		}),
	)

	if !errors.Is(err, ErrInvalidUserProfileURL) {
		t.Fatalf("expected invalid profile url error, got %v", err)
	}
}

func TestAdapterToDomainReturnsInvalidProfileURLError(t *testing.T) {
	adapter := NewAdapter()

	_, err := adapter.ToDomain(
		validUserInput(func(input *UserInput) {
			input.ProfileURL = "not a url"
		}),
	)

	if !errors.Is(err, ErrInvalidUserProfileURL) {
		t.Fatalf("expected invalid profile url error, got %v", err)
	}
}

func TestAdapterToDomainReturnsInvalidProfileURLWithoutHostError(t *testing.T) {
	adapter := NewAdapter()

	_, err := adapter.ToDomain(
		validUserInput(func(input *UserInput) {
			input.ProfileURL = "https://"
		}),
	)

	if !errors.Is(err, ErrInvalidUserProfileURL) {
		t.Fatalf("expected invalid profile url error, got %v", err)
	}
}

func TestAdapterToDomainReturnsInvalidCreatedOnError(t *testing.T) {
	adapter := NewAdapter()

	_, err := adapter.ToDomain(
		validUserInput(func(input *UserInput) {
			input.CreatedOn = time.Time{}
		}),
	)

	if !errors.Is(err, ErrInvalidUserCreatedOn) {
		t.Fatalf("expected invalid created on error, got %v", err)
	}
}

func validUserInput(
	mutate func(input *UserInput),
) UserInput {
	input := UserInput{
		Identifier:  uuid.New(),
		Platform:    platforms.NewMockPlatform("Reddit", "reddit"),
		ExternalID:  "123",
		Handle:      "steve",
		DisplayName: "Steve Rodrigue",
		ProfileURL:  "https://reddit.com/u/steve",
		CreatedOn:   time.Now().UTC(),
	}

	if mutate != nil {
		mutate(&input)
	}

	return input
}
