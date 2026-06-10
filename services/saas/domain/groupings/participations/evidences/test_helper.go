package evidences

import (
	"context"
	"sort"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/participatables"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations"
)

func NewMockEvidence() Evidence {
	return &MockEvidence{
		ID:              uuid.New(),
		DetectedOnValue: time.Now().UTC(),
	}
}

func NewMockEvidenceWithID(
	id uuid.UUID,
) Evidence {
	return &MockEvidence{
		ID:              id,
		DetectedOnValue: time.Now().UTC(),
	}
}

func NewMockEvidenceRepository() *MockEvidenceRepository {
	return &MockEvidenceRepository{
		Items: map[uuid.UUID]Evidence{},
	}
}

func NewMockEvidenceCalculator() *MockEvidenceCalculator {
	return &MockEvidenceCalculator{}
}

func NewMockEvidenceAdapter() *MockEvidenceAdapter {
	return &MockEvidenceAdapter{}
}

type MockEvidence struct {
	ID uuid.UUID

	ParticipationValue participations.Participation

	ParticipantValue participatables.Participatable
	TargetValue      participatables.Participatable

	PostValue  posts.Post
	ScoreValue float64

	DetectedOnValue time.Time
}

func (evidence *MockEvidence) Identifier() uuid.UUID {
	return evidence.ID
}

func (evidence *MockEvidence) Participation() participations.Participation {
	return evidence.ParticipationValue
}

func (evidence *MockEvidence) Participant() participatables.Participatable {
	return evidence.ParticipantValue
}

func (evidence *MockEvidence) Target() participatables.Participatable {
	return evidence.TargetValue
}

func (evidence *MockEvidence) Post() posts.Post {
	return evidence.PostValue
}

func (evidence *MockEvidence) Score() float64 {
	return evidence.ScoreValue
}

func (evidence *MockEvidence) DetectedOn() time.Time {
	return evidence.DetectedOnValue
}

type MockEvidenceAdapter struct {
	ToDomainCalls int
	ToDomainErr   error
	ToDomainValue Evidence

	LastInput EvidenceInput
}

func (adapter *MockEvidenceAdapter) ToDomain(
	input EvidenceInput,
) (Evidence, error) {
	adapter.ToDomainCalls++
	adapter.LastInput = input

	if adapter.ToDomainErr != nil {
		return nil, adapter.ToDomainErr
	}

	if adapter.ToDomainValue != nil {
		return adapter.ToDomainValue, nil
	}

	return &MockEvidence{
		ID:                 input.Identifier,
		ParticipationValue: input.Participation,
		ParticipantValue:   input.Participant,
		TargetValue:        input.Target,
		PostValue:          input.Post,
		ScoreValue:         input.Score,
		DetectedOnValue:    input.DetectedOn,
	}, nil
}

type MockEvidenceRepository struct {
	SaveCalls int
	SaveErr   error

	Items map[uuid.UUID]Evidence

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

	LastContext       context.Context
	LastEvidence      Evidence
	LastID            uuid.UUID
	LastParticipation uuid.UUID
	LastPost          uuid.UUID
	LastParticipant   participatables.Participatable
	LastTarget        participatables.Participatable
}

func (repository *MockEvidenceRepository) Save(
	ctx context.Context,
	evidence Evidence,
) error {
	repository.SaveCalls++
	repository.LastContext = ctx
	repository.LastEvidence = evidence

	if repository.Items != nil && evidence != nil {
		repository.Items[evidence.Identifier()] = evidence
	}

	return repository.SaveErr
}

func (repository *MockEvidenceRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (Evidence, error) {
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

func (repository *MockEvidenceRepository) FindByParticipation(
	ctx context.Context,
	participation uuid.UUID,
) ([]Evidence, error) {
	repository.FindByParticipationCalls++
	repository.LastContext = ctx
	repository.LastParticipation = participation

	if repository.FindByParticipationErr != nil {
		return nil, repository.FindByParticipationErr
	}

	if repository.FindByParticipationValue != nil {
		return repository.FindByParticipationValue, nil
	}

	out := []Evidence{}

	for _, evidence := range repository.Items {
		if evidence.Participation() == nil {
			continue
		}

		if evidence.Participation().Identifier() == participation {
			out = append(out, evidence)
		}
	}

	return out, nil
}

func (repository *MockEvidenceRepository) FindByPost(
	ctx context.Context,
	post uuid.UUID,
) ([]Evidence, error) {
	repository.FindByPostCalls++
	repository.LastContext = ctx
	repository.LastPost = post

	if repository.FindByPostErr != nil {
		return nil, repository.FindByPostErr
	}

	if repository.FindByPostValue != nil {
		return repository.FindByPostValue, nil
	}

	out := []Evidence{}

	for _, evidence := range repository.Items {
		if evidence.Post() == nil {
			continue
		}

		if evidence.Post().Identifier() == post {
			out = append(out, evidence)
		}
	}

	return out, nil
}

func (repository *MockEvidenceRepository) FindByParticipant(
	ctx context.Context,
	participant participatables.Participatable,
) ([]Evidence, error) {
	repository.FindByParticipantCalls++
	repository.LastContext = ctx
	repository.LastParticipant = participant

	if repository.FindByParticipantErr != nil {
		return nil, repository.FindByParticipantErr
	}

	if repository.FindByParticipantValue != nil {
		return repository.FindByParticipantValue, nil
	}

	out := []Evidence{}

	for _, evidence := range repository.Items {
		if evidence.Participant() == nil ||
			participant == nil {
			continue
		}

		if evidence.Participant().Identifier() == participant.Identifier() &&
			evidence.Participant().ParticipationKind() == participant.ParticipationKind() {
			out = append(out, evidence)
		}
	}

	return out, nil
}

func (repository *MockEvidenceRepository) FindByTarget(
	ctx context.Context,
	target participatables.Participatable,
) ([]Evidence, error) {
	repository.FindByTargetCalls++
	repository.LastContext = ctx
	repository.LastTarget = target

	if repository.FindByTargetErr != nil {
		return nil, repository.FindByTargetErr
	}

	if repository.FindByTargetValue != nil {
		return repository.FindByTargetValue, nil
	}

	out := []Evidence{}

	for _, evidence := range repository.Items {
		if evidence.Target() == nil ||
			target == nil {
			continue
		}

		if evidence.Target().Identifier() == target.Identifier() &&
			evidence.Target().ParticipationKind() == target.ParticipationKind() {
			out = append(out, evidence)
		}
	}

	return out, nil
}

func (repository *MockEvidenceRepository) sortedEvidence() []Evidence {
	out := make([]Evidence, 0, len(repository.Items))

	for _, evidence := range repository.Items {
		out = append(out, evidence)
	}

	sort.Slice(out, func(left int, right int) bool {
		return out[left].Identifier().String() <
			out[right].Identifier().String()
	})

	return out
}

type MockEvidenceCalculator struct {
	CalculateCalls int
	CalculateErr   error
	CalculateValue []Evidence

	LastContext       context.Context
	LastParticipation participations.Participation
}

func (calculator *MockEvidenceCalculator) Calculate(
	ctx context.Context,
	participation participations.Participation,
) ([]Evidence, error) {
	calculator.CalculateCalls++
	calculator.LastContext = ctx
	calculator.LastParticipation = participation

	return calculator.CalculateValue,
		calculator.CalculateErr
}
