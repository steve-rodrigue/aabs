package postgresql

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/participatables"
	domain_platforms "github.com/steve-rodrigue/aabs/services/saas/domain/entities/platforms"
)

func TestNewPlatformRepository(t *testing.T) {
	pool := &pgxpool.Pool{}
	adapter := domain_platforms.NewAdapter()

	repository := NewPlatformRepository(pool, adapter)

	if repository == nil {
		t.Fatalf("expected repository")
	}
}

func TestPlatformRepositorySaveAndFindByID(t *testing.T) {
	fixture := newPlatformRepositoryFixture(t)

	platform := newTestPlatform(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000001"),
		"Reddit",
		"reddit",
		"https://reddit.com",
	)

	err := fixture.repository.Save(fixture.ctx, platform)

	if err != nil {
		t.Fatal(err)
	}

	result, err := fixture.repository.FindByID(
		fixture.ctx,
		platform.Identifier(),
	)

	if err != nil {
		t.Fatal(err)
	}

	assertPlatform(t, result, platform)
}

func TestPlatformRepositorySaveUpdatesExistingPlatform(t *testing.T) {
	fixture := newPlatformRepositoryFixture(t)

	id := mustParseUUID("00000000-0000-0000-0000-000000000001")

	first := newTestPlatform(
		t,
		id,
		"Reddit",
		"reddit",
		"https://reddit.com",
	)

	second := newTestPlatform(
		t,
		id,
		"Reddit Updated",
		"reddit-updated",
		"https://www.reddit.com",
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

	assertPlatform(t, result, second)
}

func TestPlatformRepositoryFindByIDReturnsNilWhenNotFound(t *testing.T) {
	fixture := newPlatformRepositoryFixture(t)

	result, err := fixture.repository.FindByID(
		fixture.ctx,
		uuid.New(),
	)

	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Fatalf("expected nil platform")
	}
}

func TestPlatformRepositoryFindByHandle(t *testing.T) {
	fixture := newPlatformRepositoryFixture(t)

	platform := newTestPlatform(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000001"),
		"Reddit",
		"reddit",
		"https://reddit.com",
	)

	if err := fixture.repository.Save(fixture.ctx, platform); err != nil {
		t.Fatal(err)
	}

	result, err := fixture.repository.FindByHandle(
		fixture.ctx,
		"reddit",
	)

	if err != nil {
		t.Fatal(err)
	}

	assertPlatform(t, result, platform)
}

func TestPlatformRepositoryFindByName(t *testing.T) {
	fixture := newPlatformRepositoryFixture(t)

	platform := newTestPlatform(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000001"),
		"Reddit",
		"reddit",
		"https://reddit.com",
	)

	if err := fixture.repository.Save(fixture.ctx, platform); err != nil {
		t.Fatal(err)
	}

	result, err := fixture.repository.FindByName(
		fixture.ctx,
		"Reddit",
	)

	if err != nil {
		t.Fatal(err)
	}

	assertPlatform(t, result, platform)
}

func TestPlatformRepositoryFind(t *testing.T) {
	fixture := newPlatformRepositoryFixture(t)

	first := newTestPlatform(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000001"),
		"Reddit",
		"reddit",
		"https://reddit.com",
	)

	second := newTestPlatform(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000002"),
		"YouTube",
		"youtube",
		"https://youtube.com",
	)

	third := newTestPlatform(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000003"),
		"TikTok",
		"tiktok",
		"https://tiktok.com",
	)

	savePlatforms(t, fixture, first, second, third)

	result, err := fixture.repository.Find(fixture.ctx, 1, 2)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 platforms, got %d", len(result))
	}

	assertPlatform(t, result[0], second)
	assertPlatform(t, result[1], third)
}

func TestPlatformRepositoryFindAfterWithoutCursor(t *testing.T) {
	fixture := newPlatformRepositoryFixture(t)

	first := newTestPlatform(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000001"),
		"Reddit",
		"reddit",
		"https://reddit.com",
	)

	second := newTestPlatform(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000002"),
		"YouTube",
		"youtube",
		"https://youtube.com",
	)

	savePlatforms(t, fixture, first, second)

	result, err := fixture.repository.FindAfter(
		fixture.ctx,
		uuid.Nil,
		1,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 platform, got %d", len(result))
	}

	assertPlatform(t, result[0], first)
}

func TestPlatformRepositoryFindAfterWithCursor(t *testing.T) {
	fixture := newPlatformRepositoryFixture(t)

	first := newTestPlatform(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000001"),
		"Reddit",
		"reddit",
		"https://reddit.com",
	)

	second := newTestPlatform(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000002"),
		"YouTube",
		"youtube",
		"https://youtube.com",
	)

	third := newTestPlatform(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000003"),
		"TikTok",
		"tiktok",
		"https://tiktok.com",
	)

	savePlatforms(t, fixture, first, second, third)

	result, err := fixture.repository.FindAfter(
		fixture.ctx,
		first.Identifier(),
		2,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 platforms, got %d", len(result))
	}

	assertPlatform(t, result[0], second)
	assertPlatform(t, result[1], third)
}

func TestPlatformRepositoryCount(t *testing.T) {
	fixture := newPlatformRepositoryFixture(t)

	savePlatforms(
		t,
		fixture,
		newTestPlatform(
			t,
			mustParseUUID("00000000-0000-0000-0000-000000000001"),
			"Reddit",
			"reddit",
			"https://reddit.com",
		),
		newTestPlatform(
			t,
			mustParseUUID("00000000-0000-0000-0000-000000000002"),
			"YouTube",
			"youtube",
			"https://youtube.com",
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

type platformRepositoryFixture struct {
	ctx        context.Context
	pool       *pgxpool.Pool
	repository domain_platforms.Repository
}

func newPlatformRepositoryFixture(t *testing.T) *platformRepositoryFixture {
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

	fixture := &platformRepositoryFixture{
		ctx:  ctx,
		pool: pool,
		repository: NewPlatformRepository(
			pool,
			domain_platforms.NewAdapter(),
		),
	}

	createPlatformsTable(t, fixture)
	truncatePlatformsTable(t, fixture)

	t.Cleanup(func() {
		truncatePlatformsTable(t, fixture)
		pool.Close()
	})

	return fixture
}

func createPlatformsTable(
	t *testing.T,
	fixture *platformRepositoryFixture,
) {
	t.Helper()

	_, err := fixture.pool.Exec(
		fixture.ctx,
		`
		CREATE TABLE IF NOT EXISTS platforms (
			identifier UUID PRIMARY KEY,
			participation_kind TEXT NOT NULL,
			name TEXT NOT NULL,
			handle TEXT NOT NULL UNIQUE,
			base_url TEXT NOT NULL,
			created_on TIMESTAMPTZ NOT NULL
		)
		`,
	)

	if err != nil {
		t.Fatal(err)
	}
}

func truncatePlatformsTable(
	t *testing.T,
	fixture *platformRepositoryFixture,
) {
	t.Helper()

	_, err := fixture.pool.Exec(
		fixture.ctx,
		`
		TRUNCATE TABLE platforms
		`,
	)

	if err != nil {
		t.Fatal(err)
	}
}

func savePlatforms(
	t *testing.T,
	fixture *platformRepositoryFixture,
	platforms ...domain_platforms.Platform,
) {
	t.Helper()

	for _, platform := range platforms {
		if err := fixture.repository.Save(fixture.ctx, platform); err != nil {
			t.Fatal(err)
		}
	}
}

func assertPlatform(
	t *testing.T,
	result domain_platforms.Platform,
	expected domain_platforms.Platform,
) {
	t.Helper()

	if result == nil {
		t.Fatalf("expected platform")
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

	if result.Name() != expected.Name() {
		t.Fatalf("expected name %q, got %q", expected.Name(), result.Name())
	}

	if result.Handle() != expected.Handle() {
		t.Fatalf("expected handle %q, got %q", expected.Handle(), result.Handle())
	}

	if result.BaseURL() != expected.BaseURL() {
		t.Fatalf("expected base url %q, got %q", expected.BaseURL(), result.BaseURL())
	}
}

func mustParseUUID(value string) uuid.UUID {
	id, err := uuid.Parse(value)
	if err != nil {
		panic(err)
	}

	return id
}

func newTestPlatform(
	t *testing.T,
	id uuid.UUID,
	name string,
	handle string,
	baseURL string,
) domain_platforms.Platform {
	t.Helper()

	platform, err := domain_platforms.NewAdapter().ToDomain(
		domain_platforms.PlatformInput{
			Identifier:        id,
			ParticipationKind: participatables.PlatformKind,
			Name:              name,
			Handle:            handle,
			BaseURL:           baseURL,
			CreatedOn:         time.Now().UTC(),
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	return platform
}
