package clusters

import (
	"github.com/google/uuid"
	domain_clusters "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters/clusterables"
)

// Application represents the cluster application
type Application interface {
	BuildForTarget(target clusterables.Clusterable, members []clusterables.Clusterable) ([]domain_clusters.Cluster, error)
	FindByID(id uuid.UUID) (domain_clusters.Cluster, error)
	FindByTarget(target clusterables.Clusterable) ([]domain_clusters.Cluster, error)
	FindByMember(member clusterables.Clusterable) ([]domain_clusters.Cluster, error)
	RebuildAll() error
	RebuildPostClusters() error
	RebuildCampaignClusters() error
	RebuildTopicClusters() error
	RebuildNarrativeClusters() error
}
