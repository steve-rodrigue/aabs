package contents

import (
	"time"

	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/posts/contents/replies"
	"github.com/steve-rodrigue/aabs/services/saas/domain/posts/contents/threads"
)

type MockContent struct {
	TextValue string
}

func (content *MockContent) Identifier() uuid.UUID {
	return uuid.New()
}

func (content *MockContent) IsReply() bool {
	return false
}

func (content *MockContent) Reply() replies.Reply {
	return nil
}

func (content *MockContent) IsThread() bool {
	return false
}

func (content *MockContent) Thread() threads.Thread {
	return nil
}

func (content *MockContent) Text() string {
	return content.TextValue
}

func (content *MockContent) CreatedAt() time.Time {
	return time.Time{}
}
