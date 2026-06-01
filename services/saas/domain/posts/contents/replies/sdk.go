package replies

import (
	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/posts/contents/threads"
)

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
