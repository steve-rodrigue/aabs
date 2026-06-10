package clusters

import (
	"context"
	"errors"
	"math"
	"testing"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/clusterables"
)

func TestNewDetector(t *testing.T) {
	detector := NewDetector(
		NewMockClusterAdapter(),
		clusterables.NewMockComparableRepository(),
	)

	if detector == nil {
		t.Fatalf("expected detector")
	}
}

func TestDetectorDetect(t *testing.T) {
	ctx := context.Background()

	adapter := NewMockClusterAdapter()
	comparableRepository := clusterables.NewMockComparableRepository()

	target := clusterables.NewMockClusterable(clusterables.PostKind)

	first := clusterables.NewMockClusterable(clusterables.PostKind)
	second := clusterables.NewMockClusterable(clusterables.PostKind)

	firstComparable := clusterables.NewMockComparableWithID(
		first.Identifier(),
		clusterables.PostKind,
		[]float32{1, 0},
	)

	secondComparable := clusterables.NewMockComparableWithID(
		second.Identifier(),
		clusterables.PostKind,
		[]float32{0, 1},
	)

	comparableRepository.Items[first.Identifier()] = firstComparable
	comparableRepository.Items[second.Identifier()] = secondComparable

	detector := NewDetector(
		adapter,
		comparableRepository,
	)

	result, err := detector.Detect(
		ctx,
		target,
		[]clusterables.Clusterable{
			first,
			second,
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	if comparableRepository.FindByIDCalls != 2 {
		t.Fatalf(
			"expected 2 find by id calls, got %d",
			comparableRepository.FindByIDCalls,
		)
	}

	if comparableRepository.LastContext != ctx {
		t.Fatalf("expected context to be passed")
	}

	if adapter.ToDomainCalls != 1 {
		t.Fatalf(
			"expected 1 adapter call, got %d",
			adapter.ToDomainCalls,
		)
	}

	input := adapter.LastInput

	if input.Target.Identifier != target.Identifier() {
		t.Fatalf("expected target identifier")
	}

	if input.Target.ClusterKind != target.ClusterKind() {
		t.Fatalf("expected target kind")
	}

	if len(input.MemberIDs) != 2 {
		t.Fatalf("expected 2 member ids, got %d", len(input.MemberIDs))
	}

	if input.MemberIDs[0] != first.Identifier() {
		t.Fatalf("expected first member id")
	}

	if input.MemberIDs[1] != second.Identifier() {
		t.Fatalf("expected second member id")
	}

	if input.MemberKind != clusterables.PostKind {
		t.Fatalf("expected member kind %s, got %s", clusterables.PostKind, input.MemberKind)
	}

	assertFloat32Slice(
		t,
		input.Centroid,
		[]float32{
			0.5,
			0.5,
		},
	)

	expectedConfidence := (1 / math.Sqrt2) +
		(1 / math.Sqrt2)
	expectedConfidence = expectedConfidence / 2

	if math.Abs(input.ConfidenceScore-expectedConfidence) > 0.000001 {
		t.Fatalf(
			"expected confidence score %f, got %f",
			expectedConfidence,
			input.ConfidenceScore,
		)
	}

	if input.Identifier == uuid.Nil {
		t.Fatalf("expected generated cluster identifier")
	}

	if input.CreatedOn.IsZero() {
		t.Fatalf("expected created on")
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 cluster, got %d", len(result))
	}

	if result[0] != adapter.ToDomainValue && adapter.ToDomainValue != nil {
		t.Fatalf("expected adapter cluster result")
	}
}

func TestDetectorDetectReturnsEmptyWhenMembersAreEmpty(t *testing.T) {
	ctx := context.Background()

	adapter := NewMockClusterAdapter()
	comparableRepository := clusterables.NewMockComparableRepository()

	detector := NewDetector(
		adapter,
		comparableRepository,
	)

	result, err := detector.Detect(
		ctx,
		clusterables.NewMockClusterable(clusterables.PostKind),
		[]clusterables.Clusterable{},
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 0 {
		t.Fatalf("expected empty result")
	}

	if comparableRepository.FindByIDCalls != 0 {
		t.Fatalf("expected no comparable lookup")
	}

	if adapter.ToDomainCalls != 0 {
		t.Fatalf("expected no adapter call")
	}
}

func TestDetectorDetectReturnsInvalidTargetError(t *testing.T) {
	detector := NewDetector(
		NewMockClusterAdapter(),
		clusterables.NewMockComparableRepository(),
	)

	_, err := detector.Detect(
		context.Background(),
		nil,
		[]clusterables.Clusterable{
			clusterables.NewMockClusterable(clusterables.PostKind),
		},
	)

	if !errors.Is(err, ErrInvalidClusterDetectorTarget) {
		t.Fatalf("expected invalid target error, got %v", err)
	}
}

func TestDetectorDetectReturnsInvalidMemberError(t *testing.T) {
	detector := NewDetector(
		NewMockClusterAdapter(),
		clusterables.NewMockComparableRepository(),
	)

	_, err := detector.Detect(
		context.Background(),
		clusterables.NewMockClusterable(clusterables.PostKind),
		[]clusterables.Clusterable{
			nil,
		},
	)

	if !errors.Is(err, ErrInvalidClusterDetectorMember) {
		t.Fatalf("expected invalid member error, got %v", err)
	}
}

func TestDetectorDetectReturnsInvalidMemberKindError(t *testing.T) {
	first := clusterables.NewMockClusterable(clusterables.PostKind)
	second := clusterables.NewMockClusterable(clusterables.UserKind)

	comparableRepository := clusterables.NewMockComparableRepository()
	comparableRepository.Items[first.Identifier()] =
		clusterables.NewMockComparableWithID(
			first.Identifier(),
			clusterables.PostKind,
			[]float32{1, 0},
		)

	detector := NewDetector(
		NewMockClusterAdapter(),
		comparableRepository,
	)

	_, err := detector.Detect(
		context.Background(),
		clusterables.NewMockClusterable(clusterables.PostKind),
		[]clusterables.Clusterable{
			first,
			second,
		},
	)

	if !errors.Is(err, ErrInvalidClusterDetectorMemberKind) {
		t.Fatalf("expected invalid member kind error, got %v", err)
	}
}

func TestDetectorDetectReturnsComparableRepositoryError(t *testing.T) {
	comparableRepository := clusterables.NewMockComparableRepository()
	comparableRepository.FindByIDErr = errTest

	detector := NewDetector(
		NewMockClusterAdapter(),
		comparableRepository,
	)

	_, err := detector.Detect(
		context.Background(),
		clusterables.NewMockClusterable(clusterables.PostKind),
		[]clusterables.Clusterable{
			clusterables.NewMockClusterable(clusterables.PostKind),
		},
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected comparable repository error, got %v", err)
	}
}

func TestDetectorDetectReturnsInvalidComparableError(t *testing.T) {
	comparableRepository := clusterables.NewMockComparableRepository()
	comparableRepository.FindByIDValue = nil

	detector := NewDetector(
		NewMockClusterAdapter(),
		comparableRepository,
	)

	_, err := detector.Detect(
		context.Background(),
		clusterables.NewMockClusterable(clusterables.PostKind),
		[]clusterables.Clusterable{
			clusterables.NewMockClusterable(clusterables.PostKind),
		},
	)

	if !errors.Is(err, ErrInvalidClusterDetectorComparable) {
		t.Fatalf("expected invalid comparable error, got %v", err)
	}
}

func TestDetectorDetectReturnsInvalidVectorSizeErrorWhenVectorIsEmpty(t *testing.T) {
	member := clusterables.NewMockClusterable(clusterables.PostKind)

	comparableRepository := clusterables.NewMockComparableRepository()
	comparableRepository.Items[member.Identifier()] =
		clusterables.NewMockComparableWithID(
			member.Identifier(),
			clusterables.PostKind,
			[]float32{},
		)

	detector := NewDetector(
		NewMockClusterAdapter(),
		comparableRepository,
	)

	_, err := detector.Detect(
		context.Background(),
		clusterables.NewMockClusterable(clusterables.PostKind),
		[]clusterables.Clusterable{
			member,
		},
	)

	if !errors.Is(err, ErrInvalidClusterDetectorVectorSize) {
		t.Fatalf("expected invalid vector size error, got %v", err)
	}
}

func TestDetectorDetectReturnsInvalidVectorSizeErrorWhenVectorsMismatch(t *testing.T) {
	first := clusterables.NewMockClusterable(clusterables.PostKind)
	second := clusterables.NewMockClusterable(clusterables.PostKind)

	comparableRepository := clusterables.NewMockComparableRepository()
	comparableRepository.Items[first.Identifier()] =
		clusterables.NewMockComparableWithID(
			first.Identifier(),
			clusterables.PostKind,
			[]float32{1, 0},
		)

	comparableRepository.Items[second.Identifier()] =
		clusterables.NewMockComparableWithID(
			second.Identifier(),
			clusterables.PostKind,
			[]float32{1, 0, 0},
		)

	detector := NewDetector(
		NewMockClusterAdapter(),
		comparableRepository,
	)

	_, err := detector.Detect(
		context.Background(),
		clusterables.NewMockClusterable(clusterables.PostKind),
		[]clusterables.Clusterable{
			first,
			second,
		},
	)

	if !errors.Is(err, ErrInvalidClusterDetectorVectorSize) {
		t.Fatalf("expected invalid vector size error, got %v", err)
	}
}

func TestDetectorDetectReturnsAdapterError(t *testing.T) {
	member := clusterables.NewMockClusterable(clusterables.PostKind)

	comparableRepository := clusterables.NewMockComparableRepository()
	comparableRepository.Items[member.Identifier()] =
		clusterables.NewMockComparableWithID(
			member.Identifier(),
			clusterables.PostKind,
			[]float32{1, 0},
		)

	adapter := NewMockClusterAdapter()
	adapter.ToDomainErr = errTest

	detector := NewDetector(
		adapter,
		comparableRepository,
	)

	_, err := detector.Detect(
		context.Background(),
		clusterables.NewMockClusterable(clusterables.PostKind),
		[]clusterables.Clusterable{
			member,
		},
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected adapter error, got %v", err)
	}
}

func TestCalculateCentroid(t *testing.T) {
	centroid, err := calculateCentroid(
		[]clusterables.Comparable{
			clusterables.NewMockComparable(
				clusterables.PostKind,
				[]float32{1, 0},
			),
			clusterables.NewMockComparable(
				clusterables.PostKind,
				[]float32{0, 1},
			),
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	assertFloat32Slice(
		t,
		centroid,
		[]float32{
			0.5,
			0.5,
		},
	)
}

func TestCalculateCentroidReturnsInvalidMemberErrorWhenMembersAreEmpty(t *testing.T) {
	_, err := calculateCentroid(
		[]clusterables.Comparable{},
	)

	if !errors.Is(err, ErrInvalidClusterDetectorMember) {
		t.Fatalf("expected invalid member error, got %v", err)
	}
}

func TestCalculateConfidenceScore(t *testing.T) {
	score, err := calculateConfidenceScore(
		[]float32{
			0.5,
			0.5,
		},
		[]clusterables.Comparable{
			clusterables.NewMockComparable(
				clusterables.PostKind,
				[]float32{1, 0},
			),
			clusterables.NewMockComparable(
				clusterables.PostKind,
				[]float32{0, 1},
			),
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	expected := (1 / math.Sqrt2) +
		(1 / math.Sqrt2)
	expected = expected / 2

	if math.Abs(score-expected) > 0.000001 {
		t.Fatalf("expected score %f, got %f", expected, score)
	}
}

func TestCalculateConfidenceScoreReturnsInvalidVectorSizeErrorWhenCentroidIsEmpty(t *testing.T) {
	_, err := calculateConfidenceScore(
		[]float32{},
		[]clusterables.Comparable{
			clusterables.NewMockComparable(
				clusterables.PostKind,
				[]float32{1, 0},
			),
		},
	)

	if !errors.Is(err, ErrInvalidClusterDetectorVectorSize) {
		t.Fatalf("expected invalid vector size error, got %v", err)
	}
}

func TestCalculateConfidenceScoreReturnsInvalidVectorSizeErrorWhenMembersAreEmpty(t *testing.T) {
	_, err := calculateConfidenceScore(
		[]float32{
			1,
			0,
		},
		[]clusterables.Comparable{},
	)

	if !errors.Is(err, ErrInvalidClusterDetectorVectorSize) {
		t.Fatalf("expected invalid vector size error, got %v", err)
	}
}

func TestCosineSimilarity(t *testing.T) {
	result := cosineSimilarity(
		[]float32{
			1,
			0,
		},
		[]float32{
			1,
			0,
		},
	)

	if result != 1 {
		t.Fatalf("expected 1, got %f", result)
	}
}

func TestCosineSimilarityReturnsZeroWhenSourceMagnitudeIsZero(t *testing.T) {
	result := cosineSimilarity(
		[]float32{
			0,
			0,
		},
		[]float32{
			1,
			0,
		},
	)

	if result != 0 {
		t.Fatalf("expected 0, got %f", result)
	}
}

func TestCosineSimilarityReturnsZeroWhenTargetMagnitudeIsZero(t *testing.T) {
	result := cosineSimilarity(
		[]float32{
			1,
			0,
		},
		[]float32{
			0,
			0,
		},
	)

	if result != 0 {
		t.Fatalf("expected 0, got %f", result)
	}
}

func assertFloat32Slice(
	t *testing.T,
	result []float32,
	expected []float32,
) {
	t.Helper()

	if len(result) != len(expected) {
		t.Fatalf(
			"expected length %d, got %d",
			len(expected),
			len(result),
		)
	}

	for index := range expected {
		if math.Abs(float64(result[index]-expected[index])) > 0.000001 {
			t.Fatalf(
				"expected value[%d] %f, got %f",
				index,
				expected[index],
				result[index],
			)
		}
	}
}
