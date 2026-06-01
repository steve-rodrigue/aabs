package participations

import (
	"github.com/google/uuid"
	domain_participations "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
)

type MockParticipationsApplication struct {
	RebuildParticipationsCalls int
	RebuildParticipationsErr   error
}

func (application *MockParticipationsApplication) FindByID(id uuid.UUID) (domain_participations.Participation, error) {
	return nil, nil
}

func (application *MockParticipationsApplication) FindByParticipant(
	participant participatables.Participatable,
) ([]domain_participations.Participation, error) {
	return nil, nil
}

func (application *MockParticipationsApplication) FindByTarget(
	target participatables.Participatable,
) ([]domain_participations.Participation, error) {
	return nil, nil
}

func (application *MockParticipationsApplication) FindBetween(
	participant participatables.Participatable,
	target participatables.Participatable,
) (domain_participations.Participation, error) {
	return nil, nil
}

func (application *MockParticipationsApplication) RebuildParticipations() error {
	application.RebuildParticipationsCalls++

	return application.RebuildParticipationsErr
}
