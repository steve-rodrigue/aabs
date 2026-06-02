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

func TestFindAll(t *testing.T) {
	fixture := newApplicationFixture()

	platform := domain_platforms.NewMockPlatform("Platform", "platform")
	fixture.repository.Items[platform.Identifier()] = platform

	result, err := fixture.application.FindAll()

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindAllCalls != 1 {
		t.Fatalf("expected 1 find all call")
	}

	if len(result) != 1 || result[0] != platform {
		t.Fatalf("expected platform result")
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
