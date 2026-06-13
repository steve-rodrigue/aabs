package functional

import (
	"context"
	"fmt"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	hatchetsdk "github.com/hatchet-dev/hatchet/sdks/go"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/applications"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/clients"
	infrastructure_applications "github.com/steve-rodrigue/aabs/services/posts/restapis/infrastructure/applications"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/servers"
)

type applicationFixture struct {
	Context context.Context
	Pool    *pgxpool.Pool
	Server  *httptest.Server

	Application applications.Application
}

func newApplicationFixture(
	t *testing.T,
) *applicationFixture {
	t.Helper()

	dsn := os.Getenv("POSTS_POSTGRES_TEST_DSN")
	if dsn == "" {
		t.Skip("POSTS_POSTGRES_TEST_DSN not set")
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}

	fixture := &applicationFixture{
		Context: ctx,
		Pool:    pool,
	}

	fixture.resetSchema(t)

	hatchetClient := newTestHatchetClient(t)
	serverApplication := infrastructure_applications.New(
		pool,
		hatchetClient,
	)

	fixture.Server = httptest.NewServer(
		servers.New(serverApplication),
	)

	fixture.Application = clients.New(
		fixture.Server.URL,
		fixture.Server.Client(),
	)

	t.Cleanup(func() {
		fixture.Server.Close()
		fixture.truncateTables(t)
		pool.Close()
	})

	return fixture
}

func newTestHatchetClient(t *testing.T) *hatchetsdk.Client {
	t.Helper()

	token := os.Getenv("HATCHET_TEST_CLIENT_TOKEN")
	fmt.Printf("voila: %s", token)
	hostPort := os.Getenv("HATCHET_TEST_HOST_PORT")
	tlsStrategy := os.Getenv("HATCHET_TEST_TLS_STRATEGY")

	if token == "" {
		t.Skip("HATCHET_TEST_CLIENT_TOKEN not set")
	}

	if hostPort == "" {
		t.Skip("HATCHET_TEST_HOST_PORT not set")
	}

	if tlsStrategy == "" {
		tlsStrategy = "none"
	}

	t.Setenv("HATCHET_CLIENT_TOKEN", token)
	t.Setenv("HATCHET_CLIENT_HOST_PORT", hostPort)
	t.Setenv("HATCHET_CLIENT_TLS_STRATEGY", tlsStrategy)

	client, err := hatchetsdk.NewClient()
	if err != nil {
		t.Fatal(err)
	}

	return client
}

func (fixture *applicationFixture) createTables(
	t *testing.T,
) {
	t.Helper()

	if err := infrastructure_applications.Install(
		fixture.Context,
		fixture.Pool,
	); err != nil {
		t.Fatal(err)
	}
}

func (fixture *applicationFixture) resetSchema(
	t *testing.T,
) {
	t.Helper()

	_, err := fixture.Pool.Exec(
		fixture.Context,
		`
		DROP TABLE IF EXISTS
			post_content_replies,
			post_content_threads,
			post_communities,
			posts,
			post_contents,
			community_moderators,
			communities,
			users,
			platforms
		CASCADE
		`,
	)
	if err != nil {
		t.Fatal(err)
	}

	if err := infrastructure_applications.Install(
		fixture.Context,
		fixture.Pool,
	); err != nil {
		t.Fatal(err)
	}
}

func (fixture *applicationFixture) truncateTables(
	t *testing.T,
) {
	t.Helper()

	_, err := fixture.Pool.Exec(
		fixture.Context,
		`
		TRUNCATE TABLE
			post_content_replies,
			post_content_threads,
			post_communities,
			community_moderators,
			posts,
			post_contents,
			communities,
			users,
			platforms
		CASCADE
		`,
	)

	if err != nil {
		t.Fatal(err)
	}
}
