package redis

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/relatables"
)

func TestNewRelationshipRelatableRepository(t *testing.T) {
	client := redis.NewClient(&redis.Options{})
	adapter := relatables.NewMockRelatableAdapter()

	repository := NewRelationshipRelatableRepository(client, adapter)

	if repository == nil {
		t.Fatalf("expected repository")
	}
}

func TestRelationshipRelatableRepositorySaveAndFind(t *testing.T) {
	fixture := newRelationshipRelatableRepositoryFixture(t)

	first := relatables.NewMockRelatable(
		mustParseUUID("00000000-0000-0000-0000-000000000001"),
		relatables.UserKind,
	)

	second := relatables.NewMockRelatable(
		mustParseUUID("00000000-0000-0000-0000-000000000002"),
		relatables.PostKind,
	)

	if err := fixture.repository.Save(fixture.ctx, first); err != nil {
		t.Fatal(err)
	}

	if err := fixture.repository.Save(fixture.ctx, second); err != nil {
		t.Fatal(err)
	}

	result, err := fixture.repository.Find(fixture.ctx, 0, 10)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 relatables, got %d", len(result))
	}

	assertRelatable(t, result[0], first)
	assertRelatable(t, result[1], second)
}

func TestRelationshipRelatableRepositoryFindWithPagination(t *testing.T) {
	fixture := newRelationshipRelatableRepositoryFixture(t)

	first := relatables.NewMockRelatable(
		mustParseUUID("00000000-0000-0000-0000-000000000001"),
		relatables.UserKind,
	)

	second := relatables.NewMockRelatable(
		mustParseUUID("00000000-0000-0000-0000-000000000002"),
		relatables.PostKind,
	)

	third := relatables.NewMockRelatable(
		mustParseUUID("00000000-0000-0000-0000-000000000003"),
		relatables.TopicKind,
	)

	saveRelatables(t, fixture, first, second, third)

	result, err := fixture.repository.Find(fixture.ctx, 1, 2)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 relatables, got %d", len(result))
	}

	assertRelatable(t, result[0], second)
	assertRelatable(t, result[1], third)
}

func TestRelationshipRelatableRepositoryFindAfterWithoutCursor(t *testing.T) {
	fixture := newRelationshipRelatableRepositoryFixture(t)

	first := relatables.NewMockRelatable(
		mustParseUUID("00000000-0000-0000-0000-000000000001"),
		relatables.UserKind,
	)

	second := relatables.NewMockRelatable(
		mustParseUUID("00000000-0000-0000-0000-000000000002"),
		relatables.PostKind,
	)

	saveRelatables(t, fixture, first, second)

	result, err := fixture.repository.FindAfter(
		fixture.ctx,
		uuid.Nil,
		1,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 relatable, got %d", len(result))
	}

	assertRelatable(t, result[0], first)
}

func TestRelationshipRelatableRepositoryFindAfterWithCursor(t *testing.T) {
	fixture := newRelationshipRelatableRepositoryFixture(t)

	first := relatables.NewMockRelatable(
		mustParseUUID("00000000-0000-0000-0000-000000000001"),
		relatables.UserKind,
	)

	second := relatables.NewMockRelatable(
		mustParseUUID("00000000-0000-0000-0000-000000000002"),
		relatables.UserKind,
	)

	third := relatables.NewMockRelatable(
		mustParseUUID("00000000-0000-0000-0000-000000000003"),
		relatables.UserKind,
	)

	saveRelatables(t, fixture, first, second, third)

	result, err := fixture.repository.FindAfter(
		fixture.ctx,
		first.Identifier(),
		2,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 relatables, got %d", len(result))
	}

	assertRelatable(t, result[0], second)
	assertRelatable(t, result[1], third)
}

func TestRelationshipRelatableRepositoryFindAfterReturnsEmptyWhenCursorMissing(t *testing.T) {
	fixture := newRelationshipRelatableRepositoryFixture(t)

	result, err := fixture.repository.FindAfter(
		fixture.ctx,
		uuid.New(),
		10,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 0 {
		t.Fatalf("expected empty result, got %d", len(result))
	}
}

func TestRelationshipRelatableRepositoryCount(t *testing.T) {
	fixture := newRelationshipRelatableRepositoryFixture(t)

	saveRelatables(
		t,
		fixture,
		relatables.NewMockRelatable(
			mustParseUUID("00000000-0000-0000-0000-000000000001"),
			relatables.UserKind,
		),
		relatables.NewMockRelatable(
			mustParseUUID("00000000-0000-0000-0000-000000000002"),
			relatables.PostKind,
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

func TestRelationshipRelatableRepositoryFindByKind(t *testing.T) {
	fixture := newRelationshipRelatableRepositoryFixture(t)

	user := relatables.NewMockRelatable(
		mustParseUUID("00000000-0000-0000-0000-000000000001"),
		relatables.UserKind,
	)

	post := relatables.NewMockRelatable(
		mustParseUUID("00000000-0000-0000-0000-000000000002"),
		relatables.PostKind,
	)

	saveRelatables(t, fixture, user, post)

	result, err := fixture.repository.FindByKind(
		fixture.ctx,
		relatables.UserKind,
		0,
		10,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 relatable, got %d", len(result))
	}

	assertRelatable(t, result[0], user)
}

func TestRelationshipRelatableRepositoryCountByKind(t *testing.T) {
	fixture := newRelationshipRelatableRepositoryFixture(t)

	saveRelatables(
		t,
		fixture,
		relatables.NewMockRelatable(
			mustParseUUID("00000000-0000-0000-0000-000000000001"),
			relatables.UserKind,
		),
		relatables.NewMockRelatable(
			mustParseUUID("00000000-0000-0000-0000-000000000002"),
			relatables.UserKind,
		),
		relatables.NewMockRelatable(
			mustParseUUID("00000000-0000-0000-0000-000000000003"),
			relatables.PostKind,
		),
	)

	result, err := fixture.repository.CountByKind(
		fixture.ctx,
		relatables.UserKind,
	)

	if err != nil {
		t.Fatal(err)
	}

	if result != 2 {
		t.Fatalf("expected count 2, got %d", result)
	}
}

func TestRelationshipRelatableRepositoryDelete(t *testing.T) {
	fixture := newRelationshipRelatableRepositoryFixture(t)

	relatable := relatables.NewMockRelatable(
		mustParseUUID("00000000-0000-0000-0000-000000000001"),
		relatables.UserKind,
	)

	if err := fixture.repository.Save(fixture.ctx, relatable); err != nil {
		t.Fatal(err)
	}

	if err := fixture.repository.Delete(fixture.ctx, relatable); err != nil {
		t.Fatal(err)
	}

	count, err := fixture.repository.Count(fixture.ctx)

	if err != nil {
		t.Fatal(err)
	}

	if count != 0 {
		t.Fatalf("expected count 0, got %d", count)
	}

	kindCount, err := fixture.repository.CountByKind(
		fixture.ctx,
		relatables.UserKind,
	)

	if err != nil {
		t.Fatal(err)
	}

	if kindCount != 0 {
		t.Fatalf("expected kind count 0, got %d", kindCount)
	}
}

func TestRelationshipRelatableRepositoryDeleteByID(t *testing.T) {
	fixture := newRelationshipRelatableRepositoryFixture(t)

	relatable := relatables.NewMockRelatable(
		mustParseUUID("00000000-0000-0000-0000-000000000001"),
		relatables.UserKind,
	)

	if err := fixture.repository.Save(fixture.ctx, relatable); err != nil {
		t.Fatal(err)
	}

	if err := fixture.repository.DeleteByID(
		fixture.ctx,
		relatable.Identifier(),
	); err != nil {
		t.Fatal(err)
	}

	count, err := fixture.repository.Count(fixture.ctx)

	if err != nil {
		t.Fatal(err)
	}

	if count != 0 {
		t.Fatalf("expected count 0, got %d", count)
	}
}

func TestRelationshipRelatableRepositoryDeleteByIDIgnoresMissingID(t *testing.T) {
	fixture := newRelationshipRelatableRepositoryFixture(t)

	err := fixture.repository.DeleteByID(fixture.ctx, uuid.New())

	if err != nil {
		t.Fatal(err)
	}
}

type relationshipRelatableRepositoryFixture struct {
	ctx        context.Context
	client     redis.UniversalClient
	adapter    *relatables.MockRelatableAdapter
	repository relatables.Repository
}

func newRelationshipRelatableRepositoryFixture(
	t *testing.T,
) *relationshipRelatableRepositoryFixture {
	t.Helper()

	address := os.Getenv("REDIS_TEST_ADDR")
	if address == "" {
		t.Skip("REDIS_TEST_ADDR is not set")
	}

	ctx := context.Background()

	client := redis.NewClient(
		&redis.Options{
			Addr: address,
		},
	)

	if err := client.Ping(ctx).Err(); err != nil {
		t.Fatal(err)
	}

	adapter := relatables.NewMockRelatableAdapter()

	fixture := &relationshipRelatableRepositoryFixture{
		ctx:        ctx,
		client:     client,
		adapter:    adapter,
		repository: NewRelationshipRelatableRepository(client, adapter),
	}

	truncateRelationshipRelatableRepository(t, fixture)

	t.Cleanup(func() {
		truncateRelationshipRelatableRepository(t, fixture)
		_ = client.Close()
	})

	return fixture
}

func truncateRelationshipRelatableRepository(
	t *testing.T,
	fixture *relationshipRelatableRepositoryFixture,
) {
	t.Helper()

	keys := []string{
		relationshipRelatablesKey,
		relationshipRelatableKindsKey,
		"relationship:relatables:kind:" + string(relatables.UserKind),
		"relationship:relatables:kind:" + string(relatables.PostKind),
		"relationship:relatables:kind:" + string(relatables.TopicKind),
		"relationship:relatables:kind:" + string(relatables.CampaignKind),
		"relationship:relatables:kind:" + string(relatables.NarrativeKind),
	}

	if err := fixture.client.Del(fixture.ctx, keys...).Err(); err != nil {
		t.Fatal(err)
	}
}

func saveRelatables(
	t *testing.T,
	fixture *relationshipRelatableRepositoryFixture,
	items ...relatables.Relatable,
) {
	t.Helper()

	for _, item := range items {
		if err := fixture.repository.Save(fixture.ctx, item); err != nil {
			t.Fatal(err)
		}
	}
}

func assertRelatable(
	t *testing.T,
	result relatables.Relatable,
	expected relatables.Relatable,
) {
	t.Helper()

	if result == nil {
		t.Fatalf("expected relatable")
	}

	if result.Identifier() != expected.Identifier() {
		t.Fatalf(
			"expected id %s, got %s",
			expected.Identifier(),
			result.Identifier(),
		)
	}

	if result.RelationshipKind() != expected.RelationshipKind() {
		t.Fatalf(
			"expected kind %s, got %s",
			expected.RelationshipKind(),
			result.RelationshipKind(),
		)
	}
}

func mustParseUUID(value string) uuid.UUID {
	id, err := uuid.Parse(value)
	if err != nil {
		panic(err)
	}

	return id
}
