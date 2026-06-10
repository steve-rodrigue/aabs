package relationships

import (
	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/relatables"
	domain_relationships "github.com/steve-rodrigue/aabs/services/saas/domain/relationships"
	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships/builders"
	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships/comparables"
)

type applicationFixture struct {
	application Application

	repository *domain_relationships.MockRelationshipRepository
	builder    *builders.MockBuilder
	relatables *relatables.MockRelatableRepository
	candidates *relatables.MockCandidateRepository
	comparator *comparables.MockComparator
}

func newApplicationFixture() *applicationFixture {
	repository := domain_relationships.NewMockRelationshipRepository()
	builder := builders.NewMockBuilder()
	relatableRepository := relatables.NewMockRelatableRepository()
	candidateRepository := relatables.NewMockCandidateRepository()
	comparator := comparables.NewMockComparator()

	application := New(
		repository,
		builder,
		relatableRepository,
		candidateRepository,
		comparator,
		25,
	)

	return &applicationFixture{
		application: application,

		repository: repository,
		builder:    builder,
		relatables: relatableRepository,
		candidates: candidateRepository,
		comparator: comparator,
	}
}
