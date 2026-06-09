package postgresql

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	domain_relationships "github.com/steve-rodrigue/aabs/services/saas/domain/relationships"
	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships/relatables"
)

func TestNewRelationshipRepository(t *testing.T) {
	pool := &pgxpool.Pool{}

	repository := NewRelationshipRepository(
		pool,
		domain_relationships.NewMockRelationshipAdapter(),
		relatables.NewMockRelatableAdapter(),
	)

	if repository == nil {
		t.Fatalf("expected repository")
	}
}

func TestRelationshipRepositorySaveAndFindByID(t *testing.T) {
	fixture := newRelationshipRepositoryFixture(t)

	relationship := newTestRelationship(
		t,
		mustParseRelationshipUUID("00000000-0000-0000-0000-000000000001"),
		relatables.UserKind,
		relatables.PostKind,
		0.75,
	)

	if err := fixture.repository.Save(fixture.ctx, relationship); err != nil {
		t.Fatal(err)
	}

	result, err := fixture.repository.FindByID(
		fixture.ctx,
		relationship.Identifier(),
	)

	if err != nil {
		t.Fatal(err)
	}

	assertRelationship(t, result, relationship)
}

func TestRelationshipRepositorySaveUpdatesExistingRelationship(t *testing.T) {
	fixture := newRelationshipRepositoryFixture(t)

	id := mustParseRelationshipUUID("00000000-0000-0000-0000-000000000001")

	first := newTestRelationship(
		t,
		id,
		relatables.UserKind,
		relatables.PostKind,
		0.25,
	)

	second := newTestRelationship(
		t,
		id,
		relatables.TopicKind,
		relatables.NarrativeKind,
		0.95,
	)

	if err := fixture.repository.Save(fixture.ctx, first); err != nil {
		t.Fatal(err)
	}

	if err := fixture.repository.Save(fixture.ctx, second); err != nil {
		t.Fatal(err)
	}

	result, err := fixture.repository.FindByID(fixture.ctx, id)

	if err != nil {
		t.Fatal(err)
	}

	assertRelationship(t, result, second)
}

func TestRelationshipRepositoryFindByIDReturnsNilWhenNotFound(t *testing.T) {
	fixture := newRelationshipRepositoryFixture(t)

	result, err := fixture.repository.FindByID(
		fixture.ctx,
		uuid.New(),
	)

	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Fatalf("expected nil relationship")
	}
}

func TestRelationshipRepositoryFind(t *testing.T) {
	fixture := newRelationshipRepositoryFixture(t)

	first := newTestRelationship(
		t,
		mustParseRelationshipUUID("00000000-0000-0000-0000-000000000001"),
		relatables.UserKind,
		relatables.PostKind,
		0.1,
	)

	second := newTestRelationship(
		t,
		mustParseRelationshipUUID("00000000-0000-0000-0000-000000000002"),
		relatables.UserKind,
		relatables.TopicKind,
		0.2,
	)

	third := newTestRelationship(
		t,
		mustParseRelationshipUUID("00000000-0000-0000-0000-000000000003"),
		relatables.PostKind,
		relatables.TopicKind,
		0.3,
	)

	saveRelationships(t, fixture, first, second, third)

	result, err := fixture.repository.Find(fixture.ctx, 1, 2)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 relationships, got %d", len(result))
	}

	assertRelationship(t, result[0], second)
	assertRelationship(t, result[1], third)
}

func TestRelationshipRepositoryFindAfterWithoutCursor(t *testing.T) {
	fixture := newRelationshipRepositoryFixture(t)

	first := newTestRelationship(
		t,
		mustParseRelationshipUUID("00000000-0000-0000-0000-000000000001"),
		relatables.UserKind,
		relatables.PostKind,
		0.1,
	)

	second := newTestRelationship(
		t,
		mustParseRelationshipUUID("00000000-0000-0000-0000-000000000002"),
		relatables.UserKind,
		relatables.TopicKind,
		0.2,
	)

	saveRelationships(t, fixture, first, second)

	result, err := fixture.repository.FindAfter(
		fixture.ctx,
		uuid.Nil,
		1,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 relationship, got %d", len(result))
	}

	assertRelationship(t, result[0], first)
}

func TestRelationshipRepositoryFindAfterWithCursor(t *testing.T) {
	fixture := newRelationshipRepositoryFixture(t)

	first := newTestRelationship(
		t,
		mustParseRelationshipUUID("00000000-0000-0000-0000-000000000001"),
		relatables.UserKind,
		relatables.PostKind,
		0.1,
	)

	second := newTestRelationship(
		t,
		mustParseRelationshipUUID("00000000-0000-0000-0000-000000000002"),
		relatables.UserKind,
		relatables.TopicKind,
		0.2,
	)

	third := newTestRelationship(
		t,
		mustParseRelationshipUUID("00000000-0000-0000-0000-000000000003"),
		relatables.PostKind,
		relatables.TopicKind,
		0.3,
	)

	saveRelationships(t, fixture, first, second, third)

	result, err := fixture.repository.FindAfter(
		fixture.ctx,
		first.Identifier(),
		2,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 relationships, got %d", len(result))
	}

	assertRelationship(t, result[0], second)
	assertRelationship(t, result[1], third)
}

func TestRelationshipRepositoryCount(t *testing.T) {
	fixture := newRelationshipRepositoryFixture(t)

	saveRelationships(
		t,
		fixture,
		newTestRelationship(
			t,
			mustParseRelationshipUUID("00000000-0000-0000-0000-000000000001"),
			relatables.UserKind,
			relatables.PostKind,
			0.1,
		),
		newTestRelationship(
			t,
			mustParseRelationshipUUID("00000000-0000-0000-0000-000000000002"),
			relatables.UserKind,
			relatables.TopicKind,
			0.2,
		),
	)

	result, err := fixture.repository.Count(fixture.ctx)

	if err != nil {
		t.Fatal(err)
	}

	if result != 2 {
		t.Fatalf("expected count 2, got %d", result)
	}
}

func TestRelationshipRepositoryFindBySourceID(t *testing.T) {
	fixture := newRelationshipRepositoryFixture(t)

	sourceID := uuid.New()

	first := newTestRelationshipWithRelatables(
		t,
		mustParseRelationshipUUID("00000000-0000-0000-0000-000000000001"),
		relatables.NewMockRelatable(sourceID, relatables.UserKind),
		relatables.NewMockRelatable(uuid.New(), relatables.PostKind),
		0.5,
	)

	second := newTestRelationship(
		t,
		mustParseRelationshipUUID("00000000-0000-0000-0000-000000000002"),
		relatables.TopicKind,
		relatables.PostKind,
		0.6,
	)

	saveRelationships(t, fixture, first, second)

	result, err := fixture.repository.FindBySourceID(fixture.ctx, sourceID)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 relationship, got %d", len(result))
	}

	assertRelationship(t, result[0], first)
}

func TestRelationshipRepositoryFindByTargetID(t *testing.T) {
	fixture := newRelationshipRepositoryFixture(t)

	targetID := uuid.New()

	first := newTestRelationshipWithRelatables(
		t,
		mustParseRelationshipUUID("00000000-0000-0000-0000-000000000001"),
		relatables.NewMockRelatable(uuid.New(), relatables.UserKind),
		relatables.NewMockRelatable(targetID, relatables.PostKind),
		0.5,
	)

	second := newTestRelationship(
		t,
		mustParseRelationshipUUID("00000000-0000-0000-0000-000000000002"),
		relatables.TopicKind,
		relatables.PostKind,
		0.6,
	)

	saveRelationships(t, fixture, first, second)

	result, err := fixture.repository.FindByTargetID(fixture.ctx, targetID)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 relationship, got %d", len(result))
	}

	assertRelationship(t, result[0], first)
}

func TestRelationshipRepositoryFindBySource(t *testing.T) {
	fixture := newRelationshipRepositoryFixture(t)

	source := relatables.NewMockRelatable(uuid.New(), relatables.UserKind)

	first := newTestRelationshipWithRelatables(
		t,
		mustParseRelationshipUUID("00000000-0000-0000-0000-000000000001"),
		source,
		relatables.NewMockRelatable(uuid.New(), relatables.PostKind),
		0.5,
	)

	wrongKindSameID := newTestRelationshipWithRelatables(
		t,
		mustParseRelationshipUUID("00000000-0000-0000-0000-000000000002"),
		relatables.NewMockRelatable(source.Identifier(), relatables.PostKind),
		relatables.NewMockRelatable(uuid.New(), relatables.TopicKind),
		0.6,
	)

	saveRelationships(t, fixture, first, wrongKindSameID)

	result, err := fixture.repository.FindBySource(fixture.ctx, source)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 relationship, got %d", len(result))
	}

	assertRelationship(t, result[0], first)
}

func TestRelationshipRepositoryFindByTarget(t *testing.T) {
	fixture := newRelationshipRepositoryFixture(t)

	target := relatables.NewMockRelatable(uuid.New(), relatables.PostKind)

	first := newTestRelationshipWithRelatables(
		t,
		mustParseRelationshipUUID("00000000-0000-0000-0000-000000000001"),
		relatables.NewMockRelatable(uuid.New(), relatables.UserKind),
		target,
		0.5,
	)

	wrongKindSameID := newTestRelationshipWithRelatables(
		t,
		mustParseRelationshipUUID("00000000-0000-0000-0000-000000000002"),
		relatables.NewMockRelatable(uuid.New(), relatables.UserKind),
		relatables.NewMockRelatable(target.Identifier(), relatables.TopicKind),
		0.6,
	)

	saveRelationships(t, fixture, first, wrongKindSameID)

	result, err := fixture.repository.FindByTarget(fixture.ctx, target)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 relationship, got %d", len(result))
	}

	assertRelationship(t, result[0], first)
}

func TestRelationshipRepositoryFindBetween(t *testing.T) {
	fixture := newRelationshipRepositoryFixture(t)

	source := relatables.NewMockRelatable(uuid.New(), relatables.UserKind)
	target := relatables.NewMockRelatable(uuid.New(), relatables.PostKind)

	relationship := newTestRelationshipWithRelatables(
		t,
		mustParseRelationshipUUID("00000000-0000-0000-0000-000000000001"),
		source,
		target,
		0.9,
	)

	saveRelationships(t, fixture, relationship)

	result, err := fixture.repository.FindBetween(
		fixture.ctx,
		source,
		target,
	)

	if err != nil {
		t.Fatal(err)
	}

	assertRelationship(t, result, relationship)
}

func TestRelationshipRepositoryFindBetweenReturnsNilWhenNotFound(t *testing.T) {
	fixture := newRelationshipRepositoryFixture(t)

	result, err := fixture.repository.FindBetween(
		fixture.ctx,
		relatables.NewMockRelatable(uuid.New(), relatables.UserKind),
		relatables.NewMockRelatable(uuid.New(), relatables.PostKind),
	)

	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Fatalf("expected nil relationship")
	}
}

type relationshipRepositoryFixture struct {
	ctx        context.Context
	pool       *pgxpool.Pool
	repository domain_relationships.Repository
}

func newRelationshipRepositoryFixture(t *testing.T) *relationshipRepositoryFixture {
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

	fixture := &relationshipRepositoryFixture{
		ctx:  ctx,
		pool: pool,
		repository: NewRelationshipRepository(
			pool,
			domain_relationships.NewAdapter(),
			relatables.NewAdapter(),
		),
	}

	createRelationshipsTable(t, fixture)
	truncateRelationshipsTable(t, fixture)

	t.Cleanup(func() {
		truncateRelationshipsTable(t, fixture)
		pool.Close()
	})

	return fixture
}

func createRelationshipsTable(
	t *testing.T,
	fixture *relationshipRepositoryFixture,
) {
	t.Helper()

	_, err := fixture.pool.Exec(
		fixture.ctx,
		`
		CREATE TABLE IF NOT EXISTS relationships (
			identifier UUID PRIMARY KEY,
			source_id UUID NOT NULL,
			source_kind TEXT NOT NULL,
			target_id UUID NOT NULL,
			target_kind TEXT NOT NULL,
			similarity DOUBLE PRECISION NOT NULL,
			created_on TIMESTAMPTZ NOT NULL,

			UNIQUE (source_id, source_kind, target_id, target_kind)
		);

		CREATE INDEX IF NOT EXISTS relationships_source_idx
		ON relationships (source_id, source_kind);

		CREATE INDEX IF NOT EXISTS relationships_target_idx
		ON relationships (target_id, target_kind);

		CREATE INDEX IF NOT EXISTS relationships_identifier_idx
		ON relationships (identifier);
		`,
	)

	if err != nil {
		t.Fatal(err)
	}
}

func truncateRelationshipsTable(
	t *testing.T,
	fixture *relationshipRepositoryFixture,
) {
	t.Helper()

	_, err := fixture.pool.Exec(
		fixture.ctx,
		`
		TRUNCATE TABLE relationships
		`,
	)

	if err != nil {
		t.Fatal(err)
	}
}

func saveRelationships(
	t *testing.T,
	fixture *relationshipRepositoryFixture,
	relationships ...domain_relationships.Relationship,
) {
	t.Helper()

	for _, relationship := range relationships {
		if err := fixture.repository.Save(fixture.ctx, relationship); err != nil {
			t.Fatal(err)
		}
	}
}

func newTestRelationship(
	t *testing.T,
	id uuid.UUID,
	sourceKind relatables.Kind,
	targetKind relatables.Kind,
	similarity float64,
) domain_relationships.Relationship {
	t.Helper()

	return newTestRelationshipWithRelatables(
		t,
		id,
		relatables.NewMockRelatable(uuid.New(), sourceKind),
		relatables.NewMockRelatable(uuid.New(), targetKind),
		similarity,
	)
}

func newTestRelationshipWithRelatables(
	t *testing.T,
	id uuid.UUID,
	source relatables.Relatable,
	target relatables.Relatable,
	similarity float64,
) domain_relationships.Relationship {
	t.Helper()

	relationship, err := domain_relationships.NewAdapter().ToDomain(
		domain_relationships.RelationshipInput{
			Identifier: id,
			Source:     source,
			Target:     target,
			Similarity: similarity,
			CreatedOn:  time.Now().UTC(),
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	return relationship
}

func assertRelationship(
	t *testing.T,
	result domain_relationships.Relationship,
	expected domain_relationships.Relationship,
) {
	t.Helper()

	if result == nil {
		t.Fatalf("expected relationship")
	}

	if result.Identifier() != expected.Identifier() {
		t.Fatalf("expected id %s, got %s", expected.Identifier(), result.Identifier())
	}

	if result.Source().Identifier() != expected.Source().Identifier() {
		t.Fatalf("expected source id %s, got %s", expected.Source().Identifier(), result.Source().Identifier())
	}

	if result.Source().RelationshipKind() != expected.Source().RelationshipKind() {
		t.Fatalf("expected source kind %s, got %s", expected.Source().RelationshipKind(), result.Source().RelationshipKind())
	}

	if result.Target().Identifier() != expected.Target().Identifier() {
		t.Fatalf("expected target id %s, got %s", expected.Target().Identifier(), result.Target().Identifier())
	}

	if result.Target().RelationshipKind() != expected.Target().RelationshipKind() {
		t.Fatalf("expected target kind %s, got %s", expected.Target().RelationshipKind(), result.Target().RelationshipKind())
	}

	if result.Similarity() != expected.Similarity() {
		t.Fatalf("expected similarity %f, got %f", expected.Similarity(), result.Similarity())
	}

	if result.CreatedOn().IsZero() {
		t.Fatalf("expected created on")
	}
}

func mustParseRelationshipUUID(value string) uuid.UUID {
	id, err := uuid.Parse(value)
	if err != nil {
		panic(err)
	}

	return id
}
