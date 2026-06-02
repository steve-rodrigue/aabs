package evidences

import (
	"github.com/google/uuid"

	domain_evidences "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/evidences"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
)

func NewMockEvidencesApplication() *MockEvidencesApplication {
	return &MockEvidencesApplication{}
}

type MockEvidencesApplication struct {
	FindByIDCalls int
	FindByIDErr   error
	FindByIDValue domain_evidences.Evidence

	FindByParticipationCalls int
	FindByParticipationErr   error
	FindByParticipationValue []domain_evidences.Evidence

	FindByPostCalls int
	FindByPostErr   error
	FindByPostValue []domain_evidences.Evidence

	FindByParticipantCalls int
	FindByParticipantErr   error
	FindByParticipantValue []domain_evidences.Evidence

	FindByTargetCalls int
	FindByTargetErr   error
	FindByTargetValue []domain_evidences.Evidence
}

func (application *MockEvidencesApplication) FindByID(
	id uuid.UUID,
) (domain_evidences.Evidence, error) {
	application.FindByIDCalls++

	return application.FindByIDValue, application.FindByIDErr
}

func (application *MockEvidencesApplication) FindByParticipation(
	participation uuid.UUID,
) ([]domain_evidences.Evidence, error) {
	application.FindByParticipationCalls++

	return application.FindByParticipationValue, application.FindByParticipationErr
}

func (application *MockEvidencesApplication) FindByPost(
	post uuid.UUID,
) ([]domain_evidences.Evidence, error) {
	application.FindByPostCalls++

	return application.FindByPostValue, application.FindByPostErr
}

func (application *MockEvidencesApplication) FindByParticipant(
	participant participatables.Participatable,
) ([]domain_evidences.Evidence, error) {
	application.FindByParticipantCalls++

	return application.FindByParticipantValue, application.FindByParticipantErr
}

func (application *MockEvidencesApplication) FindByTarget(
	target participatables.Participatable,
) ([]domain_evidences.Evidence, error) {
	application.FindByTargetCalls++

	return application.FindByTargetValue, application.FindByTargetErr
}
