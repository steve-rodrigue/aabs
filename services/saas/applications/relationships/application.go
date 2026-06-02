package relationships

import (
	"github.com/google/uuid"

	domain_relationships "github.com/steve-rodrigue/aabs/services/saas/domain/relationships"
	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships/relatables"
)

type application struct {
	repository          domain_relationships.Repository
	builder             domain_relationships.Builder
	relatableRepository relatables.Repository
}

func createApplication(
	repository domain_relationships.Repository,
	builder domain_relationships.Builder,
	relatableRepository relatables.Repository,
) Application {
	return &application{
		repository:          repository,
		builder:             builder,
		relatableRepository: relatableRepository,
	}
}

// Build builds relationships between a source and targets
func (app *application) Build(
	source relatables.Relatable,
	targets []relatables.Relatable,
) ([]domain_relationships.Relationship, error) {
	return app.builder.Build(source, targets)
}

// Sync stores relationships
func (app *application) Sync(
	relationships []domain_relationships.Relationship,
) error {
	for _, relationship := range relationships {
		if err := app.repository.Save(relationship); err != nil {
			return err
		}
	}

	return nil
}

// FindAll finds all relationships
func (app *application) FindAll() ([]domain_relationships.Relationship, error) {
	return app.repository.FindAll()
}

// RelationshipsBySource finds relationships by source id
func (app *application) RelationshipsBySource(
	id uuid.UUID,
) ([]domain_relationships.Relationship, error) {
	return app.repository.FindBySourceID(id)
}

// RelationshipsByTarget finds relationships by target id
func (app *application) RelationshipsByTarget(
	id uuid.UUID,
) ([]domain_relationships.Relationship, error) {
	return app.repository.FindByTargetID(id)
}

// RebuildRelationships rebuilds all relationships between relatable entities
func (app *application) RebuildRelationships() error {
	relats, err := app.relatableRepository.FindAll()
	if err != nil {
		return err
	}

	for _, source := range relats {
		targets := make([]relatables.Relatable, 0, len(relats)-1)

		for _, target := range relats {
			if source.Identifier() == target.Identifier() {
				continue
			}

			targets = append(targets, target)
		}

		relationships, err := app.builder.Build(source, targets)
		if err != nil {
			return err
		}

		if err := app.Sync(relationships); err != nil {
			return err
		}
	}

	return nil
}
