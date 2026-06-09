package clusters

import (
	"context"

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
	rebuildBatchSize int,
	candidateAmount int,
) Application {
	return createApplication(
		repository,
		detector,
		clusterableRepository,
		candidateRepository,
		rebuildBatchSize,
		candidateAmount,
	)
}

// Application represents the cluster application
type Application interface {
	BuildForTarget(
		ctx context.Context,
		target clusterables.Clusterable,
		members []clusterables.Clusterable,
	) ([]domain_clusters.Cluster, error)

	FindByID(
		ctx context.Context,
		id uuid.UUID,
	) (domain_clusters.Cluster, error)

	FindByTarget(
		ctx context.Context,
		target clusterables.Clusterable,
	) ([]domain_clusters.Cluster, error)

	FindByMember(
		ctx context.Context,
		member clusterables.Clusterable,
	) ([]domain_clusters.Cluster, error)

	RebuildAll(ctx context.Context) error

	RebuildPostClusters(ctx context.Context) error
	RebuildUserClusters(ctx context.Context) error
	RebuildCommunityClusters(ctx context.Context) error
	RebuildPlatformClusters(ctx context.Context) error
	RebuildCampaignClusters(ctx context.Context) error
	RebuildTopicClusters(ctx context.Context) error
	RebuildNarrativeClusters(ctx context.Context) error
}
