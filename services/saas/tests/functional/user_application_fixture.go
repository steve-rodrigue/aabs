package functional

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	application_platforms "github.com/steve-rodrigue/aabs/services/saas/applications/entities/platforms"
	application_users "github.com/steve-rodrigue/aabs/services/saas/applications/entities/users"
	"github.com/steve-rodrigue/aabs/services/saas/infrastructure/postgresql"

	domain_platforms "github.com/steve-rodrigue/aabs/services/saas/domain/entities/platforms"
	domain_users "github.com/steve-rodrigue/aabs/services/saas/domain/entities/users"
)

type userApplicationFixture struct {
	Context context.Context
	Pool    *pgxpool.Pool

	PlatformRepository  domain_platforms.Repository
	PlatformApplication application_platforms.Application

	UserRepository  domain_users.Repository
	UserApplication application_users.Application
}

func newUserApplicationFixture(
	t *testing.T,
) *userApplicationFixture {
	t.Helper()

	dsn := os.Getenv("POSTGRES_TEST_DSN")
	if dsn == "" {
		t.Skip("POSTGRES_TEST_DSN not set")
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}

	fixture := &userApplicationFixture{
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

	t.Cleanup(func() {
		fixture.truncateTables(t)
		pool.Close()
	})

	return fixture
}

func (fixture *userApplicationFixture) createTables(
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
}

func (fixture *userApplicationFixture) truncateTables(
	t *testing.T,
) {
	t.Helper()

	_, err := fixture.Pool.Exec(
		fixture.Context,
		`
		TRUNCATE TABLE
			users,
			platforms
		`,
	)
	if err != nil {
		t.Fatal(err)
	}
}
