package communities

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	domain_communities "github.com/steve-rodrigue/aabs/services/posts/restapis/domain/communities"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/platforms"
)

var errTest = errors.New("test error")

func TestSave(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	community := domain_communities.NewMockCommunity("Community", "Text")

	err := fixture.application.Save(ctx, community)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.SaveCalls != 1 {
		t.Fatalf("expected 1 save call, got %d", fixture.repository.SaveCalls)
	}

	if fixture.repository.LastContext != ctx {
		t.Fatalf("expected context to be passed")
	}

	if fixture.repository.LastSaved != community {
		t.Fatalf("expected community to be passed")
	}
}

func TestSaveReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	fixture.repository.SaveErr = errTest

	err := fixture.application.Save(
		ctx,
		domain_communities.NewMockCommunity("Community", "Text"),
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected save error, got %v", err)
	}
}

func TestFindByID(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	community := domain_communities.NewMockCommunity("Community", "Text")
	fixture.repository.Items[community.Identifier()] = community

	result, err := fixture.application.FindByID(ctx, community.Identifier())

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindByIDCalls != 1 {
		t.Fatalf("expected 1 find by id call")
	}

	if fixture.repository.LastContext != ctx {
		t.Fatalf("expected context to be passed")
	}

	if fixture.repository.LastID != community.Identifier() {
		t.Fatalf("expected id to be passed")
	}

	if result != community {
		t.Fatalf("expected community result")
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

func TestFindByHandle(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	platform := platforms.NewMockPlatform("Platform", "platform")
	community := domain_communities.NewMockCommunity("Community", "Text")

	fixture.repository.FindByHandleValue = community

	result, err := fixture.application.FindByHandle(ctx, platform, "community")

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindByHandleCalls != 1 {
		t.Fatalf("expected 1 find by handle call")
	}

	if fixture.repository.LastContext != ctx {
		t.Fatalf("expected context to be passed")
	}

	if fixture.repository.LastPlatform != platform {
		t.Fatalf("expected platform to be passed")
	}

	if fixture.repository.LastHandle != "community" {
		t.Fatalf("expected handle to be passed")
	}

	if result != community {
		t.Fatalf("expected community result")
	}
}

func TestFindByHandleReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	fixture.repository.FindByHandleErr = errTest

	_, err := fixture.application.FindByHandle(
		ctx,
		platforms.NewMockPlatform("Platform", "platform"),
		"community",
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected find by handle error, got %v", err)
	}
}

func TestFind(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	community := domain_communities.NewMockCommunity("Community", "Text")
	fixture.repository.FindValue = []domain_communities.Community{
		community,
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

	if len(result) != 1 || result[0] != community {
		t.Fatalf("expected community result")
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
	community := domain_communities.NewMockCommunity("Community", "Text")

	fixture.repository.FindAfterValue = []domain_communities.Community{
		community,
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

	if len(result) != 1 || result[0] != community {
		t.Fatalf("expected community result")
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

func TestFindByPlatform(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	platform := platforms.NewMockPlatform("Platform", "platform")
	community := domain_communities.NewMockCommunity("Community", "Text")

	fixture.repository.FindByPlatformValue = []domain_communities.Community{
		community,
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

	if len(result) != 1 || result[0] != community {
		t.Fatalf("expected community result")
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
