package clusters

import (
	"github.com/google/uuid"

	domain_clusters "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters/clusterables"
)

// New creates a new cluster application
func New(
	repository domain_clusters.Repository,
	detector domain_clusters.Detector,
	clusterableRepository clusterables.Repository,
	candidateRepository clusterables.CandidateRepository,
) Application {
	return createApplication(
		repository,
		detector,
		clusterableRepository,
		candidateRepository,
	)
}

// Application represents the cluster application
type Application interface {
	BuildForTarget(target clusterables.Clusterable, members []clusterables.Clusterable) ([]domain_clusters.Cluster, error)

	FindByID(id uuid.UUID) (domain_clusters.Cluster, error)
	FindByTarget(target clusterables.Clusterable) ([]domain_clusters.Cluster, error)
	FindByMember(member clusterables.Clusterable) ([]domain_clusters.Cluster, error)

	RebuildAll() error

	RebuildPostClusters() error
	RebuildUserClusters() error
	RebuildCommunityClusters() error
	RebuildPlatformClusters() error
	RebuildCampaignClusters() error
	RebuildTopicClusters() error
	RebuildNarrativeClusters() error
}
