package replies

import (
	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/posts/contents/threads"
)

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
