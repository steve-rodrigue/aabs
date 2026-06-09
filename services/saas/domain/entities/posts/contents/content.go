package contents

import (
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts/contents/replies"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts/contents/threads"
)

type content struct {
	identifier uuid.UUID
	reply      replies.Reply
	thread     threads.Thread
	createdAt  time.Time
}

func (content *content) Identifier() uuid.UUID {
	return content.identifier
}

func (content *content) IsReply() bool {
	return content.reply != nil
}

func (content *content) Reply() replies.Reply {
	return content.reply
}

func (content *content) IsThread() bool {
	return content.thread != nil
}

func (content *content) Thread() threads.Thread {
	return content.thread
}

func (content *content) Text() string {
	if content.IsReply() {
		return content.reply.Text()
	}

	if content.IsThread() {
		return content.thread.Text()
	}

	return ""
}

func (content *content) CreatedAt() time.Time {
	return content.createdAt
}
