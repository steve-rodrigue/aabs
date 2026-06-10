package evidences

import (
	"context"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/participatables"
	domain_evidences "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/evidences"
)

type application struct {
	repository domain_evidences.Repository
}

func createApplication(
	repository domain_evidences.Repository,
) Application {
	return &application{
		repository: repository,
	}
}

// FindByID finds evidence by id
func (app *application) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (domain_evidences.Evidence, error) {
	return app.repository.FindByID(ctx, id)
}

// FindByParticipation finds evidences by participation id
func (app *application) FindByParticipation(
	ctx context.Context,
	participation uuid.UUID,
) ([]domain_evidences.Evidence, error) {
	return app.repository.FindByParticipation(ctx, participation)
}

// FindByPost finds evidences by post id
func (app *application) FindByPost(
	ctx context.Context,
	post uuid.UUID,
) ([]domain_evidences.Evidence, error) {
	return app.repository.FindByPost(ctx, post)
}

// FindByParticipant finds evidences by participant
func (app *application) FindByParticipant(
	ctx context.Context,
	participant participatables.Participatable,
) ([]domain_evidences.Evidence, error) {
	return app.repository.FindByParticipant(ctx, participant)
}

// FindByTarget finds evidences by target
func (app *application) FindByTarget(
	ctx context.Context,
	target participatables.Participatable,
) ([]domain_evidences.Evidence, error) {
	return app.repository.FindByTarget(ctx, target)
}
