package postgresql

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/participatables"
	domain_communities "github.com/steve-rodrigue/aabs/services/saas/domain/entities/communities"
	domain_platforms "github.com/steve-rodrigue/aabs/services/saas/domain/entities/platforms"
	domain_posts "github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts/contents"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts/contents/replies"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts/contents/threads"
	domain_users "github.com/steve-rodrigue/aabs/services/saas/domain/entities/users"
)

func TestNewPostRepository(t *testing.T) {
	pool := &pgxpool.Pool{}
	adapter := domain_posts.NewMockPostAdapter()
	users := domain_users.NewMockUserRepository()

	repository := NewPostRepository(pool, adapter, users)

	if repository == nil {
		t.Fatalf("expected repository")
	}
}

func TestPostRepositorySaveAndFindByIDWithThreadContent(t *testing.T) {
	fixture := newPostRepositoryFixture(t)

	creator := newTestRepositoryUser()
	fixture.users.Items[creator.Identifier()] = creator

	communityID := uuid.New()

	post := newTestThreadPost(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000001"),
		creator,
		[]uuid.UUID{communityID},
		"Thread title",
		"Thread text",
	)

	if err := fixture.repository.Save(fixture.ctx, post); err != nil {
		t.Fatal(err)
	}

	result, err := fixture.repository.FindByID(
		fixture.ctx,
		post.Identifier(),
	)

	if err != nil {
		t.Fatal(err)
	}

	assertPost(t, result, post)

	if !result.Content().IsThread() {
		t.Fatalf("expected thread content")
	}

	if result.Content().Text() != "Thread text" {
		t.Fatalf("expected thread text, got %q", result.Content().Text())
	}
}

func TestPostRepositorySaveAndFindByIDWithReplyContent(t *testing.T) {
	fixture := newPostRepositoryFixture(t)

	creator := newTestRepositoryUser()
	fixture.users.Items[creator.Identifier()] = creator

	post := newTestReplyPost(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000001"),
		creator,
		[]uuid.UUID{uuid.New()},
		"Reply text",
	)

	if err := fixture.repository.Save(fixture.ctx, post); err != nil {
		t.Fatal(err)
	}

	result, err := fixture.repository.FindByID(
		fixture.ctx,
		post.Identifier(),
	)

	if err != nil {
		t.Fatal(err)
	}

	assertPost(t, result, post)

	if !result.Content().IsReply() {
		t.Fatalf("expected reply content")
	}

	if result.Content().Text() != "Reply text" {
		t.Fatalf("expected reply text, got %q", result.Content().Text())
	}
}

func TestPostRepositorySaveUpdatesExistingPost(t *testing.T) {
	fixture := newPostRepositoryFixture(t)

	creator := newTestRepositoryUser()
	fixture.users.Items[creator.Identifier()] = creator

	id := mustParseUUID("00000000-0000-0000-0000-000000000001")

	first := newTestThreadPost(
		t,
		id,
		creator,
		[]uuid.UUID{uuid.New()},
		"Old title",
		"Old text",
	)

	second := newTestThreadPost(
		t,
		id,
		creator,
		[]uuid.UUID{uuid.New(), uuid.New()},
		"Updated title",
		"Updated text",
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

	assertPost(t, result, second)
}

func TestPostRepositoryFindByIDReturnsNilWhenNotFound(t *testing.T) {
	fixture := newPostRepositoryFixture(t)

	result, err := fixture.repository.FindByID(
		fixture.ctx,
		uuid.New(),
	)

	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Fatalf("expected nil post")
	}
}

func TestPostRepositoryFind(t *testing.T) {
	fixture := newPostRepositoryFixture(t)

	creator := newTestRepositoryUser()
	fixture.users.Items[creator.Identifier()] = creator

	first := newTestThreadPost(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000001"),
		creator,
		nil,
		"First",
		"First text",
	)

	second := newTestThreadPost(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000002"),
		creator,
		nil,
		"Second",
		"Second text",
	)

	third := newTestThreadPost(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000003"),
		creator,
		nil,
		"Third",
		"Third text",
	)

	savePosts(t, fixture, first, second, third)

	result, err := fixture.repository.Find(fixture.ctx, 1, 2)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 posts, got %d", len(result))
	}

	assertPost(t, result[0], second)
	assertPost(t, result[1], third)
}

func TestPostRepositoryFindAfterWithoutCursor(t *testing.T) {
	fixture := newPostRepositoryFixture(t)

	creator := newTestRepositoryUser()
	fixture.users.Items[creator.Identifier()] = creator

	first := newTestThreadPost(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000001"),
		creator,
		nil,
		"First",
		"First text",
	)

	second := newTestThreadPost(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000002"),
		creator,
		nil,
		"Second",
		"Second text",
	)

	savePosts(t, fixture, first, second)

	result, err := fixture.repository.FindAfter(
		fixture.ctx,
		uuid.Nil,
		1,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 post, got %d", len(result))
	}

	assertPost(t, result[0], first)
}

func TestPostRepositoryFindAfterWithCursor(t *testing.T) {
	fixture := newPostRepositoryFixture(t)

	creator := newTestRepositoryUser()
	fixture.users.Items[creator.Identifier()] = creator

	first := newTestThreadPost(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000001"),
		creator,
		nil,
		"First",
		"First text",
	)

	second := newTestThreadPost(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000002"),
		creator,
		nil,
		"Second",
		"Second text",
	)

	third := newTestThreadPost(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000003"),
		creator,
		nil,
		"Third",
		"Third text",
	)

	savePosts(t, fixture, first, second, third)

	result, err := fixture.repository.FindAfter(
		fixture.ctx,
		first.Identifier(),
		2,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 posts, got %d", len(result))
	}

	assertPost(t, result[0], second)
	assertPost(t, result[1], third)
}

func TestPostRepositoryCount(t *testing.T) {
	fixture := newPostRepositoryFixture(t)

	creator := newTestRepositoryUser()
	fixture.users.Items[creator.Identifier()] = creator

	savePosts(
		t,
		fixture,
		newTestThreadPost(
			t,
			mustParseUUID("00000000-0000-0000-0000-000000000001"),
			creator,
			nil,
			"First",
			"First text",
		),
		newTestThreadPost(
			t,
			mustParseUUID("00000000-0000-0000-0000-000000000002"),
			creator,
			nil,
			"Second",
			"Second text",
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

func TestPostRepositoryFindByUser(t *testing.T) {
	fixture := newPostRepositoryFixture(t)

	user := newTestRepositoryUser()
	otherUser := newTestRepositoryUser()

	fixture.users.Items[user.Identifier()] = user
	fixture.users.Items[otherUser.Identifier()] = otherUser

	first := newTestThreadPost(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000001"),
		user,
		nil,
		"First",
		"First text",
	)

	second := newTestThreadPost(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000002"),
		otherUser,
		nil,
		"Second",
		"Second text",
	)

	savePosts(t, fixture, first, second)

	result, err := fixture.repository.FindByUser(fixture.ctx, user)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 post, got %d", len(result))
	}

	assertPost(t, result[0], first)
}

func TestPostRepositoryFindByCommunity(t *testing.T) {
	fixture := newPostRepositoryFixture(t)

	creator := newTestRepositoryUser()
	fixture.users.Items[creator.Identifier()] = creator

	community := domain_communities.NewMockCommunity("Community", "Text")
	otherCommunityID := uuid.New()

	first := newTestThreadPost(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000001"),
		creator,
		[]uuid.UUID{community.Identifier()},
		"First",
		"First text",
	)

	second := newTestThreadPost(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000002"),
		creator,
		[]uuid.UUID{otherCommunityID},
		"Second",
		"Second text",
	)

	savePosts(t, fixture, first, second)

	result, err := fixture.repository.FindByCommunity(
		fixture.ctx,
		community,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 post, got %d", len(result))
	}

	assertPost(t, result[0], first)
}

func TestPostRepositoryFindByPlatform(t *testing.T) {
	fixture := newPostRepositoryFixture(t)

	platform := domain_platforms.NewMockPlatform("Reddit", "reddit")
	otherPlatform := domain_platforms.NewMockPlatform("YouTube", "youtube")

	user := newTestRepositoryUserWithPlatform(platform)
	otherUser := newTestRepositoryUserWithPlatform(otherPlatform)

	fixture.users.Items[user.Identifier()] = user
	fixture.users.Items[otherUser.Identifier()] = otherUser

	first := newTestThreadPost(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000001"),
		user,
		nil,
		"First",
		"First text",
	)

	second := newTestThreadPost(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000002"),
		otherUser,
		nil,
		"Second",
		"Second text",
	)

	savePosts(t, fixture, first, second)

	result, err := fixture.repository.FindByPlatform(
		fixture.ctx,
		platform,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 post, got %d", len(result))
	}

	assertPost(t, result[0], first)
}

type postRepositoryFixture struct {
	ctx        context.Context
	pool       *pgxpool.Pool
	adapter    *domain_posts.MockPostAdapter
	users      *domain_users.MockUserRepository
	repository domain_posts.Repository
}

func newPostRepositoryFixture(t *testing.T) *postRepositoryFixture {
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

	adapter := domain_posts.NewMockPostAdapter()
	users := domain_users.NewMockUserRepository()

	fixture := &postRepositoryFixture{
		ctx:        ctx,
		pool:       pool,
		adapter:    adapter,
		users:      users,
		repository: NewPostRepository(pool, adapter, users),
	}

	createPostTables(t, fixture)
	truncatePostTables(t, fixture)

	t.Cleanup(func() {
		truncatePostTables(t, fixture)
		pool.Close()
	})

	return fixture
}

func createPostTables(
	t *testing.T,
	fixture *postRepositoryFixture,
) {
	t.Helper()

	_, err := fixture.pool.Exec(
		fixture.ctx,
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

func truncatePostTables(
	t *testing.T,
	fixture *postRepositoryFixture,
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
			posts,
			users
		`,
	)

	if err != nil {
		t.Fatal(err)
	}
}

func savePosts(
	t *testing.T,
	fixture *postRepositoryFixture,
	posts ...domain_posts.Post,
) {
	t.Helper()

	for _, post := range posts {
		if err := fixture.repository.Save(fixture.ctx, post); err != nil {
			t.Fatal(err)
		}
	}
}

func assertPost(
	t *testing.T,
	result domain_posts.Post,
	expected domain_posts.Post,
) {
	t.Helper()

	if result == nil {
		t.Fatalf("expected post")
	}

	if result.Identifier() != expected.Identifier() {
		t.Fatalf("expected id %s, got %s", expected.Identifier(), result.Identifier())
	}

	if result.Creator() != expected.Creator() {
		t.Fatalf("expected creator")
	}

	if result.Content().Text() != expected.Content().Text() {
		t.Fatalf(
			"expected content text %q, got %q",
			expected.Content().Text(),
			result.Content().Text(),
		)
	}

	expectedCommunities := expected.CommunityIDs()
	resultCommunities := result.CommunityIDs()

	if len(resultCommunities) != len(expectedCommunities) {
		t.Fatalf(
			"expected %d communities, got %d",
			len(expectedCommunities),
			len(resultCommunities),
		)
	}

	for index, expectedCommunityID := range expectedCommunities {
		if resultCommunities[index] != expectedCommunityID {
			t.Fatalf("expected community id at index %d", index)
		}
	}
}

func newTestThreadPost(
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

func newTestReplyPost(
	t *testing.T,
	id uuid.UUID,
	creator domain_users.User,
	communityIDs []uuid.UUID,
	text string,
) domain_posts.Post {
	t.Helper()

	thread := threads.NewMockThread("Target thread", "Target text")

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
				Reply: &replies.ReplyInput{
					Identifier: uuid.New(),
					Target: replies.TargetInput{
						Thread: thread,
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

func newTestRepositoryUser() domain_users.User {
	return newTestRepositoryUserWithPlatform(
		domain_platforms.NewMockPlatform("Platform", "platform"),
	)
}

func newTestRepositoryUserWithPlatform(
	platform domain_platforms.Platform,
) domain_users.User {
	return &testRepositoryUser{
		id:       uuid.New(),
		platform: platform,
	}
}

type testRepositoryUser struct {
	id       uuid.UUID
	platform domain_platforms.Platform
}

func (user *testRepositoryUser) Identifier() uuid.UUID {
	return user.id
}

func (user *testRepositoryUser) ParticipationKind() participatables.Kind {
	return participatables.UserKind
}

func (user *testRepositoryUser) Platform() domain_platforms.Platform {
	return user.platform
}

func (user *testRepositoryUser) ExternalID() string {
	return ""
}

func (user *testRepositoryUser) Handle() string {
	return ""
}

func (user *testRepositoryUser) DisplayName() string {
	return ""
}

func (user *testRepositoryUser) ProfileURL() string {
	return ""
}

func (user *testRepositoryUser) CreatedOn() time.Time {
	return time.Now().UTC()
}
