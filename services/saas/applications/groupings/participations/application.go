package participations

import (
	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/applications/groupings/participations/evidences"
	domain_participations "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations"
	domain_evidences "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/evidences"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
)

type application struct {
	repository domain_participations.Repository
	calculator domain_participations.Calculator

	participatableRepository participatables.Repository

	evidenceCalculator  domain_evidences.Calculator
	evidenceApplication evidences.Application
}

func createApplication(
	repository domain_participations.Repository,
	calculator domain_participations.Calculator,
	participatableRepository participatables.Repository,
	evidenceRepository domain_evidences.Repository,
	evidenceCalculator domain_evidences.Calculator,
	evidenceApplication evidences.Application,
) Application {
	return &application{
		repository: repository,
		calculator: calculator,

		participatableRepository: participatableRepository,

		evidenceRepository:  evidenceRepository,
		evidenceCalculator:  evidenceCalculator,
		evidenceApplication: evidenceApplication,
	}
}

// Evidences returns the evidence application
func (app *application) Evidences() evidences.Application {
	return app.evidenceApplication
}

// FindByID finds a participation by id
func (app *application) FindByID(
	id uuid.UUID,
) (domain_participations.Participation, error) {
	return app.repository.FindByID(id)
}

// FindByParticipant finds participations by participant
func (app *application) FindByParticipant(
	participant participatables.Participatable,
) ([]domain_participations.Participation, error) {
	return app.repository.FindByParticipant(participant)
}

// FindByTarget finds participations by target
func (app *application) FindByTarget(
	target participatables.Participatable,
) ([]domain_participations.Participation, error) {
	return app.repository.FindByTarget(target)
}

// FindBetween finds a participation between a participant and a target
func (app *application) FindBetween(
	participant participatables.Participatable,
	target participatables.Participatable,
) (domain_participations.Participation, error) {
	return app.repository.FindBetween(participant, target)
}

// RebuildParticipations rebuilds all participations and evidences
func (app *application) RebuildParticipations() error {
	participants, err := app.participatableRepository.FindAllParticipants()
	if err != nil {
		return err
	}

	targets, err := app.participatableRepository.FindAllTargets()
	if err != nil {
		return err
	}

	for _, participant := range participants {
		for _, target := range targets {
			if participant.Identifier() == target.Identifier() {
				continue
			}

			participation, err := app.calculator.Calculate(
				participant,
				target,
			)
			if err != nil {
				return err
			}

			if err := app.repository.Save(participation); err != nil {
				return err
			}

			evidences, err := app.evidenceCalculator.Calculate(participation)
			if err != nil {
				return err
			}

			for _, evidence := range evidences {
				if err := app.evidenceRepository.Save(evidence); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
