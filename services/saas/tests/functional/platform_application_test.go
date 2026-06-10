package functional

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/participatables"
	domain_platforms "github.com/steve-rodrigue/aabs/services/saas/domain/entities/platforms"
)

func TestPlatformApplicationSaveAndFindByID(t *testing.T) {
	fixture := newPlatformApplicationFixture(t)

	platform := newTestPlatform(t, "Reddit", "reddit", "https://reddit.com")

	err := fixture.PlatformApplication.Save(
		fixture.Context,
		platform,
	)
	if err != nil {
		t.Fatal(err)
	}

	result, err := fixture.PlatformApplication.FindByID(
		fixture.Context,
		platform.Identifier(),
	)
	if err != nil {
		t.Fatal(err)
	}

	assertPlatform(t, result, platform)
}

func TestPlatformApplicationSaveUpdatesExistingPlatform(t *testing.T) {
	fixture := newPlatformApplicationFixture(t)

	platform := newTestPlatform(t, "Reddit", "reddit", "https://reddit.com")

	err := fixture.PlatformApplication.Save(
		fixture.Context,
		platform,
	)
	if err != nil {
		t.Fatal(err)
	}

	updated, err := domain_platforms.NewAdapter().ToDomain(
		domain_platforms.PlatformInput{
			Identifier:        platform.Identifier(),
			ParticipationKind: participatables.PlatformKind,
			Name:              "Reddit Updated",
			Handle:            "reddit-updated",
			BaseURL:           "https://www.reddit.com",
			CreatedOn:         platform.CreatedOn(),
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	err = fixture.PlatformApplication.Save(
		fixture.Context,
		updated,
	)
	if err != nil {
		t.Fatal(err)
	}

	result, err := fixture.PlatformApplication.FindByID(
		fixture.Context,
		platform.Identifier(),
	)
	if err != nil {
		t.Fatal(err)
	}

	assertPlatform(t, result, updated)
}

func TestPlatformApplicationFindByIDReturnsNilWhenNotFound(t *testing.T) {
	fixture := newPlatformApplicationFixture(t)

	result, err := fixture.PlatformApplication.FindByID(
		fixture.Context,
		uuid.New(),
	)
	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Fatalf("expected nil platform")
	}
}

func TestPlatformApplicationFindByHandle(t *testing.T) {
	fixture := newPlatformApplicationFixture(t)

	platform := newTestPlatform(t, "Reddit", "reddit", "https://reddit.com")

	err := fixture.PlatformApplication.Save(
		fixture.Context,
		platform,
	)
	if err != nil {
		t.Fatal(err)
	}

	result, err := fixture.PlatformApplication.FindByHandle(
		fixture.Context,
		platform.Handle(),
	)
	if err != nil {
		t.Fatal(err)
	}

	assertPlatform(t, result, platform)
}

func TestPlatformApplicationFindByHandleReturnsNilWhenNotFound(t *testing.T) {
	fixture := newPlatformApplicationFixture(t)

	result, err := fixture.PlatformApplication.FindByHandle(
		fixture.Context,
		"missing",
	)
	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Fatalf("expected nil platform")
	}
}

func TestPlatformApplicationFind(t *testing.T) {
	fixture := newPlatformApplicationFixture(t)

	first := newTestPlatform(t, "Reddit", "reddit", "https://reddit.com")
	second := newTestPlatform(t, "X", "x", "https://x.com")
	third := newTestPlatform(t, "Facebook", "facebook", "https://facebook.com")

	savePlatforms(t, fixture, first, second, third)

	result, err := fixture.PlatformApplication.Find(
		fixture.Context,
		0,
		2,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 platforms, got %d", len(result))
	}
}

func TestPlatformApplicationFindReturnsEmptyWhenIndexIsOutOfRange(t *testing.T) {
	fixture := newPlatformApplicationFixture(t)

	platform := newTestPlatform(t, "Reddit", "reddit", "https://reddit.com")

	savePlatforms(t, fixture, platform)

	result, err := fixture.PlatformApplication.Find(
		fixture.Context,
		10,
		2,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 0 {
		t.Fatalf("expected empty result, got %d", len(result))
	}
}

func TestPlatformApplicationFindAfter(t *testing.T) {
	fixture := newPlatformApplicationFixture(t)

	first := newTestPlatform(t, "Reddit", "reddit", "https://reddit.com")
	second := newTestPlatform(t, "X", "x", "https://x.com")
	third := newTestPlatform(t, "Facebook", "facebook", "https://facebook.com")

	savePlatforms(t, fixture, first, second, third)

	all, err := fixture.PlatformApplication.Find(
		fixture.Context,
		0,
		10,
	)
	if err != nil {
		t.Fatal(err)
	}

	result, err := fixture.PlatformApplication.FindAfter(
		fixture.Context,
		all[0].Identifier(),
		10,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 platforms, got %d", len(result))
	}

	if result[0].Identifier() != all[1].Identifier() {
		t.Fatalf("expected second platform first")
	}
}

func TestPlatformApplicationFindAfterWithNilCursor(t *testing.T) {
	fixture := newPlatformApplicationFixture(t)

	first := newTestPlatform(t, "Reddit", "reddit", "https://reddit.com")
	second := newTestPlatform(t, "X", "x", "https://x.com")

	savePlatforms(t, fixture, first, second)

	result, err := fixture.PlatformApplication.FindAfter(
		fixture.Context,
		uuid.Nil,
		10,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 platforms, got %d", len(result))
	}
}

func TestPlatformApplicationCount(t *testing.T) {
	fixture := newPlatformApplicationFixture(t)

	first := newTestPlatform(t, "Reddit", "reddit", "https://reddit.com")
	second := newTestPlatform(t, "X", "x", "https://x.com")

	savePlatforms(t, fixture, first, second)

	count, err := fixture.PlatformApplication.Count(
		fixture.Context,
	)
	if err != nil {
		t.Fatal(err)
	}

	if count != 2 {
		t.Fatalf("expected count 2, got %d", count)
	}
}

func savePlatforms(
	t *testing.T,
	fixture *platformApplicationFixture,
	platforms ...domain_platforms.Platform,
) {
	t.Helper()

	for _, platform := range platforms {
		err := fixture.PlatformApplication.Save(
			fixture.Context,
			platform,
		)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func newTestPlatform(
	t *testing.T,
	name string,
	handle string,
	baseURL string,
) domain_platforms.Platform {
	t.Helper()

	platform, err := domain_platforms.NewAdapter().ToDomain(
		domain_platforms.PlatformInput{
			Identifier:        uuid.New(),
			ParticipationKind: participatables.PlatformKind,
			Name:              name,
			Handle:            handle,
			BaseURL:           baseURL,
			CreatedOn:         time.Now().UTC(),
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	return platform
}

func assertPlatform(
	t *testing.T,
	result domain_platforms.Platform,
	expected domain_platforms.Platform,
) {
	t.Helper()

	if result == nil {
		t.Fatalf("expected platform")
	}

	if result.Identifier() != expected.Identifier() {
		t.Fatalf("expected identifier")
	}

	if result.ParticipationKind() != expected.ParticipationKind() {
		t.Fatalf("expected participation kind")
	}

	if result.Name() != expected.Name() {
		t.Fatalf("expected name %s, got %s", expected.Name(), result.Name())
	}

	if result.Handle() != expected.Handle() {
		t.Fatalf("expected handle %s, got %s", expected.Handle(), result.Handle())
	}

	if result.BaseURL() != expected.BaseURL() {
		t.Fatalf("expected base url %s, got %s", expected.BaseURL(), result.BaseURL())
	}
}
