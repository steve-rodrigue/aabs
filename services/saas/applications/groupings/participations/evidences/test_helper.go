package evidences

import (
	"context"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/participatables"
	domain_evidences "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/evidences"
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

	LastContext       context.Context
	LastID            uuid.UUID
	LastParticipation uuid.UUID
	LastPost          uuid.UUID
	LastParticipant   participatables.Participatable
	LastTarget        participatables.Participatable
}

func (application *MockEvidencesApplication) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (domain_evidences.Evidence, error) {
	application.FindByIDCalls++
	application.LastContext = ctx
	application.LastID = id

	return application.FindByIDValue, application.FindByIDErr
}

func (application *MockEvidencesApplication) FindByParticipation(
	ctx context.Context,
	participation uuid.UUID,
) ([]domain_evidences.Evidence, error) {
	application.FindByParticipationCalls++
	application.LastContext = ctx
	application.LastParticipation = participation

	return application.FindByParticipationValue,
		application.FindByParticipationErr
}

func (application *MockEvidencesApplication) FindByPost(
	ctx context.Context,
	post uuid.UUID,
) ([]domain_evidences.Evidence, error) {
	application.FindByPostCalls++
	application.LastContext = ctx
	application.LastPost = post

	return application.FindByPostValue,
		application.FindByPostErr
}

func (application *MockEvidencesApplication) FindByParticipant(
	ctx context.Context,
	participant participatables.Participatable,
) ([]domain_evidences.Evidence, error) {
	application.FindByParticipantCalls++
	application.LastContext = ctx
	application.LastParticipant = participant

	return application.FindByParticipantValue,
		application.FindByParticipantErr
}

func (application *MockEvidencesApplication) FindByTarget(
	ctx context.Context,
	target participatables.Participatable,
) ([]domain_evidences.Evidence, error) {
	application.FindByTargetCalls++
	application.LastContext = ctx
	application.LastTarget = target

	return application.FindByTargetValue,
		application.FindByTargetErr
}
