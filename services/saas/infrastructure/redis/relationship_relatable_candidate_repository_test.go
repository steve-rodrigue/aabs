package redis

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships/relatables"
)

func TestNewRelationshipRelatableCandidateRepository(t *testing.T) {
	client := redis.NewClient(&redis.Options{})
	adapter := relatables.NewMockRelatableAdapter()

	repository := NewRelationshipRelatableCandidateRepository(client, adapter)

	if repository == nil {
		t.Fatalf("expected repository")
	}
}

func TestRelationshipRelatableCandidateRepositoryFindCandidates(t *testing.T) {
	fixture := newRelationshipRelatableCandidateRepositoryFixture(t)

	source := relatables.NewMockRelatable(
		mustParseUUID("00000000-0000-0000-0000-000000000001"),
		relatables.UserKind,
	)

	first := relatables.NewMockRelatable(
		mustParseUUID("00000000-0000-0000-0000-000000000002"),
		relatables.UserKind,
	)

	second := relatables.NewMockRelatable(
		mustParseUUID("00000000-0000-0000-0000-000000000003"),
		relatables.UserKind,
	)

	otherKind := relatables.NewMockRelatable(
		mustParseUUID("00000000-0000-0000-0000-000000000004"),
		relatables.PostKind,
	)

	saveCandidateRelatables(t, fixture, source, first, second, otherKind)

	result, err := fixture.repository.FindCandidates(
		fixture.ctx,
		source,
		2,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 candidates, got %d", len(result))
	}

	assertRelatable(t, result[0], first)
	assertRelatable(t, result[1], second)
}

func TestRelationshipRelatableCandidateRepositoryFindCandidatesSkipsSource(t *testing.T) {
	fixture := newRelationshipRelatableCandidateRepositoryFixture(t)

	source := relatables.NewMockRelatable(
		mustParseUUID("00000000-0000-0000-0000-000000000001"),
		relatables.UserKind,
	)

	candidate := relatables.NewMockRelatable(
		mustParseUUID("00000000-0000-0000-0000-000000000002"),
		relatables.UserKind,
	)

	saveCandidateRelatables(t, fixture, source, candidate)

	result, err := fixture.repository.FindCandidates(
		fixture.ctx,
		source,
		1,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 candidate, got %d", len(result))
	}

	assertRelatable(t, result[0], candidate)
}

func TestRelationshipRelatableCandidateRepositoryFindCandidatesReturnsEmptyWhenAmountIsZero(t *testing.T) {
	fixture := newRelationshipRelatableCandidateRepositoryFixture(t)

	source := relatables.NewMockRelatable(
		uuid.New(),
		relatables.UserKind,
	)

	result, err := fixture.repository.FindCandidates(
		fixture.ctx,
		source,
		0,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 0 {
		t.Fatalf("expected empty result, got %d", len(result))
	}
}

func TestRelationshipRelatableCandidateRepositoryFindCandidatesReturnsEmptyWhenAmountIsNegative(t *testing.T) {
	fixture := newRelationshipRelatableCandidateRepositoryFixture(t)

	source := relatables.NewMockRelatable(
		uuid.New(),
		relatables.UserKind,
	)

	result, err := fixture.repository.FindCandidates(
		fixture.ctx,
		source,
		-1,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 0 {
		t.Fatalf("expected empty result, got %d", len(result))
	}
}

func TestRelationshipRelatableCandidateRepositoryFindCandidatesReturnsEmptyWhenNoCandidatesExist(t *testing.T) {
	fixture := newRelationshipRelatableCandidateRepositoryFixture(t)

	source := relatables.NewMockRelatable(
		uuid.New(),
		relatables.UserKind,
	)

	result, err := fixture.repository.FindCandidates(
		fixture.ctx,
		source,
		10,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 0 {
		t.Fatalf("expected empty result, got %d", len(result))
	}
}

type relationshipRelatableCandidateRepositoryFixture struct {
	ctx                 context.Context
	client              redis.UniversalClient
	adapter             *relatables.MockRelatableAdapter
	relatableRepository relatables.Repository
	repository          relatables.CandidateRepository
}

func newRelationshipRelatableCandidateRepositoryFixture(
	t *testing.T,
) *relationshipRelatableCandidateRepositoryFixture {
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

	fixture := &relationshipRelatableCandidateRepositoryFixture{
		ctx:                 ctx,
		client:              client,
		adapter:             adapter,
		relatableRepository: NewRelationshipRelatableRepository(client, adapter),
		repository:          NewRelationshipRelatableCandidateRepository(client, adapter),
	}

	truncateRelationshipRelatableCandidateRepository(t, fixture)

	t.Cleanup(func() {
		truncateRelationshipRelatableCandidateRepository(t, fixture)
		_ = client.Close()
	})

	return fixture
}

func truncateRelationshipRelatableCandidateRepository(
	t *testing.T,
	fixture *relationshipRelatableCandidateRepositoryFixture,
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

func saveCandidateRelatables(
	t *testing.T,
	fixture *relationshipRelatableCandidateRepositoryFixture,
	items ...relatables.Relatable,
) {
	t.Helper()

	for _, item := range items {
		if err := fixture.relatableRepository.Save(fixture.ctx, item); err != nil {
			t.Fatal(err)
		}
	}
}
