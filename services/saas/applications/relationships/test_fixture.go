package relationships

import (
	domain_relationships "github.com/steve-rodrigue/aabs/services/saas/domain/relationships"
	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships/relatables"
)

type applicationFixture struct {
	application Application

	repository          *domain_relationships.MockRelationshipRepository
	builder             *domain_relationships.MockRelationshipBuilder
	relatableRepository *relatables.MockRelatableRepository
}

func newApplicationFixture() *applicationFixture {
	repository := domain_relationships.NewMockRelationshipRepository()
	builder := domain_relationships.NewMockRelationshipBuilder()
	relatableRepository := relatables.NewMockRelatableRepository()

	application := New(
		repository,
		builder,
		relatableRepository,
	)

	return &applicationFixture{
		application:         application,
		repository:          repository,
		builder:             builder,
		relatableRepository: relatableRepository,
	}
}
