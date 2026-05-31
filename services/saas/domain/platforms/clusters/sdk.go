package clusters

import (
	"time"

	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/platforms"
	"github.com/steve-rodrigue/aabs/services/saas/domain/posts"
)

// Cluster represents a semantic cluster of posts on a platform
type Cluster interface {
	Identifier() uuid.UUID
	Platform() platforms.Platform
	MemberPostIDs() []uuid.UUID
	ConfidenceScore() float64
	Centroid() []float32
	CreatedOn() time.Time
}

// Repository represents a platform cluster repository
type Repository interface {
	Save(cluster Cluster) error
	FindByID(id uuid.UUID) (Cluster, error)
	FindByPlatform(platform uuid.UUID) ([]Cluster, error)
}

// Detector detects post clusters inside a platform
type Detector interface {
	Detect(platform platforms.Platform, posts []posts.Post) ([]Cluster, error)
}
