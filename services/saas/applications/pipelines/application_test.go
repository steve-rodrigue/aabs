package pipelines

import (
	"errors"
	"testing"

	domain_posts "github.com/steve-rodrigue/aabs/services/saas/domain/posts"
	domain_searches "github.com/steve-rodrigue/aabs/services/saas/domain/searches"
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

	if fixture.searches.IndexCalls != 0 {
		t.Fatalf("expected search index not to be called")
	}

	if fixture.searches.SearchCalls != 0 {
		t.Fatalf("expected search not to be called")
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

	if fixture.searches.IndexCalls != 0 {
		t.Fatalf("expected index not to be called")
	}

	if fixture.searches.SearchCalls != 0 {
		t.Fatalf("expected search not to be called")
	}

	if fixture.clusters.RebuildPostClustersCalls != 0 {
		t.Fatalf("expected rebuild not to be called")
	}
}

func TestProcessPostIndexesSearchablePost(t *testing.T) {
	fixture := newApplicationFixture()

	post := &mockSearchablePost{
		Post: domain_posts.NewMockPost("hello world"),
	}

	err := fixture.application.ProcessPost(post)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.posts.SaveCalls != 1 {
		t.Fatalf("expected 1 post save")
	}

	if fixture.searches.IndexCalls != 1 {
		t.Fatalf("expected 1 index call")
	}

	if fixture.searches.LastIndexed != post {
		t.Fatalf("expected indexed searchable post")
	}

	if fixture.clusters.RebuildPostClustersCalls != 1 {
		t.Fatalf("expected rebuild to be called")
	}
}

func TestProcessPostReturnsIndexError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.searches.IndexErr = errTest

	post := &mockSearchablePost{
		Post: domain_posts.NewMockPost("hello world"),
	}

	err := fixture.application.ProcessPost(post)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected index error, got %v", err)
	}

	if fixture.searches.IndexCalls != 1 {
		t.Fatalf("expected 1 index call")
	}

	if fixture.clusters.RebuildPostClustersCalls != 0 {
		t.Fatalf("expected rebuild not to be called")
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

	if fixture.searches.IndexCalls != 0 {
		t.Fatalf("expected no index calls")
	}

	if fixture.searches.SearchCalls != 0 {
		t.Fatalf("expected no search calls")
	}

	if fixture.clusters.RebuildPostClustersCalls != 2 {
		t.Fatalf("expected 2 rebuilds, got %d", fixture.clusters.RebuildPostClustersCalls)
	}
}

func TestProcessPostsStopsOnProcessPostError(t *testing.T) {
	fixture := newApplicationFixture()

	fixture.posts.SaveErr = errTest

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

	if fixture.clusters.RebuildPostClustersCalls != 0 {
		t.Fatalf("expected rebuild not to be called")
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

type mockSearchablePost struct {
	domain_posts.Post
}

func (post *mockSearchablePost) SearchKind() domain_searches.Kind {
	return domain_searches.PostKind
}

func (post *mockSearchablePost) SearchTitle() string {
	if post.Content().IsThread() {
		return post.Content().Thread().Title()
	}

	return ""
}

func (post *mockSearchablePost) SearchText() string {
	return post.Content().Text()
}
