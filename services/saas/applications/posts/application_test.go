package posts

import (
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/communities"
	"github.com/steve-rodrigue/aabs/services/saas/domain/platforms"
	domain_posts "github.com/steve-rodrigue/aabs/services/saas/domain/posts"
	"github.com/steve-rodrigue/aabs/services/saas/domain/users"
)

var errTest = errors.New("test error")

func TestSave(t *testing.T) {
	fixture := newApplicationFixture()
	post := domain_posts.NewMockPost("hello")

	err := fixture.application.Save(post)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.SaveCalls != 1 {
		t.Fatalf("expected 1 save call, got %d", fixture.repository.SaveCalls)
	}
}

func TestSaveReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.repository.SaveErr = errTest

	err := fixture.application.Save(domain_posts.NewMockPost("hello"))

	if !errors.Is(err, errTest) {
		t.Fatalf("expected save error, got %v", err)
	}
}

func TestFindByID(t *testing.T) {
	fixture := newApplicationFixture()

	id := uuid.New()
	post := domain_posts.NewMockPost("hello")

	fixture.repository.Items[id] = post

	result, err := fixture.application.FindByID(id)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindByIDCalls != 1 {
		t.Fatalf("expected 1 find by id call")
	}

	if result != post {
		t.Fatalf("expected post result")
	}
}

func TestFindByIDReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.repository.FindByIDErr = errTest

	_, err := fixture.application.FindByID(uuid.New())

	if !errors.Is(err, errTest) {
		t.Fatalf("expected find by id error, got %v", err)
	}
}

func TestFindAll(t *testing.T) {
	fixture := newApplicationFixture()

	post := domain_posts.NewMockPost("hello")
	fixture.repository.Items[post.Identifier()] = post

	result, err := fixture.application.FindAll()

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindAllCalls != 1 {
		t.Fatalf("expected 1 find all call")
	}

	if len(result) != 1 || result[0] != post {
		t.Fatalf("expected post result")
	}
}

func TestFindAllReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.repository.FindAllErr = errTest

	_, err := fixture.application.FindAll()

	if !errors.Is(err, errTest) {
		t.Fatalf("expected find all error, got %v", err)
	}
}

func TestFindByUser(t *testing.T) {
	fixture := newApplicationFixture()

	user := users.NewMockUser("@user", "User")
	post := domain_posts.NewMockPostWithUser("hello", user)

	fixture.repository.Items[post.Identifier()] = post

	result, err := fixture.application.FindByUser(user)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindByUserCalls != 1 {
		t.Fatalf("expected 1 find by user call")
	}

	if len(result) != 1 || result[0] != post {
		t.Fatalf("expected post result")
	}
}

func TestFindByUserReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.repository.FindByUserErr = errTest

	_, err := fixture.application.FindByUser(
		users.NewMockUser("@user", "User"),
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected find by user error, got %v", err)
	}
}

func TestFindByCommunity(t *testing.T) {
	fixture := newApplicationFixture()

	community := communities.NewMockCommunity("Community", "Text")
	post := domain_posts.NewMockPostWithCommunities(
		"hello",
		[]uuid.UUID{community.Identifier()},
	)

	fixture.repository.Items[post.Identifier()] = post

	result, err := fixture.application.FindByCommunity(community)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindByCommunityCalls != 1 {
		t.Fatalf("expected 1 find by community call")
	}

	if len(result) != 1 || result[0] != post {
		t.Fatalf("expected post result")
	}
}

func TestFindByCommunityReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.repository.FindByCommunityErr = errTest

	_, err := fixture.application.FindByCommunity(
		communities.NewMockCommunity("Community", "Text"),
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected find by community error, got %v", err)
	}
}

func TestFindByPlatform(t *testing.T) {
	fixture := newApplicationFixture()

	platform := platforms.NewMockPlatform("Platform", "platform")
	user := users.NewMockUser("@user", "User")
	post := domain_posts.NewMockPostWithUser("hello", user)

	fixture.repository.FindByPlatformValue = []domain_posts.Post{
		post,
	}

	result, err := fixture.application.FindByPlatform(platform)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindByPlatformCalls != 1 {
		t.Fatalf("expected 1 find by platform call")
	}

	if len(result) != 1 || result[0] != post {
		t.Fatalf("expected post result")
	}
}

func TestFindByPlatformReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.repository.FindByPlatformErr = errTest

	_, err := fixture.application.FindByPlatform(
		platforms.NewMockPlatform("Platform", "platform"),
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected find by platform error, got %v", err)
	}
}
