package contents

import (
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts/contents/replies"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts/contents/threads"
)

func NewMockContent(
	text string,
) Content {
	return &MockContent{
		id:        uuid.New(),
		TextValue: text,
		createdAt: time.Now().UTC(),
	}
}

func NewMockContentAdapter() *MockContentAdapter {
	return &MockContentAdapter{}
}

type MockContent struct {
	id uuid.UUID

	reply  replies.Reply
	thread threads.Thread

	TextValue string

	createdAt time.Time
}

func (content *MockContent) Identifier() uuid.UUID {
	return content.id
}

func (content *MockContent) IsReply() bool {
	return content.reply != nil
}

func (content *MockContent) Reply() replies.Reply {
	return content.reply
}

func (content *MockContent) IsThread() bool {
	return content.thread != nil
}

func (content *MockContent) Thread() threads.Thread {
	return content.thread
}

func (content *MockContent) Text() string {
	if content.TextValue != "" {
		return content.TextValue
	}

	if content.reply != nil {
		return content.reply.Text()
	}

	if content.thread != nil {
		return content.thread.Text()
	}

	return ""
}

func (content *MockContent) CreatedAt() time.Time {
	return content.createdAt
}

type MockContentAdapter struct {
	ToDomainCalls int
	ToDomainErr   error
	ToDomainValue Content
	ToDomainNil   bool

	LastInput ContentInput
}

func (adapter *MockContentAdapter) ToDomain(
	input ContentInput,
) (Content, error) {
	adapter.ToDomainCalls++
	adapter.LastInput = input

	if adapter.ToDomainErr != nil {
		return nil, adapter.ToDomainErr
	}

	if adapter.ToDomainNil {
		return nil, nil
	}

	if adapter.ToDomainValue != nil {
		return adapter.ToDomainValue, nil
	}

	content := &MockContent{
		id:        input.Identifier,
		createdAt: input.CreatedAt,
	}

	if input.Reply != nil {
		content.reply = replies.NewMockReply(input.Reply.Text)
		content.TextValue = input.Reply.Text
	}

	if input.Thread != nil {
		content.thread = threads.NewMockThread(input.Thread.Title, input.Thread.Text)
		content.TextValue = input.Thread.Text
	}

	return content, nil
}
