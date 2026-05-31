package application

import (
	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/communities"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/campaigns"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/narratives"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/topics"
	"github.com/steve-rodrigue/aabs/services/saas/domain/platforms"
	"github.com/steve-rodrigue/aabs/services/saas/domain/posts"
	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships"
	"github.com/steve-rodrigue/aabs/services/saas/domain/scores"
	"github.com/steve-rodrigue/aabs/services/saas/domain/users"
)

// Application represents the public application API.
type Application interface {
	// Ingestion
	ProcessPost(post posts.Post) error
	ProcessPosts(posts []posts.Post) error
	Rebuild() error

	// Posts
	Post(id uuid.UUID) (posts.Post, error)
	Posts() ([]posts.Post, error)

	// Campaigns
	Campaign(id uuid.UUID) (campaigns.Campaign, error)
	Campaigns() ([]campaigns.Campaign, error)
	FindCampaignsByUser(user users.User) ([]campaigns.Campaign, error)
	FindCampaignsByCommunity(community communities.Community) ([]campaigns.Campaign, error)
	FindCampaignsByPlatform(platform platforms.Platform) ([]campaigns.Campaign, error)

	// Topics
	Topic(id uuid.UUID) (topics.Topic, error)
	Topics() ([]topics.Topic, error)
	FindTopicsByUser(user users.User) ([]topics.Topic, error)
	FindTopicsByCommunity(community communities.Community) ([]topics.Topic, error)

	// Narratives
	Narrative(id uuid.UUID) (narratives.Narrative, error)
	Narratives() ([]narratives.Narrative, error)
	FindNarrativesByUser(user users.User) ([]narratives.Narrative, error)
	FindNarrativesByCommunity(community communities.Community) ([]narratives.Narrative, error)

	// Users
	User(id uuid.UUID) (users.User, error)
	Users() ([]users.User, error)

	// Communities
	Community(id uuid.UUID) (communities.Community, error)
	Communities() ([]communities.Community, error)

	// Relationships
	Relationships() ([]relationships.Relationship, error)
	RelationshipsBySource(id uuid.UUID) ([]relationships.Relationship, error)
	RelationshipsByTarget(id uuid.UUID) ([]relationships.Relationship, error)

	// Scores
	LatestScore(id uuid.UUID) (scores.Score, error)
	ScoreHistory(id uuid.UUID) ([]scores.Score, error)

	// Search
	SearchPosts(query string, limit int) ([]posts.Post, error)
	SearchCampaigns(query string, limit int) ([]campaigns.Campaign, error)
	SearchTopics(query string, limit int) ([]topics.Topic, error)
	SearchNarratives(query string, limit int) ([]narratives.Narrative, error)
	SearchUsers(query string, limit int) ([]users.User, error)
	SearchCommunities(query string, limit int) ([]communities.Community, error)
	SearchRelationships(query string, limit int) ([]relationships.Relationship, error)

	// Analytics
	RecalculateScores() error
	RebuildCampaigns() error
	RebuildTopics() error
	RebuildNarratives() error
	RebuildRelationships() error
	RebuildParticipations() error
}
