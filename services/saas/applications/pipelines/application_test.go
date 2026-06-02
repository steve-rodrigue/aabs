package pipelines

import (
	"errors"
	"testing"

	domain_posts "github.com/steve-rodrigue/aabs/services/saas/domain/posts"
)

var errTest = errors.New("test error")

func TestProcessPostSuccess(t *testing.T) {
	fixture := newApplicationFixture()
	post := domain_posts.NewMockPost("hello world")

	err := fixture.application.ProcessPost(post)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.posts.SaveCalls != 1 {
		t.Fatalf("expected 1 post save, got %d", fixture.posts.SaveCalls)
	}

	if fixture.searches.IndexPostCalls != 1 {
		t.Fatalf("expected 1 index post call, got %d", fixture.searches.IndexPostCalls)
	}

	if fixture.searches.LastPost != post {
		t.Fatalf("expected indexed post to be the processed post")
	}

	if fixture.searches.SearchPostsCalls != 1 {
		t.Fatalf("expected 1 search posts call, got %d", fixture.searches.SearchPostsCalls)
	}

	if fixture.searches.LastQuery != "hello world" {
		t.Fatalf("expected search query %q, got %q", "hello world", fixture.searches.LastQuery)
	}

	if fixture.clusters.RebuildPostClustersCalls != 1 {
		t.Fatalf("expected post clusters rebuild")
	}

	if fixture.campaigns.RebuildCampaignsCalls != 1 {
		t.Fatalf("expected campaigns rebuild")
	}

	if fixture.topics.RebuildTopicsCalls != 1 {
		t.Fatalf("expected topics rebuild")
	}

	if fixture.narratives.RebuildNarrativesCalls != 1 {
		t.Fatalf("expected narratives rebuild")
	}

	if fixture.participations.RebuildParticipationsCalls != 1 {
		t.Fatalf("expected participations rebuild")
	}

	if fixture.relationships.RebuildRelationshipsCalls != 1 {
		t.Fatalf("expected relationships rebuild")
	}

	if fixture.scores.RecalculateScoresCalls != 1 {
		t.Fatalf("expected scores recalculation")
	}
}

func TestProcessPostReturnsSaveError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.posts.SaveErr = errTest

	err := fixture.application.ProcessPost(domain_posts.NewMockPost("hello"))

	if !errors.Is(err, errTest) {
		t.Fatalf("expected save error, got %v", err)
	}

	if fixture.searches.IndexPostCalls != 0 {
		t.Fatalf("expected index post not to be called")
	}

	if fixture.searches.SearchPostsCalls != 0 {
		t.Fatalf("expected search not to be called")
	}
}

func TestProcessPostReturnsIndexPostError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.searches.IndexPostErr = errTest

	err := fixture.application.ProcessPost(domain_posts.NewMockPost("hello"))

	if !errors.Is(err, errTest) {
		t.Fatalf("expected index post error, got %v", err)
	}

	if fixture.searches.IndexPostCalls != 1 {
		t.Fatalf("expected 1 index post call, got %d", fixture.searches.IndexPostCalls)
	}

	if fixture.searches.SearchPostsCalls != 0 {
		t.Fatalf("expected search not to be called")
	}

	if fixture.clusters.RebuildPostClustersCalls != 0 {
		t.Fatalf("expected rebuild not to be called")
	}
}

func TestProcessPostReturnsSearchError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.searches.SearchPostsErr = errTest

	err := fixture.application.ProcessPost(domain_posts.NewMockPost("hello"))

	if !errors.Is(err, errTest) {
		t.Fatalf("expected search error, got %v", err)
	}

	if fixture.clusters.RebuildPostClustersCalls != 0 {
		t.Fatalf("expected rebuild not to be called")
	}

	if fixture.searches.IndexPostCalls != 1 {
		t.Fatalf("expected index post to be called before search")
	}
}

func TestProcessPostsSuccess(t *testing.T) {
	fixture := newApplicationFixture()

	err := fixture.application.ProcessPosts([]domain_posts.Post{
		domain_posts.NewMockPost("one"),
		domain_posts.NewMockPost("two"),
	})

	if err != nil {
		t.Fatal(err)
	}

	if fixture.posts.SaveCalls != 2 {
		t.Fatalf("expected 2 saves, got %d", fixture.posts.SaveCalls)
	}

	if fixture.searches.SearchPostsCalls != 2 {
		t.Fatalf("expected 2 search calls, got %d", fixture.searches.SearchPostsCalls)
	}

	if fixture.clusters.RebuildPostClustersCalls != 2 {
		t.Fatalf("expected 2 rebuilds, got %d", fixture.clusters.RebuildPostClustersCalls)
	}

	if fixture.searches.IndexPostCalls != 2 {
		t.Fatalf("expected 2 index post calls, got %d", fixture.searches.IndexPostCalls)
	}
}

func TestProcessPostsStopsOnError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.searches.SearchPostsErr = errTest

	err := fixture.application.ProcessPosts([]domain_posts.Post{
		domain_posts.NewMockPost("one"),
		domain_posts.NewMockPost("two"),
	})

	if !errors.Is(err, errTest) {
		t.Fatalf("expected error, got %v", err)
	}

	if fixture.posts.SaveCalls != 1 {
		t.Fatalf("expected processing to stop after first post")
	}

	if fixture.searches.IndexPostCalls != 1 {
		t.Fatalf("expected 1 index post call, got %d", fixture.searches.IndexPostCalls)
	}
}

func TestRebuildReturnsErrors(t *testing.T) {
	tests := []struct {
		name      string
		configure func(*applicationFixture)
	}{
		{
			name: "clusters",
			configure: func(fixture *applicationFixture) {
				fixture.clusters.RebuildPostClustersErr = errTest
			},
		},
		{
			name: "campaigns",
			configure: func(fixture *applicationFixture) {
				fixture.campaigns.RebuildCampaignsErr = errTest
			},
		},
		{
			name: "topics",
			configure: func(fixture *applicationFixture) {
				fixture.topics.RebuildTopicsErr = errTest
			},
		},
		{
			name: "narratives",
			configure: func(fixture *applicationFixture) {
				fixture.narratives.RebuildNarrativesErr = errTest
			},
		},
		{
			name: "participations",
			configure: func(fixture *applicationFixture) {
				fixture.participations.RebuildParticipationsErr = errTest
			},
		},
		{
			name: "relationships",
			configure: func(fixture *applicationFixture) {
				fixture.relationships.RebuildRelationshipsErr = errTest
			},
		},
		{
			name: "scores",
			configure: func(fixture *applicationFixture) {
				fixture.scores.RecalculateScoresErr = errTest
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fixture := newApplicationFixture()
			test.configure(fixture)

			err := fixture.application.Rebuild()

			if !errors.Is(err, errTest) {
				t.Fatalf("expected rebuild error, got %v", err)
			}
		})
	}
}
