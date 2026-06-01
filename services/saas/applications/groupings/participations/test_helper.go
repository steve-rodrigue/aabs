package participations

import (
	"github.com/google/uuid"

	domain_participations "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
)

func NewMockParticipationsApplication() *MockParticipationsApplication {
	return &MockParticipationsApplication{}
}

type MockParticipationsApplication struct {
	FindByIDCalls int
	FindByIDErr   error
	FindByIDValue domain_participations.Participation

	FindByParticipantCalls int
	FindByParticipantErr   error
	FindByParticipantValue []domain_participations.Participation

	FindByTargetCalls int
	FindByTargetErr   error
	FindByTargetValue []domain_participations.Participation

	FindBetweenCalls int
	FindBetweenErr   error
	FindBetweenValue domain_participations.Participation

	RebuildParticipationsCalls int
	RebuildParticipationsErr   error
}

func (application *MockParticipationsApplication) FindByID(
	id uuid.UUID,
) (domain_participations.Participation, error) {
	application.FindByIDCalls++

	return application.FindByIDValue, application.FindByIDErr
}

func (application *MockParticipationsApplication) FindByParticipant(
	participant participatables.Participatable,
) ([]domain_participations.Participation, error) {
	application.FindByParticipantCalls++

	return application.FindByParticipantValue, application.FindByParticipantErr
}

func (application *MockParticipationsApplication) FindByTarget(
	target participatables.Participatable,
) ([]domain_participations.Participation, error) {
	application.FindByTargetCalls++

	return application.FindByTargetValue, application.FindByTargetErr
}

func (application *MockParticipationsApplication) FindBetween(
	participant participatables.Participatable,
	target participatables.Participatable,
) (domain_participations.Participation, error) {
	application.FindBetweenCalls++

	return application.FindBetweenValue, application.FindBetweenErr
}

func (application *MockParticipationsApplication) RebuildParticipations() error {
	application.RebuildParticipationsCalls++

	return application.RebuildParticipationsErr
}
