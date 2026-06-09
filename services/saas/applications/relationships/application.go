package relationships

import (
	"context"

	"github.com/google/uuid"

	domain_relationships "github.com/steve-rodrigue/aabs/services/saas/domain/relationships"
	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships/builders"
	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships/comparables"
	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships/relatables"
)

type application struct {
	repository          domain_relationships.Repository
	builder             builders.Builder
	relatableRepository relatables.Repository
	candidateRepository relatables.CandidateRepository
	comparator          comparables.Comparator
	rebuildBatchSize    int
}

func createApplication(
	repository domain_relationships.Repository,
	builder builders.Builder,
	relatableRepository relatables.Repository,
	candidateRepository relatables.CandidateRepository,
	comparator comparables.Comparator,
	rebuildBatchSize int,
) Application {
	return &application{
		repository:          repository,
		builder:             builder,
		relatableRepository: relatableRepository,
		candidateRepository: candidateRepository,
		comparator:          comparator,
		rebuildBatchSize:    rebuildBatchSize,
	}
}

func (app *application) Build(
	source relatables.Relatable,
	targets []relatables.Relatable,
) ([]domain_relationships.Relationship, error) {
	return app.builder.Build(source, targets)
}

func (app *application) Sync(
	ctx context.Context,
	relationships []domain_relationships.Relationship,
) error {
	for _, relationship := range relationships {
		if err := app.repository.Save(ctx, relationship); err != nil {
			return err
		}
	}

	return nil
}

func (app *application) Find(
	ctx context.Context,
	index int,
	amount int,
) ([]domain_relationships.Relationship, error) {
	return app.repository.Find(ctx, index, amount)
}

func (app *application) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (domain_relationships.Relationship, error) {
	return app.repository.FindByID(ctx, id)
}

func (app *application) FindAfter(
	ctx context.Context,
	cursor uuid.UUID,
	amount int,
) ([]domain_relationships.Relationship, error) {
	return app.repository.FindAfter(ctx, cursor, amount)
}

func (app *application) Count(
	ctx context.Context,
) (int64, error) {
	return app.repository.Count(ctx)
}

func (app *application) RelationshipsBySource(
	ctx context.Context,
	id uuid.UUID,
) ([]domain_relationships.Relationship, error) {
	return app.repository.FindBySourceID(ctx, id)
}

func (app *application) RelationshipsByTarget(
	ctx context.Context,
	id uuid.UUID,
) ([]domain_relationships.Relationship, error) {
	return app.repository.FindByTargetID(ctx, id)
}

func (app *application) RebuildRelationships(
	ctx context.Context,
) error {
	cursor := uuid.Nil

	for {
		sources, err := app.relatableRepository.FindAfter(
			ctx,
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
				ctx,
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

			if err := app.Sync(ctx, relationships); err != nil {
				return err
			}
		}

		cursor = sources[len(sources)-1].Identifier()
	}
}
