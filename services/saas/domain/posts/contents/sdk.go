package contents

import (
	"time"

	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/posts/contents/replies"
	"github.com/steve-rodrigue/aabs/services/saas/domain/posts/contents/threads"
)

// ContentInput represents a content input
type ContentInput struct {
	Identifier uuid.UUID
	Reply      *replies.ReplyInput
	Thread     *threads.ThreadInput
	CreatedAt  time.Time
}

// Adapter represents a content adapter
type Adapter interface {
	ToDomain(input ContentInput) (Content, error)
}

// Content represents a post's content
type Content interface {
	Identifier() uuid.UUID
	IsReply() bool
	Reply() replies.Reply
	IsThread() bool
	Thread() threads.Thread
	Text() string
	CreatedAt() time.Time
}
