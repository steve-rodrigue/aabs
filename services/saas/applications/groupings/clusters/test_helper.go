package clusters

import (
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
}

func (application *MockClustersApplication) BuildForTarget(
	target clusterables.Clusterable,
	members []clusterables.Clusterable,
) ([]domain_clusters.Cluster, error) {
	application.BuildForTargetCalls++

	return application.BuildForTargetValue,
		application.BuildForTargetErr
}

func (application *MockClustersApplication) FindByID(
	id uuid.UUID,
) (domain_clusters.Cluster, error) {
	application.FindByIDCalls++

	return application.FindByIDValue,
		application.FindByIDErr
}

func (application *MockClustersApplication) FindByTarget(
	target clusterables.Clusterable,
) ([]domain_clusters.Cluster, error) {
	application.FindByTargetCalls++

	return application.FindByTargetValue,
		application.FindByTargetErr
}

func (application *MockClustersApplication) FindByMember(
	member clusterables.Clusterable,
) ([]domain_clusters.Cluster, error) {
	application.FindByMemberCalls++

	return application.FindByMemberValue,
		application.FindByMemberErr
}

func (application *MockClustersApplication) RebuildAll() error {
	application.RebuildAllCalls++

	return application.RebuildAllErr
}

func (application *MockClustersApplication) RebuildPostClusters() error {
	application.RebuildPostClustersCalls++

	return application.RebuildPostClustersErr
}

func (application *MockClustersApplication) RebuildUserClusters() error {
	application.RebuildUserClustersCalls++

	return application.RebuildUserClustersErr
}

func (application *MockClustersApplication) RebuildCommunityClusters() error {
	application.RebuildCommunityClustersCalls++

	return application.RebuildCommunityClustersErr
}

func (application *MockClustersApplication) RebuildPlatformClusters() error {
	application.RebuildPlatformClustersCalls++

	return application.RebuildPlatformClustersErr
}

func (application *MockClustersApplication) RebuildCampaignClusters() error {
	application.RebuildCampaignClustersCalls++

	return application.RebuildCampaignClustersErr
}

func (application *MockClustersApplication) RebuildTopicClusters() error {
	application.RebuildTopicClustersCalls++

	return application.RebuildTopicClustersErr
}

func (application *MockClustersApplication) RebuildNarrativeClusters() error {
	application.RebuildNarrativeClustersCalls++

	return application.RebuildNarrativeClustersErr
}
