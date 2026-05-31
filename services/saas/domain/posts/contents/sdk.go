package contents

import (
	"time"

	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/posts/contents/replies"
	"github.com/steve-rodrigue/aabs/services/saas/domain/posts/contents/threads"
)

// Content represents a post's content
type Content interface {
	Identifier() uuid.UUID
	IsReply() bool
	Reply() replies.Reply
	IsThread() bool
	Thread() threads.Thread
	CreatedAt() time.Time
}
