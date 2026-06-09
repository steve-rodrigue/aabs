package milvus

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"

	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters/clusterables"
)

const testVectorDimensions = 3

func TestNewGroupingsClustersClusterablesComparableRepository(t *testing.T) {
	fixture := newComparableRepositoryFixture(t)

	repository := NewGroupingsClustersClusterablesComparableRepository(
		fixture.client,
		fixture.adapter,
		fixture.collection,
	)

	if repository == nil {
		t.Fatalf("expected repository")
	}
}

func TestGroupingsClustersClusterablesComparableRepositoryFindByID(t *testing.T) {
	fixture := newComparableRepositoryFixture(t)

	id := uuid.New()

	insertComparable(
		t,
		fixture,
		id,
		clusterables.PostKind,
		[]float32{1, 0, 0},
	)

	result, err := fixture.repository.FindByID(
		fixture.ctx,
		id,
	)

	if err != nil {
		t.Fatal(err)
	}

	assertComparable(
		t,
		result,
		id,
		clusterables.PostKind,
		[]float32{1, 0, 0},
	)
}

func TestGroupingsClustersClusterablesComparableRepositoryFindByIDReturnsNilWhenNotFound(t *testing.T) {
	fixture := newComparableRepositoryFixture(t)

	result, err := fixture.repository.FindByID(
		fixture.ctx,
		uuid.New(),
	)

	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Fatalf("expected nil comparable")
	}
}

func TestGroupingsClustersClusterablesComparableRepositoryFindByKind(t *testing.T) {
	fixture := newComparableRepositoryFixture(t)

	insertComparable(
		t,
		fixture,
		uuid.New(),
		clusterables.PostKind,
		[]float32{1, 0, 0},
	)

	insertComparable(
		t,
		fixture,
		uuid.New(),
		clusterables.PostKind,
		[]float32{0.9, 0.1, 0},
	)

	insertComparable(
		t,
		fixture,
		uuid.New(),
		clusterables.UserKind,
		[]float32{0, 1, 0},
	)

	result, err := fixture.repository.FindByKind(
		fixture.ctx,
		clusterables.PostKind,
		0,
		10,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 comparables, got %d", len(result))
	}

	assertAllKind(t, result, clusterables.PostKind)
}

func TestGroupingsClustersClusterablesComparableRepositoryFindByKindWithOffsetAndLimit(t *testing.T) {
	fixture := newComparableRepositoryFixture(t)

	insertComparable(
		t,
		fixture,
		uuid.New(),
		clusterables.PostKind,
		[]float32{1, 0, 0},
	)

	insertComparable(
		t,
		fixture,
		uuid.New(),
		clusterables.PostKind,
		[]float32{0.9, 0.1, 0},
	)

	insertComparable(
		t,
		fixture,
		uuid.New(),
		clusterables.PostKind,
		[]float32{0.8, 0.2, 0},
	)

	result, err := fixture.repository.FindByKind(
		fixture.ctx,
		clusterables.PostKind,
		1,
		1,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 comparable, got %d", len(result))
	}

	assertAllKind(t, result, clusterables.PostKind)
}

func TestGroupingsClustersClusterablesComparableRepositoryFindNearest(t *testing.T) {
	fixture := newComparableRepositoryFixture(t)

	targetID := uuid.New()
	closeID := uuid.New()

	insertComparable(
		t,
		fixture,
		targetID,
		clusterables.PostKind,
		[]float32{1, 0, 0},
	)

	insertComparable(
		t,
		fixture,
		closeID,
		clusterables.PostKind,
		[]float32{0.95, 0.05, 0},
	)

	insertComparable(
		t,
		fixture,
		uuid.New(),
		clusterables.PostKind,
		[]float32{0, 1, 0},
	)

	insertComparable(
		t,
		fixture,
		uuid.New(),
		clusterables.UserKind,
		[]float32{1, 0, 0},
	)

	target := clusterables.NewMockComparableWithID(
		targetID,
		clusterables.PostKind,
		[]float32{1, 0, 0},
	)

	result, err := fixture.repository.FindNearest(
		fixture.ctx,
		target,
		clusterables.PostKind,
		3,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) == 0 {
		t.Fatalf("expected nearest comparables")
	}

	assertAllKind(t, result, clusterables.PostKind)

	if result[0].Identifier() != targetID &&
		result[0].Identifier() != closeID {
		t.Fatalf("expected closest result to be target or close vector")
	}
}

func TestGroupingsClustersClusterablesComparableRepositoryFindNearestFiltersByKind(t *testing.T) {
	fixture := newComparableRepositoryFixture(t)

	targetID := uuid.New()
	userID := uuid.New()

	insertComparable(
		t,
		fixture,
		targetID,
		clusterables.PostKind,
		[]float32{1, 0, 0},
	)

	insertComparable(
		t,
		fixture,
		userID,
		clusterables.UserKind,
		[]float32{1, 0, 0},
	)

	target := clusterables.NewMockComparableWithID(
		targetID,
		clusterables.PostKind,
		[]float32{1, 0, 0},
	)

	result, err := fixture.repository.FindNearest(
		fixture.ctx,
		target,
		clusterables.UserKind,
		10,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 comparable, got %d", len(result))
	}

	assertComparable(
		t,
		result[0],
		userID,
		clusterables.UserKind,
		[]float32{1, 0, 0},
	)
}

func TestGroupingsClustersClusterablesComparableRepositoryFindNearestReturnsEmptyWhenAmountIsZero(t *testing.T) {
	fixture := newComparableRepositoryFixture(t)

	target := clusterables.NewMockComparableWithID(
		uuid.New(),
		clusterables.PostKind,
		[]float32{1, 0, 0},
	)

	result, err := fixture.repository.FindNearest(
		fixture.ctx,
		target,
		clusterables.PostKind,
		0,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 0 {
		t.Fatalf("expected empty result, got %d", len(result))
	}
}

type comparableRepositoryFixture struct {
	ctx        context.Context
	client     client.Client
	adapter    clusterables.ComparableAdapter
	collection string
	repository clusterables.ComparableRepository
}

func newComparableRepositoryFixture(t *testing.T) *comparableRepositoryFixture {
	t.Helper()

	address := os.Getenv("MILVUS_TEST_ADDR")
	if address == "" {
		t.Skip("MILVUS_TEST_ADDR is not set")
	}

	ctx := context.Background()

	milvusClient, err := client.NewGrpcClient(ctx, address)
	if err != nil {
		t.Fatal(err)
	}

	collection := "test_clusterable_comparables_" +
		strings.ReplaceAll(uuid.New().String(), "-", "_")

	adapter := clusterables.NewComparableAdapter(
		clusterables.NewAdapter(),
	)

	fixture := &comparableRepositoryFixture{
		ctx:        ctx,
		client:     milvusClient,
		adapter:    adapter,
		collection: collection,
	}

	createComparableCollection(t, fixture)

	fixture.repository = NewGroupingsClustersClusterablesComparableRepository(
		milvusClient,
		adapter,
		collection,
	)

	t.Cleanup(func() {
		_ = milvusClient.DropCollection(ctx, collection)
		milvusClient.Close()
	})

	return fixture
}

func createComparableCollection(
	t *testing.T,
	fixture *comparableRepositoryFixture,
) {
	t.Helper()

	schema := entity.NewSchema().
		WithName(fixture.collection).
		WithDescription("test clusterable comparable vectors").
		WithField(
			entity.NewField().
				WithName(groupingsClustersClusterablesCandidatesComparableIdentifierField).
				WithDataType(entity.FieldTypeVarChar).
				WithMaxLength(36).
				WithIsPrimaryKey(true),
		).
		WithField(
			entity.NewField().
				WithName(groupingsClustersClusterablesCandidatesComparableKindField).
				WithDataType(entity.FieldTypeVarChar).
				WithMaxLength(64),
		).
		WithField(
			entity.NewField().
				WithName(groupingsClustersClusterablesCandidatesComparableVectorField).
				WithDataType(entity.FieldTypeFloatVector).
				WithDim(testVectorDimensions),
		)

	if err := fixture.client.CreateCollection(
		fixture.ctx,
		schema,
		2,
	); err != nil {
		t.Fatal(err)
	}

	index, err := entity.NewIndexFlat(entity.COSINE)
	if err != nil {
		t.Fatal(err)
	}

	if err := fixture.client.CreateIndex(
		fixture.ctx,
		fixture.collection,
		groupingsClustersClusterablesCandidatesComparableVectorField,
		index,
		false,
	); err != nil {
		t.Fatal(err)
	}

	if err := fixture.client.LoadCollection(
		fixture.ctx,
		fixture.collection,
		false,
	); err != nil {
		t.Fatal(err)
	}
}

func insertComparable(
	t *testing.T,
	fixture *comparableRepositoryFixture,
	id uuid.UUID,
	kind clusterables.Kind,
	vector []float32,
) {
	t.Helper()

	_, err := fixture.client.Insert(
		fixture.ctx,
		fixture.collection,
		"",
		entity.NewColumnVarChar(
			groupingsClustersClusterablesCandidatesComparableIdentifierField,
			[]string{id.String()},
		),
		entity.NewColumnVarChar(
			groupingsClustersClusterablesCandidatesComparableKindField,
			[]string{string(kind)},
		),
		entity.NewColumnFloatVector(
			groupingsClustersClusterablesCandidatesComparableVectorField,
			testVectorDimensions,
			[][]float32{vector},
		),
	)
	if err != nil {
		t.Fatal(err)
	}

	if err := fixture.client.Flush(
		fixture.ctx,
		fixture.collection,
		false,
	); err != nil {
		t.Fatal(err)
	}
}

func assertComparable(
	t *testing.T,
	result clusterables.Comparable,
	id uuid.UUID,
	kind clusterables.Kind,
	vector []float32,
) {
	t.Helper()

	if result == nil {
		t.Fatalf("expected comparable")
	}

	if result.Identifier() != id {
		t.Fatalf("expected id %s, got %s", id, result.Identifier())
	}

	if result.ClusterKind() != kind {
		t.Fatalf("expected kind %s, got %s", kind, result.ClusterKind())
	}

	resultVector := result.Vector()

	if len(resultVector) != len(vector) {
		t.Fatalf(
			"expected vector length %d, got %d",
			len(vector),
			len(resultVector),
		)
	}

	for index := range vector {
		if resultVector[index] != vector[index] {
			t.Fatalf(
				"expected vector[%d] %f, got %f",
				index,
				vector[index],
				resultVector[index],
			)
		}
	}
}

func assertAllKind(
	t *testing.T,
	items []clusterables.Comparable,
	kind clusterables.Kind,
) {
	t.Helper()

	for _, item := range items {
		if item.ClusterKind() != kind {
			t.Fatalf(
				"expected kind %s, got %s",
				kind,
				item.ClusterKind(),
			)
		}
	}
}
