package participations

import (
	"context"
	"sort"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/participatables"
)

func NewMockParticipation() Participation {
	return &MockParticipation{
		ID:              uuid.New(),
		DetectedOnValue: time.Now().UTC(),
	}
}

func NewMockParticipationWithID(
	id uuid.UUID,
) Participation {
	return &MockParticipation{
		ID:              id,
		DetectedOnValue: time.Now().UTC(),
	}
}

func NewMockParticipationBetween(
	participant participatables.Participatable,
	target participatables.Participatable,
) Participation {
	return &MockParticipation{
		ID:               uuid.New(),
		ParticipantValue: participant,
		TargetValue:      target,
		DetectedOnValue:  time.Now().UTC(),
	}
}

func NewMockParticipationWithParticipantAndTarget(
	participant participatables.Participatable,
	target participatables.Participatable,
) Participation {
	return NewMockParticipationBetween(
		participant,
		target,
	)
}

func NewMockParticipationRepository() *MockParticipationRepository {
	return &MockParticipationRepository{
		Items: map[uuid.UUID]Participation{},
	}
}

func NewMockParticipationCounter() *MockParticipationCounter {
	return &MockParticipationCounter{}
}

func NewMockParticipationCalculator() *MockParticipationCalculator {
	return &MockParticipationCalculator{}
}

func NewMockParticipationAdapter() *MockParticipationAdapter {
	return &MockParticipationAdapter{}
}

type MockParticipation struct {
	ID uuid.UUID

	ParticipantValue participatables.Participatable
	TargetValue      participatables.Participatable

	PostCountValue      int
	TotalPostCountValue int
	PercentageValue     float64

	DetectedOnValue time.Time
}

func (participation *MockParticipation) Identifier() uuid.UUID {
	return participation.ID
}

func (participation *MockParticipation) Participant() participatables.Participatable {
	return participation.ParticipantValue
}

func (participation *MockParticipation) Target() participatables.Participatable {
	return participation.TargetValue
}

func (participation *MockParticipation) PostCount() int {
	return participation.PostCountValue
}

func (participation *MockParticipation) TotalPostCount() int {
	return participation.TotalPostCountValue
}

func (participation *MockParticipation) Percentage() float64 {
	return participation.PercentageValue
}

func (participation *MockParticipation) DetectedOn() time.Time {
	return participation.DetectedOnValue
}

type MockParticipationAdapter struct {
	ToDomainCalls int
	ToDomainErr   error
	ToDomainValue Participation

	LastInput ParticipationInput
}

func (adapter *MockParticipationAdapter) ToDomain(
	input ParticipationInput,
) (Participation, error) {
	adapter.ToDomainCalls++
	adapter.LastInput = input

	if adapter.ToDomainErr != nil {
		return nil, adapter.ToDomainErr
	}

	if adapter.ToDomainValue != nil {
		return adapter.ToDomainValue, nil
	}

	return &MockParticipation{
		ID:                  input.Identifier,
		ParticipantValue:    input.Participant,
		TargetValue:         input.Target,
		PostCountValue:      input.PostCount,
		TotalPostCountValue: input.TotalPostCount,
		PercentageValue:     input.Percentage,
		DetectedOnValue:     input.DetectedOn,
	}, nil
}

type MockParticipationRepository struct {
	SaveCalls int
	SaveErr   error

	Items map[uuid.UUID]Participation

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

	LastContext     context.Context
	LastSaved       Participation
	LastID          uuid.UUID
	LastParticipant participatables.Participatable
	LastTarget      participatables.Participatable
}

func (repository *MockParticipationRepository) Save(
	ctx context.Context,
	participation Participation,
) error {
	repository.SaveCalls++
	repository.LastContext = ctx
	repository.LastSaved = participation

	if repository.Items != nil && participation != nil {
		repository.Items[participation.Identifier()] = participation
	}

	return repository.SaveErr
}

func (repository *MockParticipationRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (Participation, error) {
	repository.FindByIDCalls++
	repository.LastContext = ctx
	repository.LastID = id

	if repository.FindByIDErr != nil {
		return nil, repository.FindByIDErr
	}

	if repository.FindByIDValue != nil {
		return repository.FindByIDValue, nil
	}

	if repository.Items == nil {
		return nil, nil
	}

	return repository.Items[id], nil
}

func (repository *MockParticipationRepository) FindByParticipant(
	ctx context.Context,
	participant participatables.Participatable,
) ([]Participation, error) {
	repository.FindByParticipantCalls++
	repository.LastContext = ctx
	repository.LastParticipant = participant

	if repository.FindByParticipantErr != nil {
		return nil, repository.FindByParticipantErr
	}

	if repository.FindByParticipantValue != nil {
		return repository.FindByParticipantValue, nil
	}

	out := []Participation{}

	for _, participation := range repository.Items {
		if participation.Participant() == nil ||
			participant == nil {
			continue
		}

		if participation.Participant().Identifier() == participant.Identifier() &&
			participation.Participant().ParticipationKind() == participant.ParticipationKind() {
			out = append(out, participation)
		}
	}

	return out, nil
}

func (repository *MockParticipationRepository) FindByTarget(
	ctx context.Context,
	target participatables.Participatable,
) ([]Participation, error) {
	repository.FindByTargetCalls++
	repository.LastContext = ctx
	repository.LastTarget = target

	if repository.FindByTargetErr != nil {
		return nil, repository.FindByTargetErr
	}

	if repository.FindByTargetValue != nil {
		return repository.FindByTargetValue, nil
	}

	out := []Participation{}

	for _, participation := range repository.Items {
		if participation.Target() == nil ||
			target == nil {
			continue
		}

		if participation.Target().Identifier() == target.Identifier() &&
			participation.Target().ParticipationKind() == target.ParticipationKind() {
			out = append(out, participation)
		}
	}

	return out, nil
}

func (repository *MockParticipationRepository) FindBetween(
	ctx context.Context,
	participant participatables.Participatable,
	target participatables.Participatable,
) (Participation, error) {
	repository.FindBetweenCalls++
	repository.LastContext = ctx
	repository.LastParticipant = participant
	repository.LastTarget = target

	if repository.FindBetweenErr != nil {
		return nil, repository.FindBetweenErr
	}

	if repository.FindBetweenValue != nil {
		return repository.FindBetweenValue, nil
	}

	for _, participation := range repository.Items {
		if participation.Participant() == nil ||
			participation.Target() == nil ||
			participant == nil ||
			target == nil {
			continue
		}

		if participation.Participant().Identifier() == participant.Identifier() &&
			participation.Participant().ParticipationKind() == participant.ParticipationKind() &&
			participation.Target().Identifier() == target.Identifier() &&
			participation.Target().ParticipationKind() == target.ParticipationKind() {
			return participation, nil
		}
	}

	return nil, nil
}

func (repository *MockParticipationRepository) sortedParticipations() []Participation {
	out := make([]Participation, 0, len(repository.Items))

	for _, participation := range repository.Items {
		out = append(out, participation)
	}

	sort.Slice(out, func(left int, right int) bool {
		return out[left].Identifier().String() <
			out[right].Identifier().String()
	})

	return out
}

type MockParticipationCounter struct {
	CountByParticipantAndTargetCalls int
	CountByParticipantAndTargetErr   error
	CountByParticipantAndTargetValue int
	CountByTargetCalls               int
	CountByTargetErr                 error
	CountByTargetValue               int
	LastContext                      context.Context
	LastParticipant                  participatables.Participatable
	LastTarget                       participatables.Participatable
}

func (counter *MockParticipationCounter) CountByParticipantAndTarget(
	ctx context.Context,
	participant participatables.Participatable,
	target participatables.Participatable,
) (int, error) {
	counter.CountByParticipantAndTargetCalls++
	counter.LastContext = ctx
	counter.LastParticipant = participant
	counter.LastTarget = target
	return counter.CountByParticipantAndTargetValue,
		counter.CountByParticipantAndTargetErr

}

func (counter *MockParticipationCounter) CountByTarget(
	ctx context.Context,
	target participatables.Participatable,

) (int, error) {
	counter.CountByTargetCalls++
	counter.LastContext = ctx
	counter.LastTarget = target
	return counter.CountByTargetValue,
		counter.CountByTargetErr

}

type MockParticipationCalculator struct {
	CalculateCalls int
	CalculateErr   error
	CalculateValue Participation

	LastContext     context.Context
	LastParticipant participatables.Participatable
	LastTarget      participatables.Participatable
}

func (calculator *MockParticipationCalculator) Calculate(
	ctx context.Context,
	participant participatables.Participatable,
	target participatables.Participatable,
) (Participation, error) {
	calculator.CalculateCalls++
	calculator.LastContext = ctx
	calculator.LastParticipant = participant
	calculator.LastTarget = target

	return calculator.CalculateValue,
		calculator.CalculateErr
}
