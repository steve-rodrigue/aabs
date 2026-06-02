package relationships

import (
	"github.com/google/uuid"

	relationship_comparables "github.com/steve-rodrigue/aabs/services/saas/applications/relationships/comparables"
	domain_relationships "github.com/steve-rodrigue/aabs/services/saas/domain/relationships"
	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships/relatables"
)

func NewMockRelationshipsApplication() *MockRelationshipsApplication {
	return &MockRelationshipsApplication{
		ComparablesValue: relationship_comparables.NewMockComparablesApplication(),
	}
}

type MockRelationshipsApplication struct {
	ComparablesCalls int
	ComparablesValue relationship_comparables.Application

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
}

func (application *MockRelationshipsApplication) Comparables() relationship_comparables.Application {
	application.ComparablesCalls++

	return application.ComparablesValue
}

func (application *MockRelationshipsApplication) Build(
	source relatables.Relatable,
	targets []relatables.Relatable,
) ([]domain_relationships.Relationship, error) {
	application.BuildCalls++

	return application.BuildValue, application.BuildErr
}

func (application *MockRelationshipsApplication) Sync(
	relationships []domain_relationships.Relationship,
) error {
	application.SyncCalls++

	return application.SyncErr
}

func (application *MockRelationshipsApplication) Find(
	index int,
	amount int,
) ([]domain_relationships.Relationship, error) {
	application.FindCalls++

	return application.FindValue, application.FindErr
}

func (application *MockRelationshipsApplication) FindByID(
	id uuid.UUID,
) (domain_relationships.Relationship, error) {
	application.FindByIDCalls++

	return application.FindByIDValue, application.FindByIDErr
}

func (application *MockRelationshipsApplication) FindAfter(
	cursor uuid.UUID,
	amount int,
) ([]domain_relationships.Relationship, error) {
	application.FindAfterCalls++

	return application.FindAfterValue, application.FindAfterErr
}

func (application *MockRelationshipsApplication) Count() (int64, error) {
	application.CountCalls++

	return application.CountValue, application.CountErr
}

func (application *MockRelationshipsApplication) RelationshipsBySource(
	id uuid.UUID,
) ([]domain_relationships.Relationship, error) {
	application.RelationshipsBySourceCalls++

	return application.RelationshipsBySourceValue, application.RelationshipsBySourceErr
}

func (application *MockRelationshipsApplication) RelationshipsByTarget(
	id uuid.UUID,
) ([]domain_relationships.Relationship, error) {
	application.RelationshipsByTargetCalls++

	return application.RelationshipsByTargetValue, application.RelationshipsByTargetErr
}

func (application *MockRelationshipsApplication) RebuildRelationships() error {
	application.RebuildRelationshipsCalls++

	return application.RebuildRelationshipsErr
}
