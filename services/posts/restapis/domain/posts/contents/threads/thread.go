package threads

import (
	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/users"
)

type thread struct {
	identifier uuid.UUID
	creator    users.User
	title      string
	text       string
}

func (thread *thread) Identifier() uuid.UUID {
	return thread.identifier
}

func (thread *thread) Creator() users.User {
	return thread.creator
}

func (thread *thread) Title() string {
	return thread.title
}

func (thread *thread) Text() string {
	return thread.text
}
