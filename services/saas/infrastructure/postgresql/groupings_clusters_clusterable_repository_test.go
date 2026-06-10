package postgresql

import (
	"context"
	"os"
	"sort"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/clusterables"
)

func TestNewGroupingsClustersClusterablesRepository(t *testing.T) {
	fixture := newGroupingsClustersClusterableRepositoryFixture(t)

	repository := NewGroupingsClustersClusterablesRepository(
		fixture.pool,
		fixture.adapter,
	)

	if repository == nil {
		t.Fatalf("expected repository")
	}
}

func TestGroupingsClustersClusterableRepositoryFindByKind(t *testing.T) {
	fixture := newGroupingsClustersClusterableRepositoryFixture(t)

	first := uuid.New()
	second := uuid.New()

	insertGroupingsClusterable(
		t,
		fixture,
		first,
		clusterables.PostKind,
	)

	insertGroupingsClusterable(
		t,
		fixture,
		second,
		clusterables.PostKind,
	)

	insertGroupingsClusterable(
		t,
		fixture,
		uuid.New(),
		clusterables.UserKind,
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
		t.Fatalf("expected 2 clusterables, got %d", len(result))
	}

	assertClusterablesKind(t, result, clusterables.PostKind)
}

func TestGroupingsClustersClusterableRepositoryFindByKindUsesOffsetAndLimit(t *testing.T) {
	fixture := newGroupingsClustersClusterableRepositoryFixture(t)

	ids := []uuid.UUID{
		uuid.New(),
		uuid.New(),
		uuid.New(),
	}

	sort.Slice(ids, func(left int, right int) bool {
		return ids[left].String() < ids[right].String()
	})

	for _, id := range ids {
		insertGroupingsClusterable(
			t,
			fixture,
			id,
			clusterables.PostKind,
		)
	}

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
		t.Fatalf("expected 1 clusterable, got %d", len(result))
	}

	if result[0].Identifier() != ids[1] {
		t.Fatalf(
			"expected id %s, got %s",
			ids[1],
			result[0].Identifier(),
		)
	}
}

func TestGroupingsClustersClusterableRepositoryFindByKindAfterWithNilCursor(t *testing.T) {
	fixture := newGroupingsClustersClusterableRepositoryFixture(t)

	id := uuid.New()

	insertGroupingsClusterable(
		t,
		fixture,
		id,
		clusterables.PostKind,
	)

	result, err := fixture.repository.FindByKindAfter(
		fixture.ctx,
		clusterables.PostKind,
		uuid.Nil,
		10,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 clusterable, got %d", len(result))
	}

	if result[0].Identifier() != id {
		t.Fatalf(
			"expected id %s, got %s",
			id,
			result[0].Identifier(),
		)
	}
}

func TestGroupingsClustersClusterableRepositoryFindByKindAfter(t *testing.T) {
	fixture := newGroupingsClustersClusterableRepositoryFixture(t)

	ids := []uuid.UUID{
		uuid.New(),
		uuid.New(),
		uuid.New(),
	}

	sort.Slice(ids, func(left int, right int) bool {
		return ids[left].String() < ids[right].String()
	})

	for _, id := range ids {
		insertGroupingsClusterable(
			t,
			fixture,
			id,
			clusterables.PostKind,
		)
	}

	result, err := fixture.repository.FindByKindAfter(
		fixture.ctx,
		clusterables.PostKind,
		ids[0],
		10,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 clusterables, got %d", len(result))
	}

	if result[0].Identifier() != ids[1] {
		t.Fatalf(
			"expected first id %s, got %s",
			ids[1],
			result[0].Identifier(),
		)
	}

	if result[1].Identifier() != ids[2] {
		t.Fatalf(
			"expected second id %s, got %s",
			ids[2],
			result[1].Identifier(),
		)
	}
}

func TestGroupingsClustersClusterableRepositoryFindByKindAfterUsesLimit(t *testing.T) {
	fixture := newGroupingsClustersClusterableRepositoryFixture(t)

	ids := []uuid.UUID{
		uuid.New(),
		uuid.New(),
		uuid.New(),
	}

	sort.Slice(ids, func(left int, right int) bool {
		return ids[left].String() < ids[right].String()
	})

	for _, id := range ids {
		insertGroupingsClusterable(
			t,
			fixture,
			id,
			clusterables.PostKind,
		)
	}

	result, err := fixture.repository.FindByKindAfter(
		fixture.ctx,
		clusterables.PostKind,
		ids[0],
		1,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 clusterable, got %d", len(result))
	}

	if result[0].Identifier() != ids[1] {
		t.Fatalf(
			"expected id %s, got %s",
			ids[1],
			result[0].Identifier(),
		)
	}
}

func TestGroupingsClustersClusterableRepositoryCountByKind(t *testing.T) {
	fixture := newGroupingsClustersClusterableRepositoryFixture(t)

	insertGroupingsClusterable(
		t,
		fixture,
		uuid.New(),
		clusterables.PostKind,
	)

	insertGroupingsClusterable(
		t,
		fixture,
		uuid.New(),
		clusterables.PostKind,
	)

	insertGroupingsClusterable(
		t,
		fixture,
		uuid.New(),
		clusterables.UserKind,
	)

	count, err := fixture.repository.CountByKind(
		fixture.ctx,
		clusterables.PostKind,
	)

	if err != nil {
		t.Fatal(err)
	}

	if count != 2 {
		t.Fatalf("expected count 2, got %d", count)
	}
}

type groupingsClustersClusterableRepositoryFixture struct {
	ctx        context.Context
	pool       *pgxpool.Pool
	adapter    clusterables.Adapter
	repository clusterables.Repository
}

func newGroupingsClustersClusterableRepositoryFixture(
	t *testing.T,
) *groupingsClustersClusterableRepositoryFixture {
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

	adapter := clusterables.NewAdapter()

	fixture := &groupingsClustersClusterableRepositoryFixture{
		ctx:     ctx,
		pool:    pool,
		adapter: adapter,
	}

	createGroupingsClusterablesTable(t, fixture)
	truncateGroupingsClusterablesTable(t, fixture)

	fixture.repository = NewGroupingsClustersClusterablesRepository(
		pool,
		adapter,
	)

	t.Cleanup(func() {
		truncateGroupingsClusterablesTable(t, fixture)
		pool.Close()
	})

	return fixture
}

func createGroupingsClusterablesTable(
	t *testing.T,
	fixture *groupingsClustersClusterableRepositoryFixture,
) {
	t.Helper()

	_, err := fixture.pool.Exec(
		fixture.ctx,
		`
		CREATE TABLE IF NOT EXISTS groupings_clusterables (
			identifier UUID NOT NULL,
			kind TEXT NOT NULL,
			PRIMARY KEY (identifier, kind)
		)
		`,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func truncateGroupingsClusterablesTable(
	t *testing.T,
	fixture *groupingsClustersClusterableRepositoryFixture,
) {
	t.Helper()

	_, err := fixture.pool.Exec(
		fixture.ctx,
		`
		TRUNCATE TABLE groupings_clusterables
		`,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func insertGroupingsClusterable(
	t *testing.T,
	fixture *groupingsClustersClusterableRepositoryFixture,
	id uuid.UUID,
	kind clusterables.Kind,
) {
	t.Helper()

	_, err := fixture.pool.Exec(
		fixture.ctx,
		`
		INSERT INTO groupings_clusterables (
			identifier,
			kind
		)
		VALUES ($1, $2)
		`,
		id,
		string(kind),
	)
	if err != nil {
		t.Fatal(err)
	}
}

func assertClusterablesKind(
	t *testing.T,
	items []clusterables.Clusterable,
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
