package application

import (
	"testing"
)

func TestApplication(t *testing.T) {
	fixture := newApplicationFixture()

	if fixture.application.Pipeline() != fixture.pipeline {
		t.Fatalf("expected pipeline application")
	}

	if fixture.application.Posts() != fixture.posts {
		t.Fatalf("expected posts application")
	}

	if fixture.application.Users() != fixture.users {
		t.Fatalf("expected users application")
	}

	if fixture.application.Communities() != fixture.communities {
		t.Fatalf("expected communities application")
	}

	if fixture.application.Platforms() != fixture.platforms {
		t.Fatalf("expected platforms application")
	}

	if fixture.application.Groupings() != fixture.groupings {
		t.Fatalf("expected groupings application")
	}

	if fixture.application.Relationships() != fixture.relationships {
		t.Fatalf("expected relationships application")
	}

	if fixture.application.Scores() != fixture.scores {
		t.Fatalf("expected scores application")
	}

	if fixture.application.Searches() != fixture.searches {
		t.Fatalf("expected searches application")
	}
}
