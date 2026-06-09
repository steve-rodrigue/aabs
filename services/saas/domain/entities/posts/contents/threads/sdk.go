package threads

import (
	"errors"

	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/users"
)

var (
	ErrInvalidThreadIdentifier = errors.New("invalid thread identifier")
	ErrInvalidThreadCreator    = errors.New("invalid thread creator")
	ErrInvalidThreadTitle      = errors.New("invalid thread title")
	ErrInvalidThreadText       = errors.New("invalid thread text")
)

// NewAdapter creates a new thread adapter
func NewAdapter() Adapter {
	return &adapter{}
}

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
