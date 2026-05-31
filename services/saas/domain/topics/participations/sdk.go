package participations

import (
	"time"

	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/topics"
	"github.com/steve-rodrigue/aabs/services/saas/domain/topics/clusters"
	"github.com/steve-rodrigue/aabs/services/saas/domain/users"
)

// Participation represents a user participation in a topic
type Participation interface {
	Identifier() uuid.UUID
	Topic() topics.Topic
	Cluster() clusters.Cluster
	User() users.User
	PostCount() int
	TotalUserPostCount() int
	Percentage() float64
	DetectedOn() time.Time
}

// Repository represents a participation repository
type Repository interface {
	Save(participation Participation) error
	FindByTopic(topic uuid.UUID) ([]Participation, error)
	FindByCluster(cluster uuid.UUID) ([]Participation, error)
	FindByUser(user uuid.UUID) ([]Participation, error)
}

// Calculator represents a participation calculator
type Calculator interface {
	Calculate(cluster clusters.Cluster) ([]Participation, error)
}
