package replies

import (
	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts/contents/threads"
)

func NewMockReply(
	text string,
) Reply {
	return &MockReply{
		id:   uuid.New(),
		text: text,
	}
}

func NewMockReplyWithReplyTarget(
	text string,
	target Reply,
) Reply {
	return &MockReply{
		id: uuid.New(),
		target: &MockTarget{
			reply: target,
		},
		text: text,
	}
}

func NewMockReplyWithThreadTarget(
	text string,
	target threads.Thread,
) Reply {
	return &MockReply{
		id: uuid.New(),
		target: &MockTarget{
			thread: target,
		},
		text: text,
	}
}

type MockReply struct {
	id     uuid.UUID
	target Target
	text   string
}

func (reply *MockReply) Identifier() uuid.UUID {
	return reply.id
}

func (reply *MockReply) Target() Target {
	return reply.target
}

func (reply *MockReply) Text() string {
	return reply.text
}

func NewMockTargetWithReply(
	reply Reply,
) Target {
	return &MockTarget{
		reply: reply,
	}
}

func NewMockTargetWithThread(
	thread threads.Thread,
) Target {
	return &MockTarget{
		thread: thread,
	}
}

type MockTarget struct {
	reply  Reply
	thread threads.Thread
}

func (target *MockTarget) IsReply() bool {
	return target.reply != nil
}

func (target *MockTarget) Reply() Reply {
	return target.reply
}

func (target *MockTarget) IsThread() bool {
	return target.thread != nil
}

func (target *MockTarget) Thread() threads.Thread {
	return target.thread
}

func NewMockReplyAdapter() *MockReplyAdapter {
	return &MockReplyAdapter{}
}

type MockReplyAdapter struct {
	ToDomainCalls int
	ToDomainErr   error
	ToDomainValue Reply

	LastInput ReplyInput
}

func (adapter *MockReplyAdapter) ToDomain(
	input ReplyInput,
) (Reply, error) {
	adapter.ToDomainCalls++
	adapter.LastInput = input

	if adapter.ToDomainErr != nil {
		return nil, adapter.ToDomainErr
	}

	if adapter.ToDomainValue != nil {
		return adapter.ToDomainValue, nil
	}

	target := &MockTarget{
		reply:  input.Target.Reply,
		thread: input.Target.Thread,
	}

	return &MockReply{
		id:     input.Identifier,
		target: target,
		text:   input.Text,
	}, nil
}
