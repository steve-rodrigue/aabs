package participations

import (
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
)

func NewMockParticipation() Participation {
	return &MockParticipation{
		id: uuid.New(),
	}
}

func NewMockParticipationRepository() *MockParticipationRepository {
	return &MockParticipationRepository{}
}

func NewMockParticipationCalculator() *MockParticipationCalculator {
	return &MockParticipationCalculator{}
}

type MockParticipation struct {
	id uuid.UUID

	participant participatables.Participatable
	target      participatables.Participatable

	postCount      int
	totalPostCount int
	percentage     float64
}

func (participation *MockParticipation) Identifier() uuid.UUID {
	return participation.id
}

func (participation *MockParticipation) Participant() participatables.Participatable {
	return participation.participant
}

func (participation *MockParticipation) Target() participatables.Participatable {
	return participation.target
}

func (participation *MockParticipation) PostCount() int {
	return participation.postCount
}

func (participation *MockParticipation) TotalPostCount() int {
	return participation.totalPostCount
}

func (participation *MockParticipation) Percentage() float64 {
	return participation.percentage
}

func (participation *MockParticipation) DetectedOn() time.Time {
	return time.Time{}
}

type MockParticipationRepository struct {
	SaveCalls int
	SaveErr   error

	FindByIDCalls int
	FindByIDErr   error
	FindByIDValue Participation

	FindByParticipantCalls int
	FindByParticipantErr   error
	FindByParticipantValue []Participation

	FindByTargetCalls int
	FindByTargetErr   error
	FindByTargetValue []Participation

	FindBetweenCalls int
	FindBetweenErr   error
	FindBetweenValue Participation
}

func (repository *MockParticipationRepository) Save(
	participation Participation,
) error {
	repository.SaveCalls++

	return repository.SaveErr
}

func (repository *MockParticipationRepository) FindByID(
	id uuid.UUID,
) (Participation, error) {
	repository.FindByIDCalls++

	return repository.FindByIDValue, repository.FindByIDErr
}

func (repository *MockParticipationRepository) FindByParticipant(
	participant participatables.Participatable,
) ([]Participation, error) {
	repository.FindByParticipantCalls++

	return repository.FindByParticipantValue, repository.FindByParticipantErr
}

func (repository *MockParticipationRepository) FindByTarget(
	target participatables.Participatable,
) ([]Participation, error) {
	repository.FindByTargetCalls++

	return repository.FindByTargetValue, repository.FindByTargetErr
}

func (repository *MockParticipationRepository) FindBetween(
	participant participatables.Participatable,
	target participatables.Participatable,
) (Participation, error) {
	repository.FindBetweenCalls++

	return repository.FindBetweenValue, repository.FindBetweenErr
}

type MockParticipationCalculator struct {
	CalculateCalls int
	CalculateErr   error

	LastParticipant participatables.Participatable
	LastTarget      participatables.Participatable

	Value Participation
}

func (calculator *MockParticipationCalculator) Calculate(
	participant participatables.Participatable,
	target participatables.Participatable,
) (Participation, error) {
	calculator.CalculateCalls++
	calculator.LastParticipant = participant
	calculator.LastTarget = target

	return calculator.Value, calculator.CalculateErr
}
