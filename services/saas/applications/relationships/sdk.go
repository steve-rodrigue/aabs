package relationships

import (
	"github.com/google/uuid"

	relationship_comparables "github.com/steve-rodrigue/aabs/services/saas/applications/relationships/comparables"
	domain_relationships "github.com/steve-rodrigue/aabs/services/saas/domain/relationships"
	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships/relatables"
)

// New creates a new relationships application
func New(
	repository domain_relationships.Repository,
	builder domain_relationships.Builder,
	relatableRepository relatables.Repository,
	candidateRepository relatables.CandidateRepository,
	comparables relationship_comparables.Application,
	rebuildBatchSize int,
) Application {
	return createApplication(
		repository,
		builder,
		relatableRepository,
		candidateRepository,
		comparables,
		rebuildBatchSize,
	)
}

// Application represents the relationships application
type Application interface {
	Comparables() relationship_comparables.Application

	Build(source relatables.Relatable, targets []relatables.Relatable) ([]domain_relationships.Relationship, error)
	Sync(relationships []domain_relationships.Relationship) error

	Find(index int, amount int) ([]domain_relationships.Relationship, error)
	FindByID(id uuid.UUID) (domain_relationships.Relationship, error)
	FindAfter(cursor uuid.UUID, amount int) ([]domain_relationships.Relationship, error)
	Count() (int64, error)

	RelationshipsBySource(id uuid.UUID) ([]domain_relationships.Relationship, error)
	RelationshipsByTarget(id uuid.UUID) ([]domain_relationships.Relationship, error)

	RebuildRelationships() error
}
