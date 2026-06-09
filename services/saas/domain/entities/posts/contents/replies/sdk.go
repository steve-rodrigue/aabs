package replies

import (
	"errors"

	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts/contents/threads"
)

var (
	ErrInvalidReplyIdentifier = errors.New("invalid reply identifier")
	ErrInvalidReplyTarget     = errors.New("invalid reply target")
	ErrInvalidReplyText       = errors.New("invalid reply text")
)

// NewAdapter creates a new reply adapter
func NewAdapter() Adapter {
	return &adapter{}
}

// ReplyInput represents the reply input
type ReplyInput struct {
	Identifier uuid.UUID
	Target     TargetInput
	Text       string
}

// TargetInput represents the target input
type TargetInput struct {
	Reply  Reply
	Thread threads.Thread
}

// Adapter represents a reply adapter
type Adapter interface {
	ToDomain(input ReplyInput) (Reply, error)
}

// Reply represents a reply
type Reply interface {
	Identifier() uuid.UUID
	Target() Target
	Text() string
}

// Target represents a reply target
type Target interface {
	IsReply() bool
	Reply() Reply
	IsThread() bool
	Thread() threads.Thread
}
