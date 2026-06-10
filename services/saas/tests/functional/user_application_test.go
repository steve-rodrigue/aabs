package functional

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/participatables"
	domain_platforms "github.com/steve-rodrigue/aabs/services/saas/domain/entities/platforms"
	domain_users "github.com/steve-rodrigue/aabs/services/saas/domain/entities/users"
)

func TestUserApplicationSaveAndFindByID(t *testing.T) {
	fixture := newUserApplicationFixture(t)

	platform := newTestUserPlatform(t, fixture)
	user := newTestUser(t, platform, "123", "steve", "Steve Rodrigue")

	saveUsers(t, fixture, user)

	result, err := fixture.UserApplication.FindByID(
		fixture.Context,
		user.Identifier(),
	)
	if err != nil {
		t.Fatal(err)
	}

	assertUser(t, result, user)
}

func TestUserApplicationSaveUpdatesExistingUser(t *testing.T) {
	fixture := newUserApplicationFixture(t)

	platform := newTestUserPlatform(t, fixture)
	user := newTestUser(t, platform, "123", "steve", "Steve Rodrigue")

	saveUsers(t, fixture, user)

	updated, err := domain_users.NewAdapter().ToDomain(
		domain_users.UserInput{
			Identifier:  user.Identifier(),
			Platform:    platform,
			ExternalID:  "123-updated",
			Handle:      "steve-updated",
			DisplayName: "Steve Updated",
			ProfileURL:  "https://reddit.com/u/steve-updated",
			CreatedOn:   user.CreatedOn(),
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	saveUsers(t, fixture, updated)

	result, err := fixture.UserApplication.FindByID(
		fixture.Context,
		user.Identifier(),
	)
	if err != nil {
		t.Fatal(err)
	}

	assertUser(t, result, updated)
}

func TestUserApplicationFindByIDReturnsNilWhenNotFound(t *testing.T) {
	fixture := newUserApplicationFixture(t)

	result, err := fixture.UserApplication.FindByID(
		fixture.Context,
		uuid.New(),
	)
	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Fatalf("expected nil user")
	}
}

func TestUserApplicationFindByExternalID(t *testing.T) {
	fixture := newUserApplicationFixture(t)

	platform := newTestUserPlatform(t, fixture)
	otherPlatform := newTestPlatformForUser(t, fixture, "X", "x", "https://x.com")

	expected := newTestUser(t, platform, "123", "steve", "Steve Rodrigue")
	other := newTestUser(t, otherPlatform, "123", "steve", "Steve Rodrigue")

	saveUsers(t, fixture, expected, other)

	result, err := fixture.UserApplication.FindByExternalID(
		fixture.Context,
		platform,
		expected.ExternalID(),
	)
	if err != nil {
		t.Fatal(err)
	}

	assertUser(t, result, expected)
}

func TestUserApplicationFindByExternalIDReturnsNilWhenNotFound(t *testing.T) {
	fixture := newUserApplicationFixture(t)

	platform := newTestUserPlatform(t, fixture)

	result, err := fixture.UserApplication.FindByExternalID(
		fixture.Context,
		platform,
		"missing",
	)
	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Fatalf("expected nil user")
	}
}

func TestUserApplicationFindByHandle(t *testing.T) {
	fixture := newUserApplicationFixture(t)

	platform := newTestUserPlatform(t, fixture)
	otherPlatform := newTestPlatformForUser(t, fixture, "X", "x", "https://x.com")

	expected := newTestUser(t, platform, "123", "steve", "Steve Rodrigue")
	other := newTestUser(t, otherPlatform, "456", "steve", "Steve Rodrigue")

	saveUsers(t, fixture, expected, other)

	result, err := fixture.UserApplication.FindByHandle(
		fixture.Context,
		platform,
		expected.Handle(),
	)
	if err != nil {
		t.Fatal(err)
	}

	assertUser(t, result, expected)
}

func TestUserApplicationFindByHandleReturnsNilWhenNotFound(t *testing.T) {
	fixture := newUserApplicationFixture(t)

	platform := newTestUserPlatform(t, fixture)

	result, err := fixture.UserApplication.FindByHandle(
		fixture.Context,
		platform,
		"missing",
	)
	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Fatalf("expected nil user")
	}
}

func TestUserApplicationFind(t *testing.T) {
	fixture := newUserApplicationFixture(t)

	platform := newTestUserPlatform(t, fixture)

	first := newTestUser(t, platform, "1", "first", "First User")
	second := newTestUser(t, platform, "2", "second", "Second User")
	third := newTestUser(t, platform, "3", "third", "Third User")

	saveUsers(t, fixture, first, second, third)

	result, err := fixture.UserApplication.Find(
		fixture.Context,
		0,
		2,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 users, got %d", len(result))
	}
}

func TestUserApplicationFindReturnsEmptyWhenIndexIsOutOfRange(t *testing.T) {
	fixture := newUserApplicationFixture(t)

	platform := newTestUserPlatform(t, fixture)
	user := newTestUser(t, platform, "1", "first", "First User")

	saveUsers(t, fixture, user)

	result, err := fixture.UserApplication.Find(
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

func TestUserApplicationFindAfter(t *testing.T) {
	fixture := newUserApplicationFixture(t)

	platform := newTestUserPlatform(t, fixture)

	first := newTestUser(t, platform, "1", "first", "First User")
	second := newTestUser(t, platform, "2", "second", "Second User")
	third := newTestUser(t, platform, "3", "third", "Third User")

	saveUsers(t, fixture, first, second, third)

	all, err := fixture.UserApplication.Find(
		fixture.Context,
		0,
		10,
	)
	if err != nil {
		t.Fatal(err)
	}

	result, err := fixture.UserApplication.FindAfter(
		fixture.Context,
		all[0].Identifier(),
		10,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 users, got %d", len(result))
	}

	if result[0].Identifier() != all[1].Identifier() {
		t.Fatalf("expected second user first")
	}
}

func TestUserApplicationFindAfterWithNilCursor(t *testing.T) {
	fixture := newUserApplicationFixture(t)

	platform := newTestUserPlatform(t, fixture)

	first := newTestUser(t, platform, "1", "first", "First User")
	second := newTestUser(t, platform, "2", "second", "Second User")

	saveUsers(t, fixture, first, second)

	result, err := fixture.UserApplication.FindAfter(
		fixture.Context,
		uuid.Nil,
		10,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 users, got %d", len(result))
	}
}

func TestUserApplicationCount(t *testing.T) {
	fixture := newUserApplicationFixture(t)

	platform := newTestUserPlatform(t, fixture)

	first := newTestUser(t, platform, "1", "first", "First User")
	second := newTestUser(t, platform, "2", "second", "Second User")

	saveUsers(t, fixture, first, second)

	count, err := fixture.UserApplication.Count(
		fixture.Context,
	)
	if err != nil {
		t.Fatal(err)
	}

	if count != 2 {
		t.Fatalf("expected count 2, got %d", count)
	}
}

func newTestUserPlatform(
	t *testing.T,
	fixture *userApplicationFixture,
) domain_platforms.Platform {
	t.Helper()

	return newTestPlatformForUser(
		t,
		fixture,
		"Reddit",
		"reddit",
		"https://reddit.com",
	)
}

func newTestPlatformForUser(
	t *testing.T,
	fixture *userApplicationFixture,
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

	err = fixture.PlatformApplication.Save(
		fixture.Context,
		platform,
	)
	if err != nil {
		t.Fatal(err)
	}

	return platform
}

func newTestUser(
	t *testing.T,
	platform domain_platforms.Platform,
	externalID string,
	handle string,
	displayName string,
) domain_users.User {
	t.Helper()

	user, err := domain_users.NewAdapter().ToDomain(
		domain_users.UserInput{
			Identifier:  uuid.New(),
			Platform:    platform,
			ExternalID:  externalID,
			Handle:      handle,
			DisplayName: displayName,
			ProfileURL:  "https://reddit.com/u/" + handle,
			CreatedOn:   time.Now().UTC(),
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	return user
}

func saveUsers(
	t *testing.T,
	fixture *userApplicationFixture,
	users ...domain_users.User,
) {
	t.Helper()

	for _, user := range users {
		err := fixture.UserApplication.Save(
			fixture.Context,
			user,
		)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func assertUser(
	t *testing.T,
	result domain_users.User,
	expected domain_users.User,
) {
	t.Helper()

	if result == nil {
		t.Fatalf("expected user")
	}

	if result.Identifier() != expected.Identifier() {
		t.Fatalf("expected identifier")
	}

	if result.ParticipationKind() != expected.ParticipationKind() {
		t.Fatalf("expected participation kind")
	}

	if result.Platform() == nil {
		t.Fatalf("expected platform")
	}

	if result.Platform().Identifier() != expected.Platform().Identifier() {
		t.Fatalf("expected platform identifier")
	}

	if result.ExternalID() != expected.ExternalID() {
		t.Fatalf("expected external id %s, got %s", expected.ExternalID(), result.ExternalID())
	}

	if result.Handle() != expected.Handle() {
		t.Fatalf("expected handle %s, got %s", expected.Handle(), result.Handle())
	}

	if result.DisplayName() != expected.DisplayName() {
		t.Fatalf("expected display name %s, got %s", expected.DisplayName(), result.DisplayName())
	}

	if result.ProfileURL() != expected.ProfileURL() {
		t.Fatalf("expected profile url %s, got %s", expected.ProfileURL(), result.ProfileURL())
	}
}
