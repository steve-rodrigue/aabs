package relationships

import (
	"context"

	"github.com/google/uuid"

	domain_relationships "github.com/steve-rodrigue/aabs/services/saas/domain/relationships"
	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships/builders"
	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships/comparables"
	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships/relatables"
)

// New creates a new relationships application
func New(
	repository domain_relationships.Repository,
	builder builders.Builder,
	relatableRepository relatables.Repository,
	candidateRepository relatables.CandidateRepository,
	comparator comparables.Comparator,
	rebuildBatchSize int,
) Application {
	return createApplication(
		repository,
		builder,
		relatableRepository,
		candidateRepository,
		comparator,
		rebuildBatchSize,
	)
}

// Application represents the relationships application
type Application interface {
	Build(
		source relatables.Relatable,
		targets []relatables.Relatable,
	) ([]domain_relationships.Relationship, error)

	Sync(
		ctx context.Context,
		relationships []domain_relationships.Relationship,
	) error

	Find(
		ctx context.Context,
		index int,
		amount int,
	) ([]domain_relationships.Relationship, error)

	FindByID(
		ctx context.Context,
		id uuid.UUID,
	) (domain_relationships.Relationship, error)

	FindAfter(
		ctx context.Context,
		cursor uuid.UUID,
		amount int,
	) ([]domain_relationships.Relationship, error)

	Count(
		ctx context.Context,
	) (int64, error)

	RelationshipsBySource(
		ctx context.Context,
		id uuid.UUID,
	) ([]domain_relationships.Relationship, error)

	RelationshipsByTarget(
		ctx context.Context,
		id uuid.UUID,
	) ([]domain_relationships.Relationship, error)

	RebuildRelationships(
		ctx context.Context,
	) error
}
