package relationships

import (
	"github.com/google/uuid"
	domain_relationships "github.com/steve-rodrigue/aabs/services/saas/domain/relationships"
	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships/relatables"
)

// New creates a new relationships application
func New(
	repository domain_relationships.Repository,
	builder domain_relationships.Builder,
	relatableRepository relatables.Repository,
) Application {
	return createApplication(
		repository,
		builder,
		relatableRepository,
	)
}

// Application represents the relationships application
type Application interface {
	Build(source relatables.Relatable, targets []relatables.Relatable) ([]domain_relationships.Relationship, error)
	Sync(relationships []domain_relationships.Relationship) error

	FindAll() ([]domain_relationships.Relationship, error)
	RelationshipsBySource(id uuid.UUID) ([]domain_relationships.Relationship, error)
	RelationshipsByTarget(id uuid.UUID) ([]domain_relationships.Relationship, error)

	RebuildRelationships() error
}
