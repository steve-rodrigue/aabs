package relationships

import (
	domain_relationships "github.com/steve-rodrigue/aabs/services/saas/domain/relationships"
	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships/relatables"
)

// Builder builds graph relationships between relatable entities
type Builder interface {
	Build(source relatables.Relatable, targets []relatables.Relatable) ([]domain_relationships.Relationship, error)
}

// Synchronizer persists relationships to the graph store
type Synchronizer interface {
	Sync(relationships []domain_relationships.Relationship) error
}
