package topics

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters/clusterables"
)

func TestNewAdapter(t *testing.T) {
	adapter := NewAdapter()

	if adapter == nil {
		t.Fatalf("expected adapter")
	}
}

func TestAdapterToDomain(t *testing.T) {
	adapter := NewAdapter()

	cluster := clusters.NewMockCluster(
		clusterables.NewMockClusterable(clusterables.TopicKind),
		clusterables.PostKind,
		[]uuid.UUID{uuid.New()},
	)

	parent := &topic{
		identifier: uuid.New(),
		cluster:    cluster,
		name:       "Parent Topic",
		createdOn:  time.Now().UTC(),
	}

	id := uuid.New()
	createdOn := time.Now()

	result, err := adapter.ToDomain(
		TopicInput{
			Identifier:  id,
			Cluster:     cluster,
			Name:        "  AI Spam  ",
			Description: "  Posts about AI spam  ",
			Parent:      parent,
			CreatedOn:   createdOn,
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	if result.Identifier() != id {
		t.Fatalf("expected identifier %s, got %s", id, result.Identifier())
	}

	if result.Cluster() != cluster {
		t.Fatalf("expected cluster")
	}

	if result.Name() != "AI Spam" {
		t.Fatalf("expected trimmed name AI Spam, got %s", result.Name())
	}

	if result.Description() != "Posts about AI spam" {
		t.Fatalf("expected trimmed description, got %s", result.Description())
	}

	if !result.HasParent() {
		t.Fatalf("expected parent")
	}

	if result.Parent() != parent {
		t.Fatalf("expected parent")
	}

	if !result.CreatedOn().Equal(createdOn.UTC()) {
		t.Fatalf("expected UTC created on")
	}
}

func TestAdapterToDomainWithoutParent(t *testing.T) {
	adapter := NewAdapter()

	result, err := adapter.ToDomain(
		validTopicInput(func(input *TopicInput) {
			input.Parent = nil
		}),
	)

	if err != nil {
		t.Fatal(err)
	}

	if result.HasParent() {
		t.Fatalf("expected no parent")
	}

	if result.Parent() != nil {
		t.Fatalf("expected nil parent")
	}
}

func TestAdapterToDomainAllowsEmptyDescription(t *testing.T) {
	adapter := NewAdapter()

	result, err := adapter.ToDomain(
		validTopicInput(func(input *TopicInput) {
			input.Description = "   "
		}),
	)

	if err != nil {
		t.Fatal(err)
	}

	if result.Description() != "" {
		t.Fatalf("expected empty description, got %s", result.Description())
	}
}

func TestAdapterToDomainReturnsInvalidIdentifierError(t *testing.T) {
	assertAdapterError(
		t,
		func(input *TopicInput) {
			input.Identifier = uuid.Nil
		},
		ErrInvalidTopicIdentifier,
	)
}

func TestAdapterToDomainReturnsInvalidClusterError(t *testing.T) {
	assertAdapterError(
		t,
		func(input *TopicInput) {
			input.Cluster = nil
		},
		ErrInvalidTopicCluster,
	)
}

func TestAdapterToDomainReturnsInvalidNameErrorWhenEmpty(t *testing.T) {
	assertAdapterError(
		t,
		func(input *TopicInput) {
			input.Name = ""
		},
		ErrInvalidTopicName,
	)
}

func TestAdapterToDomainReturnsInvalidNameErrorWhenWhitespace(t *testing.T) {
	assertAdapterError(
		t,
		func(input *TopicInput) {
			input.Name = "   "
		},
		ErrInvalidTopicName,
	)
}

func TestAdapterToDomainReturnsInvalidCreatedOnError(t *testing.T) {
	assertAdapterError(
		t,
		func(input *TopicInput) {
			input.CreatedOn = time.Time{}
		},
		ErrInvalidTopicCreatedOn,
	)
}

func assertAdapterError(
	t *testing.T,
	mutate func(input *TopicInput),
	expected error,
) {
	t.Helper()

	adapter := NewAdapter()

	_, err := adapter.ToDomain(
		validTopicInput(mutate),
	)

	if !errors.Is(err, expected) {
		t.Fatalf("expected error %v, got %v", expected, err)
	}
}

func validTopicInput(
	mutate func(input *TopicInput),
) TopicInput {
	input := TopicInput{
		Identifier: uuid.New(),
		Cluster: clusters.NewMockCluster(
			clusterables.NewMockClusterable(clusterables.TopicKind),
			clusterables.PostKind,
			[]uuid.UUID{uuid.New()},
		),
		Name:        "AI Spam",
		Description: "Posts about AI spam",
		CreatedOn:   time.Now(),
	}

	if mutate != nil {
		mutate(&input)
	}

	return input
}
