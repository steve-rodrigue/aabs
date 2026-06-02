package evidences

import (
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
	"github.com/steve-rodrigue/aabs/services/saas/domain/posts"
)

func NewMockEvidence() Evidence {
	return &MockEvidence{
		id: uuid.New(),
	}
}

func NewMockEvidenceRepository() *MockEvidenceRepository {
	return &MockEvidenceRepository{}
}

func NewMockEvidenceCalculator() *MockEvidenceCalculator {
	return &MockEvidenceCalculator{}
}

type MockEvidence struct {
	id uuid.UUID

	participation participations.Participation

	participant participatables.Participatable
	target      participatables.Participatable

	post  posts.Post
	score float64
}

func (evidence *MockEvidence) Identifier() uuid.UUID {
	return evidence.id
}

func (evidence *MockEvidence) Participation() participations.Participation {
	return evidence.participation
}

func (evidence *MockEvidence) Participant() participatables.Participatable {
	return evidence.participant
}

func (evidence *MockEvidence) Target() participatables.Participatable {
	return evidence.target
}

func (evidence *MockEvidence) Post() posts.Post {
	return evidence.post
}

func (evidence *MockEvidence) Score() float64 {
	return evidence.score
}

func (evidence *MockEvidence) DetectedOn() time.Time {
	return time.Time{}
}

type MockEvidenceRepository struct {
	SaveCalls int
	SaveErr   error

	FindByIDCalls int
	FindByIDErr   error
	FindByIDValue Evidence

	FindByParticipationCalls int
	FindByParticipationErr   error
	FindByParticipationValue []Evidence

	FindByPostCalls int
	FindByPostErr   error
	FindByPostValue []Evidence

	FindByParticipantCalls int
	FindByParticipantErr   error
	FindByParticipantValue []Evidence

	FindByTargetCalls int
	FindByTargetErr   error
	FindByTargetValue []Evidence
}

func (repository *MockEvidenceRepository) Save(
	evidence Evidence,
) error {
	repository.SaveCalls++

	return repository.SaveErr
}

func (repository *MockEvidenceRepository) FindByID(
	id uuid.UUID,
) (Evidence, error) {
	repository.FindByIDCalls++

	return repository.FindByIDValue, repository.FindByIDErr
}

func (repository *MockEvidenceRepository) FindByParticipation(
	participation uuid.UUID,
) ([]Evidence, error) {
	repository.FindByParticipationCalls++

	return repository.FindByParticipationValue, repository.FindByParticipationErr
}

func (repository *MockEvidenceRepository) FindByPost(
	post uuid.UUID,
) ([]Evidence, error) {
	repository.FindByPostCalls++

	return repository.FindByPostValue, repository.FindByPostErr
}

func (repository *MockEvidenceRepository) FindByParticipant(
	participant participatables.Participatable,
) ([]Evidence, error) {
	repository.FindByParticipantCalls++

	return repository.FindByParticipantValue, repository.FindByParticipantErr
}

func (repository *MockEvidenceRepository) FindByTarget(
	target participatables.Participatable,
) ([]Evidence, error) {
	repository.FindByTargetCalls++

	return repository.FindByTargetValue, repository.FindByTargetErr
}

type MockEvidenceCalculator struct {
	CalculateCalls int
	CalculateErr   error

	LastParticipation participations.Participation

	Value []Evidence
}

func (calculator *MockEvidenceCalculator) Calculate(
	participation participations.Participation,
) ([]Evidence, error) {
	calculator.CalculateCalls++
	calculator.LastParticipation = participation

	return calculator.Value, calculator.CalculateErr
}
