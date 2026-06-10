package milvus

import (
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/clusterables"
)

const (
	groupingsClustersClusterablesCandidatesComparableIdentifierField = "identifier"
	groupingsClustersClusterablesCandidatesComparableKindField       = "kind"
	groupingsClustersClusterablesCandidatesComparableVectorField     = "vector"
)

// NewGroupingsClustersClusterablesComparableRepository creates a new milvus comparable repository
func NewGroupingsClustersClusterablesComparableRepository(
	client client.Client,
	adapter clusterables.ComparableAdapter,
	collection string,
) clusterables.ComparableRepository {
	return &groupingsClustersClusterablesComparableRepository{
		client:     client,
		adapter:    adapter,
		collection: collection,
	}
}
