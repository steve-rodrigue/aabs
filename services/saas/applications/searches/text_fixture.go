package searches

import (
	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/communities"
	"github.com/steve-rodrigue/aabs/services/saas/domain/embeddings"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/campaigns"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/narratives"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/topics"
	"github.com/steve-rodrigue/aabs/services/saas/domain/posts"
	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships"
	domain_searches "github.com/steve-rodrigue/aabs/services/saas/domain/searches"
	"github.com/steve-rodrigue/aabs/services/saas/domain/users"
)

func newFixture() *fixture {
	embedder := &embeddings.MockEmbedder{
		Vector: embeddings.Vector{1, 2, 3},
	}

	searchRepository := domain_searches.NewMockSearchRepository()

	posts := &posts.MockPostRepository{
		Items: map[uuid.UUID]posts.Post{},
	}

	users := &users.MockUserRepository{
		Items: map[uuid.UUID]users.User{},
	}

	communities := &communities.MockCommunityRepository{
		Items: map[uuid.UUID]communities.Community{},
	}

	campaigns := &campaigns.MockCampaignRepository{
		Items: map[uuid.UUID]campaigns.Campaign{},
	}

	topics := &topics.MockTopicRepository{
		Items: map[uuid.UUID]topics.Topic{},
	}

	narratives := &narratives.MockNarrativeRepository{
		Items: map[uuid.UUID]narratives.Narrative{},
	}

	relationships := &relationships.MockRelationshipRepository{
		Items: map[uuid.UUID]relationships.Relationship{},
	}

	return &fixture{
		app: New(
			embedder,
			searchRepository,
			posts,
			users,
			communities,
			campaigns,
			topics,
			narratives,
			relationships,
		),

		embedder:         embedder,
		searchRepository: searchRepository,
		posts:            posts,
		users:            users,
		communities:      communities,
		campaigns:        campaigns,
		topics:           topics,
		narratives:       narratives,
		relationships:    relationships,
	}
}

type fixture struct {
	app Application

	embedder         *embeddings.MockEmbedder
	searchRepository *domain_searches.MockSearchRepository
	posts            *posts.MockPostRepository
	users            *users.MockUserRepository
	communities      *communities.MockCommunityRepository
	campaigns        *campaigns.MockCampaignRepository
	topics           *topics.MockTopicRepository
	narratives       *narratives.MockNarrativeRepository
	relationships    *relationships.MockRelationshipRepository
}
