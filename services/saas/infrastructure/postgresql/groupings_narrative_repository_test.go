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
	domain_narratives "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/narratives"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
)

func TestNewGroupingsNarrativeRepository(t *testing.T) {
	fixture := newGroupingsNarrativeRepositoryFixture(t)

	repository := NewGroupingsNarrativeRepository(
		fixture.pool,
		fixture.adapter,
		fixture.clusters,
	)

	if repository == nil {
		t.Fatalf("expected repository")
	}
}

func TestGroupingsNarrativeRepositorySaveAndFindByID(t *testing.T) {
	fixture := newGroupingsNarrativeRepositoryFixture(t)

	narrative := newTestNarrative(t, fixture, "Narrative A", "Description A")

	if err := fixture.repository.Save(fixture.ctx, narrative); err != nil {
		t.Fatal(err)
	}

	result, err := fixture.repository.FindByID(
		fixture.ctx,
		narrative.Identifier(),
	)
	if err != nil {
		t.Fatal(err)
	}

	assertNarrative(t, result, narrative)
}

func TestGroupingsNarrativeRepositoryFindByIDReturnsNilWhenNotFound(t *testing.T) {
	fixture := newGroupingsNarrativeRepositoryFixture(t)

	result, err := fixture.repository.FindByID(
		fixture.ctx,
		uuid.New(),
	)

	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Fatalf("expected nil narrative")
	}
}

func TestGroupingsNarrativeRepositorySaveUpdatesExistingNarrative(t *testing.T) {
	fixture := newGroupingsNarrativeRepositoryFixture(t)

	narrative := newTestNarrative(t, fixture, "Narrative A", "Description A")

	if err := fixture.repository.Save(fixture.ctx, narrative); err != nil {
		t.Fatal(err)
	}

	updated := &domain_narratives.MockNarrative{
		ID:                     narrative.Identifier(),
		ParticipationKindValue: participatables.NarrativeKind,
		ClusterValue:           narrative.Cluster(),
		NameValue:              "Narrative B",
		DescriptionValue:       "Description B",
		CreatedOnValue:         narrative.CreatedOn(),
	}

	if err := fixture.repository.Save(fixture.ctx, updated); err != nil {
		t.Fatal(err)
	}

	result, err := fixture.repository.FindByID(
		fixture.ctx,
		narrative.Identifier(),
	)
	if err != nil {
		t.Fatal(err)
	}

	assertNarrative(t, result, updated)
}

func TestGroupingsNarrativeRepositoryFindByName(t *testing.T) {
	fixture := newGroupingsNarrativeRepositoryFixture(t)

	expected := newTestNarrative(t, fixture, "Narrative A", "Description A")
	other := newTestNarrative(t, fixture, "Narrative B", "Description B")

	saveNarratives(t, fixture, expected, other)

	result, err := fixture.repository.FindByName(
		fixture.ctx,
		expected.Name(),
	)
	if err != nil {
		t.Fatal(err)
	}

	assertNarrative(t, result, expected)
}

func TestGroupingsNarrativeRepositoryFindByNameReturnsNilWhenNotFound(t *testing.T) {
	fixture := newGroupingsNarrativeRepositoryFixture(t)

	result, err := fixture.repository.FindByName(
		fixture.ctx,
		"missing",
	)

	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Fatalf("expected nil narrative")
	}
}

func TestGroupingsNarrativeRepositoryFind(t *testing.T) {
	fixture := newGroupingsNarrativeRepositoryFixture(t)

	first := newTestNarrative(t, fixture, "Narrative A", "Description A")
	second := newTestNarrative(t, fixture, "Narrative B", "Description B")
	third := newTestNarrative(t, fixture, "Narrative C", "Description C")

	saveNarratives(t, fixture, first, second, third)

	result, err := fixture.repository.Find(
		fixture.ctx,
		0,
		2,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 narratives, got %d", len(result))
	}
}

func TestGroupingsNarrativeRepositoryFindAfter(t *testing.T) {
	fixture := newGroupingsNarrativeRepositoryFixture(t)

	first := newTestNarrative(t, fixture, "Narrative A", "Description A")
	second := newTestNarrative(t, fixture, "Narrative B", "Description B")
	third := newTestNarrative(t, fixture, "Narrative C", "Description C")

	narratives := []domain_narratives.Narrative{
		first,
		second,
		third,
	}

	sortNarrativesByID(narratives)

	saveNarratives(t, fixture, narratives...)

	result, err := fixture.repository.FindAfter(
		fixture.ctx,
		narratives[0].Identifier(),
		10,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 narratives, got %d", len(result))
	}

	if result[0].Identifier() != narratives[1].Identifier() {
		t.Fatalf("expected second narrative first")
	}
}

func TestGroupingsNarrativeRepositoryFindAfterWithNilCursor(t *testing.T) {
	fixture := newGroupingsNarrativeRepositoryFixture(t)

	first := newTestNarrative(t, fixture, "Narrative A", "Description A")
	second := newTestNarrative(t, fixture, "Narrative B", "Description B")

	saveNarratives(t, fixture, first, second)

	result, err := fixture.repository.FindAfter(
		fixture.ctx,
		uuid.Nil,
		10,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 narratives, got %d", len(result))
	}
}

func TestGroupingsNarrativeRepositoryCount(t *testing.T) {
	fixture := newGroupingsNarrativeRepositoryFixture(t)

	first := newTestNarrative(t, fixture, "Narrative A", "Description A")
	second := newTestNarrative(t, fixture, "Narrative B", "Description B")

	saveNarratives(t, fixture, first, second)

	count, err := fixture.repository.Count(fixture.ctx)
	if err != nil {
		t.Fatal(err)
	}

	if count != 2 {
		t.Fatalf("expected count 2, got %d", count)
	}
}

type groupingsNarrativeRepositoryFixture struct {
	ctx context.Context

	pool *pgxpool.Pool

	adapter    domain_narratives.Adapter
	clusters   *domain_clusters.MockClusterRepository
	repository domain_narratives.Repository
}

func newGroupingsNarrativeRepositoryFixture(
	t *testing.T,
) *groupingsNarrativeRepositoryFixture {
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

	fixture := &groupingsNarrativeRepositoryFixture{
		ctx:      ctx,
		pool:     pool,
		adapter:  domain_narratives.NewAdapter(),
		clusters: domain_clusters.NewMockClusterRepository(),
	}

	createGroupingsNarrativeTable(t, fixture)
	truncateGroupingsNarrativeTable(t, fixture)

	fixture.repository = NewGroupingsNarrativeRepository(
		pool,
		fixture.adapter,
		fixture.clusters,
	)

	t.Cleanup(func() {
		truncateGroupingsNarrativeTable(t, fixture)
		pool.Close()
	})

	return fixture
}

func createGroupingsNarrativeTable(
	t *testing.T,
	fixture *groupingsNarrativeRepositoryFixture,
) {
	t.Helper()

	_, err := fixture.pool.Exec(
		fixture.ctx,
		`
		CREATE TABLE IF NOT EXISTS groupings_narratives (
			identifier UUID PRIMARY KEY,
			participation_kind TEXT NOT NULL,
			cluster_id UUID NOT NULL,
			name TEXT NOT NULL,
			description TEXT NOT NULL,
			created_on TIMESTAMPTZ NOT NULL
		)
		`,
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = fixture.pool.Exec(
		fixture.ctx,
		`
		CREATE INDEX IF NOT EXISTS groupings_narratives_name_idx
			ON groupings_narratives (name)
		`,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func truncateGroupingsNarrativeTable(
	t *testing.T,
	fixture *groupingsNarrativeRepositoryFixture,
) {
	t.Helper()

	_, err := fixture.pool.Exec(
		fixture.ctx,
		`
		TRUNCATE TABLE groupings_narratives
		`,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func saveNarratives(
	t *testing.T,
	fixture *groupingsNarrativeRepositoryFixture,
	narratives ...domain_narratives.Narrative,
) {
	t.Helper()

	for _, narrative := range narratives {
		if err := fixture.repository.Save(fixture.ctx, narrative); err != nil {
			t.Fatal(err)
		}
	}
}

func newTestNarrative(
	t *testing.T,
	fixture *groupingsNarrativeRepositoryFixture,
	name string,
	description string,
) domain_narratives.Narrative {
	t.Helper()

	cluster := domain_clusters.NewMockCluster(
		clusterables.NewMockClusterable(clusterables.NarrativeKind),
		clusterables.PostKind,
		[]uuid.UUID{uuid.New()},
	)

	fixture.clusters.Items[cluster.Identifier()] = cluster

	narrative, err := fixture.adapter.ToDomain(
		domain_narratives.NarrativeInput{
			Identifier:        uuid.New(),
			ParticipationKind: participatables.NarrativeKind,
			Cluster:           cluster,
			Name:              name,
			Description:       description,
			CreatedOn:         time.Now().UTC().Truncate(time.Microsecond),
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	return narrative
}

func assertNarrative(
	t *testing.T,
	result domain_narratives.Narrative,
	expected domain_narratives.Narrative,
) {
	t.Helper()

	if result == nil {
		t.Fatalf("expected narrative")
	}

	if result.Identifier() != expected.Identifier() {
		t.Fatalf(
			"expected id %s, got %s",
			expected.Identifier(),
			result.Identifier(),
		)
	}

	if result.ParticipationKind() != expected.ParticipationKind() {
		t.Fatalf(
			"expected participation kind %s, got %s",
			expected.ParticipationKind(),
			result.ParticipationKind(),
		)
	}

	if result.Cluster().Identifier() != expected.Cluster().Identifier() {
		t.Fatalf("expected cluster")
	}

	if result.Name() != expected.Name() {
		t.Fatalf(
			"expected name %s, got %s",
			expected.Name(),
			result.Name(),
		)
	}

	if result.Description() != expected.Description() {
		t.Fatalf(
			"expected description %s, got %s",
			expected.Description(),
			result.Description(),
		)
	}

	if !result.CreatedOn().Equal(expected.CreatedOn()) {
		t.Fatalf(
			"expected created on %s, got %s",
			expected.CreatedOn(),
			result.CreatedOn(),
		)
	}
}

func sortNarrativesByID(
	narratives []domain_narratives.Narrative,
) {
	sort.Slice(narratives, func(left int, right int) bool {
		return narratives[left].Identifier().String() <
			narratives[right].Identifier().String()
	})
}
