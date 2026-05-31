package clusters

import (
	"time"

	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/posts"
	"github.com/steve-rodrigue/aabs/services/saas/domain/topics"
)

// Cluster represents a topic cluster
type Cluster interface {
	Identifier() uuid.UUID
	Topic() topics.Topic
	MemberPostIDs() []uuid.UUID
	ConfidenceScore() float64
	Centroid() []float32
	CreatedOn() time.Time
}

// Repository represents a cluster repository
type Repository interface {
	Save(cluster Cluster) error
	FindByTopic(topic uuid.UUID) ([]Cluster, error)
	FindByID(id uuid.UUID) (Cluster, error)
}

// Detector represents a clusters detector for a topic
type Detector interface {
	Detect(topic topics.Topic, posts []posts.Post) ([]Cluster, error)
}
