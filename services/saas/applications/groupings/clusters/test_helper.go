package clusters

import (
	"context"

	"github.com/google/uuid"

	domain_clusters "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters/clusterables"
)

func NewMockClustersApplication() *MockClustersApplication {
	return &MockClustersApplication{}
}

type MockClustersApplication struct {
	BuildForTargetCalls int
	BuildForTargetErr   error
	BuildForTargetValue []domain_clusters.Cluster

	FindByIDCalls int
	FindByIDErr   error
	FindByIDValue domain_clusters.Cluster

	FindByTargetCalls int
	FindByTargetErr   error
	FindByTargetValue []domain_clusters.Cluster

	FindByMemberCalls int
	FindByMemberErr   error
	FindByMemberValue []domain_clusters.Cluster

	RebuildAllCalls int
	RebuildAllErr   error

	RebuildPostClustersCalls int
	RebuildPostClustersErr   error

	RebuildUserClustersCalls int
	RebuildUserClustersErr   error

	RebuildCommunityClustersCalls int
	RebuildCommunityClustersErr   error

	RebuildPlatformClustersCalls int
	RebuildPlatformClustersErr   error

	RebuildCampaignClustersCalls int
	RebuildCampaignClustersErr   error

	RebuildTopicClustersCalls int
	RebuildTopicClustersErr   error

	RebuildNarrativeClustersCalls int
	RebuildNarrativeClustersErr   error

	LastContext context.Context
	LastID      uuid.UUID
	LastTarget  clusterables.Clusterable
	LastMember  clusterables.Clusterable
	LastMembers []clusterables.Clusterable
}

func (application *MockClustersApplication) BuildForTarget(
	ctx context.Context,
	target clusterables.Clusterable,
	members []clusterables.Clusterable,
) ([]domain_clusters.Cluster, error) {
	application.BuildForTargetCalls++
	application.LastContext = ctx
	application.LastTarget = target
	application.LastMembers = members

	return application.BuildForTargetValue,
		application.BuildForTargetErr
}

func (application *MockClustersApplication) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (domain_clusters.Cluster, error) {
	application.FindByIDCalls++
	application.LastContext = ctx
	application.LastID = id

	return application.FindByIDValue,
		application.FindByIDErr
}

func (application *MockClustersApplication) FindByTarget(
	ctx context.Context,
	target clusterables.Clusterable,
) ([]domain_clusters.Cluster, error) {
	application.FindByTargetCalls++
	application.LastContext = ctx
	application.LastTarget = target

	return application.FindByTargetValue,
		application.FindByTargetErr
}

func (application *MockClustersApplication) FindByMember(
	ctx context.Context,
	member clusterables.Clusterable,
) ([]domain_clusters.Cluster, error) {
	application.FindByMemberCalls++
	application.LastContext = ctx
	application.LastMember = member

	return application.FindByMemberValue,
		application.FindByMemberErr
}

func (application *MockClustersApplication) RebuildAll(
	ctx context.Context,
) error {
	application.RebuildAllCalls++
	application.LastContext = ctx

	return application.RebuildAllErr
}

func (application *MockClustersApplication) RebuildPostClusters(
	ctx context.Context,
) error {
	application.RebuildPostClustersCalls++
	application.LastContext = ctx

	return application.RebuildPostClustersErr
}

func (application *MockClustersApplication) RebuildUserClusters(
	ctx context.Context,
) error {
	application.RebuildUserClustersCalls++
	application.LastContext = ctx

	return application.RebuildUserClustersErr
}

func (application *MockClustersApplication) RebuildCommunityClusters(
	ctx context.Context,
) error {
	application.RebuildCommunityClustersCalls++
	application.LastContext = ctx

	return application.RebuildCommunityClustersErr
}

func (application *MockClustersApplication) RebuildPlatformClusters(
	ctx context.Context,
) error {
	application.RebuildPlatformClustersCalls++
	application.LastContext = ctx

	return application.RebuildPlatformClustersErr
}

func (application *MockClustersApplication) RebuildCampaignClusters(
	ctx context.Context,
) error {
	application.RebuildCampaignClustersCalls++
	application.LastContext = ctx

	return application.RebuildCampaignClustersErr
}

func (application *MockClustersApplication) RebuildTopicClusters(
	ctx context.Context,
) error {
	application.RebuildTopicClustersCalls++
	application.LastContext = ctx

	return application.RebuildTopicClustersErr
}

func (application *MockClustersApplication) RebuildNarrativeClusters(
	ctx context.Context,
) error {
	application.RebuildNarrativeClustersCalls++
	application.LastContext = ctx

	return application.RebuildNarrativeClustersErr
}
