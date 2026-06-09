package clusterables

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestNewCandidateRepository(t *testing.T) {
	repository := NewCandidateRepository(
		NewMockComparableRepository(),
	)

	if repository == nil {
		t.Fatalf("expected repository")
	}
}

func TestCandidateRepositoryFindCandidates(t *testing.T) {
	ctx := context.Background()

	comparableRepository := NewMockComparableRepository()
	repository := NewCandidateRepository(comparableRepository)

	targetID := uuid.New()

	target := NewMockClusterableWithID(
		targetID,
		PostKind,
	)

	targetComparable := NewMockComparableWithID(
		targetID,
		PostKind,
		[]float32{1, 0, 0},
	)

	first := NewMockComparableWithID(
		uuid.New(),
		PostKind,
		[]float32{0.9, 0.1, 0},
	)

	second := NewMockComparableWithID(
		uuid.New(),
		PostKind,
		[]float32{0.8, 0.2, 0},
	)

	comparableRepository.FindByIDValue = targetComparable
	comparableRepository.FindNearestValue = []Comparable{
		targetComparable,
		first,
		second,
	}

	result, err := repository.FindCandidates(
		ctx,
		target,
		PostKind,
		2,
	)

	if err != nil {
		t.Fatal(err)
	}

	if comparableRepository.FindByIDCalls != 1 {
		t.Fatalf("expected 1 find by id call")
	}

	if comparableRepository.LastContext != ctx {
		t.Fatalf("expected context to be passed to find by id")
	}

	if comparableRepository.LastID != targetID {
		t.Fatalf("expected target id to be passed")
	}

	if comparableRepository.FindNearestCalls != 1 {
		t.Fatalf("expected 1 find nearest call")
	}

	if comparableRepository.LastTarget != targetComparable {
		t.Fatalf("expected target comparable to be passed")
	}

	if comparableRepository.LastKind != PostKind {
		t.Fatalf("expected kind %s, got %s", PostKind, comparableRepository.LastKind)
	}

	if comparableRepository.LastAmount != 3 {
		t.Fatalf("expected amount 3, got %d", comparableRepository.LastAmount)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 candidates, got %d", len(result))
	}

	if result[0] != first {
		t.Fatalf("expected first candidate")
	}

	if result[1] != second {
		t.Fatalf("expected second candidate")
	}
}

func TestCandidateRepositoryFindCandidatesDoesNotRemoveDifferentKindSameID(t *testing.T) {
	ctx := context.Background()

	comparableRepository := NewMockComparableRepository()
	repository := NewCandidateRepository(comparableRepository)

	targetID := uuid.New()

	target := NewMockClusterableWithID(
		targetID,
		PostKind,
	)

	targetComparable := NewMockComparableWithID(
		targetID,
		PostKind,
		[]float32{1, 0, 0},
	)

	differentKindSameID := NewMockComparableWithID(
		targetID,
		UserKind,
		[]float32{1, 0, 0},
	)

	comparableRepository.FindByIDValue = targetComparable
	comparableRepository.FindNearestValue = []Comparable{
		differentKindSameID,
	}

	result, err := repository.FindCandidates(
		ctx,
		target,
		UserKind,
		1,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 candidate, got %d", len(result))
	}

	if result[0] != differentKindSameID {
		t.Fatalf("expected different kind same id candidate")
	}
}

func TestCandidateRepositoryFindCandidatesLimitsResults(t *testing.T) {
	ctx := context.Background()

	comparableRepository := NewMockComparableRepository()
	repository := NewCandidateRepository(comparableRepository)

	target := NewMockClusterable(PostKind)
	targetComparable := NewMockComparableWithID(
		target.Identifier(),
		PostKind,
		[]float32{1, 0, 0},
	)

	first := NewMockComparable(PostKind, []float32{0.9, 0.1, 0})
	second := NewMockComparable(PostKind, []float32{0.8, 0.2, 0})

	comparableRepository.FindByIDValue = targetComparable
	comparableRepository.FindNearestValue = []Comparable{
		first,
		second,
	}

	result, err := repository.FindCandidates(
		ctx,
		target,
		PostKind,
		1,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 candidate, got %d", len(result))
	}

	if result[0] != first {
		t.Fatalf("expected first candidate")
	}
}

func TestCandidateRepositoryFindCandidatesReturnsEmptyWhenAmountIsZero(t *testing.T) {
	comparableRepository := NewMockComparableRepository()
	repository := NewCandidateRepository(comparableRepository)

	result, err := repository.FindCandidates(
		context.Background(),
		NewMockClusterable(PostKind),
		PostKind,
		0,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 0 {
		t.Fatalf("expected empty result")
	}

	if comparableRepository.FindByIDCalls != 0 {
		t.Fatalf("expected no find by id call")
	}

	if comparableRepository.FindNearestCalls != 0 {
		t.Fatalf("expected no find nearest call")
	}
}

func TestCandidateRepositoryFindCandidatesReturnsEmptyWhenAmountIsNegative(t *testing.T) {
	comparableRepository := NewMockComparableRepository()
	repository := NewCandidateRepository(comparableRepository)

	result, err := repository.FindCandidates(
		context.Background(),
		NewMockClusterable(PostKind),
		PostKind,
		-1,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 0 {
		t.Fatalf("expected empty result")
	}

	if comparableRepository.FindByIDCalls != 0 {
		t.Fatalf("expected no find by id call")
	}
}

func TestCandidateRepositoryFindCandidatesReturnsFindByIDError(t *testing.T) {
	ctx := context.Background()

	comparableRepository := NewMockComparableRepository()
	comparableRepository.FindByIDErr = errTest

	repository := NewCandidateRepository(comparableRepository)

	_, err := repository.FindCandidates(
		ctx,
		NewMockClusterable(PostKind),
		PostKind,
		1,
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected find by id error, got %v", err)
	}

	if comparableRepository.FindNearestCalls != 0 {
		t.Fatalf("expected no find nearest call")
	}
}

func TestCandidateRepositoryFindCandidatesReturnsEmptyWhenTargetComparableIsNil(t *testing.T) {
	ctx := context.Background()

	comparableRepository := NewMockComparableRepository()
	comparableRepository.FindByIDValue = nil

	repository := NewCandidateRepository(comparableRepository)

	result, err := repository.FindCandidates(
		ctx,
		NewMockClusterable(PostKind),
		PostKind,
		1,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 0 {
		t.Fatalf("expected empty result")
	}

	if comparableRepository.FindNearestCalls != 0 {
		t.Fatalf("expected no find nearest call")
	}
}

func TestCandidateRepositoryFindCandidatesReturnsFindNearestError(t *testing.T) {
	ctx := context.Background()

	comparableRepository := NewMockComparableRepository()
	comparableRepository.FindByIDValue = NewMockComparable(
		PostKind,
		[]float32{1, 0, 0},
	)
	comparableRepository.FindNearestErr = errTest

	repository := NewCandidateRepository(comparableRepository)

	_, err := repository.FindCandidates(
		ctx,
		NewMockClusterable(PostKind),
		PostKind,
		1,
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected find nearest error, got %v", err)
	}
}

func TestCandidateRepositoryFindCandidatesReturnsEmptyWhenNearestIsEmpty(t *testing.T) {
	ctx := context.Background()

	comparableRepository := NewMockComparableRepository()
	comparableRepository.FindByIDValue = NewMockComparable(
		PostKind,
		[]float32{1, 0, 0},
	)
	comparableRepository.FindNearestValue = []Comparable{}

	repository := NewCandidateRepository(comparableRepository)

	result, err := repository.FindCandidates(
		ctx,
		NewMockClusterable(PostKind),
		PostKind,
		1,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 0 {
		t.Fatalf("expected empty result")
	}
}
