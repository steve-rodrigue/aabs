package relationships

import (
	"github.com/google/uuid"
	domain_relationships "github.com/steve-rodrigue/aabs/services/saas/domain/relationships"
	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships/relatables"
)

// Application represents the relationships application
type Application interface {
	Build(source relatables.Relatable, targets []relatables.Relatable) ([]domain_relationships.Relationship, error)
	Sync(relationships []domain_relationships.Relationship) error

	FindAll() ([]domain_relationships.Relationship, error)
	RelationshipsBySource(id uuid.UUID) ([]domain_relationships.Relationship, error)
	RelationshipsByTarget(id uuid.UUID) ([]domain_relationships.Relationship, error)

	RebuildRelationships() error
}
