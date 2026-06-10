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
