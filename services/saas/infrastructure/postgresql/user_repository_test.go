package postgresql

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	domain_platforms "github.com/steve-rodrigue/aabs/services/saas/domain/entities/platforms"
	domain_users "github.com/steve-rodrigue/aabs/services/saas/domain/entities/users"
)

func TestNewUserRepository(t *testing.T) {
	pool := &pgxpool.Pool{}
	adapter := domain_users.NewMockUserAdapter()
	platforms := domain_platforms.NewMockPlatformRepository()

	repository := NewUserRepository(pool, adapter, platforms)

	if repository == nil {
		t.Fatalf("expected repository")
	}
}

func TestUserRepositorySaveAndFindByID(t *testing.T) {
	fixture := newUserRepositoryFixture(t)

	platform := domain_platforms.NewMockPlatform("Reddit", "reddit")
	fixture.platforms.Items[platform.Identifier()] = platform

	user := newTestUser(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000001"),
		platform,
		"external-1",
		"steve-rodrigue",
		"Steve Rodrigue",
		"https://reddit.com/u/steve-rodrigue",
	)

	if err := fixture.repository.Save(fixture.ctx, user); err != nil {
		t.Fatal(err)
	}

	result, err := fixture.repository.FindByID(
		fixture.ctx,
		user.Identifier(),
	)

	if err != nil {
		t.Fatal(err)
	}

	assertUser(t, result, user)

	if fixture.adapter.ToDomainCalls != 1 {
		t.Fatalf("expected 1 adapter call, got %d", fixture.adapter.ToDomainCalls)
	}

	if fixture.platforms.FindByIDCalls != 1 {
		t.Fatalf("expected 1 platform lookup, got %d", fixture.platforms.FindByIDCalls)
	}
}

func TestUserRepositorySaveUpdatesExistingUser(t *testing.T) {
	fixture := newUserRepositoryFixture(t)

	platform := domain_platforms.NewMockPlatform("Reddit", "reddit")
	fixture.platforms.Items[platform.Identifier()] = platform

	id := mustParseUUID("00000000-0000-0000-0000-000000000001")

	first := newTestUser(
		t,
		id,
		platform,
		"external-1",
		"steve-rodrigue",
		"Steve Rodrigue",
		"https://reddit.com/u/steve-rodrigue",
	)

	second := newTestUser(
		t,
		id,
		platform,
		"external-2",
		"steve-rodrigue-updated",
		"Steve Updated",
		"https://reddit.com/u/steve-rodrigue-updated",
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

	assertUser(t, result, second)
}

func TestUserRepositoryFindByIDReturnsNilWhenNotFound(t *testing.T) {
	fixture := newUserRepositoryFixture(t)

	result, err := fixture.repository.FindByID(
		fixture.ctx,
		uuid.New(),
	)

	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Fatalf("expected nil user")
	}
}

func TestUserRepositoryFindByPlatformAndExternalID(t *testing.T) {
	fixture := newUserRepositoryFixture(t)

	platform := domain_platforms.NewMockPlatform("Reddit", "reddit")
	fixture.platforms.Items[platform.Identifier()] = platform

	user := newTestUser(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000001"),
		platform,
		"external-1",
		"steve-rodrigue",
		"Steve Rodrigue",
		"https://reddit.com/u/steve-rodrigue",
	)

	if err := fixture.repository.Save(fixture.ctx, user); err != nil {
		t.Fatal(err)
	}

	result, err := fixture.repository.FindByPlatformAndExternalID(
		fixture.ctx,
		platform,
		"external-1",
	)

	if err != nil {
		t.Fatal(err)
	}

	assertUser(t, result, user)
}

func TestUserRepositoryFindByPlatformAndHandle(t *testing.T) {
	fixture := newUserRepositoryFixture(t)

	platform := domain_platforms.NewMockPlatform("Reddit", "reddit")
	fixture.platforms.Items[platform.Identifier()] = platform

	user := newTestUser(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000001"),
		platform,
		"external-1",
		"steve-rodrigue",
		"Steve Rodrigue",
		"https://reddit.com/u/steve-rodrigue",
	)

	if err := fixture.repository.Save(fixture.ctx, user); err != nil {
		t.Fatal(err)
	}

	result, err := fixture.repository.FindByPlatformAndHandle(
		fixture.ctx,
		platform,
		"steve-rodrigue",
	)

	if err != nil {
		t.Fatal(err)
	}

	assertUser(t, result, user)
}

func TestUserRepositoryFind(t *testing.T) {
	fixture := newUserRepositoryFixture(t)

	platform := domain_platforms.NewMockPlatform("Reddit", "reddit")
	fixture.platforms.Items[platform.Identifier()] = platform

	first := newTestUser(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000001"),
		platform,
		"external-1",
		"first",
		"First User",
		"https://reddit.com/u/first",
	)

	second := newTestUser(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000002"),
		platform,
		"external-2",
		"second",
		"Second User",
		"https://reddit.com/u/second",
	)

	third := newTestUser(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000003"),
		platform,
		"external-3",
		"third",
		"Third User",
		"https://reddit.com/u/third",
	)

	saveUsers(t, fixture, first, second, third)

	result, err := fixture.repository.Find(fixture.ctx, 1, 2)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 users, got %d", len(result))
	}

	assertUser(t, result[0], second)
	assertUser(t, result[1], third)
}

func TestUserRepositoryFindAfterWithoutCursor(t *testing.T) {
	fixture := newUserRepositoryFixture(t)

	platform := domain_platforms.NewMockPlatform("Reddit", "reddit")
	fixture.platforms.Items[platform.Identifier()] = platform

	first := newTestUser(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000001"),
		platform,
		"external-1",
		"first",
		"First User",
		"https://reddit.com/u/first",
	)

	second := newTestUser(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000002"),
		platform,
		"external-2",
		"second",
		"Second User",
		"https://reddit.com/u/second",
	)

	saveUsers(t, fixture, first, second)

	result, err := fixture.repository.FindAfter(
		fixture.ctx,
		uuid.Nil,
		1,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 user, got %d", len(result))
	}

	assertUser(t, result[0], first)
}

func TestUserRepositoryFindAfterWithCursor(t *testing.T) {
	fixture := newUserRepositoryFixture(t)

	platform := domain_platforms.NewMockPlatform("Reddit", "reddit")
	fixture.platforms.Items[platform.Identifier()] = platform

	first := newTestUser(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000001"),
		platform,
		"external-1",
		"first",
		"First User",
		"https://reddit.com/u/first",
	)

	second := newTestUser(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000002"),
		platform,
		"external-2",
		"second",
		"Second User",
		"https://reddit.com/u/second",
	)

	third := newTestUser(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000003"),
		platform,
		"external-3",
		"third",
		"Third User",
		"https://reddit.com/u/third",
	)

	saveUsers(t, fixture, first, second, third)

	result, err := fixture.repository.FindAfter(
		fixture.ctx,
		first.Identifier(),
		2,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 users, got %d", len(result))
	}

	assertUser(t, result[0], second)
	assertUser(t, result[1], third)
}

func TestUserRepositoryCount(t *testing.T) {
	fixture := newUserRepositoryFixture(t)

	platform := domain_platforms.NewMockPlatform("Reddit", "reddit")
	fixture.platforms.Items[platform.Identifier()] = platform

	saveUsers(
		t,
		fixture,
		newTestUser(
			t,
			mustParseUUID("00000000-0000-0000-0000-000000000001"),
			platform,
			"external-1",
			"first",
			"First User",
			"https://reddit.com/u/first",
		),
		newTestUser(
			t,
			mustParseUUID("00000000-0000-0000-0000-000000000002"),
			platform,
			"external-2",
			"second",
			"Second User",
			"https://reddit.com/u/second",
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

type userRepositoryFixture struct {
	ctx        context.Context
	pool       *pgxpool.Pool
	adapter    *domain_users.MockUserAdapter
	platforms  *domain_platforms.MockPlatformRepository
	repository domain_users.Repository
}

func newUserRepositoryFixture(t *testing.T) *userRepositoryFixture {
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

	adapter := domain_users.NewMockUserAdapter()
	platforms := domain_platforms.NewMockPlatformRepository()

	fixture := &userRepositoryFixture{
		ctx:        ctx,
		pool:       pool,
		adapter:    adapter,
		platforms:  platforms,
		repository: NewUserRepository(pool, adapter, platforms),
	}

	createUsersTable(t, fixture)
	truncateUsersTable(t, fixture)

	t.Cleanup(func() {
		truncateUsersTable(t, fixture)
		pool.Close()
	})

	return fixture
}

func createUsersTable(
	t *testing.T,
	fixture *userRepositoryFixture,
) {
	t.Helper()

	_, err := fixture.pool.Exec(
		fixture.ctx,
		`
		CREATE TABLE IF NOT EXISTS users (
			identifier UUID PRIMARY KEY,
			platform_id UUID NOT NULL,
			external_id TEXT NOT NULL,
			handle TEXT NOT NULL,
			display_name TEXT NOT NULL,
			profile_url TEXT NOT NULL,
			created_on TIMESTAMPTZ NOT NULL,

			UNIQUE (platform_id, external_id),
			UNIQUE (platform_id, handle)
		)
		`,
	)

	if err != nil {
		t.Fatal(err)
	}
}

func truncateUsersTable(
	t *testing.T,
	fixture *userRepositoryFixture,
) {
	t.Helper()

	_, err := fixture.pool.Exec(
		fixture.ctx,
		`
		TRUNCATE TABLE users
		`,
	)

	if err != nil {
		t.Fatal(err)
	}
}

func saveUsers(
	t *testing.T,
	fixture *userRepositoryFixture,
	users ...domain_users.User,
) {
	t.Helper()

	for _, user := range users {
		if err := fixture.repository.Save(fixture.ctx, user); err != nil {
			t.Fatal(err)
		}
	}
}

func assertUser(
	t *testing.T,
	result domain_users.User,
	expected domain_users.User,
) {
	t.Helper()

	if result == nil {
		t.Fatalf("expected user")
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

	if result.Platform() != expected.Platform() {
		t.Fatalf("expected platform")
	}

	if result.ExternalID() != expected.ExternalID() {
		t.Fatalf(
			"expected external id %q, got %q",
			expected.ExternalID(),
			result.ExternalID(),
		)
	}

	if result.Handle() != expected.Handle() {
		t.Fatalf("expected handle %q, got %q", expected.Handle(), result.Handle())
	}

	if result.DisplayName() != expected.DisplayName() {
		t.Fatalf(
			"expected display name %q, got %q",
			expected.DisplayName(),
			result.DisplayName(),
		)
	}

	if result.ProfileURL() != expected.ProfileURL() {
		t.Fatalf(
			"expected profile url %q, got %q",
			expected.ProfileURL(),
			result.ProfileURL(),
		)
	}
}

func newTestUser(
	t *testing.T,
	id uuid.UUID,
	platform domain_platforms.Platform,
	externalID string,
	handle string,
	displayName string,
	profileURL string,
) domain_users.User {
	t.Helper()

	user, err := domain_users.NewAdapter().ToDomain(
		domain_users.UserInput{
			Identifier:  id,
			Platform:    platform,
			ExternalID:  externalID,
			Handle:      handle,
			DisplayName: displayName,
			ProfileURL:  profileURL,
			CreatedOn:   time.Now().UTC(),
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	return user
}
