package functional

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	application_communities "github.com/steve-rodrigue/aabs/services/saas/applications/entities/communities"
	application_platforms "github.com/steve-rodrigue/aabs/services/saas/applications/entities/platforms"
	application_users "github.com/steve-rodrigue/aabs/services/saas/applications/entities/users"
	"github.com/steve-rodrigue/aabs/services/saas/infrastructure/postgresql"

	domain_communities "github.com/steve-rodrigue/aabs/services/saas/domain/entities/communities"
	domain_platforms "github.com/steve-rodrigue/aabs/services/saas/domain/entities/platforms"
	domain_users "github.com/steve-rodrigue/aabs/services/saas/domain/entities/users"
)

type communityApplicationFixture struct {
	Context context.Context
	Pool    *pgxpool.Pool

	PlatformRepository  domain_platforms.Repository
	PlatformApplication application_platforms.Application

	UserRepository  domain_users.Repository
	UserApplication application_users.Application

	CommunityRepository  domain_communities.Repository
	CommunityApplication application_communities.Application
}

func newCommunityApplicationFixture(
	t *testing.T,
) *communityApplicationFixture {
	t.Helper()

	dsn := os.Getenv("POSTGRES_TEST_DSN")
	if dsn == "" {
		t.Skip("POSTGRES_TEST_DSN not set")
	}

	ctx := context.Background()

	pool, err := pgxpool.New(
		ctx,
		dsn,
	)
	if err != nil {
		t.Fatal(err)
	}

	fixture := &communityApplicationFixture{
		Context: ctx,
		Pool:    pool,
	}

	fixture.createTables(t)
	fixture.truncateTables(t)

	fixture.PlatformRepository = postgresql.NewPlatformRepository(
		pool,
		domain_platforms.NewAdapter(),
	)

	fixture.PlatformApplication = application_platforms.New(
		fixture.PlatformRepository,
	)

	fixture.UserRepository = postgresql.NewUserRepository(
		pool,
		domain_users.NewAdapter(),
		fixture.PlatformRepository,
	)

	fixture.UserApplication = application_users.New(
		fixture.UserRepository,
	)

	fixture.CommunityRepository = postgresql.NewCommunityRepository(
		pool,
		domain_communities.NewAdapter(),
		fixture.PlatformRepository,
		fixture.UserRepository,
	)

	fixture.CommunityApplication = application_communities.New(
		fixture.CommunityRepository,
	)

	t.Cleanup(func() {
		fixture.truncateTables(t)
		pool.Close()
	})

	return fixture
}
func (fixture *communityApplicationFixture) createTables(
	t *testing.T,
) {
	t.Helper()

	_, err := fixture.Pool.Exec(
		fixture.Context,
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

	_, err = fixture.Pool.Exec(
		fixture.Context,
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

	_, err = fixture.Pool.Exec(
		fixture.Context,
		`
		CREATE TABLE IF NOT EXISTS communities (
			identifier UUID PRIMARY KEY,
			platform_id UUID NOT NULL,
			handle TEXT NOT NULL,
			title TEXT NOT NULL,
			text TEXT NOT NULL,
			created_on TIMESTAMPTZ NOT NULL,
			UNIQUE (platform_id, handle)
		)
		`,
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = fixture.Pool.Exec(
		fixture.Context,
		`
		CREATE TABLE IF NOT EXISTS community_moderators (
			community_id UUID NOT NULL,
			user_id UUID NOT NULL,
			PRIMARY KEY (community_id, user_id)
		)
		`,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func (fixture *communityApplicationFixture) truncateTables(
	t *testing.T,
) {
	t.Helper()

	_, err := fixture.Pool.Exec(
		fixture.Context,
		`
		TRUNCATE TABLE
			communities_moderators,
			communities,
			users,
			platforms
		`,
	)
	if err != nil {
		t.Fatal(err)
	}
}
