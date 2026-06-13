package postgresql

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	domain_posts "github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts/contents"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts/contents/replies"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts/contents/threads"
	domain_users "github.com/steve-rodrigue/aabs/services/posts/restapis/domain/users"
)

func TestNewPostService(t *testing.T) {
	pool := &pgxpool.Pool{}

	service := NewPostService(pool)

	if service == nil {
		t.Fatalf("expected service")
	}
}

func TestPostServiceSaveThreadPost(t *testing.T) {
	fixture := newPostServiceFixture(t)

	creator := domain_users.NewMockUser("steve", "Steve")
	post := newPostServiceThreadPost(
		t,
		creator,
		[]uuid.UUID{uuid.New()},
		"Thread title",
		"Thread text",
	)

	err := fixture.service.Save(fixture.ctx, post)
	if err != nil {
		t.Fatal(err)
	}

	var count int

	err = fixture.pool.QueryRow(
		fixture.ctx,
		`
		SELECT COUNT(*)
		FROM posts
		WHERE identifier = $1
		AND creator_id = $2
		AND content_id = $3
		`,
		post.Identifier(),
		post.Creator().Identifier(),
		post.Content().Identifier(),
	).Scan(&count)
	if err != nil {
		t.Fatal(err)
	}

	if count != 1 {
		t.Fatalf("expected post to be saved")
	}

	err = fixture.pool.QueryRow(
		fixture.ctx,
		`
		SELECT COUNT(*)
		FROM post_contents
		WHERE identifier = $1
		AND kind = 'thread'
		`,
		post.Content().Identifier(),
	).Scan(&count)
	if err != nil {
		t.Fatal(err)
	}

	if count != 1 {
		t.Fatalf("expected thread content to be saved")
	}

	err = fixture.pool.QueryRow(
		fixture.ctx,
		`
		SELECT COUNT(*)
		FROM post_content_threads
		WHERE content_id = $1
		AND title = $2
		AND text = $3
		`,
		post.Content().Identifier(),
		"Thread title",
		"Thread text",
	).Scan(&count)
	if err != nil {
		t.Fatal(err)
	}

	if count != 1 {
		t.Fatalf("expected thread to be saved")
	}
}

func TestPostServiceSaveReplyPost(t *testing.T) {
	fixture := newPostServiceFixture(t)

	creator := domain_users.NewMockUser("steve", "Steve")

	target := newPostServiceThreadPost(
		t,
		creator,
		nil,
		"Target title",
		"Target text",
	)

	err := fixture.service.Save(fixture.ctx, target)
	if err != nil {
		t.Fatal(err)
	}

	reply := newPostServiceReplyPost(
		t,
		creator,
		[]uuid.UUID{uuid.New()},
		target.Content().Thread(),
		"Reply text",
	)

	err = fixture.service.Save(fixture.ctx, reply)
	if err != nil {
		t.Fatal(err)
	}

	var count int

	err = fixture.pool.QueryRow(
		fixture.ctx,
		`
		SELECT COUNT(*)
		FROM post_contents
		WHERE identifier = $1
		AND kind = 'reply'
		`,
		reply.Content().Identifier(),
	).Scan(&count)
	if err != nil {
		t.Fatal(err)
	}

	if count != 1 {
		t.Fatalf("expected reply content to be saved")
	}

	err = fixture.pool.QueryRow(
		fixture.ctx,
		`
		SELECT COUNT(*)
		FROM post_content_replies
		WHERE content_id = $1
		AND target_thread_id = $2
		AND text = $3
		`,
		reply.Content().Identifier(),
		target.Content().Thread().Identifier(),
		"Reply text",
	).Scan(&count)
	if err != nil {
		t.Fatal(err)
	}

	if count != 1 {
		t.Fatalf("expected reply to be saved")
	}
}

func TestPostServiceSaveUpdatesExistingPost(t *testing.T) {
	fixture := newPostServiceFixture(t)

	creator := domain_users.NewMockUser("steve", "Steve")
	id := uuid.New()

	firstCommunity := uuid.New()
	secondCommunity := uuid.New()

	first := newPostServiceThreadPostWithID(
		t,
		id,
		creator,
		[]uuid.UUID{firstCommunity},
		"Old title",
		"Old text",
	)

	second := newPostServiceThreadPostWithID(
		t,
		id,
		creator,
		[]uuid.UUID{secondCommunity},
		"New title",
		"New text",
	)

	err := fixture.service.Save(fixture.ctx, first)
	if err != nil {
		t.Fatal(err)
	}

	err = fixture.service.Save(fixture.ctx, second)
	if err != nil {
		t.Fatal(err)
	}

	var count int

	err = fixture.pool.QueryRow(
		fixture.ctx,
		`
		SELECT COUNT(*)
		FROM post_communities
		WHERE post_id = $1
		AND community_id = $2
		`,
		id,
		firstCommunity,
	).Scan(&count)
	if err != nil {
		t.Fatal(err)
	}

	if count != 0 {
		t.Fatalf("expected old community to be removed")
	}

	err = fixture.pool.QueryRow(
		fixture.ctx,
		`
		SELECT COUNT(*)
		FROM post_communities
		WHERE post_id = $1
		AND community_id = $2
		`,
		id,
		secondCommunity,
	).Scan(&count)
	if err != nil {
		t.Fatal(err)
	}

	if count != 1 {
		t.Fatalf("expected new community to be saved")
	}
}

type postServiceFixture struct {
	ctx     context.Context
	pool    *pgxpool.Pool
	service domain_posts.Service
}

func newPostServiceFixture(t *testing.T) *postServiceFixture {
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

	fixture := &postServiceFixture{
		ctx:     ctx,
		pool:    pool,
		service: NewPostService(pool),
	}

	createPostServiceTables(t, fixture)
	truncatePostServiceTables(t, fixture)

	t.Cleanup(func() {
		truncatePostServiceTables(t, fixture)
		pool.Close()
	})

	return fixture
}

func createPostServiceTables(
	t *testing.T,
	fixture *postServiceFixture,
) {
	t.Helper()

	_, err := fixture.pool.Exec(
		fixture.ctx,
		`
			DROP TABLE IF EXISTS post_content_replies CASCADE;
			DROP TABLE IF EXISTS post_content_threads CASCADE;
			DROP TABLE IF EXISTS post_contents CASCADE;
			DROP TABLE IF EXISTS post_communities CASCADE;
			DROP TABLE IF EXISTS posts CASCADE;
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

func truncatePostServiceTables(
	t *testing.T,
	fixture *postServiceFixture,
) {
	t.Helper()

	_, err := fixture.pool.Exec(
		fixture.ctx,
		`
		TRUNCATE TABLE
			post_content_replies,
			post_content_threads,
			post_contents,
			post_communities,
			posts
		`,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func newPostServiceThreadPost(
	t *testing.T,
	creator domain_users.User,
	communityIDs []uuid.UUID,
	title string,
	text string,
) domain_posts.Post {
	t.Helper()

	return newPostServiceThreadPostWithID(
		t,
		uuid.New(),
		creator,
		communityIDs,
		title,
		text,
	)
}

func newPostServiceThreadPostWithID(
	t *testing.T,
	id uuid.UUID,
	creator domain_users.User,
	communityIDs []uuid.UUID,
	title string,
	text string,
) domain_posts.Post {
	t.Helper()

	contentAdapter := contents.NewAdapter(
		replies.NewAdapter(),
		threads.NewAdapter(),
	)

	post, err := domain_posts.NewAdapter(contentAdapter).ToDomain(
		domain_posts.PostInput{
			Identifier:   id,
			CommunityIDs: communityIDs,
			Creator:      creator,
			Content: contents.ContentInput{
				Identifier: uuid.New(),
				Thread: &threads.ThreadInput{
					Identifier: uuid.New(),
					Creator:    creator,
					Title:      title,
					Text:       text,
				},
				CreatedAt: time.Now().UTC(),
			},
			CreatedOn: time.Now().UTC(),
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	return post
}

func newPostServiceReplyPost(
	t *testing.T,
	creator domain_users.User,
	communityIDs []uuid.UUID,
	target threads.Thread,
	text string,
) domain_posts.Post {
	t.Helper()

	contentAdapter := contents.NewAdapter(
		replies.NewAdapter(),
		threads.NewAdapter(),
	)

	post, err := domain_posts.NewAdapter(contentAdapter).ToDomain(
		domain_posts.PostInput{
			Identifier:   uuid.New(),
			CommunityIDs: communityIDs,
			Creator:      creator,
			Content: contents.ContentInput{
				Identifier: uuid.New(),
				Reply: &replies.ReplyInput{
					Identifier: uuid.New(),
					Target: replies.TargetInput{
						Thread: target,
					},
					Text: text,
				},
				CreatedAt: time.Now().UTC(),
			},
			CreatedOn: time.Now().UTC(),
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	return post
}
