package functional

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	application_posts "github.com/steve-rodrigue/aabs/services/saas/applications/entities/posts"
	"github.com/steve-rodrigue/aabs/services/saas/infrastructure/postgresql"

	domain_posts "github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts/contents"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts/contents/replies"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts/contents/threads"
	domain_users "github.com/steve-rodrigue/aabs/services/saas/domain/entities/users"
)

type postApplicationFixture struct {
	Context context.Context
	Pool    *pgxpool.Pool

	UserRepository *domain_users.MockUserRepository

	PostRepository  domain_posts.Repository
	PostApplication application_posts.Application
}

func newPostApplicationFixture(
	t *testing.T,
) *postApplicationFixture {
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

	fixture := &postApplicationFixture{
		Context:        ctx,
		Pool:           pool,
		UserRepository: domain_users.NewMockUserRepository(),
	}

	fixture.createTables(t)
	fixture.truncateTables(t)

	contentAdapter := contents.NewAdapter(
		replies.NewAdapter(),
		threads.NewAdapter(),
	)

	fixture.PostRepository = postgresql.NewPostRepository(
		pool,
		domain_posts.NewAdapter(contentAdapter),
		fixture.UserRepository,
	)

	fixture.PostApplication = application_posts.New(
		fixture.PostRepository,
	)

	t.Cleanup(func() {
		fixture.truncateTables(t)
		pool.Close()
	})

	return fixture
}

func (fixture *postApplicationFixture) createTables(
	t *testing.T,
) {
	t.Helper()

	_, err := fixture.Pool.Exec(
		fixture.Context,
		`
		CREATE TABLE IF NOT EXISTS users (
			identifier UUID PRIMARY KEY,
			platform_id UUID NULL
		);

		CREATE TABLE IF NOT EXISTS posts (
			identifier UUID PRIMARY KEY,
			creator_id UUID NOT NULL,
			content_id UUID NOT NULL,
			created_on TIMESTAMPTZ NOT NULL
		);

		CREATE TABLE IF NOT EXISTS post_communities (
			post_id UUID NOT NULL,
			community_id UUID NOT NULL,
			PRIMARY KEY (post_id, community_id)
		);

		CREATE TABLE IF NOT EXISTS post_contents (
			identifier UUID PRIMARY KEY,
			kind TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL
		);

		CREATE TABLE IF NOT EXISTS post_content_threads (
			content_id UUID PRIMARY KEY,
			identifier UUID NOT NULL UNIQUE,
			creator_id UUID NOT NULL,
			title TEXT NOT NULL,
			text TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS post_content_replies (
			content_id UUID PRIMARY KEY,
			identifier UUID NOT NULL UNIQUE,
			target_reply_id UUID NULL,
			target_thread_id UUID NULL,
			text TEXT NOT NULL
		);
		`,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func (fixture *postApplicationFixture) truncateTables(
	t *testing.T,
) {
	t.Helper()

	_, err := fixture.Pool.Exec(
		fixture.Context,
		`
		TRUNCATE TABLE
			post_content_replies,
			post_content_threads,
			post_contents,
			post_communities,
			posts,
			users
		`,
	)
	if err != nil {
		t.Fatal(err)
	}
}
