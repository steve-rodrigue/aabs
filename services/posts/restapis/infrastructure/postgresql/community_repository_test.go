package postgresql

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	domain_communities "github.com/steve-rodrigue/aabs/services/posts/restapis/domain/communities"
	domain_platforms "github.com/steve-rodrigue/aabs/services/posts/restapis/domain/platforms"
	domain_users "github.com/steve-rodrigue/aabs/services/posts/restapis/domain/users"
)

func TestNewCommunityRepository(t *testing.T) {
	pool := &pgxpool.Pool{}
	adapter := domain_communities.NewMockCommunityAdapter()
	platforms := domain_platforms.NewMockPlatformRepository()
	users := domain_users.NewMockUserRepository()

	repository := NewCommunityRepository(pool, adapter, platforms, users)

	if repository == nil {
		t.Fatalf("expected repository")
	}
}

func TestCommunityRepositorySaveAndFindByID(t *testing.T) {
	fixture := newCommunityRepositoryFixture(t)

	platform := domain_platforms.NewMockPlatform("Reddit", "reddit")
	moderator := domain_users.NewMockUser("@mod", "Moderator")

	fixture.platforms.Items[platform.Identifier()] = platform
	fixture.users.Items[moderator.Identifier()] = moderator

	community := newTestCommunity(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000001"),
		platform,
		"aabs",
		"AABS",
		"Anti-AI Bot Spam",
		[]domain_users.User{moderator},
	)

	if err := fixture.repository.Save(fixture.ctx, community); err != nil {
		t.Fatal(err)
	}

	result, err := fixture.repository.FindByID(
		fixture.ctx,
		community.Identifier(),
	)

	if err != nil {
		t.Fatal(err)
	}

	assertCommunity(t, result, community)

	if fixture.adapter.ToDomainCalls != 1 {
		t.Fatalf("expected 1 adapter call, got %d", fixture.adapter.ToDomainCalls)
	}

	if fixture.platforms.FindByIDCalls != 1 {
		t.Fatalf("expected 1 platform lookup, got %d", fixture.platforms.FindByIDCalls)
	}

	if fixture.users.FindByIDCalls != 1 {
		t.Fatalf("expected 1 moderator lookup, got %d", fixture.users.FindByIDCalls)
	}
}

func TestCommunityRepositorySaveUpdatesExistingCommunity(t *testing.T) {
	fixture := newCommunityRepositoryFixture(t)

	platform := domain_platforms.NewMockPlatform("Reddit", "reddit")
	firstModerator := domain_users.NewMockUser("@first", "First Moderator")
	secondModerator := domain_users.NewMockUser("@second", "Second Moderator")

	fixture.platforms.Items[platform.Identifier()] = platform
	fixture.users.Items[firstModerator.Identifier()] = firstModerator
	fixture.users.Items[secondModerator.Identifier()] = secondModerator

	id := mustParseUUID("00000000-0000-0000-0000-000000000001")

	first := newTestCommunity(
		t,
		id,
		platform,
		"aabs",
		"AABS",
		"Old text",
		[]domain_users.User{firstModerator},
	)

	second := newTestCommunity(
		t,
		id,
		platform,
		"aabs-updated",
		"AABS Updated",
		"Updated text",
		[]domain_users.User{secondModerator},
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

	assertCommunity(t, result, second)
}

func TestCommunityRepositoryFindByIDReturnsNilWhenNotFound(t *testing.T) {
	fixture := newCommunityRepositoryFixture(t)

	result, err := fixture.repository.FindByID(
		fixture.ctx,
		uuid.New(),
	)

	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Fatalf("expected nil community")
	}
}

func TestCommunityRepositoryFindByHandle(t *testing.T) {
	fixture := newCommunityRepositoryFixture(t)

	platform := domain_platforms.NewMockPlatform("Reddit", "reddit")
	fixture.platforms.Items[platform.Identifier()] = platform

	community := newTestCommunity(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000001"),
		platform,
		"aabs",
		"AABS",
		"Anti-AI Bot Spam",
		nil,
	)

	if err := fixture.repository.Save(fixture.ctx, community); err != nil {
		t.Fatal(err)
	}

	result, err := fixture.repository.FindByHandle(
		fixture.ctx,
		platform,
		"aabs",
	)

	if err != nil {
		t.Fatal(err)
	}

	assertCommunity(t, result, community)
}

func TestCommunityRepositoryFind(t *testing.T) {
	fixture := newCommunityRepositoryFixture(t)

	platform := domain_platforms.NewMockPlatform("Reddit", "reddit")
	fixture.platforms.Items[platform.Identifier()] = platform

	first := newTestCommunity(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000001"),
		platform,
		"first",
		"First",
		"First text",
		nil,
	)

	second := newTestCommunity(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000002"),
		platform,
		"second",
		"Second",
		"Second text",
		nil,
	)

	third := newTestCommunity(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000003"),
		platform,
		"third",
		"Third",
		"Third text",
		nil,
	)

	saveCommunities(t, fixture, first, second, third)

	result, err := fixture.repository.Find(fixture.ctx, 1, 2)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 communities, got %d", len(result))
	}

	assertCommunity(t, result[0], second)
	assertCommunity(t, result[1], third)
}

func TestCommunityRepositoryFindAfterWithoutCursor(t *testing.T) {
	fixture := newCommunityRepositoryFixture(t)

	platform := domain_platforms.NewMockPlatform("Reddit", "reddit")
	fixture.platforms.Items[platform.Identifier()] = platform

	first := newTestCommunity(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000001"),
		platform,
		"first",
		"First",
		"First text",
		nil,
	)

	second := newTestCommunity(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000002"),
		platform,
		"second",
		"Second",
		"Second text",
		nil,
	)

	saveCommunities(t, fixture, first, second)

	result, err := fixture.repository.FindAfter(
		fixture.ctx,
		uuid.Nil,
		1,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 community, got %d", len(result))
	}

	assertCommunity(t, result[0], first)
}

func TestCommunityRepositoryFindAfterWithCursor(t *testing.T) {
	fixture := newCommunityRepositoryFixture(t)

	platform := domain_platforms.NewMockPlatform("Reddit", "reddit")
	fixture.platforms.Items[platform.Identifier()] = platform

	first := newTestCommunity(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000001"),
		platform,
		"first",
		"First",
		"First text",
		nil,
	)

	second := newTestCommunity(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000002"),
		platform,
		"second",
		"Second",
		"Second text",
		nil,
	)

	third := newTestCommunity(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000003"),
		platform,
		"third",
		"Third",
		"Third text",
		nil,
	)

	saveCommunities(t, fixture, first, second, third)

	result, err := fixture.repository.FindAfter(
		fixture.ctx,
		first.Identifier(),
		2,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 communities, got %d", len(result))
	}

	assertCommunity(t, result[0], second)
	assertCommunity(t, result[1], third)
}

func TestCommunityRepositoryFindByPlatform(t *testing.T) {
	fixture := newCommunityRepositoryFixture(t)

	platform := domain_platforms.NewMockPlatform("Reddit", "reddit")
	otherPlatform := domain_platforms.NewMockPlatform("YouTube", "youtube")

	fixture.platforms.Items[platform.Identifier()] = platform
	fixture.platforms.Items[otherPlatform.Identifier()] = otherPlatform

	first := newTestCommunity(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000001"),
		platform,
		"first",
		"First",
		"First text",
		nil,
	)

	second := newTestCommunity(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000002"),
		platform,
		"second",
		"Second",
		"Second text",
		nil,
	)

	other := newTestCommunity(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000003"),
		otherPlatform,
		"other",
		"Other",
		"Other text",
		nil,
	)

	saveCommunities(t, fixture, first, second, other)

	result, err := fixture.repository.FindByPlatform(
		fixture.ctx,
		platform,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 communities, got %d", len(result))
	}

	assertCommunity(t, result[0], first)
	assertCommunity(t, result[1], second)
}

func TestCommunityRepositoryCount(t *testing.T) {
	fixture := newCommunityRepositoryFixture(t)

	platform := domain_platforms.NewMockPlatform("Reddit", "reddit")
	fixture.platforms.Items[platform.Identifier()] = platform

	saveCommunities(
		t,
		fixture,
		newTestCommunity(
			t,
			mustParseUUID("00000000-0000-0000-0000-000000000001"),
			platform,
			"first",
			"First",
			"First text",
			nil,
		),
		newTestCommunity(
			t,
			mustParseUUID("00000000-0000-0000-0000-000000000002"),
			platform,
			"second",
			"Second",
			"Second text",
			nil,
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

type communityRepositoryFixture struct {
	ctx        context.Context
	pool       *pgxpool.Pool
	adapter    *domain_communities.MockCommunityAdapter
	platforms  *domain_platforms.MockPlatformRepository
	users      *domain_users.MockUserRepository
	repository domain_communities.Repository
}

func newCommunityRepositoryFixture(t *testing.T) *communityRepositoryFixture {
	t.Helper()

	dsn := os.Getenv("POSTS_POSTGRES_TEST_DSN")
	if dsn == "" {
		t.Skip("POSTS_POSTGRES_TEST_DSN is not set")
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}

	adapter := domain_communities.NewMockCommunityAdapter()
	platforms := domain_platforms.NewMockPlatformRepository()
	users := domain_users.NewMockUserRepository()

	fixture := &communityRepositoryFixture{
		ctx:        ctx,
		pool:       pool,
		adapter:    adapter,
		platforms:  platforms,
		users:      users,
		repository: NewCommunityRepository(pool, adapter, platforms, users),
	}

	createCommunitiesTables(t, fixture)
	truncateCommunitiesTables(t, fixture)

	t.Cleanup(func() {
		truncateCommunitiesTables(t, fixture)
		pool.Close()
	})

	return fixture
}

func createCommunitiesTables(
	t *testing.T,
	fixture *communityRepositoryFixture,
) {
	t.Helper()

	_, err := fixture.pool.Exec(
		fixture.ctx,
		`
			DROP TABLE IF EXISTS communities CASCADE;
			DROP TABLE IF EXISTS community_moderators CASCADE
		`,
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = fixture.pool.Exec(
		fixture.ctx,
		`
		CREATE TABLE IF NOT EXISTS communities (
			identifier UUID PRIMARY KEY,
			platform_id UUID NOT NULL,
			handle TEXT NOT NULL,
			title TEXT NOT NULL,
			text TEXT NOT NULL,
			created_on TIMESTAMPTZ NOT NULL,

			UNIQUE (platform_id, handle)
		);

		CREATE TABLE IF NOT EXISTS community_moderators (
			community_id UUID NOT NULL REFERENCES communities(identifier) ON DELETE CASCADE,
			user_id UUID NOT NULL,

			PRIMARY KEY (community_id, user_id)
		);
		`,
	)

	if err != nil {
		t.Fatal(err)
	}
}

func truncateCommunitiesTables(
	t *testing.T,
	fixture *communityRepositoryFixture,
) {
	t.Helper()

	_, err := fixture.pool.Exec(
		fixture.ctx,
		`
		TRUNCATE TABLE community_moderators, communities
		`,
	)

	if err != nil {
		t.Fatal(err)
	}
}

func saveCommunities(
	t *testing.T,
	fixture *communityRepositoryFixture,
	communities ...domain_communities.Community,
) {
	t.Helper()

	for _, community := range communities {
		if err := fixture.repository.Save(fixture.ctx, community); err != nil {
			t.Fatal(err)
		}
	}
}

func assertCommunity(
	t *testing.T,
	result domain_communities.Community,
	expected domain_communities.Community,
) {
	t.Helper()

	if result == nil {
		t.Fatalf("expected community")
	}

	if result.Identifier() != expected.Identifier() {
		t.Fatalf("expected id %s, got %s", expected.Identifier(), result.Identifier())
	}

	if result.Platform() != expected.Platform() {
		t.Fatalf("expected platform")
	}

	if result.Handle() != expected.Handle() {
		t.Fatalf("expected handle %q, got %q", expected.Handle(), result.Handle())
	}

	if result.Title() != expected.Title() {
		t.Fatalf("expected title %q, got %q", expected.Title(), result.Title())
	}

	if result.Text() != expected.Text() {
		t.Fatalf("expected text %q, got %q", expected.Text(), result.Text())
	}

	if len(result.Moderators()) != len(expected.Moderators()) {
		t.Fatalf(
			"expected %d moderators, got %d",
			len(expected.Moderators()),
			len(result.Moderators()),
		)
	}

	for index, moderator := range expected.Moderators() {
		if result.Moderators()[index] != moderator {
			t.Fatalf("expected moderator at index %d", index)
		}
	}
}

func newTestCommunity(
	t *testing.T,
	id uuid.UUID,
	platform domain_platforms.Platform,
	handle string,
	title string,
	text string,
	moderators []domain_users.User,
) domain_communities.Community {
	t.Helper()

	community, err := domain_communities.NewAdapter().ToDomain(
		domain_communities.CommunityInput{
			Identifier: id,
			Platform:   platform,
			Handle:     handle,
			Title:      title,
			Text:       text,
			CreatedOn:  time.Now().UTC(),
			Moderators: moderators,
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	return community
}
