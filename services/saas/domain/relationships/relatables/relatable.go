package relatables

import "github.com/google/uuid"

type relatable struct {
	identifier       uuid.UUID
	relationshipKind Kind
}

func (relatable *relatable) Identifier() uuid.UUID {
	return relatable.identifier
}

func (relatable *relatable) RelationshipKind() Kind {
	return relatable.relationshipKind
}
