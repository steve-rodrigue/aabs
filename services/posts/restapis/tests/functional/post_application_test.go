package functional

import (
	"testing"
	"time"

	"github.com/google/uuid"

	domain_communities "github.com/steve-rodrigue/aabs/services/posts/restapis/domain/communities"
	domain_platforms "github.com/steve-rodrigue/aabs/services/posts/restapis/domain/platforms"
	domain_posts "github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts/contents"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts/contents/replies"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts/contents/threads"
	domain_users "github.com/steve-rodrigue/aabs/services/posts/restapis/domain/users"
)

func TestPostApplicationSaveAndFindByIDWithThread(t *testing.T) {
	fixture := newApplicationFixture(t)

	creator := newPostApplicationUser(t, fixture)
	communityID := newPostApplicationCommunityID(t, fixture, creator.Platform(), "thread-community")

	post := newPostApplicationThreadPost(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000001"),
		creator,
		[]uuid.UUID{communityID},
		"Thread title",
		"Thread text",
	)

	saveApplicationPosts(t, fixture, post)

	result, err := fixture.Application.Posts().FindByID(
		fixture.Context,
		post.Identifier(),
	)
	if err != nil {
		t.Fatal(err)
	}

	assertApplicationPost(t, result, post)

	if !result.Content().IsThread() {
		t.Fatalf("expected thread content")
	}

	if result.Content().Text() != "Thread text" {
		t.Fatalf("expected thread text")
	}
}

func TestPostApplicationSaveAndFindByIDWithReply(t *testing.T) {
	fixture := newApplicationFixture(t)

	creator := newPostApplicationUser(t, fixture)

	parent := newPostApplicationThreadPost(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000001"),
		creator,
		nil,
		"Parent title",
		"Parent text",
	)

	saveApplicationPosts(t, fixture, parent)

	communityID := newPostApplicationCommunityID(t, fixture, creator.Platform(), "reply-community")

	reply := newPostApplicationReplyPostWithTargetThread(
		t,
		mustParseUUID("00000000-0000-0000-0000-000000000002"),
		creator,
		[]uuid.UUID{communityID},
		parent.Content().Thread(),
		"Reply text",
	)

	saveApplicationPosts(t, fixture, reply)

	result, err := fixture.Application.Posts().FindByID(
		fixture.Context,
		reply.Identifier(),
	)
	if err != nil {
		t.Fatal(err)
	}

	assertApplicationPost(t, result, reply)

	if !result.Content().IsReply() {
		t.Fatalf("expected reply content")
	}

	if result.Content().Text() != "Reply text" {
		t.Fatalf("expected reply text")
	}
}

func TestPostApplicationSaveUpdatesExistingPost(t *testing.T) {
	fixture := newApplicationFixture(t)

	creator := newPostApplicationUser(t, fixture)

	firstCommunityID := newPostApplicationCommunityID(t, fixture, creator.Platform(), "old-community")
	secondCommunityID := newPostApplicationCommunityID(t, fixture, creator.Platform(), "new-community")
	thirdCommunityID := newPostApplicationCommunityID(t, fixture, creator.Platform(), "another-community")

	id := mustParseUUID("00000000-0000-0000-0000-000000000001")

	first := newPostApplicationThreadPost(
		t,
		id,
		creator,
		[]uuid.UUID{firstCommunityID},
		"Old title",
		"Old text",
	)

	second := newPostApplicationThreadPost(
		t,
		id,
		creator,
		[]uuid.UUID{secondCommunityID, thirdCommunityID},
		"Updated title",
		"Updated text",
	)

	saveApplicationPosts(t, fixture, first, second)

	result, err := fixture.Application.Posts().FindByID(
		fixture.Context,
		id,
	)
	if err != nil {
		t.Fatal(err)
	}

	assertApplicationPost(t, result, second)
}

func TestPostApplicationFindByIDReturnsNilWhenNotFound(t *testing.T) {
	fixture := newApplicationFixture(t)

	result, err := fixture.Application.Posts().FindByID(
		fixture.Context,
		uuid.New(),
	)
	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Fatalf("expected nil post")
	}
}

func TestPostApplicationFind(t *testing.T) {
	fixture := newApplicationFixture(t)

	creator := newPostApplicationUser(t, fixture)

	first := newPostApplicationThreadPost(t, mustParseUUID("00000000-0000-0000-0000-000000000001"), creator, nil, "First", "First text")
	second := newPostApplicationThreadPost(t, mustParseUUID("00000000-0000-0000-0000-000000000002"), creator, nil, "Second", "Second text")
	third := newPostApplicationThreadPost(t, mustParseUUID("00000000-0000-0000-0000-000000000003"), creator, nil, "Third", "Third text")

	saveApplicationPosts(t, fixture, first, second, third)

	result, err := fixture.Application.Posts().Find(
		fixture.Context,
		1,
		2,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 posts, got %d", len(result))
	}

	assertApplicationPost(t, result[0], second)
	assertApplicationPost(t, result[1], third)
}

func TestPostApplicationFindAfterWithNilCursor(t *testing.T) {
	fixture := newApplicationFixture(t)

	creator := newPostApplicationUser(t, fixture)

	first := newPostApplicationThreadPost(t, mustParseUUID("00000000-0000-0000-0000-000000000001"), creator, nil, "First", "First text")
	second := newPostApplicationThreadPost(t, mustParseUUID("00000000-0000-0000-0000-000000000002"), creator, nil, "Second", "Second text")

	saveApplicationPosts(t, fixture, first, second)

	result, err := fixture.Application.Posts().FindAfter(
		fixture.Context,
		uuid.Nil,
		1,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 post, got %d", len(result))
	}

	assertApplicationPost(t, result[0], first)
}

func TestPostApplicationFindAfterWithCursor(t *testing.T) {
	fixture := newApplicationFixture(t)

	creator := newPostApplicationUser(t, fixture)

	first := newPostApplicationThreadPost(t, mustParseUUID("00000000-0000-0000-0000-000000000001"), creator, nil, "First", "First text")
	second := newPostApplicationThreadPost(t, mustParseUUID("00000000-0000-0000-0000-000000000002"), creator, nil, "Second", "Second text")
	third := newPostApplicationThreadPost(t, mustParseUUID("00000000-0000-0000-0000-000000000003"), creator, nil, "Third", "Third text")

	saveApplicationPosts(t, fixture, first, second, third)

	result, err := fixture.Application.Posts().FindAfter(
		fixture.Context,
		first.Identifier(),
		2,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 posts, got %d", len(result))
	}

	assertApplicationPost(t, result[0], second)
	assertApplicationPost(t, result[1], third)
}

func TestPostApplicationCount(t *testing.T) {
	fixture := newApplicationFixture(t)

	creator := newPostApplicationUser(t, fixture)

	first := newPostApplicationThreadPost(t, mustParseUUID("00000000-0000-0000-0000-000000000001"), creator, nil, "First", "First text")
	second := newPostApplicationThreadPost(t, mustParseUUID("00000000-0000-0000-0000-000000000002"), creator, nil, "Second", "Second text")

	saveApplicationPosts(t, fixture, first, second)

	count, err := fixture.Application.Posts().Count(
		fixture.Context,
	)
	if err != nil {
		t.Fatal(err)
	}

	if count != 2 {
		t.Fatalf("expected count 2, got %d", count)
	}
}

func TestPostApplicationFindByCriteriaWithUserIDs(t *testing.T) {
	fixture := newApplicationFixture(t)

	user := newPostApplicationUser(t, fixture)
	otherUser := newPostApplicationUser(t, fixture)

	first := newPostApplicationThreadPost(t, mustParseUUID("00000000-0000-0000-0000-000000000001"), user, nil, "First", "First text")
	second := newPostApplicationThreadPost(t, mustParseUUID("00000000-0000-0000-0000-000000000002"), otherUser, nil, "Second", "Second text")

	saveApplicationPosts(t, fixture, first, second)

	result, err := fixture.Application.Posts().FindByCriteria(
		fixture.Context,
		domain_posts.Criteria{
			UserIDs: []uuid.UUID{user.Identifier()},
		},
		0,
		10,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 post, got %d", len(result))
	}

	assertApplicationPost(t, result[0], first)
}

func TestPostApplicationFindByCriteriaWithCommunityIDs(t *testing.T) {
	fixture := newApplicationFixture(t)

	creator := newPostApplicationUser(t, fixture)

	communityID := newPostApplicationCommunityID(t, fixture, creator.Platform(), "first-community")
	otherCommunityID := newPostApplicationCommunityID(t, fixture, creator.Platform(), "second-community")

	first := newPostApplicationThreadPost(t, mustParseUUID("00000000-0000-0000-0000-000000000001"), creator, []uuid.UUID{communityID}, "First", "First text")
	second := newPostApplicationThreadPost(t, mustParseUUID("00000000-0000-0000-0000-000000000002"), creator, []uuid.UUID{otherCommunityID}, "Second", "Second text")

	saveApplicationPosts(t, fixture, first, second)

	result, err := fixture.Application.Posts().FindByCriteria(
		fixture.Context,
		domain_posts.Criteria{
			CommunityIDs: []uuid.UUID{communityID},
		},
		0,
		10,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 post, got %d", len(result))
	}

	assertApplicationPost(t, result[0], first)
}

func TestPostApplicationFindByCriteriaWithPlatformIDs(t *testing.T) {
	fixture := newApplicationFixture(t)

	platform := newPostApplicationPlatform(
		t,
		fixture,
		"Reddit",
		"reddit",
		"https://reddit.com",
	)

	otherPlatform := newPostApplicationPlatform(
		t,
		fixture,
		"YouTube",
		"youtube",
		"https://youtube.com",
	)

	user := newPostApplicationUserWithPlatform(t, fixture, platform)
	otherUser := newPostApplicationUserWithPlatform(t, fixture, otherPlatform)

	first := newPostApplicationThreadPost(t, mustParseUUID("00000000-0000-0000-0000-000000000001"), user, nil, "First", "First text")
	second := newPostApplicationThreadPost(t, mustParseUUID("00000000-0000-0000-0000-000000000002"), otherUser, nil, "Second", "Second text")

	saveApplicationPosts(t, fixture, first, second)

	result, err := fixture.Application.Posts().FindByCriteria(
		fixture.Context,
		domain_posts.Criteria{
			PlatformIDs: []uuid.UUID{platform.Identifier()},
		},
		0,
		10,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 post, got %d", len(result))
	}

	assertApplicationPost(t, result[0], first)
}

func TestPostApplicationFindByCriteriaAfter(t *testing.T) {
	fixture := newApplicationFixture(t)

	user := newPostApplicationUser(t, fixture)

	first := newPostApplicationThreadPost(t, mustParseUUID("00000000-0000-0000-0000-000000000001"), user, nil, "First", "First text")
	second := newPostApplicationThreadPost(t, mustParseUUID("00000000-0000-0000-0000-000000000002"), user, nil, "Second", "Second text")
	third := newPostApplicationThreadPost(t, mustParseUUID("00000000-0000-0000-0000-000000000003"), user, nil, "Third", "Third text")

	saveApplicationPosts(t, fixture, first, second, third)

	result, err := fixture.Application.Posts().FindByCriteriaAfter(
		fixture.Context,
		domain_posts.Criteria{
			UserIDs: []uuid.UUID{user.Identifier()},
		},
		first.Identifier(),
		2,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 posts, got %d", len(result))
	}

	assertApplicationPost(t, result[0], second)
	assertApplicationPost(t, result[1], third)
}

func TestPostApplicationCountByCriteria(t *testing.T) {
	fixture := newApplicationFixture(t)

	user := newPostApplicationUser(t, fixture)
	otherUser := newPostApplicationUser(t, fixture)

	first := newPostApplicationThreadPost(t, mustParseUUID("00000000-0000-0000-0000-000000000001"), user, nil, "First", "First text")
	second := newPostApplicationThreadPost(t, mustParseUUID("00000000-0000-0000-0000-000000000002"), user, nil, "Second", "Second text")
	third := newPostApplicationThreadPost(t, mustParseUUID("00000000-0000-0000-0000-000000000003"), otherUser, nil, "Third", "Third text")

	saveApplicationPosts(t, fixture, first, second, third)

	count, err := fixture.Application.Posts().CountByCriteria(
		fixture.Context,
		domain_posts.Criteria{
			UserIDs: []uuid.UUID{user.Identifier()},
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	if count != 2 {
		t.Fatalf("expected count 2, got %d", count)
	}
}

func saveApplicationPosts(
	t *testing.T,
	fixture *applicationFixture,
	posts ...domain_posts.Post,
) {
	t.Helper()

	for _, post := range posts {
		err := fixture.Application.Posts().Save(
			fixture.Context,
			post,
		)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func newPostApplicationPlatform(
	t *testing.T,
	fixture *applicationFixture,
	name string,
	handle string,
	baseURL string,
) domain_platforms.Platform {
	t.Helper()

	platform, err := domain_platforms.NewAdapter().ToDomain(
		domain_platforms.PlatformInput{
			Identifier: uuid.New(),
			Name:       name,
			Handle:     handle,
			BaseURL:    baseURL,
			CreatedOn:  time.Now().UTC(),
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	err = fixture.Application.Platforms().Save(fixture.Context, platform)
	if err != nil {
		t.Fatal(err)
	}

	return platform
}

func newPostApplicationUser(
	t *testing.T,
	fixture *applicationFixture,
) domain_users.User {
	t.Helper()

	return newPostApplicationUserWithPlatform(
		t,
		fixture,
		newPostApplicationPlatform(
			t,
			fixture,
			"Platform",
			"platform-"+uuid.NewString(),
			"https://platform.test",
		),
	)
}

func newPostApplicationUserWithPlatform(
	t *testing.T,
	fixture *applicationFixture,
	platform domain_platforms.Platform,
) domain_users.User {
	t.Helper()

	user, err := domain_users.NewAdapter().ToDomain(
		domain_users.UserInput{
			Identifier:  uuid.New(),
			Platform:    platform,
			ExternalID:  uuid.NewString(),
			Handle:      "user-" + uuid.NewString(),
			DisplayName: "Test User",
			ProfileURL:  "https://platform.test/user",
			CreatedOn:   time.Now().UTC(),
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	err = fixture.Application.Users().Save(fixture.Context, user)
	if err != nil {
		t.Fatal(err)
	}

	return user
}

func newPostApplicationCommunityID(
	t *testing.T,
	fixture *applicationFixture,
	platform domain_platforms.Platform,
	handle string,
) uuid.UUID {
	t.Helper()

	community, err := domain_communities.NewAdapter().ToDomain(
		domain_communities.CommunityInput{
			Identifier: uuid.New(),
			Platform:   platform,
			Handle:     handle + "-" + uuid.NewString(),
			Title:      "Test Community",
			Text:       "Test community text",
			CreatedOn:  time.Now().UTC(),
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	err = fixture.Application.Communities().Save(fixture.Context, community)
	if err != nil {
		t.Fatal(err)
	}

	return community.Identifier()
}

func newPostApplicationThreadPost(
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

func newPostApplicationReplyPostWithTargetThread(
	t *testing.T,
	id uuid.UUID,
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
			Identifier:   id,
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

func assertApplicationPost(
	t *testing.T,
	result domain_posts.Post,
	expected domain_posts.Post,
) {
	t.Helper()

	if result == nil {
		t.Fatalf("expected post")
	}

	if result.Identifier() != expected.Identifier() {
		t.Fatalf("expected id")
	}

	if result.Creator() == nil ||
		result.Creator().Identifier() != expected.Creator().Identifier() {
		t.Fatalf("expected creator")
	}

	if result.Content().Text() != expected.Content().Text() {
		t.Fatalf("expected text %q, got %q", expected.Content().Text(), result.Content().Text())
	}

	expectedCommunities := expected.CommunityIDs()
	resultCommunities := result.CommunityIDs()

	if len(resultCommunities) != len(expectedCommunities) {
		t.Fatalf("expected %d communities, got %d", len(expectedCommunities), len(resultCommunities))
	}

	for index, expectedCommunityID := range expectedCommunities {
		if resultCommunities[index] != expectedCommunityID {
			t.Fatalf("expected community id at index %d", index)
		}
	}
}
