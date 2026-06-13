package contents

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts/contents/replies"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts/contents/threads"
)

var (
	ErrInvalidContentIdentifier = errors.New("invalid content identifier")
	ErrInvalidContentTarget     = errors.New("invalid content target")
	ErrInvalidContentCreatedAt  = errors.New("invalid content created at")
)

// NewAdapter creates a new content adapter
func NewAdapter(
	replies replies.Adapter,
	threads threads.Adapter,
) Adapter {
	return &adapter{
		replies: replies,
		threads: threads,
	}
}

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
