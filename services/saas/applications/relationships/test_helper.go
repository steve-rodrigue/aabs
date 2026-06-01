package relationships

import (
	"github.com/google/uuid"
	domain_relationships "github.com/steve-rodrigue/aabs/services/saas/domain/relationships"
	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships/relatables"
)

type MockRelationshipsApplication struct {
	RebuildRelationshipsCalls int
	RebuildRelationshipsErr   error
}

func (application *MockRelationshipsApplication) Build(
	source relatables.Relatable,
	targets []relatables.Relatable,
) ([]domain_relationships.Relationship, error) {
	return nil, nil
}

func (application *MockRelationshipsApplication) Sync(relationships []domain_relationships.Relationship) error {
	return nil
}

func (application *MockRelationshipsApplication) FindAll() ([]domain_relationships.Relationship, error) {
	return nil, nil
}

func (application *MockRelationshipsApplication) RelationshipsBySource(id uuid.UUID) ([]domain_relationships.Relationship, error) {
	return nil, nil
}

func (application *MockRelationshipsApplication) RelationshipsByTarget(id uuid.UUID) ([]domain_relationships.Relationship, error) {
	return nil, nil
}

func (application *MockRelationshipsApplication) RebuildRelationships() error {
	application.RebuildRelationshipsCalls++

	return application.RebuildRelationshipsErr
}
