package relatables

import (
	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/campaigns"
	"github.com/steve-rodrigue/aabs/services/saas/domain/narratives"
	"github.com/steve-rodrigue/aabs/services/saas/domain/posts"
	"github.com/steve-rodrigue/aabs/services/saas/domain/topics"
	"github.com/steve-rodrigue/aabs/services/saas/domain/users"
)

// Relatable represents a relationship relatable
type Relatable interface {
	Identifier() uuid.UUID
	IsCampaign() bool
	Campaign() campaigns.Campaign
	IsTopic() bool
	Topic() topics.Topic
	IsUser() bool
	User() users.User
	IsPost() bool
	Post() posts.Post
	IsNarrative() bool
	Narrative() narratives.Narrative
}
