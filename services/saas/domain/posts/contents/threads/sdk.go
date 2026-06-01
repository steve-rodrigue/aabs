package threads

import (
	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/users"
)

// ThreadInput represents the thread input
type ThreadInput struct {
	Identifier uuid.UUID
	Creator    users.User
	Title      string
	Text       string
}

// Adapter represents a thread adapter
type Adapter interface {
	ToDomain(input ThreadInput) (Thread, error)
}

// Thread represents a thread
type Thread interface {
	Identifier() uuid.UUID
	Creator() users.User
	Title() string
	Text() string
}
