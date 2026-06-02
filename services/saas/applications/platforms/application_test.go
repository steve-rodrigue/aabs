package platforms

import (
	"errors"
	"testing"

	"github.com/google/uuid"

	domain_platforms "github.com/steve-rodrigue/aabs/services/saas/domain/platforms"
)

var errTest = errors.New("test error")

func TestSave(t *testing.T) {
	fixture := newApplicationFixture()
	platform := domain_platforms.NewMockPlatform("Platform", "platform")

	err := fixture.application.Save(platform)

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
		domain_platforms.NewMockPlatform("Platform", "platform"),
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected save error, got %v", err)
	}
}

func TestFindByID(t *testing.T) {
	fixture := newApplicationFixture()

	platform := domain_platforms.NewMockPlatform("Platform", "platform")
	fixture.repository.Items[platform.Identifier()] = platform

	result, err := fixture.application.FindByID(platform.Identifier())

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindByIDCalls != 1 {
		t.Fatalf("expected 1 find by id call")
	}

	if result != platform {
		t.Fatalf("expected platform result")
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

func TestFindByHandle(t *testing.T) {
	fixture := newApplicationFixture()

	platform := domain_platforms.NewMockPlatform("Platform", "platform")
	fixture.repository.Items[platform.Identifier()] = platform

	result, err := fixture.application.FindByHandle("platform")

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindByHandleCalls != 1 {
		t.Fatalf("expected 1 find by handle call")
	}

	if result != platform {
		t.Fatalf("expected platform result")
	}
}

func TestFindByHandleReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.repository.FindByHandleErr = errTest

	_, err := fixture.application.FindByHandle("platform")

	if !errors.Is(err, errTest) {
		t.Fatalf("expected find by handle error, got %v", err)
	}
}

func TestFind(t *testing.T) {
	fixture := newApplicationFixture()

	platform := domain_platforms.NewMockPlatform("Platform", "platform")
	fixture.repository.FindValue = []domain_platforms.Platform{
		platform,
	}

	result, err := fixture.application.Find(0, 25)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindCalls != 1 {
		t.Fatalf("expected 1 find call")
	}

	if len(result) != 1 || result[0] != platform {
		t.Fatalf("expected platform result")
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
	platform := domain_platforms.NewMockPlatform("Platform", "platform")

	fixture.repository.FindAfterValue = []domain_platforms.Platform{
		platform,
	}

	result, err := fixture.application.FindAfter(cursor, 25)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindAfterCalls != 1 {
		t.Fatalf("expected 1 find after call")
	}

	if len(result) != 1 || result[0] != platform {
		t.Fatalf("expected platform result")
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
