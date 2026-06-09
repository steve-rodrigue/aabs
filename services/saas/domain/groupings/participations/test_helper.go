package participations

import (
	"context"
	"sort"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
)

func NewMockParticipation() Participation {
	return &MockParticipation{
		id:         uuid.New(),
		detectedOn: time.Now().UTC(),
	}
}

func NewMockParticipationWithID(
	id uuid.UUID,
) Participation {
	return &MockParticipation{
		id:         id,
		detectedOn: time.Now().UTC(),
	}
}

func NewMockParticipationBetween(
	participant participatables.Participatable,
	target participatables.Participatable,
) Participation {
	return &MockParticipation{
		id:          uuid.New(),
		participant: participant,
		target:      target,
		detectedOn:  time.Now().UTC(),
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

func NewMockParticipationCalculator() *MockParticipationCalculator {
	return &MockParticipationCalculator{}
}

func NewMockParticipationAdapter() *MockParticipationAdapter {
	return &MockParticipationAdapter{}
}

type MockParticipation struct {
	id uuid.UUID

	participant participatables.Participatable
	target      participatables.Participatable

	postCount      int
	totalPostCount int
	percentage     float64

	detectedOn time.Time
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
	return participation.detectedOn
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
		id:             input.Identifier,
		participant:    input.Participant,
		target:         input.Target,
		postCount:      input.PostCount,
		totalPostCount: input.TotalPostCount,
		percentage:     input.Percentage,
		detectedOn:     input.DetectedOn,
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
