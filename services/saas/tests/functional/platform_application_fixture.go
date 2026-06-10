package functional

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	application_platforms "github.com/steve-rodrigue/aabs/services/saas/applications/entities/platforms"
	"github.com/steve-rodrigue/aabs/services/saas/infrastructure/postgresql"

	domain_platforms "github.com/steve-rodrigue/aabs/services/saas/domain/entities/platforms"
)

type platformApplicationFixture struct {
	Context context.Context
	Pool    *pgxpool.Pool

	PlatformRepository  domain_platforms.Repository
	PlatformApplication application_platforms.Application
}

func newPlatformApplicationFixture(
	t *testing.T,
) *platformApplicationFixture {
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

	fixture := &platformApplicationFixture{
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

	t.Cleanup(func() {
		fixture.truncateTables(t)
		pool.Close()
	})

	return fixture
}

func (fixture *platformApplicationFixture) createTables(
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
}

func (fixture *platformApplicationFixture) truncateTables(
	t *testing.T,
) {
	t.Helper()

	_, err := fixture.Pool.Exec(
		fixture.Context,
		`
		TRUNCATE TABLE platforms
		`,
	)
	if err != nil {
		t.Fatal(err)
	}
}
