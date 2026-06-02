package users

import (
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/platforms"
	domain_users "github.com/steve-rodrigue/aabs/services/saas/domain/users"
)

var errTest = errors.New("test error")

func TestSave(t *testing.T) {
	fixture := newApplicationFixture()
	user := domain_users.NewMockUser("@user", "User")

	err := fixture.application.Save(user)

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

	err := fixture.application.Save(
		domain_users.NewMockUser("@user", "User"),
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected save error, got %v", err)
	}
}

func TestFindByID(t *testing.T) {
	fixture := newApplicationFixture()

	user := domain_users.NewMockUser("@user", "User")
	fixture.repository.Items[user.Identifier()] = user

	result, err := fixture.application.FindByID(user.Identifier())

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindByIDCalls != 1 {
		t.Fatalf("expected 1 find by id call")
	}

	if result != user {
		t.Fatalf("expected user result")
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

func TestFindByExternalID(t *testing.T) {
	fixture := newApplicationFixture()

	platform := platforms.NewMockPlatform("Platform", "platform")
	user := domain_users.NewMockUser("@user", "User")

	fixture.repository.FindByPlatformAndExternalIDValue = user

	result, err := fixture.application.FindByExternalID(platform, "external-id")

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindByPlatformAndExternalIDCalls != 1 {
		t.Fatalf("expected 1 external id lookup")
	}

	if result != user {
		t.Fatalf("expected user result")
	}
}

func TestFindByExternalIDReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.repository.FindByPlatformAndExternalIDErr = errTest

	_, err := fixture.application.FindByExternalID(
		platforms.NewMockPlatform("Platform", "platform"),
		"external-id",
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected external id error, got %v", err)
	}
}

func TestFind(t *testing.T) {
	fixture := newApplicationFixture()

	user := domain_users.NewMockUser("@user", "User")
	fixture.repository.FindValue = []domain_users.User{
		user,
	}

	result, err := fixture.application.Find(0, 25)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindCalls != 1 {
		t.Fatalf("expected 1 find call")
	}

	if len(result) != 1 || result[0] != user {
		t.Fatalf("expected user result")
	}
}

func TestFindReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.repository.FindErr = errTest

	_, err := fixture.application.Find(0, 25)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected find error, got %v", err)
	}
}

func TestFindAfter(t *testing.T) {
	fixture := newApplicationFixture()

	cursor := uuid.New()
	user := domain_users.NewMockUser("@user", "User")

	fixture.repository.FindAfterValue = []domain_users.User{
		user,
	}

	result, err := fixture.application.FindAfter(cursor, 25)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindAfterCalls != 1 {
		t.Fatalf("expected 1 find after call")
	}

	if len(result) != 1 || result[0] != user {
		t.Fatalf("expected user result")
	}
}

func TestFindAfterReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.repository.FindAfterErr = errTest

	_, err := fixture.application.FindAfter(uuid.New(), 25)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected find after error, got %v", err)
	}
}

func TestCount(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.repository.CountValue = 123

	result, err := fixture.application.Count()

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.CountCalls != 1 {
		t.Fatalf("expected 1 count call")
	}

	if result != 123 {
		t.Fatalf("expected count 123, got %d", result)
	}
}

func TestCountReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.repository.CountErr = errTest

	_, err := fixture.application.Count()

	if !errors.Is(err, errTest) {
		t.Fatalf("expected count error, got %v", err)
	}
}
