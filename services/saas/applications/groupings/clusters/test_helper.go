package clusters

import (
	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters"
	domain_clusters "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters/clusterables"
)

type MockClustersApplication struct {
	RebuildPostClustersCalls int
	RebuildPostClustersErr   error
}

func (application *MockClustersApplication) BuildForTarget(
	target clusterables.Clusterable,
	members []clusterables.Clusterable,
) ([]domain_clusters.Cluster, error) {
	return nil, nil
}

func (application *MockClustersApplication) FindByID(id uuid.UUID) (clusters.Cluster, error) {
	return nil, nil
}

func (application *MockClustersApplication) FindByTarget(target clusterables.Clusterable) ([]domain_clusters.Cluster, error) {
	return nil, nil
}

func (application *MockClustersApplication) FindByMember(member clusterables.Clusterable) ([]domain_clusters.Cluster, error) {
	return nil, nil
}

func (application *MockClustersApplication) RebuildAll() error {
	return nil
}

func (application *MockClustersApplication) RebuildPostClusters() error {
	application.RebuildPostClustersCalls++

	return application.RebuildPostClustersErr
}

func (application *MockClustersApplication) RebuildCampaignClusters() error {
	return nil
}

func (application *MockClustersApplication) RebuildTopicClusters() error {
	return nil
}

func (application *MockClustersApplication) RebuildNarrativeClusters() error {
	return nil
}
