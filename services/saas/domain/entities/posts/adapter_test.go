package posts

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts/contents"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts/contents/threads"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/users"
)

func TestNewAdapter(t *testing.T) {
	adapter := NewAdapter(contents.NewMockContentAdapter())

	if adapter == nil {
		t.Fatalf("expected adapter")
	}
}

func TestAdapterToDomain(t *testing.T) {
	contentAdapter := contents.NewMockContentAdapter()
	adapter := NewAdapter(contentAdapter)

	id := uuid.New()
	communityID := uuid.New()
	creator := users.NewMockUser("@user", "User")
	createdOn := time.Now()

	result, err := adapter.ToDomain(
		PostInput{
			Identifier: id,
			CommunityIDs: []uuid.UUID{
				communityID,
			},
			Creator: creator,
			Content: contents.ContentInput{
				Identifier: uuid.New(),
				Thread: &threads.ThreadInput{
					Identifier: uuid.New(),
					Creator:    creator,
					Title:      "Thread title",
					Text:       "Post text",
				},
				CreatedAt: time.Now().UTC(),
			},
			CreatedOn: createdOn,
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	if result.Identifier() != id {
		t.Fatalf("expected identifier %s, got %s", id, result.Identifier())
	}

	if len(result.CommunityIDs()) != 1 || result.CommunityIDs()[0] != communityID {
		t.Fatalf("expected community id")
	}

	if result.Creator() != creator {
		t.Fatalf("expected creator")
	}

	if result.Content() == nil {
		t.Fatalf("expected content")
	}

	if !result.CreatedOn().Equal(createdOn.UTC()) {
		t.Fatalf("expected created on %s, got %s", createdOn.UTC(), result.CreatedOn())
	}

	if contentAdapter.ToDomainCalls != 1 {
		t.Fatalf("expected 1 content adapter call, got %d", contentAdapter.ToDomainCalls)
	}
}

func TestAdapterToDomainCopiesCommunityIDs(t *testing.T) {
	adapter := NewAdapter(contents.NewMockContentAdapter())

	communityID := uuid.New()
	communityIDs := []uuid.UUID{
		communityID,
	}

	result, err := adapter.ToDomain(
		validPostInput(func(input *PostInput) {
			input.CommunityIDs = communityIDs
		}),
	)

	if err != nil {
		t.Fatal(err)
	}

	communityIDs[0] = uuid.New()

	if result.CommunityIDs()[0] != communityID {
		t.Fatalf("expected community ids to be copied")
	}
}

func TestAdapterToDomainReturnsInvalidIdentifierError(t *testing.T) {
	adapter := NewAdapter(contents.NewMockContentAdapter())

	_, err := adapter.ToDomain(
		validPostInput(func(input *PostInput) {
			input.Identifier = uuid.Nil
		}),
	)

	if !errors.Is(err, ErrInvalidPostIdentifier) {
		t.Fatalf("expected invalid identifier error, got %v", err)
	}
}

func TestAdapterToDomainReturnsInvalidCreatorError(t *testing.T) {
	adapter := NewAdapter(contents.NewMockContentAdapter())

	_, err := adapter.ToDomain(
		validPostInput(func(input *PostInput) {
			input.Creator = nil
		}),
	)

	if !errors.Is(err, ErrInvalidPostCreator) {
		t.Fatalf("expected invalid creator error, got %v", err)
	}
}

func TestAdapterToDomainReturnsInvalidCreatedOnError(t *testing.T) {
	adapter := NewAdapter(contents.NewMockContentAdapter())

	_, err := adapter.ToDomain(
		validPostInput(func(input *PostInput) {
			input.CreatedOn = time.Time{}
		}),
	)

	if !errors.Is(err, ErrInvalidPostCreatedOn) {
		t.Fatalf("expected invalid created on error, got %v", err)
	}
}

func TestAdapterToDomainReturnsInvalidCommunityIDError(t *testing.T) {
	adapter := NewAdapter(contents.NewMockContentAdapter())

	_, err := adapter.ToDomain(
		validPostInput(func(input *PostInput) {
			input.CommunityIDs = []uuid.UUID{
				uuid.Nil,
			}
		}),
	)

	if !errors.Is(err, ErrInvalidPostCommunityID) {
		t.Fatalf("expected invalid community id error, got %v", err)
	}
}

func TestAdapterToDomainReturnsContentAdapterError(t *testing.T) {
	contentAdapter := contents.NewMockContentAdapter()
	contentAdapter.ToDomainErr = errTest

	adapter := NewAdapter(contentAdapter)

	_, err := adapter.ToDomain(validPostInput(nil))

	if !errors.Is(err, errTest) {
		t.Fatalf("expected content adapter error, got %v", err)
	}
}

func TestAdapterToDomainReturnsInvalidPostContentError(t *testing.T) {
	contentAdapter := contents.NewMockContentAdapter()
	contentAdapter.ToDomainValue = nil

	adapter := NewAdapter(contentAdapter)

	_, err := adapter.ToDomain(validPostInput(nil))

	if !errors.Is(err, ErrInvalidPostContent) {
		t.Fatalf("expected invalid post content error, got %v", err)
	}
}

var errTest = errors.New("test error")

func validPostInput(
	mutate func(input *PostInput),
) PostInput {
	creator := users.NewMockUser("@user", "User")

	input := PostInput{
		Identifier: uuid.New(),
		CommunityIDs: []uuid.UUID{
			uuid.New(),
		},
		Creator: creator,
		Content: contents.ContentInput{
			Identifier: uuid.New(),
			Thread: &threads.ThreadInput{
				Identifier: uuid.New(),
				Creator:    creator,
				Title:      "Thread title",
				Text:       "Post text",
			},
			CreatedAt: time.Now().UTC(),
		},
		CreatedOn: time.Now().UTC(),
	}

	if mutate != nil {
		mutate(&input)
	}

	return input
}
