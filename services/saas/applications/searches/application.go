package searches

import (
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

type application struct {
	embedder         embeddings.Embedder
	searchRepository domain_searches.Repository

	postRepository         posts.Repository
	userRepository         users.Repository
	communityRepository    communities.Repository
	campaignRepository     campaigns.Repository
	topicRepository        topics.Repository
	narrativeRepository    narratives.Repository
	relationshipRepository relationships.Repository
}

func createApplication(
	embedder embeddings.Embedder,
	searchRepository domain_searches.Repository,
	postRepository posts.Repository,
	userRepository users.Repository,
	communityRepository communities.Repository,
	campaignRepository campaigns.Repository,
	topicRepository topics.Repository,
	narrativeRepository narratives.Repository,
	relationshipRepository relationships.Repository,
) Application {
	out := application{
		embedder:         embedder,
		searchRepository: searchRepository,

		postRepository:         postRepository,
		userRepository:         userRepository,
		communityRepository:    communityRepository,
		campaignRepository:     campaignRepository,
		topicRepository:        topicRepository,
		narrativeRepository:    narrativeRepository,
		relationshipRepository: relationshipRepository,
	}

	return &out
}

// IndexPost indexes a post for semantic search
func (app *application) IndexPost(post posts.Post) error {
	vector, err := app.embedder.Embed(post.Content().Text())
	if err != nil {
		return err
	}

	return app.searchRepository.Store(
		post.Identifier(),
		domain_searches.PostKind,
		vector,
	)
}

// Search searches supported entities
func (app *application) Search(query string, limit int) ([]Result, error) {
	matches, err := app.search(query, limit)
	if err != nil {
		return nil, err
	}

	results := make([]Result, 0, len(matches))

	for _, match := range matches {
		switch match.Kind() {
		case domain_searches.PostKind:
			post, err := app.postRepository.FindByID(match.Target())
			if err != nil {
				return nil, err
			}

			results = append(results, &result{
				identifier: post.Identifier(),
				kind:       PostKind,
				title:      "Post",
				text:       post.Content().Text(),
				score:      match.Similarity(),
			})

		case domain_searches.CampaignKind:
			campaign, err := app.campaignRepository.FindByID(match.Target())
			if err != nil {
				return nil, err
			}

			results = append(results, &result{
				identifier: campaign.Identifier(),
				kind:       CampaignKind,
				title:      campaign.Name(),
				text:       campaign.Description(),
				score:      match.Similarity(),
			})

		case domain_searches.TopicKind:
			topic, err := app.topicRepository.FindByID(match.Target())
			if err != nil {
				return nil, err
			}

			results = append(results, &result{
				identifier: topic.Identifier(),
				kind:       TopicKind,
				title:      topic.Name(),
				text:       topic.Description(),
				score:      match.Similarity(),
			})

		case domain_searches.NarrativeKind:
			narrative, err := app.narrativeRepository.FindByID(match.Target())
			if err != nil {
				return nil, err
			}

			results = append(results, &result{
				identifier: narrative.Identifier(),
				kind:       NarrativeKind,
				title:      narrative.Name(),
				text:       narrative.Description(),
				score:      match.Similarity(),
			})

		case domain_searches.UserKind:
			user, err := app.userRepository.FindByID(match.Target())
			if err != nil {
				return nil, err
			}

			results = append(results, &result{
				identifier: user.Identifier(),
				kind:       UserKind,
				title:      user.Handle(),
				text:       user.DisplayName(),
				score:      match.Similarity(),
			})

		case domain_searches.CommunityKind:
			community, err := app.communityRepository.FindByID(match.Target())
			if err != nil {
				return nil, err
			}

			results = append(results, &result{
				identifier: community.Identifier(),
				kind:       CommunityKind,
				title:      community.Title(),
				text:       community.Text(),
				score:      match.Similarity(),
			})

		case domain_searches.RelationshipKind:
			relationship, err := app.relationshipRepository.FindByID(match.Target())
			if err != nil {
				return nil, err
			}

			results = append(results, &result{
				identifier: relationship.Identifier(),
				kind:       RelationshipKind,
				title:      "Relationship",
				text:       "",
				score:      match.Similarity(),
			})
		}
	}

	return results, nil
}

func (app *application) search(
	query string,
	limit int,
) ([]domain_searches.Match, error) {
	vector, err := app.embedder.Embed(query)
	if err != nil {
		return nil, err
	}

	return app.searchRepository.Search(vector, limit)
}

// SearchPosts searches posts
func (app *application) SearchPosts(
	query string,
	limit int,
) ([]posts.Post, error) {
	matches, err := app.search(query, limit)
	if err != nil {
		return nil, err
	}

	out := make([]posts.Post, 0, len(matches))

	for _, match := range matches {
		post, err := app.postRepository.FindByID(match.Target())
		if err != nil {
			return nil, err
		}

		out = append(out, post)
	}

	return out, nil
}

// SearchCampaigns searches campaigns
func (app *application) SearchCampaigns(
	query string,
	limit int,
) ([]campaigns.Campaign, error) {
	matches, err := app.search(query, limit)
	if err != nil {
		return nil, err
	}

	out := make([]campaigns.Campaign, 0, len(matches))

	for _, match := range matches {
		campaign, err := app.campaignRepository.FindByID(match.Target())
		if err != nil {
			return nil, err
		}

		out = append(out, campaign)
	}

	return out, nil
}

// SearchTopics searches topics
func (app *application) SearchTopics(
	query string,
	limit int,
) ([]topics.Topic, error) {
	matches, err := app.search(query, limit)
	if err != nil {
		return nil, err
	}

	out := make([]topics.Topic, 0, len(matches))

	for _, match := range matches {
		topic, err := app.topicRepository.FindByID(match.Target())
		if err != nil {
			return nil, err
		}

		out = append(out, topic)
	}

	return out, nil
}

// SearchNarratives searches narratives
func (app *application) SearchNarratives(
	query string,
	limit int,
) ([]narratives.Narrative, error) {
	matches, err := app.search(query, limit)
	if err != nil {
		return nil, err
	}

	out := make([]narratives.Narrative, 0, len(matches))

	for _, match := range matches {
		narrative, err := app.narrativeRepository.FindByID(match.Target())
		if err != nil {
			return nil, err
		}

		out = append(out, narrative)
	}

	return out, nil
}

// SearchUsers searches users
func (app *application) SearchUsers(
	query string,
	limit int,
) ([]users.User, error) {
	matches, err := app.search(query, limit)
	if err != nil {
		return nil, err
	}

	out := make([]users.User, 0, len(matches))

	for _, match := range matches {
		user, err := app.userRepository.FindByID(match.Target())
		if err != nil {
			return nil, err
		}

		out = append(out, user)
	}

	return out, nil
}

// SearchCommunities searches communities
func (app *application) SearchCommunities(
	query string,
	limit int,
) ([]communities.Community, error) {
	matches, err := app.search(query, limit)
	if err != nil {
		return nil, err
	}

	out := make([]communities.Community, 0, len(matches))

	for _, match := range matches {
		community, err := app.communityRepository.FindByID(match.Target())
		if err != nil {
			return nil, err
		}

		out = append(out, community)
	}

	return out, nil
}

// SearchRelationships searches relationships
func (app *application) SearchRelationships(
	query string,
	limit int,
) ([]relationships.Relationship, error) {
	matches, err := app.search(query, limit)
	if err != nil {
		return nil, err
	}

	out := make([]relationships.Relationship, 0, len(matches))

	for _, match := range matches {
		relationship, err := app.relationshipRepository.FindByID(match.Target())
		if err != nil {
			return nil, err
		}

		out = append(out, relationship)
	}

	return out, nil
}
