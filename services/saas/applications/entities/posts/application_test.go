package posts

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/communities"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/platforms"
	domain_posts "github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/users"
)

var errTest = errors.New("test error")

func TestSave(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	post := domain_posts.NewMockPost("hello")

	err := fixture.application.Save(ctx, post)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.SaveCalls != 1 {
		t.Fatalf("expected 1 save call, got %d", fixture.repository.SaveCalls)
	}

	if fixture.repository.LastContext != ctx {
		t.Fatalf("expected context to be passed")
	}

	if fixture.repository.LastSaved != post {
		t.Fatalf("expected post to be passed")
	}
}

func TestSaveReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	fixture.repository.SaveErr = errTest

	err := fixture.application.Save(ctx, domain_posts.NewMockPost("hello"))

	if !errors.Is(err, errTest) {
		t.Fatalf("expected save error, got %v", err)
	}
}

func TestFindByID(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	id := uuid.New()
	post := domain_posts.NewMockPost("hello")

	fixture.repository.Items[id] = post

	result, err := fixture.application.FindByID(ctx, id)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindByIDCalls != 1 {
		t.Fatalf("expected 1 find by id call")
	}

	if fixture.repository.LastContext != ctx {
		t.Fatalf("expected context to be passed")
	}

	if fixture.repository.LastID != id {
		t.Fatalf("expected id to be passed")
	}

	if result != post {
		t.Fatalf("expected post result")
	}
}

func TestFindByIDReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	fixture.repository.FindByIDErr = errTest

	_, err := fixture.application.FindByID(ctx, uuid.New())

	if !errors.Is(err, errTest) {
		t.Fatalf("expected find by id error, got %v", err)
	}
}

func TestFind(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	post := domain_posts.NewMockPost("hello")
	fixture.repository.FindValue = []domain_posts.Post{
		post,
	}

	result, err := fixture.application.Find(ctx, 0, 25)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindCalls != 1 {
		t.Fatalf("expected 1 find call")
	}

	if fixture.repository.LastContext != ctx {
		t.Fatalf("expected context to be passed")
	}

	if fixture.repository.LastIndex != 0 {
		t.Fatalf("expected index 0")
	}

	if fixture.repository.LastAmount != 25 {
		t.Fatalf("expected amount 25")
	}

	if len(result) != 1 || result[0] != post {
		t.Fatalf("expected post result")
	}
}

func TestFindReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	fixture.repository.FindErr = errTest

	_, err := fixture.application.Find(ctx, 0, 25)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected find error, got %v", err)
	}
}

func TestFindAfter(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	cursor := uuid.New()
	post := domain_posts.NewMockPost("hello")

	fixture.repository.FindAfterValue = []domain_posts.Post{
		post,
	}

	result, err := fixture.application.FindAfter(ctx, cursor, 25)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindAfterCalls != 1 {
		t.Fatalf("expected 1 find after call")
	}

	if fixture.repository.LastContext != ctx {
		t.Fatalf("expected context to be passed")
	}

	if fixture.repository.LastCursor != cursor {
		t.Fatalf("expected cursor to be passed")
	}

	if fixture.repository.LastAmount != 25 {
		t.Fatalf("expected amount 25")
	}

	if len(result) != 1 || result[0] != post {
		t.Fatalf("expected post result")
	}
}

func TestFindAfterReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	fixture.repository.FindAfterErr = errTest

	_, err := fixture.application.FindAfter(ctx, uuid.New(), 25)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected find after error, got %v", err)
	}
}

func TestCount(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	fixture.repository.CountValue = 123

	result, err := fixture.application.Count(ctx)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.CountCalls != 1 {
		t.Fatalf("expected 1 count call")
	}

	if fixture.repository.LastContext != ctx {
		t.Fatalf("expected context to be passed")
	}

	if result != 123 {
		t.Fatalf("expected count 123, got %d", result)
	}
}

func TestCountReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	fixture.repository.CountErr = errTest

	_, err := fixture.application.Count(ctx)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected count error, got %v", err)
	}
}

func TestFindByUser(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	user := users.NewMockUser("@user", "User")
	post := domain_posts.NewMockPostWithUser("hello", user)

	fixture.repository.Items[post.Identifier()] = post

	result, err := fixture.application.FindByUser(ctx, user)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindByUserCalls != 1 {
		t.Fatalf("expected 1 find by user call")
	}

	if fixture.repository.LastContext != ctx {
		t.Fatalf("expected context to be passed")
	}

	if fixture.repository.LastUser != user {
		t.Fatalf("expected user to be passed")
	}

	if len(result) != 1 || result[0] != post {
		t.Fatalf("expected post result")
	}
}

func TestFindByUserReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	fixture.repository.FindByUserErr = errTest

	_, err := fixture.application.FindByUser(
		ctx,
		users.NewMockUser("@user", "User"),
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected find by user error, got %v", err)
	}
}

func TestFindByCommunity(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	community := communities.NewMockCommunity("Community", "Text")
	post := domain_posts.NewMockPostWithCommunities(
		"hello",
		[]uuid.UUID{community.Identifier()},
	)

	fixture.repository.Items[post.Identifier()] = post

	result, err := fixture.application.FindByCommunity(ctx, community)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindByCommunityCalls != 1 {
		t.Fatalf("expected 1 find by community call")
	}

	if fixture.repository.LastContext != ctx {
		t.Fatalf("expected context to be passed")
	}

	if fixture.repository.LastCommunity != community {
		t.Fatalf("expected community to be passed")
	}

	if len(result) != 1 || result[0] != post {
		t.Fatalf("expected post result")
	}
}

func TestFindByCommunityReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	fixture.repository.FindByCommunityErr = errTest

	_, err := fixture.application.FindByCommunity(
		ctx,
		communities.NewMockCommunity("Community", "Text"),
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected find by community error, got %v", err)
	}
}

func TestFindByPlatform(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	platform := platforms.NewMockPlatform("Platform", "platform")
	user := users.NewMockUser("@user", "User")
	post := domain_posts.NewMockPostWithUser("hello", user)

	fixture.repository.FindByPlatformValue = []domain_posts.Post{
		post,
	}

	result, err := fixture.application.FindByPlatform(ctx, platform)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindByPlatformCalls != 1 {
		t.Fatalf("expected 1 find by platform call")
	}

	if fixture.repository.LastContext != ctx {
		t.Fatalf("expected context to be passed")
	}

	if fixture.repository.LastPlatform != platform {
		t.Fatalf("expected platform to be passed")
	}

	if len(result) != 1 || result[0] != post {
		t.Fatalf("expected post result")
	}
}

func TestFindByPlatformReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	fixture.repository.FindByPlatformErr = errTest

	_, err := fixture.application.FindByPlatform(
		ctx,
		platforms.NewMockPlatform("Platform", "platform"),
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected find by platform error, got %v", err)
	}
}
