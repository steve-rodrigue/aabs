package platforms

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
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
	createdOn := time.Now()

	result, err := adapter.ToDomain(
		PlatformInput{
			Identifier:        id,
			ParticipationKind: participatables.PlatformKind,
			Name:              " Reddit ",
			Handle:            " reddit ",
			BaseURL:           " https://reddit.com ",
			CreatedOn:         createdOn,
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	if result.Identifier() != id {
		t.Fatalf("expected identifier %s, got %s", id, result.Identifier())
	}

	if result.ParticipationKind() != participatables.PlatformKind {
		t.Fatalf(
			"expected participation kind %s, got %s",
			participatables.PlatformKind,
			result.ParticipationKind(),
		)
	}

	if result.Name() != "Reddit" {
		t.Fatalf("expected trimmed name %q, got %q", "Reddit", result.Name())
	}

	if result.Handle() != "reddit" {
		t.Fatalf("expected trimmed handle %q, got %q", "reddit", result.Handle())
	}

	if result.BaseURL() != "https://reddit.com" {
		t.Fatalf(
			"expected trimmed base url %q, got %q",
			"https://reddit.com",
			result.BaseURL(),
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
		validPlatformInput(func(input *PlatformInput) {
			input.Identifier = uuid.Nil
		}),
	)

	if !errors.Is(err, ErrInvalidPlatformIdentifier) {
		t.Fatalf("expected invalid identifier error, got %v", err)
	}
}

func TestAdapterToDomainReturnsInvalidParticipationKindError(t *testing.T) {
	adapter := NewAdapter()

	_, err := adapter.ToDomain(
		validPlatformInput(func(input *PlatformInput) {
			input.ParticipationKind = participatables.UserKind
		}),
	)

	if !errors.Is(err, ErrInvalidPlatformParticipationKind) {
		t.Fatalf("expected invalid participation kind error, got %v", err)
	}
}

func TestAdapterToDomainReturnsInvalidNameError(t *testing.T) {
	adapter := NewAdapter()

	_, err := adapter.ToDomain(
		validPlatformInput(func(input *PlatformInput) {
			input.Name = "   "
		}),
	)

	if !errors.Is(err, ErrInvalidPlatformName) {
		t.Fatalf("expected invalid name error, got %v", err)
	}
}

func TestAdapterToDomainReturnsInvalidHandleError(t *testing.T) {
	adapter := NewAdapter()

	_, err := adapter.ToDomain(
		validPlatformInput(func(input *PlatformInput) {
			input.Handle = "   "
		}),
	)

	if !errors.Is(err, ErrInvalidPlatformHandle) {
		t.Fatalf("expected invalid handle error, got %v", err)
	}
}

func TestAdapterToDomainReturnsInvalidEmptyBaseURLError(t *testing.T) {
	adapter := NewAdapter()

	_, err := adapter.ToDomain(
		validPlatformInput(func(input *PlatformInput) {
			input.BaseURL = "   "
		}),
	)

	if !errors.Is(err, ErrInvalidPlatformBaseURL) {
		t.Fatalf("expected invalid base url error, got %v", err)
	}
}

func TestAdapterToDomainReturnsInvalidBaseURLError(t *testing.T) {
	adapter := NewAdapter()

	_, err := adapter.ToDomain(
		validPlatformInput(func(input *PlatformInput) {
			input.BaseURL = "not a url"
		}),
	)

	if !errors.Is(err, ErrInvalidPlatformBaseURL) {
		t.Fatalf("expected invalid base url error, got %v", err)
	}
}

func TestAdapterToDomainReturnsInvalidBaseURLWithoutHostError(t *testing.T) {
	adapter := NewAdapter()

	_, err := adapter.ToDomain(
		validPlatformInput(func(input *PlatformInput) {
			input.BaseURL = "https://"
		}),
	)

	if !errors.Is(err, ErrInvalidPlatformBaseURL) {
		t.Fatalf("expected invalid base url error, got %v", err)
	}
}

func TestAdapterToDomainReturnsInvalidCreatedOnError(t *testing.T) {
	adapter := NewAdapter()

	_, err := adapter.ToDomain(
		validPlatformInput(func(input *PlatformInput) {
			input.CreatedOn = time.Time{}
		}),
	)

	if !errors.Is(err, ErrInvalidPlatformCreatedOn) {
		t.Fatalf("expected invalid created on error, got %v", err)
	}
}

func validPlatformInput(
	mutate func(input *PlatformInput),
) PlatformInput {
	input := PlatformInput{
		Identifier:        uuid.New(),
		ParticipationKind: participatables.PlatformKind,
		Name:              "Reddit",
		Handle:            "reddit",
		BaseURL:           "https://reddit.com",
		CreatedOn:         time.Now().UTC(),
	}

	if mutate != nil {
		mutate(&input)
	}

	return input
}
