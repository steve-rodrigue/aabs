package relatables

import (
	"errors"
	"testing"

	"github.com/google/uuid"
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

	result, err := adapter.ToDomain(
		RelatableInput{
			Identifier:       id,
			RelationshipKind: PostKind,
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	if result.Identifier() != id {
		t.Fatalf("expected identifier %s, got %s", id, result.Identifier())
	}

	if result.RelationshipKind() != PostKind {
		t.Fatalf(
			"expected relationship kind %s, got %s",
			PostKind,
			result.RelationshipKind(),
		)
	}
}

func TestAdapterToDomainAcceptsCampaignKind(t *testing.T) {
	assertAdapterAcceptsKind(t, CampaignKind)
}

func TestAdapterToDomainAcceptsTopicKind(t *testing.T) {
	assertAdapterAcceptsKind(t, TopicKind)
}

func TestAdapterToDomainAcceptsUserKind(t *testing.T) {
	assertAdapterAcceptsKind(t, UserKind)
}

func TestAdapterToDomainAcceptsPostKind(t *testing.T) {
	assertAdapterAcceptsKind(t, PostKind)
}

func TestAdapterToDomainAcceptsNarrativeKind(t *testing.T) {
	assertAdapterAcceptsKind(t, NarrativeKind)
}

func TestAdapterToDomainReturnsInvalidIdentifierError(t *testing.T) {
	adapter := NewAdapter()

	_, err := adapter.ToDomain(
		validRelatableInput(func(input *RelatableInput) {
			input.Identifier = uuid.Nil
		}),
	)

	if !errors.Is(err, ErrInvalidRelatableIdentifier) {
		t.Fatalf("expected invalid identifier error, got %v", err)
	}
}

func TestAdapterToDomainReturnsInvalidRelationshipKindError(t *testing.T) {
	adapter := NewAdapter()

	_, err := adapter.ToDomain(
		validRelatableInput(func(input *RelatableInput) {
			input.RelationshipKind = Kind("")
		}),
	)

	if !errors.Is(err, ErrInvalidRelatableRelationshipKind) {
		t.Fatalf("expected invalid relationship kind error, got %v", err)
	}
}

func TestAdapterToDomainReturnsInvalidUnknownRelationshipKindError(t *testing.T) {
	adapter := NewAdapter()

	_, err := adapter.ToDomain(
		validRelatableInput(func(input *RelatableInput) {
			input.RelationshipKind = Kind("unknown")
		}),
	)

	if !errors.Is(err, ErrInvalidRelatableRelationshipKind) {
		t.Fatalf("expected invalid relationship kind error, got %v", err)
	}
}

func assertAdapterAcceptsKind(
	t *testing.T,
	kind Kind,
) {
	t.Helper()

	adapter := NewAdapter()

	result, err := adapter.ToDomain(
		RelatableInput{
			Identifier:       uuid.New(),
			RelationshipKind: kind,
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	if result.RelationshipKind() != kind {
		t.Fatalf(
			"expected relationship kind %s, got %s",
			kind,
			result.RelationshipKind(),
		)
	}
}

func validRelatableInput(
	mutate func(input *RelatableInput),
) RelatableInput {
	input := RelatableInput{
		Identifier:       uuid.New(),
		RelationshipKind: PostKind,
	}

	if mutate != nil {
		mutate(&input)
	}

	return input
}
