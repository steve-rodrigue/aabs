package relationships

import (
	relationship_comparables "github.com/steve-rodrigue/aabs/services/saas/applications/relationships/comparables"
	domain_relationships "github.com/steve-rodrigue/aabs/services/saas/domain/relationships"
	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships/relatables"
)

type applicationFixture struct {
	application Application

	repository  *domain_relationships.MockRelationshipRepository
	builder     *domain_relationships.MockRelationshipBuilder
	relatables  *relatables.MockRelatableRepository
	candidates  *relatables.MockCandidateRepository
	comparables *relationship_comparables.MockComparablesApplication
}

func newApplicationFixture() *applicationFixture {
	repository := domain_relationships.NewMockRelationshipRepository()
	builder := domain_relationships.NewMockRelationshipBuilder()
	relatableRepository := relatables.NewMockRelatableRepository()
	candidateRepository := relatables.NewMockCandidateRepository()
	comparables := relationship_comparables.NewMockComparablesApplication()

	application := New(
		repository,
		builder,
		relatableRepository,
		candidateRepository,
		comparables,
		25,
	)

	return &applicationFixture{
		application: application,

		repository:  repository,
		builder:     builder,
		relatables:  relatableRepository,
		candidates:  candidateRepository,
		comparables: comparables,
	}
}
