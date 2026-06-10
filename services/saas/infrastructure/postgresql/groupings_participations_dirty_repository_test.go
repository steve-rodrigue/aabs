package postgresql

import (
	"context"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/participatables"
	domain_dirty_participation "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/dirty"
)

func TestNewGroupingsParticipationsDirtyRepository(t *testing.T) {
	fixture := newGroupingsParticipationsDirtyRepositoryFixture(t)

	repository := NewGroupingsParticipationsDirtyRepository(
		fixture.pool,
		fixture.adapter,
	)

	if repository == nil {
		t.Fatalf("expected repository")
	}
}

func TestGroupingsParticipationsDirtyRepositorySaveAndFindByID(t *testing.T) {
	fixture := newGroupingsParticipationsDirtyRepositoryFixture(t)

	dirty := newTestDirtyParticipation(t, fixture)

	if err := fixture.repository.Save(fixture.ctx, dirty); err != nil {
		t.Fatal(err)
	}

	result, err := fixture.repository.FindByID(
		fixture.ctx,
		dirty.Identifier(),
	)
	if err != nil {
		t.Fatal(err)
	}

	assertDirtyParticipation(t, result, dirty)
}

func TestGroupingsParticipationsDirtyRepositoryFindByIDReturnsNilWhenNotFound(t *testing.T) {
	fixture := newGroupingsParticipationsDirtyRepositoryFixture(t)

	result, err := fixture.repository.FindByID(
		fixture.ctx,
		uuid.New(),
	)
	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Fatalf("expected nil dirty participation")
	}
}

func TestGroupingsParticipationsDirtyRepositorySaveUpdatesExistingBetweenPair(t *testing.T) {
	fixture := newGroupingsParticipationsDirtyRepositoryFixture(t)

	dirty := newTestDirtyParticipation(t, fixture)

	if err := fixture.repository.Save(fixture.ctx, dirty); err != nil {
		t.Fatal(err)
	}

	updated := &domain_dirty_participation.MockDirty{
		ID:               uuid.New(),
		ParticipantValue: dirty.Participant(),
		TargetValue:      dirty.Target(),
		MarkedOnValue:    time.Now().UTC().Add(time.Hour).Truncate(time.Microsecond),
	}

	if err := fixture.repository.Save(fixture.ctx, updated); err != nil {
		t.Fatal(err)
	}

	count, err := fixture.repository.Count(fixture.ctx)
	if err != nil {
		t.Fatal(err)
	}

	if count != 1 {
		t.Fatalf("expected count 1, got %d", count)
	}

	result, err := fixture.repository.FindBetween(
		fixture.ctx,
		dirty.Participant(),
		dirty.Target(),
	)
	if err != nil {
		t.Fatal(err)
	}

	if result.Identifier() != dirty.Identifier() {
		t.Fatalf("expected original identifier to remain on upsert")
	}

	if !result.MarkedOn().Equal(updated.MarkedOn()) {
		t.Fatalf(
			"expected marked on %s, got %s",
			updated.MarkedOn(),
			result.MarkedOn(),
		)
	}
}

func TestGroupingsParticipationsDirtyRepositoryDelete(t *testing.T) {
	fixture := newGroupingsParticipationsDirtyRepositoryFixture(t)

	dirty := newTestDirtyParticipation(t, fixture)

	if err := fixture.repository.Save(fixture.ctx, dirty); err != nil {
		t.Fatal(err)
	}

	if err := fixture.repository.Delete(fixture.ctx, dirty.Identifier()); err != nil {
		t.Fatal(err)
	}

	result, err := fixture.repository.FindByID(
		fixture.ctx,
		dirty.Identifier(),
	)
	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Fatalf("expected nil dirty participation")
	}
}

func TestGroupingsParticipationsDirtyRepositoryFindBetween(t *testing.T) {
	fixture := newGroupingsParticipationsDirtyRepositoryFixture(t)

	expected := newTestDirtyParticipation(t, fixture)
	other := newTestDirtyParticipation(t, fixture)

	saveDirtyParticipations(t, fixture, expected, other)

	result, err := fixture.repository.FindBetween(
		fixture.ctx,
		expected.Participant(),
		expected.Target(),
	)
	if err != nil {
		t.Fatal(err)
	}

	assertDirtyParticipation(t, result, expected)
}

func TestGroupingsParticipationsDirtyRepositoryFindBetweenReturnsNilWhenNotFound(t *testing.T) {
	fixture := newGroupingsParticipationsDirtyRepositoryFixture(t)

	participant := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.UserKind,
	)

	target := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.CampaignKind,
	)

	result, err := fixture.repository.FindBetween(
		fixture.ctx,
		participant,
		target,
	)
	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Fatalf("expected nil dirty participation")
	}
}

func TestGroupingsParticipationsDirtyRepositoryFind(t *testing.T) {
	fixture := newGroupingsParticipationsDirtyRepositoryFixture(t)

	first := newTestDirtyParticipation(t, fixture)
	second := newTestDirtyParticipation(t, fixture)
	third := newTestDirtyParticipation(t, fixture)

	saveDirtyParticipations(t, fixture, first, second, third)

	result, err := fixture.repository.Find(
		fixture.ctx,
		0,
		2,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 dirty participations, got %d", len(result))
	}
}

func TestGroupingsParticipationsDirtyRepositoryFindAfter(t *testing.T) {
	fixture := newGroupingsParticipationsDirtyRepositoryFixture(t)

	first := newTestDirtyParticipation(t, fixture)
	second := newTestDirtyParticipation(t, fixture)
	third := newTestDirtyParticipation(t, fixture)

	items := []domain_dirty_participation.Dirty{
		first,
		second,
		third,
	}

	sortDirtyParticipationsByID(items)

	saveDirtyParticipations(t, fixture, items...)

	result, err := fixture.repository.FindAfter(
		fixture.ctx,
		items[0].Identifier(),
		10,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 dirty participations, got %d", len(result))
	}

	if result[0].Identifier() != items[1].Identifier() {
		t.Fatalf("expected second dirty participation first")
	}
}

func TestGroupingsParticipationsDirtyRepositoryFindAfterWithNilCursor(t *testing.T) {
	fixture := newGroupingsParticipationsDirtyRepositoryFixture(t)

	first := newTestDirtyParticipation(t, fixture)
	second := newTestDirtyParticipation(t, fixture)

	saveDirtyParticipations(t, fixture, first, second)

	result, err := fixture.repository.FindAfter(
		fixture.ctx,
		uuid.Nil,
		10,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 dirty participations, got %d", len(result))
	}
}

func TestGroupingsParticipationsDirtyRepositoryCount(t *testing.T) {
	fixture := newGroupingsParticipationsDirtyRepositoryFixture(t)

	first := newTestDirtyParticipation(t, fixture)
	second := newTestDirtyParticipation(t, fixture)

	saveDirtyParticipations(t, fixture, first, second)

	count, err := fixture.repository.Count(fixture.ctx)
	if err != nil {
		t.Fatal(err)
	}

	if count != 2 {
		t.Fatalf("expected count 2, got %d", count)
	}
}

type groupingsParticipationsDirtyRepositoryFixture struct {
	ctx context.Context

	pool *pgxpool.Pool

	adapter    domain_dirty_participation.Adapter
	repository domain_dirty_participation.Repository
}

func newGroupingsParticipationsDirtyRepositoryFixture(
	t *testing.T,
) *groupingsParticipationsDirtyRepositoryFixture {
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

	fixture := &groupingsParticipationsDirtyRepositoryFixture{
		ctx:     ctx,
		pool:    pool,
		adapter: domain_dirty_participation.NewAdapter(),
	}

	createGroupingsParticipationsDirtyTable(t, fixture)
	truncateGroupingsParticipationsDirtyTable(t, fixture)

	fixture.repository = NewGroupingsParticipationsDirtyRepository(
		pool,
		fixture.adapter,
	)

	t.Cleanup(func() {
		truncateGroupingsParticipationsDirtyTable(t, fixture)
		pool.Close()
	})

	return fixture
}

func createGroupingsParticipationsDirtyTable(
	t *testing.T,
	fixture *groupingsParticipationsDirtyRepositoryFixture,
) {
	t.Helper()

	_, err := fixture.pool.Exec(
		fixture.ctx,
		`
		CREATE TABLE IF NOT EXISTS groupings_participations_dirty (
			identifier UUID PRIMARY KEY,
			participant_id UUID NOT NULL,
			participant_kind TEXT NOT NULL,
			target_id UUID NOT NULL,
			target_kind TEXT NOT NULL,
			marked_on TIMESTAMPTZ NOT NULL,
			UNIQUE (
				participant_id,
				participant_kind,
				target_id,
				target_kind
			)
		)
		`,
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = fixture.pool.Exec(
		fixture.ctx,
		`
		CREATE INDEX IF NOT EXISTS groupings_participations_dirty_between_idx
			ON groupings_participations_dirty (
				participant_id,
				participant_kind,
				target_id,
				target_kind
			)
		`,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func truncateGroupingsParticipationsDirtyTable(
	t *testing.T,
	fixture *groupingsParticipationsDirtyRepositoryFixture,
) {
	t.Helper()

	_, err := fixture.pool.Exec(
		fixture.ctx,
		`
		TRUNCATE TABLE groupings_participations_dirty
		`,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func saveDirtyParticipations(
	t *testing.T,
	fixture *groupingsParticipationsDirtyRepositoryFixture,
	items ...domain_dirty_participation.Dirty,
) {
	t.Helper()

	for _, item := range items {
		if err := fixture.repository.Save(fixture.ctx, item); err != nil {
			t.Fatal(err)
		}
	}
}

func newTestDirtyParticipation(
	t *testing.T,
	fixture *groupingsParticipationsDirtyRepositoryFixture,
) domain_dirty_participation.Dirty {
	t.Helper()

	participant := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.UserKind,
	)

	target := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.CampaignKind,
	)

	dirty, err := fixture.adapter.ToDomain(
		domain_dirty_participation.DirtyInput{
			Identifier:  uuid.New(),
			Participant: participant,
			Target:      target,
			MarkedOn:    time.Now().UTC().Truncate(time.Microsecond),
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	return dirty
}

func assertDirtyParticipation(
	t *testing.T,
	result domain_dirty_participation.Dirty,
	expected domain_dirty_participation.Dirty,
) {
	t.Helper()

	if result == nil {
		t.Fatalf("expected dirty participation")
	}

	if result.Identifier() != expected.Identifier() {
		t.Fatalf(
			"expected id %s, got %s",
			expected.Identifier(),
			result.Identifier(),
		)
	}

	if result.Participant().Identifier() != expected.Participant().Identifier() {
		t.Fatalf("expected participant id")
	}

	if result.Participant().ParticipationKind() != expected.Participant().ParticipationKind() {
		t.Fatalf("expected participant kind")
	}

	if result.Target().Identifier() != expected.Target().Identifier() {
		t.Fatalf("expected target id")
	}

	if result.Target().ParticipationKind() != expected.Target().ParticipationKind() {
		t.Fatalf("expected target kind")
	}

	if !result.MarkedOn().Equal(expected.MarkedOn()) {
		t.Fatalf(
			"expected marked on %s, got %s",
			expected.MarkedOn(),
			result.MarkedOn(),
		)
	}
}

func sortDirtyParticipationsByID(
	items []domain_dirty_participation.Dirty,
) {
	sort.Slice(items, func(left int, right int) bool {
		return items[left].Identifier().String() <
			items[right].Identifier().String()
	})
}
