package functional

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/participatables"
	domain_communities "github.com/steve-rodrigue/aabs/services/saas/domain/entities/communities"
	domain_platforms "github.com/steve-rodrigue/aabs/services/saas/domain/entities/platforms"
	domain_users "github.com/steve-rodrigue/aabs/services/saas/domain/entities/users"
)

func TestCommunityApplicationSaveAndFindByID(t *testing.T) {
	fixture := newCommunityApplicationFixture(t)

	platform := newTestCommunityPlatform(t, fixture)
	community := newTestCommunity(t, platform, "golang", "Golang", "Go community", nil)

	saveCommunities(t, fixture, community)

	result, err := fixture.CommunityApplication.FindByID(
		fixture.Context,
		community.Identifier(),
	)
	if err != nil {
		t.Fatal(err)
	}

	assertCommunity(t, result, community)
}

func TestCommunityApplicationSaveAndFindByIDWithModerators(t *testing.T) {
	fixture := newCommunityApplicationFixture(t)

	platform := newTestCommunityPlatform(t, fixture)

	moderator := newTestCommunityUser(
		t,
		fixture,
		platform,
		"mod-1",
		"moderator",
		"Moderator",
	)

	community := newTestCommunity(
		t,
		platform,
		"golang",
		"Golang",
		"Go community",
		[]domain_users.User{moderator},
	)

	saveCommunities(t, fixture, community)

	result, err := fixture.CommunityApplication.FindByID(
		fixture.Context,
		community.Identifier(),
	)
	if err != nil {
		t.Fatal(err)
	}

	assertCommunity(t, result, community)
}

func TestCommunityApplicationSaveUpdatesExistingCommunity(t *testing.T) {
	fixture := newCommunityApplicationFixture(t)

	platform := newTestCommunityPlatform(t, fixture)
	community := newTestCommunity(t, platform, "golang", "Golang", "Go community", nil)

	saveCommunities(t, fixture, community)

	updated, err := domain_communities.NewAdapter().ToDomain(
		domain_communities.CommunityInput{
			Identifier: community.Identifier(),
			Platform:   platform,
			Handle:     "golang-updated",
			Title:      "Golang Updated",
			Text:       "Updated Go community",
			CreatedOn:  community.CreatedOn(),
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	saveCommunities(t, fixture, updated)

	result, err := fixture.CommunityApplication.FindByID(
		fixture.Context,
		community.Identifier(),
	)
	if err != nil {
		t.Fatal(err)
	}

	assertCommunity(t, result, updated)
}

func TestCommunityApplicationFindByIDReturnsNilWhenNotFound(t *testing.T) {
	fixture := newCommunityApplicationFixture(t)

	result, err := fixture.CommunityApplication.FindByID(
		fixture.Context,
		uuid.New(),
	)
	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Fatalf("expected nil community")
	}
}

func TestCommunityApplicationFindByHandle(t *testing.T) {
	fixture := newCommunityApplicationFixture(t)

	platform := newTestCommunityPlatform(t, fixture)
	otherPlatform := newTestCommunityPlatformWithValues(
		t,
		fixture,
		"X",
		"x",
		"https://x.com",
	)

	expected := newTestCommunity(t, platform, "golang", "Golang", "Go community", nil)
	other := newTestCommunity(t, otherPlatform, "golang", "Golang Other", "Other", nil)

	saveCommunities(t, fixture, expected, other)

	result, err := fixture.CommunityApplication.FindByHandle(
		fixture.Context,
		platform,
		expected.Handle(),
	)
	if err != nil {
		t.Fatal(err)
	}

	assertCommunity(t, result, expected)
}

func TestCommunityApplicationFindByHandleReturnsNilWhenNotFound(t *testing.T) {
	fixture := newCommunityApplicationFixture(t)

	platform := newTestCommunityPlatform(t, fixture)

	result, err := fixture.CommunityApplication.FindByHandle(
		fixture.Context,
		platform,
		"missing",
	)
	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Fatalf("expected nil community")
	}
}

func TestCommunityApplicationFind(t *testing.T) {
	fixture := newCommunityApplicationFixture(t)

	platform := newTestCommunityPlatform(t, fixture)

	first := newTestCommunity(t, platform, "first", "First", "First community", nil)
	second := newTestCommunity(t, platform, "second", "Second", "Second community", nil)
	third := newTestCommunity(t, platform, "third", "Third", "Third community", nil)

	saveCommunities(t, fixture, first, second, third)

	result, err := fixture.CommunityApplication.Find(
		fixture.Context,
		0,
		2,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 communities, got %d", len(result))
	}
}

func TestCommunityApplicationFindReturnsEmptyWhenIndexIsOutOfRange(t *testing.T) {
	fixture := newCommunityApplicationFixture(t)

	platform := newTestCommunityPlatform(t, fixture)
	community := newTestCommunity(t, platform, "golang", "Golang", "Go community", nil)

	saveCommunities(t, fixture, community)

	result, err := fixture.CommunityApplication.Find(
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

func TestCommunityApplicationFindAfter(t *testing.T) {
	fixture := newCommunityApplicationFixture(t)

	platform := newTestCommunityPlatform(t, fixture)

	first := newTestCommunity(t, platform, "first", "First", "First community", nil)
	second := newTestCommunity(t, platform, "second", "Second", "Second community", nil)
	third := newTestCommunity(t, platform, "third", "Third", "Third community", nil)

	saveCommunities(t, fixture, first, second, third)

	all, err := fixture.CommunityApplication.Find(
		fixture.Context,
		0,
		10,
	)
	if err != nil {
		t.Fatal(err)
	}

	result, err := fixture.CommunityApplication.FindAfter(
		fixture.Context,
		all[0].Identifier(),
		10,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 communities, got %d", len(result))
	}

	if result[0].Identifier() != all[1].Identifier() {
		t.Fatalf("expected second community first")
	}
}

func TestCommunityApplicationFindAfterWithNilCursor(t *testing.T) {
	fixture := newCommunityApplicationFixture(t)

	platform := newTestCommunityPlatform(t, fixture)

	first := newTestCommunity(t, platform, "first", "First", "First community", nil)
	second := newTestCommunity(t, platform, "second", "Second", "Second community", nil)

	saveCommunities(t, fixture, first, second)

	result, err := fixture.CommunityApplication.FindAfter(
		fixture.Context,
		uuid.Nil,
		10,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 communities, got %d", len(result))
	}
}

func TestCommunityApplicationFindByPlatform(t *testing.T) {
	fixture := newCommunityApplicationFixture(t)

	platform := newTestCommunityPlatform(t, fixture)
	otherPlatform := newTestCommunityPlatformWithValues(
		t,
		fixture,
		"X",
		"x",
		"https://x.com",
	)

	first := newTestCommunity(t, platform, "first", "First", "First community", nil)
	second := newTestCommunity(t, platform, "second", "Second", "Second community", nil)
	other := newTestCommunity(t, otherPlatform, "third", "Third", "Third community", nil)

	saveCommunities(t, fixture, first, second, other)

	result, err := fixture.CommunityApplication.FindByPlatform(
		fixture.Context,
		platform,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 communities, got %d", len(result))
	}
}

func TestCommunityApplicationCount(t *testing.T) {
	fixture := newCommunityApplicationFixture(t)

	platform := newTestCommunityPlatform(t, fixture)

	first := newTestCommunity(t, platform, "first", "First", "First community", nil)
	second := newTestCommunity(t, platform, "second", "Second", "Second community", nil)

	saveCommunities(t, fixture, first, second)

	count, err := fixture.CommunityApplication.Count(
		fixture.Context,
	)
	if err != nil {
		t.Fatal(err)
	}

	if count != 2 {
		t.Fatalf("expected count 2, got %d", count)
	}
}

func newTestCommunityPlatform(
	t *testing.T,
	fixture *communityApplicationFixture,
) domain_platforms.Platform {
	t.Helper()

	return newTestCommunityPlatformWithValues(
		t,
		fixture,
		"Reddit",
		"reddit",
		"https://reddit.com",
	)
}

func newTestCommunityPlatformWithValues(
	t *testing.T,
	fixture *communityApplicationFixture,
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

func newTestCommunityUser(
	t *testing.T,
	fixture *communityApplicationFixture,
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

	err = fixture.UserApplication.Save(
		fixture.Context,
		user,
	)
	if err != nil {
		t.Fatal(err)
	}

	return user
}

func newTestCommunity(
	t *testing.T,
	platform domain_platforms.Platform,
	handle string,
	title string,
	text string,
	moderators []domain_users.User,
) domain_communities.Community {
	t.Helper()

	community, err := domain_communities.NewAdapter().ToDomain(
		domain_communities.CommunityInput{
			Identifier: uuid.New(),
			Platform:   platform,
			Handle:     handle,
			Title:      title,
			Text:       text,
			CreatedOn:  time.Now().UTC(),
			Moderators: moderators,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	return community
}

func saveCommunities(
	t *testing.T,
	fixture *communityApplicationFixture,
	communities ...domain_communities.Community,
) {
	t.Helper()

	for _, community := range communities {
		err := fixture.CommunityApplication.Save(
			fixture.Context,
			community,
		)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func assertCommunity(
	t *testing.T,
	result domain_communities.Community,
	expected domain_communities.Community,
) {
	t.Helper()

	if result == nil {
		t.Fatalf("expected community")
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

	if result.Handle() != expected.Handle() {
		t.Fatalf("expected handle %s, got %s", expected.Handle(), result.Handle())
	}

	if result.Title() != expected.Title() {
		t.Fatalf("expected title %s, got %s", expected.Title(), result.Title())
	}

	if result.Text() != expected.Text() {
		t.Fatalf("expected text %s, got %s", expected.Text(), result.Text())
	}

	expectedModerators := expected.Moderators()
	resultModerators := result.Moderators()

	if len(resultModerators) != len(expectedModerators) {
		t.Fatalf(
			"expected %d moderators, got %d",
			len(expectedModerators),
			len(resultModerators),
		)
	}

	expectedModeratorIDs := map[uuid.UUID]bool{}
	for _, moderator := range expectedModerators {
		expectedModeratorIDs[moderator.Identifier()] = true
	}

	for _, moderator := range resultModerators {
		if !expectedModeratorIDs[moderator.Identifier()] {
			t.Fatalf("unexpected moderator %s", moderator.Identifier())
		}
	}
}
