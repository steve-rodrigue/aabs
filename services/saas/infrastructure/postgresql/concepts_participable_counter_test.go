package postgresql

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/participatables"
)

func TestNewConceptParticipatableCounter(t *testing.T) {
	fixture := newConceptParticipatableCounterFixture(t)

	counter := NewConceptParticipatableCounter(fixture.pool)

	if counter == nil {
		t.Fatalf("expected counter")
	}
}

func TestConceptParticipatableCounterCountUserByCommunity(t *testing.T) {
	fixture := newConceptParticipatableCounterFixture(t)

	user := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.UserKind,
	)

	community := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.CommunityKind,
	)

	otherCommunity := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.CommunityKind,
	)

	insertPost(t, fixture, uuid.New(), user.Identifier(), []uuid.UUID{
		community.Identifier(),
	})

	insertPost(t, fixture, uuid.New(), user.Identifier(), []uuid.UUID{
		community.Identifier(),
		otherCommunity.Identifier(),
	})

	insertPost(t, fixture, uuid.New(), user.Identifier(), []uuid.UUID{
		otherCommunity.Identifier(),
	})

	count, err := fixture.counter.CountByParticipantAndTarget(
		fixture.ctx,
		user,
		community,
	)
	if err != nil {
		t.Fatal(err)
	}

	if count != 2 {
		t.Fatalf("expected count 2, got %d", count)
	}
}

func TestConceptParticipatableCounterCountByCommunity(t *testing.T) {
	fixture := newConceptParticipatableCounterFixture(t)

	user := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.UserKind,
	)

	community := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.CommunityKind,
	)

	otherCommunity := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.CommunityKind,
	)

	insertPost(t, fixture, uuid.New(), user.Identifier(), []uuid.UUID{
		community.Identifier(),
	})

	insertPost(t, fixture, uuid.New(), user.Identifier(), []uuid.UUID{
		community.Identifier(),
	})

	insertPost(t, fixture, uuid.New(), user.Identifier(), []uuid.UUID{
		otherCommunity.Identifier(),
	})

	count, err := fixture.counter.CountByTarget(
		fixture.ctx,
		community,
	)
	if err != nil {
		t.Fatal(err)
	}

	if count != 2 {
		t.Fatalf("expected count 2, got %d", count)
	}
}

func TestConceptParticipatableCounterCountPlatformByTarget(t *testing.T) {
	fixture := newConceptParticipatableCounterFixture(t)

	platform := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.PlatformKind,
	)

	otherPlatform := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.PlatformKind,
	)

	user := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.UserKind,
	)

	otherUser := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.UserKind,
	)

	community := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.CommunityKind,
	)

	insertUser(t, fixture, user.Identifier(), platform.Identifier())
	insertUser(t, fixture, otherUser.Identifier(), otherPlatform.Identifier())

	insertPost(t, fixture, uuid.New(), user.Identifier(), []uuid.UUID{
		community.Identifier(),
	})

	insertPost(t, fixture, uuid.New(), user.Identifier(), []uuid.UUID{
		community.Identifier(),
	})

	insertPost(t, fixture, uuid.New(), otherUser.Identifier(), []uuid.UUID{
		community.Identifier(),
	})

	count, err := fixture.counter.CountByParticipantAndTarget(
		fixture.ctx,
		platform,
		community,
	)
	if err != nil {
		t.Fatal(err)
	}

	if count != 2 {
		t.Fatalf("expected count 2, got %d", count)
	}
}

func TestConceptParticipatableCounterCountByPlatform(t *testing.T) {
	fixture := newConceptParticipatableCounterFixture(t)

	platform := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.PlatformKind,
	)

	otherPlatform := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.PlatformKind,
	)

	user := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.UserKind,
	)

	otherUser := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.UserKind,
	)

	insertUser(t, fixture, user.Identifier(), platform.Identifier())
	insertUser(t, fixture, otherUser.Identifier(), otherPlatform.Identifier())

	insertPost(t, fixture, uuid.New(), user.Identifier(), nil)
	insertPost(t, fixture, uuid.New(), user.Identifier(), nil)
	insertPost(t, fixture, uuid.New(), otherUser.Identifier(), nil)

	count, err := fixture.counter.CountByTarget(
		fixture.ctx,
		platform,
	)
	if err != nil {
		t.Fatal(err)
	}

	if count != 2 {
		t.Fatalf("expected count 2, got %d", count)
	}
}

func TestConceptParticipatableCounterCountPostByGroupingTarget(t *testing.T) {
	fixture := newConceptParticipatableCounterFixture(t)

	postID := uuid.New()

	post := participatables.NewMockParticipatable(
		postID,
		participatables.PostKind,
	)

	campaign := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.CampaignKind,
	)

	otherCampaign := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.CampaignKind,
	)

	insertPost(t, fixture, postID, uuid.New(), nil)

	insertGroupingClusterMember(
		t,
		fixture,
		campaign.Identifier(),
		participatables.CampaignKind,
		postID,
	)

	insertGroupingClusterMember(
		t,
		fixture,
		otherCampaign.Identifier(),
		participatables.CampaignKind,
		uuid.New(),
	)

	count, err := fixture.counter.CountByParticipantAndTarget(
		fixture.ctx,
		post,
		campaign,
	)
	if err != nil {
		t.Fatal(err)
	}

	if count != 1 {
		t.Fatalf("expected count 1, got %d", count)
	}
}

func TestConceptParticipatableCounterCountByGrouping(t *testing.T) {
	fixture := newConceptParticipatableCounterFixture(t)

	topic := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.TopicKind,
	)

	firstPost := uuid.New()
	secondPost := uuid.New()
	thirdPost := uuid.New()

	insertGroupingClusterMember(
		t,
		fixture,
		topic.Identifier(),
		participatables.TopicKind,
		firstPost,
	)

	insertGroupingClusterMember(
		t,
		fixture,
		topic.Identifier(),
		participatables.TopicKind,
		secondPost,
	)

	insertGroupingClusterMember(
		t,
		fixture,
		topic.Identifier(),
		participatables.TopicKind,
		thirdPost,
	)

	count, err := fixture.counter.CountByTarget(
		fixture.ctx,
		topic,
	)
	if err != nil {
		t.Fatal(err)
	}

	if count != 3 {
		t.Fatalf("expected count 3, got %d", count)
	}
}

func TestConceptParticipatableCounterReturnsInvalidKindErrorForParticipant(t *testing.T) {
	fixture := newConceptParticipatableCounterFixture(t)

	participant := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.Kind("invalid"),
	)

	target := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.UserKind,
	)

	_, err := fixture.counter.CountByParticipantAndTarget(
		fixture.ctx,
		participant,
		target,
	)

	if !errors.Is(err, ErrInvalidConceptParticipatableCounterKind) {
		t.Fatalf("expected invalid kind error, got %v", err)
	}
}

func TestConceptParticipatableCounterReturnsInvalidKindErrorForTarget(t *testing.T) {
	fixture := newConceptParticipatableCounterFixture(t)

	target := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.Kind("invalid"),
	)

	_, err := fixture.counter.CountByTarget(
		fixture.ctx,
		target,
	)

	if !errors.Is(err, ErrInvalidConceptParticipatableCounterKind) {
		t.Fatalf("expected invalid kind error, got %v", err)
	}
}

type conceptParticipatableCounterFixture struct {
	ctx context.Context

	pool    *pgxpool.Pool
	counter participatables.Counter
}

func newConceptParticipatableCounterFixture(
	t *testing.T,
) *conceptParticipatableCounterFixture {
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

	fixture := &conceptParticipatableCounterFixture{
		ctx:  ctx,
		pool: pool,
	}

	createConceptParticipatableCounterTables(t, fixture)
	truncateConceptParticipatableCounterTables(t, fixture)

	fixture.counter = NewConceptParticipatableCounter(pool)

	t.Cleanup(func() {
		truncateConceptParticipatableCounterTables(t, fixture)
		pool.Close()
	})

	return fixture
}

func createConceptParticipatableCounterTables(
	t *testing.T,
	fixture *conceptParticipatableCounterFixture,
) {
	t.Helper()

	_, err := fixture.pool.Exec(
		fixture.ctx,
		`
		CREATE TABLE IF NOT EXISTS users (
			identifier UUID PRIMARY KEY,
			platform_id UUID NOT NULL
		)
		`,
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = fixture.pool.Exec(
		fixture.ctx,
		`
		CREATE TABLE IF NOT EXISTS posts (
			identifier UUID PRIMARY KEY,
			user_id UUID NOT NULL,
			community_ids UUID[] NOT NULL DEFAULT '{}'
		)
		`,
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = fixture.pool.Exec(
		fixture.ctx,
		`
		CREATE TABLE IF NOT EXISTS groupings_clusters (
			identifier UUID PRIMARY KEY,
			target_id UUID NOT NULL,
			target_kind TEXT NOT NULL
		)
		`,
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = fixture.pool.Exec(
		fixture.ctx,
		`
		CREATE TABLE IF NOT EXISTS groupings_clusters_members (
			cluster_id UUID NOT NULL,
			member_id UUID NOT NULL,
			member_kind TEXT NOT NULL
		)
		`,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func truncateConceptParticipatableCounterTables(
	t *testing.T,
	fixture *conceptParticipatableCounterFixture,
) {
	t.Helper()

	_, err := fixture.pool.Exec(
		fixture.ctx,
		`
		TRUNCATE TABLE
			groupings_clusters_members,
			groupings_clusters,
			posts,
			users
		`,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func insertUser(
	t *testing.T,
	fixture *conceptParticipatableCounterFixture,
	id uuid.UUID,
	platform uuid.UUID,
) {
	t.Helper()

	_, err := fixture.pool.Exec(
		fixture.ctx,
		`
		INSERT INTO users (
			identifier,
			platform_id
		)
		VALUES ($1, $2)
		ON CONFLICT (identifier)
		DO UPDATE SET
			platform_id = EXCLUDED.platform_id
		`,
		id,
		platform,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func insertPost(
	t *testing.T,
	fixture *conceptParticipatableCounterFixture,
	id uuid.UUID,
	user uuid.UUID,
	communities []uuid.UUID,
) {
	t.Helper()

	if communities == nil {
		communities = []uuid.UUID{}
	}

	_, err := fixture.pool.Exec(
		fixture.ctx,
		`
		INSERT INTO posts (
			identifier,
			user_id,
			community_ids
		)
		VALUES ($1, $2, $3)
		ON CONFLICT (identifier)
		DO UPDATE SET
			user_id = EXCLUDED.user_id,
			community_ids = EXCLUDED.community_ids
		`,
		id,
		user,
		communities,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func insertGroupingClusterMember(
	t *testing.T,
	fixture *conceptParticipatableCounterFixture,
	targetID uuid.UUID,
	targetKind participatables.Kind,
	postID uuid.UUID,
) {
	t.Helper()

	clusterID := uuid.New()

	_, err := fixture.pool.Exec(
		fixture.ctx,
		`
		INSERT INTO groupings_clusters (
			identifier,
			target_id,
			target_kind
		)
		VALUES ($1, $2, $3)
		`,
		clusterID,
		targetID,
		string(targetKind),
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = fixture.pool.Exec(
		fixture.ctx,
		`
		INSERT INTO groupings_clusters_members (
			cluster_id,
			member_id,
			member_kind
		)
		VALUES ($1, $2, $3)
		`,
		clusterID,
		postID,
		string(participatables.PostKind),
	)
	if err != nil {
		t.Fatal(err)
	}
}
