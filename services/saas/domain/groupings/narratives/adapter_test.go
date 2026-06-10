package narratives

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/clusterables"
	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/participatables"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters"
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
	cluster := validNarrativeCluster()
	createdOn := time.Now()

	result, err := adapter.ToDomain(
		NarrativeInput{
			Identifier:        id,
			ParticipationKind: participatables.NarrativeKind,
			Cluster:           cluster,
			Name:              " Election Integrity ",
			Description:       " Discussion about election integrity. ",
			CreatedOn:         createdOn,
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	if result.Identifier() != id {
		t.Fatalf("expected identifier %s, got %s", id, result.Identifier())
	}

	if result.ParticipationKind() != participatables.NarrativeKind {
		t.Fatalf(
			"expected participation kind %s, got %s",
			participatables.NarrativeKind,
			result.ParticipationKind(),
		)
	}

	if result.Cluster() != cluster {
		t.Fatalf("expected cluster")
	}

	if result.Name() != "Election Integrity" {
		t.Fatalf("expected trimmed name, got %s", result.Name())
	}

	if result.Description() != "Discussion about election integrity." {
		t.Fatalf(
			"expected trimmed description, got %s",
			result.Description(),
		)
	}

	if !result.CreatedOn().Equal(createdOn.UTC()) {
		t.Fatalf("expected UTC created on")
	}
}

func TestAdapterToDomainReturnsInvalidIdentifierError(t *testing.T) {
	assertAdapterError(
		t,
		func(input *NarrativeInput) {
			input.Identifier = uuid.Nil
		},
		ErrInvalidNarrativeIdentifier,
	)
}

func TestAdapterToDomainReturnsInvalidParticipationKindErrorWhenEmpty(t *testing.T) {
	assertAdapterError(
		t,
		func(input *NarrativeInput) {
			input.ParticipationKind = ""
		},
		ErrInvalidNarrativeParticipationKind,
	)
}

func TestAdapterToDomainReturnsInvalidParticipationKindErrorWhenNotNarrative(t *testing.T) {
	assertAdapterError(
		t,
		func(input *NarrativeInput) {
			input.ParticipationKind = participatables.TopicKind
		},
		ErrInvalidNarrativeParticipationKind,
	)
}

func TestAdapterToDomainReturnsInvalidClusterError(t *testing.T) {
	assertAdapterError(
		t,
		func(input *NarrativeInput) {
			input.Cluster = nil
		},
		ErrInvalidNarrativeCluster,
	)
}

func TestAdapterToDomainReturnsInvalidNameErrorWhenEmpty(t *testing.T) {
	assertAdapterError(
		t,
		func(input *NarrativeInput) {
			input.Name = ""
		},
		ErrInvalidNarrativeName,
	)
}

func TestAdapterToDomainReturnsInvalidNameErrorWhenWhitespace(t *testing.T) {
	assertAdapterError(
		t,
		func(input *NarrativeInput) {
			input.Name = "   "
		},
		ErrInvalidNarrativeName,
	)
}

func TestAdapterToDomainReturnsInvalidDescriptionErrorWhenEmpty(t *testing.T) {
	assertAdapterError(
		t,
		func(input *NarrativeInput) {
			input.Description = ""
		},
		ErrInvalidNarrativeDescription,
	)
}

func TestAdapterToDomainReturnsInvalidDescriptionErrorWhenWhitespace(t *testing.T) {
	assertAdapterError(
		t,
		func(input *NarrativeInput) {
			input.Description = "   "
		},
		ErrInvalidNarrativeDescription,
	)
}

func TestAdapterToDomainReturnsInvalidCreatedOnError(t *testing.T) {
	assertAdapterError(
		t,
		func(input *NarrativeInput) {
			input.CreatedOn = time.Time{}
		},
		ErrInvalidNarrativeCreatedOn,
	)
}

func assertAdapterError(
	t *testing.T,
	mutate func(input *NarrativeInput),
	expected error,
) {
	t.Helper()

	_, err := NewAdapter().ToDomain(
		validNarrativeInput(mutate),
	)

	if !errors.Is(err, expected) {
		t.Fatalf("expected error %v, got %v", expected, err)
	}
}

func validNarrativeInput(
	mutate func(input *NarrativeInput),
) NarrativeInput {
	input := NarrativeInput{
		Identifier:        uuid.New(),
		ParticipationKind: participatables.NarrativeKind,
		Cluster:           validNarrativeCluster(),
		Name:              "Election Integrity",
		Description:       "Discussion about election integrity.",
		CreatedOn:         time.Now(),
	}

	if mutate != nil {
		mutate(&input)
	}

	return input
}

func validNarrativeCluster() clusters.Cluster {
	return clusters.NewMockCluster(
		clusterables.NewMockClusterable(clusterables.NarrativeKind),
		clusterables.PostKind,
		[]uuid.UUID{uuid.New()},
	)
}
