package searches

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/communities"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/campaigns"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/narratives"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/topics"
	"github.com/steve-rodrigue/aabs/services/saas/domain/posts"
	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships"
	domain_searches "github.com/steve-rodrigue/aabs/services/saas/domain/searches"
	"github.com/steve-rodrigue/aabs/services/saas/domain/users"
)

var errTest = errors.New("test error")

func TestIndexPostSuccess(t *testing.T) {
	fixture := newFixture()
	post := posts.NewMockPost("hello")

	err := fixture.app.IndexPost(post)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.embedder.EmbedCalls != 1 {
		t.Fatalf("expected 1 embed call")
	}

	if fixture.embedder.LastText != "hello" {
		t.Fatalf("expected text %q, got %q", "hello", fixture.embedder.LastText)
	}

	if fixture.searchRepository.StoreCalls != 1 {
		t.Fatalf("expected 1 store call")
	}

	if fixture.searchRepository.LastStoredTarget != post.Identifier() {
		t.Fatalf("expected stored target to be post id")
	}

	if fixture.searchRepository.LastStoredKind != domain_searches.PostKind {
		t.Fatalf("expected stored kind post")
	}
}

func TestIndexPostReturnsEmbedError(t *testing.T) {
	fixture := newFixture()
	fixture.embedder.EmbedErr = errTest

	err := fixture.app.IndexPost(posts.NewMockPost("hello"))

	if !errors.Is(err, errTest) {
		t.Fatalf("expected embed error, got %v", err)
	}

	if fixture.searchRepository.StoreCalls != 0 {
		t.Fatalf("expected store not to be called")
	}
}

func TestIndexPostReturnsStoreError(t *testing.T) {
	fixture := newFixture()
	fixture.searchRepository.StoreErr = errTest

	err := fixture.app.IndexPost(posts.NewMockPost("hello"))

	if !errors.Is(err, errTest) {
		t.Fatalf("expected store error, got %v", err)
	}
}

func TestSearchReturnsMixedResults(t *testing.T) {
	fixture := newFixture()

	post := posts.NewMockPost("post text")
	campaign := campaigns.NewMockCampaign("Campaign", "Campaign description")
	topic := topics.NewMockTopic("Topic", "Topic description")
	narrative := narratives.NewMockNarrative("Narrative", "Narrative description")
	user := users.NewMockUser("@user", "User Display")
	community := communities.NewMockCommunity("Community", "Community text")
	relationship := relationships.NewMockRelationship()

	fixture.posts.Items[post.Identifier()] = post
	fixture.campaigns.Items[campaign.Identifier()] = campaign
	fixture.topics.Items[topic.Identifier()] = topic
	fixture.narratives.Items[narrative.Identifier()] = narrative
	fixture.users.Items[user.Identifier()] = user
	fixture.communities.Items[community.Identifier()] = community
	fixture.relationships.Items[relationship.Identifier()] = relationship

	fixture.searchRepository.Matches = []domain_searches.Match{
		domain_searches.NewMockMatch(post.Identifier(), domain_searches.PostKind, 0.91),
		domain_searches.NewMockMatch(campaign.Identifier(), domain_searches.CampaignKind, 0.82),
		domain_searches.NewMockMatch(topic.Identifier(), domain_searches.TopicKind, 0.73),
		domain_searches.NewMockMatch(narrative.Identifier(), domain_searches.NarrativeKind, 0.64),
		domain_searches.NewMockMatch(user.Identifier(), domain_searches.UserKind, 0.55),
		domain_searches.NewMockMatch(community.Identifier(), domain_searches.CommunityKind, 0.46),
		domain_searches.NewMockMatch(relationship.Identifier(), domain_searches.RelationshipKind, 0.37),
	}

	results, err := fixture.app.Search("hello", 10)

	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 7 {
		t.Fatalf("expected 7 results, got %d", len(results))
	}

	assertResult(t, results[0], post.Identifier(), PostKind, "Post", "post text", 0.91)
	assertResult(t, results[1], campaign.Identifier(), CampaignKind, "Campaign", "Campaign description", 0.82)
	assertResult(t, results[2], topic.Identifier(), TopicKind, "Topic", "Topic description", 0.73)
	assertResult(t, results[3], narrative.Identifier(), NarrativeKind, "Narrative", "Narrative description", 0.64)
	assertResult(t, results[4], user.Identifier(), UserKind, "@user", "User Display", 0.55)
	assertResult(t, results[5], community.Identifier(), CommunityKind, "Community", "Community text", 0.46)
	assertResult(t, results[6], relationship.Identifier(), RelationshipKind, "Relationship", "", 0.37)
}

func TestSearchReturnsEmbedError(t *testing.T) {
	fixture := newFixture()
	fixture.embedder.EmbedErr = errTest

	_, err := fixture.app.Search("hello", 10)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected embed error, got %v", err)
	}
}

func TestSearchReturnsRepositorySearchError(t *testing.T) {
	fixture := newFixture()
	fixture.searchRepository.SearchErr = errTest

	_, err := fixture.app.Search("hello", 10)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected search error, got %v", err)
	}
}

func TestSearchPosts(t *testing.T) {
	fixture := newFixture()
	post := posts.NewMockPost("post")

	fixture.posts.Items[post.Identifier()] = post
	fixture.searchRepository.Matches = []domain_searches.Match{
		domain_searches.NewMockMatch(post.Identifier(), domain_searches.PostKind, 0.9),
	}

	results, err := fixture.app.SearchPosts("query", 10)

	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 1 || results[0] != post {
		t.Fatalf("expected post result")
	}
}

func TestSearchCampaigns(t *testing.T) {
	fixture := newFixture()
	campaign := campaigns.NewMockCampaign("Campaign", "Description")

	fixture.campaigns.Items[campaign.Identifier()] = campaign
	fixture.searchRepository.Matches = []domain_searches.Match{
		domain_searches.NewMockMatch(campaign.Identifier(), domain_searches.CampaignKind, 0.9),
	}

	results, err := fixture.app.SearchCampaigns("query", 10)

	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 1 || results[0] != campaign {
		t.Fatalf("expected campaign result")
	}
}

func TestSearchTopics(t *testing.T) {
	fixture := newFixture()
	topic := topics.NewMockTopic("Topic", "Description")

	fixture.topics.Items[topic.Identifier()] = topic
	fixture.searchRepository.Matches = []domain_searches.Match{
		domain_searches.NewMockMatch(topic.Identifier(), domain_searches.TopicKind, 0.9),
	}

	results, err := fixture.app.SearchTopics("query", 10)

	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 1 || results[0] != topic {
		t.Fatalf("expected topic result")
	}
}

func TestSearchNarratives(t *testing.T) {
	fixture := newFixture()
	narrative := narratives.NewMockNarrative("Narrative", "Description")

	fixture.narratives.Items[narrative.Identifier()] = narrative
	fixture.searchRepository.Matches = []domain_searches.Match{
		domain_searches.NewMockMatch(narrative.Identifier(), domain_searches.NarrativeKind, 0.9),
	}

	results, err := fixture.app.SearchNarratives("query", 10)

	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 1 || results[0] != narrative {
		t.Fatalf("expected narrative result")
	}
}

func TestSearchUsers(t *testing.T) {
	fixture := newFixture()
	user := users.NewMockUser("@user", "Display")

	fixture.users.Items[user.Identifier()] = user
	fixture.searchRepository.Matches = []domain_searches.Match{
		domain_searches.NewMockMatch(user.Identifier(), domain_searches.UserKind, 0.9),
	}

	results, err := fixture.app.SearchUsers("query", 10)

	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 1 || results[0] != user {
		t.Fatalf("expected user result")
	}
}

func TestSearchCommunities(t *testing.T) {
	fixture := newFixture()
	community := communities.NewMockCommunity("Community", "Text")

	fixture.communities.Items[community.Identifier()] = community
	fixture.searchRepository.Matches = []domain_searches.Match{
		domain_searches.NewMockMatch(community.Identifier(), domain_searches.CommunityKind, 0.9),
	}

	results, err := fixture.app.SearchCommunities("query", 10)

	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 1 || results[0] != community {
		t.Fatalf("expected community result")
	}
}

func TestSearchRelationships(t *testing.T) {
	fixture := newFixture()
	relationship := relationships.NewMockRelationship()

	fixture.relationships.Items[relationship.Identifier()] = relationship
	fixture.searchRepository.Matches = []domain_searches.Match{
		domain_searches.NewMockMatch(relationship.Identifier(), domain_searches.RelationshipKind, 0.9),
	}

	results, err := fixture.app.SearchRelationships("query", 10)

	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 1 || results[0] != relationship {
		t.Fatalf("expected relationship result")
	}
}

func assertResult(
	t *testing.T,
	result Result,
	id uuid.UUID,
	kind ResultKind,
	title string,
	text string,
	score float64,
) {
	t.Helper()

	if result.Identifier() != id {
		t.Fatalf("expected id %s, got %s", id, result.Identifier())
	}

	if result.Kind() != kind {
		t.Fatalf("expected kind %s, got %s", kind, result.Kind())
	}

	if result.Title() != title {
		t.Fatalf("expected title %q, got %q", title, result.Title())
	}

	if result.Text() != text {
		t.Fatalf("expected text %q, got %q", text, result.Text())
	}

	if result.Score() != score {
		t.Fatalf("expected score %f, got %f", score, result.Score())
	}
}
