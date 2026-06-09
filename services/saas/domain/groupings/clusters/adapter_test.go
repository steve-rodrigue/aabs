package clusters

import (
	"errors"
	"math"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters/clusterables"
)

var errTest = errors.New("test error")

func TestNewAdapter(t *testing.T) {
	adapter := NewAdapter(
		clusterables.NewMockClusterableAdapter(),
	)

	if adapter == nil {
		t.Fatalf("expected adapter")
	}
}

func TestAdapterToDomain(t *testing.T) {
	clusterableAdapter := clusterables.NewMockClusterableAdapter()
	adapter := NewAdapter(clusterableAdapter)

	id := uuid.New()
	targetID := uuid.New()
	memberIDs := []uuid.UUID{
		uuid.New(),
		uuid.New(),
	}
	centroid := []float32{
		0.1,
		0.2,
	}
	createdOn := time.Now()

	result, err := adapter.ToDomain(
		ClusterInput{
			Identifier: id,
			Target: clusterables.ClusterableInput{
				Identifier:  targetID,
				ClusterKind: clusterables.PostKind,
			},
			MemberIDs:       memberIDs,
			MemberKind:      clusterables.PostKind,
			ConfidenceScore: 0.8,
			Centroid:        centroid,
			CreatedOn:       createdOn,
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	if clusterableAdapter.ToDomainCalls != 1 {
		t.Fatalf(
			"expected 1 clusterable adapter call, got %d",
			clusterableAdapter.ToDomainCalls,
		)
	}

	if result.Identifier() != id {
		t.Fatalf("expected identifier %s, got %s", id, result.Identifier())
	}

	if result.Target().Identifier() != targetID {
		t.Fatalf(
			"expected target id %s, got %s",
			targetID,
			result.Target().Identifier(),
		)
	}

	if result.Target().ClusterKind() != clusterables.PostKind {
		t.Fatalf("expected target kind")
	}

	resultMemberIDs := result.MemberIDs()

	if len(resultMemberIDs) != len(memberIDs) {
		t.Fatalf("expected %d member ids, got %d", len(memberIDs), len(resultMemberIDs))
	}

	for index := range memberIDs {
		if resultMemberIDs[index] != memberIDs[index] {
			t.Fatalf("expected member id at index %d", index)
		}
	}

	if result.MemberKind() != clusterables.PostKind {
		t.Fatalf("expected member kind")
	}

	if result.ConfidenceScore() != 0.8 {
		t.Fatalf("expected confidence score 0.8, got %f", result.ConfidenceScore())
	}

	resultCentroid := result.Centroid()

	if len(resultCentroid) != len(centroid) {
		t.Fatalf("expected centroid length %d, got %d", len(centroid), len(resultCentroid))
	}

	for index := range centroid {
		if resultCentroid[index] != centroid[index] {
			t.Fatalf("expected centroid at index %d", index)
		}
	}

	if !result.CreatedOn().Equal(createdOn.UTC()) {
		t.Fatalf("expected UTC created on")
	}
}

func TestAdapterToDomainAcceptsZeroConfidenceScore(t *testing.T) {
	assertAdapterAcceptsConfidenceScore(t, 0)
}

func TestAdapterToDomainAcceptsMaximumConfidenceScore(t *testing.T) {
	assertAdapterAcceptsConfidenceScore(t, 1)
}

func TestAdapterToDomainCopiesMemberIDs(t *testing.T) {
	adapter := NewAdapter(clusterables.NewMockClusterableAdapter())

	memberIDs := []uuid.UUID{
		uuid.New(),
	}

	result, err := adapter.ToDomain(
		validClusterInput(func(input *ClusterInput) {
			input.MemberIDs = memberIDs
		}),
	)

	if err != nil {
		t.Fatal(err)
	}

	memberIDs[0] = uuid.New()

	if result.MemberIDs()[0] == memberIDs[0] {
		t.Fatalf("expected member ids copy")
	}
}

func TestAdapterToDomainCopiesCentroid(t *testing.T) {
	adapter := NewAdapter(clusterables.NewMockClusterableAdapter())

	centroid := []float32{
		0.1,
	}

	result, err := adapter.ToDomain(
		validClusterInput(func(input *ClusterInput) {
			input.Centroid = centroid
		}),
	)

	if err != nil {
		t.Fatal(err)
	}

	centroid[0] = 99

	if result.Centroid()[0] == 99 {
		t.Fatalf("expected centroid copy")
	}
}

func TestAdapterToDomainReturnsInvalidIdentifierError(t *testing.T) {
	assertAdapterError(
		t,
		func(input *ClusterInput) {
			input.Identifier = uuid.Nil
		},
		ErrInvalidClusterIdentifier,
	)
}

func TestAdapterToDomainReturnsClusterableAdapterError(t *testing.T) {
	clusterableAdapter := clusterables.NewMockClusterableAdapter()
	clusterableAdapter.ToDomainErr = errTest

	adapter := NewAdapter(clusterableAdapter)

	_, err := adapter.ToDomain(validClusterInput(nil))

	if !errors.Is(err, errTest) {
		t.Fatalf("expected clusterable adapter error, got %v", err)
	}
}

func TestAdapterToDomainReturnsInvalidMemberIDErrorWhenMembersAreEmpty(t *testing.T) {
	assertAdapterError(
		t,
		func(input *ClusterInput) {
			input.MemberIDs = []uuid.UUID{}
		},
		ErrInvalidClusterMemberID,
	)
}

func TestAdapterToDomainReturnsInvalidMemberIDErrorWhenMemberIsNil(t *testing.T) {
	assertAdapterError(
		t,
		func(input *ClusterInput) {
			input.MemberIDs = []uuid.UUID{
				uuid.Nil,
			}
		},
		ErrInvalidClusterMemberID,
	)
}

func TestAdapterToDomainReturnsInvalidMemberKindError(t *testing.T) {
	assertAdapterError(
		t,
		func(input *ClusterInput) {
			input.MemberKind = clusterables.Kind("invalid")
		},
		ErrInvalidClusterMemberKind,
	)
}

func TestAdapterToDomainReturnsInvalidConfidenceScoreErrorWhenNegative(t *testing.T) {
	assertAdapterError(
		t,
		func(input *ClusterInput) {
			input.ConfidenceScore = -0.1
		},
		ErrInvalidClusterConfidenceScore,
	)
}

func TestAdapterToDomainReturnsInvalidConfidenceScoreErrorWhenGreaterThanOne(t *testing.T) {
	assertAdapterError(
		t,
		func(input *ClusterInput) {
			input.ConfidenceScore = 1.1
		},
		ErrInvalidClusterConfidenceScore,
	)
}

func TestAdapterToDomainReturnsInvalidConfidenceScoreErrorWhenNaN(t *testing.T) {
	assertAdapterError(
		t,
		func(input *ClusterInput) {
			input.ConfidenceScore = math.NaN()
		},
		ErrInvalidClusterConfidenceScore,
	)
}

func TestAdapterToDomainReturnsInvalidConfidenceScoreErrorWhenInf(t *testing.T) {
	assertAdapterError(
		t,
		func(input *ClusterInput) {
			input.ConfidenceScore = math.Inf(1)
		},
		ErrInvalidClusterConfidenceScore,
	)
}

func TestAdapterToDomainReturnsInvalidCentroidErrorWhenNil(t *testing.T) {
	assertAdapterError(
		t,
		func(input *ClusterInput) {
			input.Centroid = nil
		},
		ErrInvalidClusterCentroid,
	)
}

func TestAdapterToDomainReturnsInvalidCentroidErrorWhenEmpty(t *testing.T) {
	assertAdapterError(
		t,
		func(input *ClusterInput) {
			input.Centroid = []float32{}
		},
		ErrInvalidClusterCentroid,
	)
}

func TestAdapterToDomainReturnsInvalidTargetErrorWhenClusterableAdapterReturnsNil(
	t *testing.T,
) {
	clusterableAdapter := clusterables.NewMockClusterableAdapter()
	clusterableAdapter.ToDomainReturnsNil = true

	adapter := NewAdapter(
		clusterableAdapter,
	)

	_, err := adapter.ToDomain(
		validClusterInput(nil),
	)

	if !errors.Is(
		err,
		ErrInvalidClusterTarget,
	) {
		t.Fatalf(
			"expected invalid cluster target error, got %v",
			err,
		)
	}
}

func TestAdapterToDomainReturnsInvalidCreatedOnError(t *testing.T) {
	assertAdapterError(
		t,
		func(input *ClusterInput) {
			input.CreatedOn = time.Time{}
		},
		ErrInvalidClusterCreatedOn,
	)
}

func assertAdapterAcceptsConfidenceScore(
	t *testing.T,
	confidenceScore float64,
) {
	t.Helper()

	adapter := NewAdapter(clusterables.NewMockClusterableAdapter())

	result, err := adapter.ToDomain(
		validClusterInput(func(input *ClusterInput) {
			input.ConfidenceScore = confidenceScore
		}),
	)

	if err != nil {
		t.Fatal(err)
	}

	if result.ConfidenceScore() != confidenceScore {
		t.Fatalf(
			"expected confidence score %f, got %f",
			confidenceScore,
			result.ConfidenceScore(),
		)
	}
}

func assertAdapterError(
	t *testing.T,
	mutate func(input *ClusterInput),
	expected error,
) {
	t.Helper()

	adapter := NewAdapter(clusterables.NewMockClusterableAdapter())

	_, err := adapter.ToDomain(
		validClusterInput(mutate),
	)

	if !errors.Is(err, expected) {
		t.Fatalf("expected error %v, got %v", expected, err)
	}
}

func validClusterInput(
	mutate func(input *ClusterInput),
) ClusterInput {
	input := ClusterInput{
		Identifier: uuid.New(),
		Target: clusterables.ClusterableInput{
			Identifier:  uuid.New(),
			ClusterKind: clusterables.PostKind,
		},
		MemberIDs: []uuid.UUID{
			uuid.New(),
			uuid.New(),
		},
		MemberKind:      clusterables.PostKind,
		ConfidenceScore: 0.8,
		Centroid: []float32{
			0.1,
			0.2,
		},
		CreatedOn: time.Now(),
	}

	if mutate != nil {
		mutate(&input)
	}

	return input
}
