package relationships

import (
	"github.com/google/uuid"

	relationship_comparables "github.com/steve-rodrigue/aabs/services/saas/applications/relationships/comparables"
	domain_relationships "github.com/steve-rodrigue/aabs/services/saas/domain/relationships"
	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships/relatables"
)

type application struct {
	repository          domain_relationships.Repository
	builder             domain_relationships.Builder
	relatableRepository relatables.Repository
	candidateRepository relatables.CandidateRepository
	comparables         relationship_comparables.Application
	rebuildBatchSize    int
}

func createApplication(
	repository domain_relationships.Repository,
	builder domain_relationships.Builder,
	relatableRepository relatables.Repository,
	candidateRepository relatables.CandidateRepository,
	comparables relationship_comparables.Application,
	rebuildBatchSize int,
) Application {
	return &application{
		repository:          repository,
		builder:             builder,
		relatableRepository: relatableRepository,
		candidateRepository: candidateRepository,
		comparables:         comparables,
		rebuildBatchSize:    rebuildBatchSize,
	}
}

func (app *application) Comparables() relationship_comparables.Application {
	return app.comparables
}

func (app *application) Build(
	source relatables.Relatable,
	targets []relatables.Relatable,
) ([]domain_relationships.Relationship, error) {
	return app.builder.Build(source, targets)
}

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

func (app *application) Find(
	index int,
	amount int,
) ([]domain_relationships.Relationship, error) {
	return app.repository.Find(index, amount)
}

func (app *application) FindByID(
	id uuid.UUID,
) (domain_relationships.Relationship, error) {
	return app.repository.FindByID(id)
}

func (app *application) FindAfter(
	cursor uuid.UUID,
	amount int,
) ([]domain_relationships.Relationship, error) {
	return app.repository.FindAfter(cursor, amount)
}

func (app *application) Count() (int64, error) {
	return app.repository.Count()
}

func (app *application) RelationshipsBySource(
	id uuid.UUID,
) ([]domain_relationships.Relationship, error) {
	return app.repository.FindBySourceID(id)
}

func (app *application) RelationshipsByTarget(
	id uuid.UUID,
) ([]domain_relationships.Relationship, error) {
	return app.repository.FindByTargetID(id)
}

func (app *application) RebuildRelationships() error {
	cursor := uuid.Nil

	for {
		sources, err := app.relatableRepository.FindAfter(
			cursor,
			app.rebuildBatchSize,
		)
		if err != nil {
			return err
		}

		if len(sources) == 0 {
			return nil
		}

		for _, source := range sources {
			targets, err := app.candidateRepository.FindCandidates(
				source,
				app.rebuildBatchSize,
			)
			if err != nil {
				return err
			}

			relationships, err := app.builder.Build(
				source,
				targets,
			)
			if err != nil {
				return err
			}

			if err := app.Sync(relationships); err != nil {
				return err
			}
		}

		cursor = sources[len(sources)-1].Identifier()
	}
}
