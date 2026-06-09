package clusterables

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
		ClusterableInput{
			Identifier:  id,
			ClusterKind: PostKind,
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	if result.Identifier() != id {
		t.Fatalf("expected identifier %s, got %s", id, result.Identifier())
	}

	if result.ClusterKind() != PostKind {
		t.Fatalf(
			"expected cluster kind %s, got %s",
			PostKind,
			result.ClusterKind(),
		)
	}
}

func TestAdapterToDomainAcceptsPostKind(t *testing.T) {
	assertAdapterAcceptsKind(t, PostKind)
}

func TestAdapterToDomainAcceptsUserKind(t *testing.T) {
	assertAdapterAcceptsKind(t, UserKind)
}

func TestAdapterToDomainAcceptsCommunityKind(t *testing.T) {
	assertAdapterAcceptsKind(t, CommunityKind)
}

func TestAdapterToDomainAcceptsPlatformKind(t *testing.T) {
	assertAdapterAcceptsKind(t, PlatformKind)
}

func TestAdapterToDomainAcceptsCampaignKind(t *testing.T) {
	assertAdapterAcceptsKind(t, CampaignKind)
}

func TestAdapterToDomainAcceptsTopicKind(t *testing.T) {
	assertAdapterAcceptsKind(t, TopicKind)
}

func TestAdapterToDomainAcceptsNarrativeKind(t *testing.T) {
	assertAdapterAcceptsKind(t, NarrativeKind)
}

func TestAdapterToDomainReturnsInvalidIdentifierError(t *testing.T) {
	adapter := NewAdapter()

	_, err := adapter.ToDomain(
		validClusterableInput(func(input *ClusterableInput) {
			input.Identifier = uuid.Nil
		}),
	)

	if !errors.Is(err, ErrInvalidClusterableIdentifier) {
		t.Fatalf("expected invalid identifier error, got %v", err)
	}
}

func TestAdapterToDomainReturnsInvalidKindErrorWhenEmpty(t *testing.T) {
	adapter := NewAdapter()

	_, err := adapter.ToDomain(
		validClusterableInput(func(input *ClusterableInput) {
			input.ClusterKind = Kind("")
		}),
	)

	if !errors.Is(err, ErrInvalidClusterableKind) {
		t.Fatalf("expected invalid kind error, got %v", err)
	}
}

func TestAdapterToDomainReturnsInvalidKindErrorWhenUnknown(t *testing.T) {
	adapter := NewAdapter()

	_, err := adapter.ToDomain(
		validClusterableInput(func(input *ClusterableInput) {
			input.ClusterKind = Kind("unknown")
		}),
	)

	if !errors.Is(err, ErrInvalidClusterableKind) {
		t.Fatalf("expected invalid kind error, got %v", err)
	}
}

func assertAdapterAcceptsKind(
	t *testing.T,
	kind Kind,
) {
	t.Helper()

	adapter := NewAdapter()

	result, err := adapter.ToDomain(
		ClusterableInput{
			Identifier:  uuid.New(),
			ClusterKind: kind,
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	if result.ClusterKind() != kind {
		t.Fatalf(
			"expected cluster kind %s, got %s",
			kind,
			result.ClusterKind(),
		)
	}
}

func validClusterableInput(
	mutate func(input *ClusterableInput),
) ClusterableInput {
	input := ClusterableInput{
		Identifier:  uuid.New(),
		ClusterKind: PostKind,
	}

	if mutate != nil {
		mutate(&input)
	}

	return input
}
