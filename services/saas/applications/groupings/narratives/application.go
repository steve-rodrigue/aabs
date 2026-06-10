package narratives

import (
	"github.com/google/uuid"

	app_participations "github.com/steve-rodrigue/aabs/services/saas/applications/groupings/participations"
	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/participatables"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/communities"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/users"
	domain_narratives "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/narratives"
)

type application struct {
	repository     domain_narratives.Repository
	participations app_participations.Application
}

func createApplication(
	repository domain_narratives.Repository,
	participations app_participations.Application,
) Application {
	return &application{
		repository:     repository,
		participations: participations,
	}
}

func (app *application) FindByID(
	id uuid.UUID,
) (domain_narratives.Narrative, error) {
	return app.repository.FindByID(id)
}

func (app *application) Find(
	index int,
	amount int,
) ([]domain_narratives.Narrative, error) {
	return app.repository.Find(index, amount)
}

func (app *application) FindAfter(
	cursor uuid.UUID,
	amount int,
) ([]domain_narratives.Narrative, error) {
	return app.repository.FindAfter(cursor, amount)
}

func (app *application) FindNarrativesByUser(
	user users.User,
) ([]domain_narratives.Narrative, error) {
	return app.findNarrativesByParticipant(user)
}

func (app *application) FindNarrativesByCommunity(
	community communities.Community,
) ([]domain_narratives.Narrative, error) {
	return app.findNarrativesByParticipant(community)
}

func (app *application) findNarrativesByParticipant(
	participant participatables.Participatable,
) ([]domain_narratives.Narrative, error) {
	participations, err := app.participations.FindByParticipant(participant)
	if err != nil {
		return nil, err
	}

	out := make([]domain_narratives.Narrative, 0, len(participations))

	for _, participation := range participations {
		if participation.Target().ParticipationKind() != participatables.NarrativeKind {
			continue
		}

		narrative, err := app.repository.FindByID(
			participation.Target().Identifier(),
		)
		if err != nil {
			return nil, err
		}

		out = append(out, narrative)
	}

	return out, nil
}

func (app *application) Count() (int64, error) {
	return app.repository.Count()
}

func (app *application) RebuildNarratives() error {
	return nil
}
