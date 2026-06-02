package clusters

import (
	domain_clusters "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters/clusterables"
)

type applicationFixture struct {
	application Application

	repository   *domain_clusters.MockClusterRepository
	detector     *domain_clusters.MockClusterDetector
	clusterables *clusterables.MockClusterableRepository
	candidates   *clusterables.MockCandidateRepository
}

func newApplicationFixture() *applicationFixture {
	repository := domain_clusters.NewMockClusterRepository()
	detector := domain_clusters.NewMockClusterDetector()
	clusterableRepository := clusterables.NewMockClusterableRepository()
	candidateRepository := clusterables.NewMockCandidateRepository()

	application := New(
		repository,
		detector,
		clusterableRepository,
		candidateRepository,
	)

	return &applicationFixture{
		application:  application,
		repository:   repository,
		detector:     detector,
		clusterables: clusterableRepository,
		candidates:   candidateRepository,
	}
}
