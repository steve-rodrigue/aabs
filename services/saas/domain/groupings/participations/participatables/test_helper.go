package participatables

import "github.com/google/uuid"

func NewMockParticipatable(
	id uuid.UUID,
	kind Kind,
) Participatable {
	return &MockParticipatable{
		id:   id,
		kind: kind,
	}
}

func NewMockParticipatableRepository() *MockParticipatableRepository {
	return &MockParticipatableRepository{}
}

type MockParticipatable struct {
	id   uuid.UUID
	kind Kind
}

func (participatable *MockParticipatable) Identifier() uuid.UUID {
	return participatable.id
}

func (participatable *MockParticipatable) ParticipationKind() Kind {
	return participatable.kind
}

type MockParticipatableRepository struct {
	FindAllParticipantsCalls int
	FindAllParticipantsErr   error
	FindAllParticipantsValue []Participatable

	FindAllTargetsCalls int
	FindAllTargetsErr   error
	FindAllTargetsValue []Participatable
}

func (repository *MockParticipatableRepository) FindAllParticipants() ([]Participatable, error) {
	repository.FindAllParticipantsCalls++

	return repository.FindAllParticipantsValue, repository.FindAllParticipantsErr
}

func (repository *MockParticipatableRepository) FindAllTargets() ([]Participatable, error) {
	repository.FindAllTargetsCalls++

	return repository.FindAllTargetsValue, repository.FindAllTargetsErr
}
