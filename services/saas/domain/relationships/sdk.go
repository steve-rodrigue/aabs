package relationships

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/relatables"
)

var (
	ErrInvalidRelationshipIdentifier = errors.New("invalid relationship identifier")
	ErrInvalidRelationshipSource     = errors.New("invalid relationship source")
	ErrInvalidRelationshipTarget     = errors.New("invalid relationship target")
	ErrInvalidRelationshipSimilarity = errors.New("invalid relationship similarity")
	ErrInvalidRelationshipCreatedOn  = errors.New("invalid relationship created on")
)

// NewAdapter creates a new relationship adapter
func NewAdapter() Adapter {
	return &adapter{}
}

// RelationshipInput represents a relationship input
type RelationshipInput struct {
	Identifier uuid.UUID
	Source     relatables.Relatable
	Target     relatables.Relatable
	Similarity float64
	CreatedOn  time.Time
}

// Adapter represents a relationship adapter
type Adapter interface {
	ToDomain(input RelationshipInput) (Relationship, error)
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
	Save(ctx context.Context, relationship Relationship) error
	FindByID(ctx context.Context, id uuid.UUID) (Relationship, error)
	Find(ctx context.Context, index int, amount int) ([]Relationship, error)
	FindAfter(ctx context.Context, cursor uuid.UUID, amount int) ([]Relationship, error)
	Count(ctx context.Context) (int64, error)
	FindBySourceID(ctx context.Context, source uuid.UUID) ([]Relationship, error)
	FindByTargetID(ctx context.Context, target uuid.UUID) ([]Relationship, error)
	FindBySource(ctx context.Context, source relatables.Relatable) ([]Relationship, error)
	FindByTarget(ctx context.Context, target relatables.Relatable) ([]Relationship, error)
	FindBetween(
		ctx context.Context,
		source relatables.Relatable,
		target relatables.Relatable,
	) (Relationship, error)
}
