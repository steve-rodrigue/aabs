package relationships

import (
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

	FindAllCalls int
	FindAllErr   error
	FindAllValue []domain_relationships.Relationship

	RelationshipsBySourceCalls int
	RelationshipsBySourceErr   error
	RelationshipsBySourceValue []domain_relationships.Relationship

	RelationshipsByTargetCalls int
	RelationshipsByTargetErr   error
	RelationshipsByTargetValue []domain_relationships.Relationship

	RebuildRelationshipsCalls int
	RebuildRelationshipsErr   error
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

func (application *MockRelationshipsApplication) FindAll() (
	[]domain_relationships.Relationship,
	error,
) {
	application.FindAllCalls++

	return application.FindAllValue, application.FindAllErr
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
