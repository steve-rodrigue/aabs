package campaigns

import (
	"errors"
	"math"
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

	id := uuid.New()
	cluster := newValidCluster()
	createdOn := time.Now()

	result, err := adapter.ToDomain(
		CampaignInput{
			Identifier:  id,
			Name:        " Campaign ",
			Description: " Description ",
			Cluster:     cluster,
			PostCount:   5,
			Confidence:  0.75,
			CreatedOn:   createdOn,
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	if result.Identifier() != id {
		t.Fatalf("expected identifier")
	}

	if result.Name() != "Campaign" {
		t.Fatalf("expected trimmed name, got %q", result.Name())
	}

	if result.Description() != "Description" {
		t.Fatalf("expected trimmed description, got %q", result.Description())
	}

	if result.Cluster() != cluster {
		t.Fatalf("expected cluster")
	}

	if result.PostCount() != 5 {
		t.Fatalf("expected post count 5")
	}

	if result.Confidence() != 0.75 {
		t.Fatalf("expected confidence 0.75")
	}

	if !result.CreatedOn().Equal(createdOn.UTC()) {
		t.Fatalf("expected UTC created on")
	}
}

func TestAdapterToDomainAllowsEmptyDescription(t *testing.T) {
	adapter := NewAdapter()

	result, err := adapter.ToDomain(
		validCampaignInput(func(input *CampaignInput) {
			input.Description = " "
		}),
	)

	if err != nil {
		t.Fatal(err)
	}

	if result.Description() != "" {
		t.Fatalf("expected empty description")
	}
}

func TestAdapterToDomainAcceptsZeroConfidence(t *testing.T) {
	assertAdapterAcceptsConfidence(t, 0)
}

func TestAdapterToDomainAcceptsOneConfidence(t *testing.T) {
	assertAdapterAcceptsConfidence(t, 1)
}

func TestAdapterToDomainReturnsInvalidIdentifierError(t *testing.T) {
	assertAdapterError(
		t,
		func(input *CampaignInput) {
			input.Identifier = uuid.Nil
		},
		ErrInvalidCampaignIdentifier,
	)
}

func TestAdapterToDomainReturnsInvalidNameError(t *testing.T) {
	assertAdapterError(
		t,
		func(input *CampaignInput) {
			input.Name = " "
		},
		ErrInvalidCampaignName,
	)
}

func TestAdapterToDomainReturnsInvalidClusterError(t *testing.T) {
	assertAdapterError(
		t,
		func(input *CampaignInput) {
			input.Cluster = nil
		},
		ErrInvalidCampaignCluster,
	)
}

func TestAdapterToDomainReturnsInvalidPostCountErrorWhenZero(t *testing.T) {
	assertAdapterError(
		t,
		func(input *CampaignInput) {
			input.PostCount = 0
		},
		ErrInvalidCampaignPostCount,
	)
}

func TestAdapterToDomainReturnsInvalidPostCountErrorWhenNegative(t *testing.T) {
	assertAdapterError(
		t,
		func(input *CampaignInput) {
			input.PostCount = -1
		},
		ErrInvalidCampaignPostCount,
	)
}

func TestAdapterToDomainReturnsInvalidConfidenceErrorWhenNegative(t *testing.T) {
	assertAdapterError(
		t,
		func(input *CampaignInput) {
			input.Confidence = -0.1
		},
		ErrInvalidCampaignConfidence,
	)
}

func TestAdapterToDomainReturnsInvalidConfidenceErrorWhenGreaterThanOne(t *testing.T) {
	assertAdapterError(
		t,
		func(input *CampaignInput) {
			input.Confidence = 1.1
		},
		ErrInvalidCampaignConfidence,
	)
}

func TestAdapterToDomainReturnsInvalidConfidenceErrorWhenNaN(t *testing.T) {
	assertAdapterError(
		t,
		func(input *CampaignInput) {
			input.Confidence = math.NaN()
		},
		ErrInvalidCampaignConfidence,
	)
}

func TestAdapterToDomainReturnsInvalidConfidenceErrorWhenInf(t *testing.T) {
	assertAdapterError(
		t,
		func(input *CampaignInput) {
			input.Confidence = math.Inf(1)
		},
		ErrInvalidCampaignConfidence,
	)
}

func TestAdapterToDomainReturnsInvalidCreatedOnError(t *testing.T) {
	assertAdapterError(
		t,
		func(input *CampaignInput) {
			input.CreatedOn = time.Time{}
		},
		ErrInvalidCampaignCreatedOn,
	)
}

func assertAdapterAcceptsConfidence(
	t *testing.T,
	confidence float64,
) {
	t.Helper()

	adapter := NewAdapter()

	result, err := adapter.ToDomain(
		validCampaignInput(func(input *CampaignInput) {
			input.Confidence = confidence
		}),
	)

	if err != nil {
		t.Fatal(err)
	}

	if result.Confidence() != confidence {
		t.Fatalf("expected confidence %f, got %f", confidence, result.Confidence())
	}
}

func assertAdapterError(
	t *testing.T,
	mutate func(input *CampaignInput),
	expected error,
) {
	t.Helper()

	adapter := NewAdapter()

	_, err := adapter.ToDomain(
		validCampaignInput(mutate),
	)

	if !errors.Is(err, expected) {
		t.Fatalf("expected error %v, got %v", expected, err)
	}
}

func validCampaignInput(
	mutate func(input *CampaignInput),
) CampaignInput {
	input := CampaignInput{
		Identifier:  uuid.New(),
		Name:        "Campaign",
		Description: "Description",
		Cluster:     newValidCluster(),
		PostCount:   3,
		Confidence:  0.8,
		CreatedOn:   time.Now().UTC(),
	}

	if mutate != nil {
		mutate(&input)
	}

	return input
}

func newValidCluster() clusters.Cluster {
	return clusters.NewMockCluster(
		clusterables.NewMockClusterable(clusterables.PostKind),
		clusterables.PostKind,
		[]uuid.UUID{uuid.New()},
	)
}
