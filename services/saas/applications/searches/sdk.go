package searches

import (
	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/communities"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/campaigns"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/narratives"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/topics"
	"github.com/steve-rodrigue/aabs/services/saas/domain/posts"
	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships"
	"github.com/steve-rodrigue/aabs/services/saas/domain/users"
)

type ResultKind string

const (
	PostKind         ResultKind = "post"
	CampaignKind     ResultKind = "campaign"
	TopicKind        ResultKind = "topic"
	NarrativeKind    ResultKind = "narrative"
	UserKind         ResultKind = "user"
	CommunityKind    ResultKind = "community"
	RelationshipKind ResultKind = "relationship"
)

// Result represents a search result
type Result interface {
	Identifier() uuid.UUID
	Kind() ResultKind
	Title() string
	Text() string
	Score() float64
}

// Application represents a search application
type Application interface {
	Search(query string, limit int) ([]Result, error)
	SearchPosts(query string, limit int) ([]posts.Post, error)
	SearchCampaigns(query string, limit int) ([]campaigns.Campaign, error)
	SearchTopics(query string, limit int) ([]topics.Topic, error)
	SearchNarratives(query string, limit int) ([]narratives.Narrative, error)
	SearchUsers(query string, limit int) ([]users.User, error)
	SearchCommunities(query string, limit int) ([]communities.Community, error)
	SearchRelationships(query string, limit int) ([]relationships.Relationship, error)
}
