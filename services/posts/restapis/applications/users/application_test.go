package users

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/platforms"
	domain_users "github.com/steve-rodrigue/aabs/services/posts/restapis/domain/users"
)

var errTest = errors.New("test error")

func TestSave(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	user := domain_users.NewMockUser("@user", "User")

	err := fixture.application.Save(ctx, user)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.SaveCalls != 1 {
		t.Fatalf("expected 1 save call, got %d", fixture.repository.SaveCalls)
	}

	if fixture.repository.LastContext != ctx {
		t.Fatalf("expected context to be passed")
	}

	if fixture.repository.LastSaved != user {
		t.Fatalf("expected user to be passed")
	}
}

func TestSaveReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	fixture.repository.SaveErr = errTest

	err := fixture.application.Save(
		ctx,
		domain_users.NewMockUser("@user", "User"),
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected save error, got %v", err)
	}
}

func TestFindByID(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	user := domain_users.NewMockUser("@user", "User")
	fixture.repository.Items[user.Identifier()] = user

	result, err := fixture.application.FindByID(ctx, user.Identifier())

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindByIDCalls != 1 {
		t.Fatalf("expected 1 find by id call")
	}

	if fixture.repository.LastContext != ctx {
		t.Fatalf("expected context to be passed")
	}

	if fixture.repository.LastID != user.Identifier() {
		t.Fatalf("expected id to be passed")
	}

	if result != user {
		t.Fatalf("expected user result")
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

func TestFindByExternalID(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	platform := platforms.NewMockPlatform("Platform", "platform")
	user := domain_users.NewMockUser("@user", "User")

	fixture.repository.FindByPlatformAndExternalIDValue = user

	result, err := fixture.application.FindByExternalID(
		ctx,
		platform,
		"external-id",
	)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindByPlatformAndExternalIDCalls != 1 {
		t.Fatalf("expected 1 external id lookup")
	}

	if fixture.repository.LastContext != ctx {
		t.Fatalf("expected context to be passed")
	}

	if fixture.repository.LastPlatform != platform {
		t.Fatalf("expected platform to be passed")
	}

	if fixture.repository.LastExternalID != "external-id" {
		t.Fatalf("expected external id to be passed")
	}

	if result != user {
		t.Fatalf("expected user result")
	}
}

func TestFindByExternalIDReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	fixture.repository.FindByPlatformAndExternalIDErr = errTest

	_, err := fixture.application.FindByExternalID(
		ctx,
		platforms.NewMockPlatform("Platform", "platform"),
		"external-id",
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected external id error, got %v", err)
	}
}

func TestFindByHandle(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	platform := platforms.NewMockPlatform("Platform", "platform")
	user := domain_users.NewMockUser("@user", "User")

	fixture.repository.FindByPlatformAndHandleValue = user

	result, err := fixture.application.FindByHandle(
		ctx,
		platform,
		"@user",
	)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindByPlatformAndHandleCalls != 1 {
		t.Fatalf("expected 1 handle lookup")
	}

	if fixture.repository.LastContext != ctx {
		t.Fatalf("expected context to be passed")
	}

	if fixture.repository.LastPlatform != platform {
		t.Fatalf("expected platform to be passed")
	}

	if fixture.repository.LastHandle != "@user" {
		t.Fatalf("expected handle to be passed")
	}

	if result != user {
		t.Fatalf("expected user result")
	}
}

func TestFindByHandleReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	fixture.repository.FindByPlatformAndHandleErr = errTest

	_, err := fixture.application.FindByHandle(
		ctx,
		platforms.NewMockPlatform("Platform", "platform"),
		"@user",
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected handle error, got %v", err)
	}
}

func TestFind(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	user := domain_users.NewMockUser("@user", "User")
	fixture.repository.FindValue = []domain_users.User{
		user,
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

	if len(result) != 1 || result[0] != user {
		t.Fatalf("expected user result")
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
	user := domain_users.NewMockUser("@user", "User")

	fixture.repository.FindAfterValue = []domain_users.User{
		user,
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

	if len(result) != 1 || result[0] != user {
		t.Fatalf("expected user result")
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
