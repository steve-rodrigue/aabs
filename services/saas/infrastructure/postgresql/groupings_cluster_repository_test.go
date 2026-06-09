package postgresql

import (
	"context"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	domain_clusters "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters/clusterables"
)

func TestNewGroupingsClusterRepository(t *testing.T) {
	fixture := newGroupingsClusterRepositoryFixture(t)

	repository := NewGroupingsClusterRepository(
		fixture.pool,
		fixture.adapter,
	)

	if repository == nil {
		t.Fatalf("expected repository")
	}
}

func TestGroupingsClusterRepositorySaveAndFindByID(t *testing.T) {
	fixture := newGroupingsClusterRepositoryFixture(t)

	cluster := newTestCluster(
		t,
		clusterables.NewMockClusterable(clusterables.PostKind),
		clusterables.PostKind,
		[]uuid.UUID{uuid.New(), uuid.New()},
	)

	if err := fixture.repository.Save(fixture.ctx, cluster); err != nil {
		t.Fatal(err)
	}

	result, err := fixture.repository.FindByID(
		fixture.ctx,
		cluster.Identifier(),
	)
	if err != nil {
		t.Fatal(err)
	}

	assertCluster(t, result, cluster)
}

func TestGroupingsClusterRepositoryFindByIDReturnsNilWhenNotFound(t *testing.T) {
	fixture := newGroupingsClusterRepositoryFixture(t)

	result, err := fixture.repository.FindByID(
		fixture.ctx,
		uuid.New(),
	)

	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Fatalf("expected nil cluster")
	}
}

func TestGroupingsClusterRepositorySaveUpdatesExistingCluster(t *testing.T) {
	fixture := newGroupingsClusterRepositoryFixture(t)

	cluster := newTestCluster(
		t,
		clusterables.NewMockClusterable(clusterables.PostKind),
		clusterables.PostKind,
		[]uuid.UUID{uuid.New()},
	)

	if err := fixture.repository.Save(fixture.ctx, cluster); err != nil {
		t.Fatal(err)
	}

	updated := &testCluster{
		id: cluster.Identifier(),
		target: clusterables.NewMockClusterableWithID(
			cluster.Target().Identifier(),
			cluster.Target().ClusterKind(),
		),
		memberIDs: []uuid.UUID{
			uuid.New(),
			uuid.New(),
		},
		memberKind:      clusterables.PostKind,
		confidenceScore: 0.55,
		centroid:        []float32{0.5, 0.5},
		createdOn:       cluster.CreatedOn(),
	}

	if err := fixture.repository.Save(fixture.ctx, updated); err != nil {
		t.Fatal(err)
	}

	result, err := fixture.repository.FindByID(
		fixture.ctx,
		cluster.Identifier(),
	)
	if err != nil {
		t.Fatal(err)
	}

	assertCluster(t, result, updated)
}

func TestGroupingsClusterRepositoryFindByTarget(t *testing.T) {
	fixture := newGroupingsClusterRepositoryFixture(t)

	target := clusterables.NewMockClusterable(clusterables.PostKind)

	first := newTestCluster(
		t,
		target,
		clusterables.PostKind,
		[]uuid.UUID{uuid.New()},
	)

	second := newTestCluster(
		t,
		target,
		clusterables.PostKind,
		[]uuid.UUID{uuid.New()},
	)

	other := newTestCluster(
		t,
		clusterables.NewMockClusterable(clusterables.UserKind),
		clusterables.UserKind,
		[]uuid.UUID{uuid.New()},
	)

	saveClusters(t, fixture, first, second, other)

	result, err := fixture.repository.FindByTarget(
		fixture.ctx,
		target.Identifier(),
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 clusters, got %d", len(result))
	}

	for _, cluster := range result {
		if cluster.Target().Identifier() != target.Identifier() {
			t.Fatalf("expected target id %s, got %s", target.Identifier(), cluster.Target().Identifier())
		}
	}
}

func TestGroupingsClusterRepositoryFindByMember(t *testing.T) {
	fixture := newGroupingsClusterRepositoryFixture(t)

	member := uuid.New()

	first := newTestCluster(
		t,
		clusterables.NewMockClusterable(clusterables.PostKind),
		clusterables.PostKind,
		[]uuid.UUID{member, uuid.New()},
	)

	second := newTestCluster(
		t,
		clusterables.NewMockClusterable(clusterables.PostKind),
		clusterables.PostKind,
		[]uuid.UUID{member},
	)

	other := newTestCluster(
		t,
		clusterables.NewMockClusterable(clusterables.PostKind),
		clusterables.PostKind,
		[]uuid.UUID{uuid.New()},
	)

	saveClusters(t, fixture, first, second, other)

	result, err := fixture.repository.FindByMember(
		fixture.ctx,
		member,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 clusters, got %d", len(result))
	}

	for _, cluster := range result {
		if !containsUUID(cluster.MemberIDs(), member) {
			t.Fatalf("expected member %s in cluster", member)
		}
	}
}

func TestGroupingsClusterRepositoryFind(t *testing.T) {
	fixture := newGroupingsClusterRepositoryFixture(t)

	clusters := []domain_clusters.Cluster{
		newTestCluster(t, clusterables.NewMockClusterable(clusterables.PostKind), clusterables.PostKind, []uuid.UUID{uuid.New()}),
		newTestCluster(t, clusterables.NewMockClusterable(clusterables.PostKind), clusterables.PostKind, []uuid.UUID{uuid.New()}),
		newTestCluster(t, clusterables.NewMockClusterable(clusterables.PostKind), clusterables.PostKind, []uuid.UUID{uuid.New()}),
	}

	sortClusters(clusters)
	saveClusters(t, fixture, clusters...)

	result, err := fixture.repository.Find(
		fixture.ctx,
		1,
		1,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 cluster, got %d", len(result))
	}

	assertCluster(t, result[0], clusters[1])
}

func TestGroupingsClusterRepositoryFindAfterWithNilCursor(t *testing.T) {
	fixture := newGroupingsClusterRepositoryFixture(t)

	cluster := newTestCluster(
		t,
		clusterables.NewMockClusterable(clusterables.PostKind),
		clusterables.PostKind,
		[]uuid.UUID{uuid.New()},
	)

	if err := fixture.repository.Save(fixture.ctx, cluster); err != nil {
		t.Fatal(err)
	}

	result, err := fixture.repository.FindAfter(
		fixture.ctx,
		uuid.Nil,
		10,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 cluster, got %d", len(result))
	}

	assertCluster(t, result[0], cluster)
}

func TestGroupingsClusterRepositoryFindAfter(t *testing.T) {
	fixture := newGroupingsClusterRepositoryFixture(t)

	clusters := []domain_clusters.Cluster{
		newTestCluster(t, clusterables.NewMockClusterable(clusterables.PostKind), clusterables.PostKind, []uuid.UUID{uuid.New()}),
		newTestCluster(t, clusterables.NewMockClusterable(clusterables.PostKind), clusterables.PostKind, []uuid.UUID{uuid.New()}),
		newTestCluster(t, clusterables.NewMockClusterable(clusterables.PostKind), clusterables.PostKind, []uuid.UUID{uuid.New()}),
	}

	sortClusters(clusters)
	saveClusters(t, fixture, clusters...)

	result, err := fixture.repository.FindAfter(
		fixture.ctx,
		clusters[0].Identifier(),
		10,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 clusters, got %d", len(result))
	}

	assertCluster(t, result[0], clusters[1])
	assertCluster(t, result[1], clusters[2])
}

func TestGroupingsClusterRepositoryCount(t *testing.T) {
	fixture := newGroupingsClusterRepositoryFixture(t)

	saveClusters(
		t,
		fixture,
		newTestCluster(t, clusterables.NewMockClusterable(clusterables.PostKind), clusterables.PostKind, []uuid.UUID{uuid.New()}),
		newTestCluster(t, clusterables.NewMockClusterable(clusterables.PostKind), clusterables.PostKind, []uuid.UUID{uuid.New()}),
	)

	count, err := fixture.repository.Count(fixture.ctx)
	if err != nil {
		t.Fatal(err)
	}

	if count != 2 {
		t.Fatalf("expected count 2, got %d", count)
	}
}

type groupingsClusterRepositoryFixture struct {
	ctx        context.Context
	pool       *pgxpool.Pool
	adapter    domain_clusters.Adapter
	repository domain_clusters.Repository
}

func newGroupingsClusterRepositoryFixture(t *testing.T) *groupingsClusterRepositoryFixture {
	t.Helper()

	dsn := os.Getenv("POSTGRES_TEST_DSN")
	if dsn == "" {
		t.Skip("POSTGRES_TEST_DSN is not set")
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}

	adapter := domain_clusters.NewAdapter(
		clusterables.NewAdapter(),
	)

	fixture := &groupingsClusterRepositoryFixture{
		ctx:     ctx,
		pool:    pool,
		adapter: adapter,
	}

	createGroupingsClustersTable(t, fixture)
	truncateGroupingsClustersTable(t, fixture)

	fixture.repository = NewGroupingsClusterRepository(
		pool,
		adapter,
	)

	t.Cleanup(func() {
		truncateGroupingsClustersTable(t, fixture)
		pool.Close()
	})

	return fixture
}

func createGroupingsClustersTable(
	t *testing.T,
	fixture *groupingsClusterRepositoryFixture,
) {
	t.Helper()

	_, err := fixture.pool.Exec(
		fixture.ctx,
		`
		CREATE TABLE IF NOT EXISTS groupings_clusters (
			identifier UUID PRIMARY KEY,
			target_id UUID NOT NULL,
			target_kind TEXT NOT NULL,
			member_ids UUID[] NOT NULL,
			member_kind TEXT NOT NULL,
			confidence_score DOUBLE PRECISION NOT NULL,
			centroid REAL[] NOT NULL,
			created_on TIMESTAMPTZ NOT NULL
		)
		`,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func truncateGroupingsClustersTable(
	t *testing.T,
	fixture *groupingsClusterRepositoryFixture,
) {
	t.Helper()

	_, err := fixture.pool.Exec(
		fixture.ctx,
		`
		TRUNCATE TABLE groupings_clusters
		`,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func saveClusters(
	t *testing.T,
	fixture *groupingsClusterRepositoryFixture,
	clusters ...domain_clusters.Cluster,
) {
	t.Helper()

	for _, cluster := range clusters {
		if err := fixture.repository.Save(fixture.ctx, cluster); err != nil {
			t.Fatal(err)
		}
	}
}

type testCluster struct {
	id uuid.UUID

	target domain_clusters_target

	memberIDs  []uuid.UUID
	memberKind clusterables.Kind

	confidenceScore float64
	centroid        []float32

	createdOn time.Time
}

type domain_clusters_target interface {
	Identifier() uuid.UUID
	ClusterKind() clusterables.Kind
}

func newTestCluster(
	t *testing.T,
	target clusterables.Clusterable,
	memberKind clusterables.Kind,
	memberIDs []uuid.UUID,
) domain_clusters.Cluster {
	t.Helper()

	return &testCluster{
		id:              uuid.New(),
		target:          target,
		memberIDs:       memberIDs,
		memberKind:      memberKind,
		confidenceScore: 0.85,
		centroid:        []float32{0.1, 0.2},
		createdOn:       time.Now().UTC().Truncate(time.Microsecond),
	}
}

func (cluster *testCluster) Identifier() uuid.UUID {
	return cluster.id
}

func (cluster *testCluster) Target() clusterables.Clusterable {
	return cluster.target
}

func (cluster *testCluster) MemberIDs() []uuid.UUID {
	out := make([]uuid.UUID, len(cluster.memberIDs))
	copy(out, cluster.memberIDs)

	return out
}

func (cluster *testCluster) MemberKind() clusterables.Kind {
	return cluster.memberKind
}

func (cluster *testCluster) ConfidenceScore() float64 {
	return cluster.confidenceScore
}

func (cluster *testCluster) Centroid() []float32 {
	out := make([]float32, len(cluster.centroid))
	copy(out, cluster.centroid)

	return out
}

func (cluster *testCluster) CreatedOn() time.Time {
	return cluster.createdOn
}

func assertCluster(
	t *testing.T,
	result domain_clusters.Cluster,
	expected domain_clusters.Cluster,
) {
	t.Helper()

	if result == nil {
		t.Fatalf("expected cluster")
	}

	if result.Identifier() != expected.Identifier() {
		t.Fatalf("expected id %s, got %s", expected.Identifier(), result.Identifier())
	}

	if result.Target().Identifier() != expected.Target().Identifier() {
		t.Fatalf("expected target id %s, got %s", expected.Target().Identifier(), result.Target().Identifier())
	}

	if result.Target().ClusterKind() != expected.Target().ClusterKind() {
		t.Fatalf("expected target kind %s, got %s", expected.Target().ClusterKind(), result.Target().ClusterKind())
	}

	if result.MemberKind() != expected.MemberKind() {
		t.Fatalf("expected member kind %s, got %s", expected.MemberKind(), result.MemberKind())
	}

	assertUUIDSlice(t, result.MemberIDs(), expected.MemberIDs())
	assertFloat32Slice(t, result.Centroid(), expected.Centroid())

	if result.ConfidenceScore() != expected.ConfidenceScore() {
		t.Fatalf("expected confidence %f, got %f", expected.ConfidenceScore(), result.ConfidenceScore())
	}

	if !result.CreatedOn().Equal(expected.CreatedOn()) {
		t.Fatalf("expected created on %s, got %s", expected.CreatedOn(), result.CreatedOn())
	}
}

func assertUUIDSlice(
	t *testing.T,
	result []uuid.UUID,
	expected []uuid.UUID,
) {
	t.Helper()

	if len(result) != len(expected) {
		t.Fatalf("expected %d ids, got %d", len(expected), len(result))
	}

	for index := range expected {
		if result[index] != expected[index] {
			t.Fatalf("expected id[%d] %s, got %s", index, expected[index], result[index])
		}
	}
}

func assertFloat32Slice(
	t *testing.T,
	result []float32,
	expected []float32,
) {
	t.Helper()

	if len(result) != len(expected) {
		t.Fatalf("expected %d floats, got %d", len(expected), len(result))
	}

	for index := range expected {
		if result[index] != expected[index] {
			t.Fatalf("expected float[%d] %f, got %f", index, expected[index], result[index])
		}
	}
}

func containsUUID(
	values []uuid.UUID,
	id uuid.UUID,
) bool {
	for _, value := range values {
		if value == id {
			return true
		}
	}

	return false
}

func sortClusters(
	clusters []domain_clusters.Cluster,
) {
	sort.Slice(clusters, func(left int, right int) bool {
		return clusters[left].Identifier().String() <
			clusters[right].Identifier().String()
	})
}
