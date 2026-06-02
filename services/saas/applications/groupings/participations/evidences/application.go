package evidences

import (
	"github.com/google/uuid"

	domain_evidences "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/evidences"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
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
	id uuid.UUID,
) (domain_evidences.Evidence, error) {
	return app.repository.FindByID(id)
}

// FindByParticipation finds evidences by participation id
func (app *application) FindByParticipation(
	participation uuid.UUID,
) ([]domain_evidences.Evidence, error) {
	return app.repository.FindByParticipation(participation)
}

// FindByPost finds evidences by post id
func (app *application) FindByPost(
	post uuid.UUID,
) ([]domain_evidences.Evidence, error) {
	return app.repository.FindByPost(post)
}

// FindByParticipant finds evidences by participant
func (app *application) FindByParticipant(
	participant participatables.Participatable,
) ([]domain_evidences.Evidence, error) {
	return app.repository.FindByParticipant(participant)
}

// FindByTarget finds evidences by target
func (app *application) FindByTarget(
	target participatables.Participatable,
) ([]domain_evidences.Evidence, error) {
	return app.repository.FindByTarget(target)
}
