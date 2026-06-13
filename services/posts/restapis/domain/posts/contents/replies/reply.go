package replies

import "github.com/google/uuid"

type reply struct {
	identifier uuid.UUID
	target     Target
	text       string
}

func (reply *reply) Identifier() uuid.UUID {
	return reply.identifier
}

func (reply *reply) Target() Target {
	return reply.target
}

func (reply *reply) Text() string {
	return reply.text
}
