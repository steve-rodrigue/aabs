package relationships

import (
	"context"

	"github.com/google/uuid"

	domain_relationships "github.com/steve-rodrigue/aabs/services/saas/domain/relationships"
	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships/relatables"
)

func NewMockRelationshipsApplication() *MockRelationshipsApplication {
	return &MockRelationshipsApplication{}
}

type MockRelationshipsApplication struct {
	BuildCalls int
	BuildErr   error
	BuildValue []domain_relationships.Relationship

	SyncCalls int
	SyncErr   error

	FindCalls int
	FindErr   error
	FindValue []domain_relationships.Relationship

	FindAfterCalls int
	FindAfterErr   error
	FindAfterValue []domain_relationships.Relationship

	CountCalls int
	CountErr   error
	CountValue int64

	FindByIDCalls int
	FindByIDErr   error
	FindByIDValue domain_relationships.Relationship

	RelationshipsBySourceCalls int
	RelationshipsBySourceErr   error
	RelationshipsBySourceValue []domain_relationships.Relationship

	RelationshipsByTargetCalls int
	RelationshipsByTargetErr   error
	RelationshipsByTargetValue []domain_relationships.Relationship

	RebuildRelationshipsCalls int
	RebuildRelationshipsErr   error

	LastContext       context.Context
	LastSource        relatables.Relatable
	LastTargets       []relatables.Relatable
	LastRelationships []domain_relationships.Relationship
	LastIndex         int
	LastAmount        int
	LastID            uuid.UUID
	LastCursor        uuid.UUID
}

func (application *MockRelationshipsApplication) Build(
	source relatables.Relatable,
	targets []relatables.Relatable,
) ([]domain_relationships.Relationship, error) {
	application.BuildCalls++
	application.LastSource = source
	application.LastTargets = targets

	return application.BuildValue, application.BuildErr
}

func (application *MockRelationshipsApplication) Sync(
	ctx context.Context,
	relationships []domain_relationships.Relationship,
) error {
	application.SyncCalls++
	application.LastContext = ctx
	application.LastRelationships = relationships

	return application.SyncErr
}

func (application *MockRelationshipsApplication) Find(
	ctx context.Context,
	index int,
	amount int,
) ([]domain_relationships.Relationship, error) {
	application.FindCalls++
	application.LastContext = ctx
	application.LastIndex = index
	application.LastAmount = amount

	return application.FindValue, application.FindErr
}

func (application *MockRelationshipsApplication) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (domain_relationships.Relationship, error) {
	application.FindByIDCalls++
	application.LastContext = ctx
	application.LastID = id

	return application.FindByIDValue, application.FindByIDErr
}

func (application *MockRelationshipsApplication) FindAfter(
	ctx context.Context,
	cursor uuid.UUID,
	amount int,
) ([]domain_relationships.Relationship, error) {
	application.FindAfterCalls++
	application.LastContext = ctx
	application.LastCursor = cursor
	application.LastAmount = amount

	return application.FindAfterValue, application.FindAfterErr
}

func (application *MockRelationshipsApplication) Count(
	ctx context.Context,
) (int64, error) {
	application.CountCalls++
	application.LastContext = ctx

	return application.CountValue, application.CountErr
}

func (application *MockRelationshipsApplication) RelationshipsBySource(
	ctx context.Context,
	id uuid.UUID,
) ([]domain_relationships.Relationship, error) {
	application.RelationshipsBySourceCalls++
	application.LastContext = ctx
	application.LastID = id

	return application.RelationshipsBySourceValue,
		application.RelationshipsBySourceErr
}

func (application *MockRelationshipsApplication) RelationshipsByTarget(
	ctx context.Context,
	id uuid.UUID,
) ([]domain_relationships.Relationship, error) {
	application.RelationshipsByTargetCalls++
	application.LastContext = ctx
	application.LastID = id

	return application.RelationshipsByTargetValue,
		application.RelationshipsByTargetErr
}

func (application *MockRelationshipsApplication) RebuildRelationships(
	ctx context.Context,
) error {
	application.RebuildRelationshipsCalls++
	application.LastContext = ctx

	return application.RebuildRelationshipsErr
}
