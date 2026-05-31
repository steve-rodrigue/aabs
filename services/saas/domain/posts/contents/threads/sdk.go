package threads

import (
	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/users"
)

// Thread represents a thread
type Thread interface {
	Identifier() uuid.UUID
	Creator() users.User
	Title() string
	Text() string
}
