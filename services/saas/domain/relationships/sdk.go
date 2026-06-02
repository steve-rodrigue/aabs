package relationships

import (
	"time"

	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships/relatables"
)

// Builder represents a relationship builder
type Builder interface {
	Build(source relatables.Relatable, targets []relatables.Relatable) ([]Relationship, error)
}

// Relationship represents a relationship
type Relationship interface {
	Identifier() uuid.UUID
	Source() relatables.Relatable
	Target() relatables.Relatable
	Similarity() float64
	CreatedOn() time.Time
}

// Repository represents a relationship repository
type Repository interface {
	Save(relationship Relationship) error

	FindByID(id uuid.UUID) (Relationship, error)

	Find(index int, amount int) ([]Relationship, error)
	FindAfter(cursor uuid.UUID, amount int) ([]Relationship, error)
	Count() (int64, error)

	FindBySourceID(source uuid.UUID) ([]Relationship, error)
	FindByTargetID(target uuid.UUID) ([]Relationship, error)

	FindBySource(source relatables.Relatable) ([]Relationship, error)
	FindByTarget(target relatables.Relatable) ([]Relationship, error)
	FindBetween(source relatables.Relatable, target relatables.Relatable) (Relationship, error)
}
