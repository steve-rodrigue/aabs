package evidences

import (
	"time"

	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/posts"
	"github.com/steve-rodrigue/aabs/services/saas/domain/topics"
	"github.com/steve-rodrigue/aabs/services/saas/domain/users"
)

// Evidence represents one user post that participated in a topic
type Evidence interface {
	Identifier() uuid.UUID
	User() users.User
	Topic() topics.Topic
	Post() posts.Post
	Score() float64
	DetectedOn() time.Time
}

// Repository represents an evidence repository
type Repository interface {
	Save(evidence Evidence) error
	FindByUser(user uuid.UUID) ([]Evidence, error)
	FindByTopic(topic uuid.UUID) ([]Evidence, error)
	FindByUserAndTopic(user uuid.UUID, topic uuid.UUID) ([]Evidence, error)
	FindByPost(post uuid.UUID) ([]Evidence, error)
}
