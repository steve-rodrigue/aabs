package clusters

import (
	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/campaigns/clusters"
	"github.com/steve-rodrigue/aabs/services/saas/domain/communities"
)

// Cluster represents a community cluster
type Cluster interface {
	Identifier() uuid.UUID
	Community() communities.Community
	Cluster() clusters.Cluster
	PostCount() int
	Confidence() float64
}

// Repository represents a community cluster repository
type Repository interface {
	Save(cluster Cluster) error
	FindByCommunity(community uuid.UUID) ([]Cluster, error)
	FindByID(id uuid.UUID) (Cluster, error)
}
