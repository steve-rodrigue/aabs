package applications

import "testing"

func TestPosts(t *testing.T) {
	fixture := newApplicationFixture()

	result := fixture.application.Posts()

	if result != fixture.posts {
		t.Fatalf("expected posts application")
	}
}

func TestUsers(t *testing.T) {
	fixture := newApplicationFixture()

	result := fixture.application.Users()

	if result != fixture.users {
		t.Fatalf("expected users application")
	}
}

func TestCommunities(t *testing.T) {
	fixture := newApplicationFixture()

	result := fixture.application.Communities()

	if result != fixture.communities {
		t.Fatalf("expected communities application")
	}
}

func TestPlatforms(t *testing.T) {
	fixture := newApplicationFixture()

	result := fixture.application.Platforms()

	if result != fixture.platforms {
		t.Fatalf("expected platforms application")
	}
}
