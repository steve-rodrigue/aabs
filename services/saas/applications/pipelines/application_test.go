package pipelines

import (
	"errors"
	"testing"

	app_searches "github.com/steve-rodrigue/aabs/services/saas/applications/searches"

	"github.com/steve-rodrigue/aabs/services/saas/applications/groupings"
	"github.com/steve-rodrigue/aabs/services/saas/applications/groupings/campaigns"
	"github.com/steve-rodrigue/aabs/services/saas/applications/groupings/clusters"
	"github.com/steve-rodrigue/aabs/services/saas/applications/groupings/narratives"
	"github.com/steve-rodrigue/aabs/services/saas/applications/groupings/participations"
	"github.com/steve-rodrigue/aabs/services/saas/applications/groupings/topics"
	"github.com/steve-rodrigue/aabs/services/saas/applications/relationships"
	"github.com/steve-rodrigue/aabs/services/saas/applications/scores"
	"github.com/steve-rodrigue/aabs/services/saas/domain/posts"
	domain_posts "github.com/steve-rodrigue/aabs/services/saas/domain/posts"
)

var errTest = errors.New("test error")

func TestProcessPostSuccess(t *testing.T) {
	fixture := newApplicationFixture()
	post := posts.NewMockPost("hello world")

	err := fixture.application.ProcessPost(post)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.postRepository.SaveCalls != 1 {
		t.Fatalf("expected 1 post save, got %d", fixture.postRepository.SaveCalls)
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
	fixture.postRepository.SaveErr = errTest

	err := fixture.application.ProcessPost(posts.NewMockPost("hello"))

	if !errors.Is(err, errTest) {
		t.Fatalf("expected save error, got %v", err)
	}

	if fixture.searches.SearchPostsCalls != 0 {
		t.Fatalf("expected search not to be called")
	}
}

func TestProcessPostReturnsSearchError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.searches.SearchPostsErr = errTest

	err := fixture.application.ProcessPost(posts.NewMockPost("hello"))

	if !errors.Is(err, errTest) {
		t.Fatalf("expected search error, got %v", err)
	}

	if fixture.clusters.RebuildPostClustersCalls != 0 {
		t.Fatalf("expected rebuild not to be called")
	}
}

func TestProcessPostsSuccess(t *testing.T) {
	fixture := newApplicationFixture()

	err := fixture.application.ProcessPosts([]domain_posts.Post{
		posts.NewMockPost("one"),
		posts.NewMockPost("two"),
	})

	if err != nil {
		t.Fatal(err)
	}

	if fixture.postRepository.SaveCalls != 2 {
		t.Fatalf("expected 2 saves, got %d", fixture.postRepository.SaveCalls)
	}

	if fixture.searches.SearchPostsCalls != 2 {
		t.Fatalf("expected 2 search calls, got %d", fixture.searches.SearchPostsCalls)
	}

	if fixture.clusters.RebuildPostClustersCalls != 2 {
		t.Fatalf("expected 2 rebuilds, got %d", fixture.clusters.RebuildPostClustersCalls)
	}
}

func TestProcessPostsStopsOnError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.searches.SearchPostsErr = errTest

	err := fixture.application.ProcessPosts([]domain_posts.Post{
		posts.NewMockPost("one"),
		posts.NewMockPost("two"),
	})

	if !errors.Is(err, errTest) {
		t.Fatalf("expected error, got %v", err)
	}

	if fixture.postRepository.SaveCalls != 1 {
		t.Fatalf("expected processing to stop after first post")
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

type applicationFixture struct {
	application Application

	postRepository *posts.MockPostRepository
	searches       *app_searches.MockSearchApplication

	groupings      *groupings.MockGroupingsApplication
	clusters       *clusters.MockClustersApplication
	campaigns      *campaigns.MockCampaignsApplication
	topics         *topics.MockTopicsApplication
	narratives     *narratives.MockNarrativesApplication
	participations *participations.MockParticipationsApplication

	relationships *relationships.MockRelationshipsApplication
	scores        *scores.MockScoresApplication
}

func newApplicationFixture() *applicationFixture {
	postRepository := &posts.MockPostRepository{}
	searches := &app_searches.MockSearchApplication{}

	clusters := &clusters.MockClustersApplication{}
	campaigns := &campaigns.MockCampaignsApplication{}
	topics := &topics.MockTopicsApplication{}
	narratives := &narratives.MockNarrativesApplication{}
	participations := &participations.MockParticipationsApplication{}

	groupings := &groupings.MockGroupingsApplication{
		ClustersIns:       clusters,
		CampaignsIns:      campaigns,
		TopicsIns:         topics,
		NarrativesIns:     narratives,
		ParticipationsIns: participations,
	}

	relationships := &relationships.MockRelationshipsApplication{}
	scores := &scores.MockScoresApplication{}

	application := New(
		postRepository,
		searches,
		groupings,
		relationships,
		scores,
	)

	return &applicationFixture{
		application: application,

		postRepository: postRepository,
		searches:       searches,

		groupings:      groupings,
		clusters:       clusters,
		campaigns:      campaigns,
		topics:         topics,
		narratives:     narratives,
		participations: participations,

		relationships: relationships,
		scores:        scores,
	}
}
