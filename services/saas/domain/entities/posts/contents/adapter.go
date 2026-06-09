package contents

import (
	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts/contents/replies"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts/contents/threads"
)

type adapter struct {
	replies replies.Adapter
	threads threads.Adapter
}

func (adapter *adapter) ToDomain(
	input ContentInput,
) (Content, error) {
	if input.Identifier == uuid.Nil {
		return nil, ErrInvalidContentIdentifier
	}

	if input.CreatedAt.IsZero() {
		return nil, ErrInvalidContentCreatedAt
	}

	hasReply := input.Reply != nil
	hasThread := input.Thread != nil

	if hasReply == hasThread {
		return nil, ErrInvalidContentTarget
	}

	var reply replies.Reply
	var thread threads.Thread
	var err error

	if hasReply {
		reply, err = adapter.replies.ToDomain(*input.Reply)
		if err != nil {
			return nil, err
		}
	}

	if hasThread {
		thread, err = adapter.threads.ToDomain(*input.Thread)
		if err != nil {
			return nil, err
		}
	}

	return &content{
		identifier: input.Identifier,
		reply:      reply,
		thread:     thread,
		createdAt:  input.CreatedAt.UTC(),
	}, nil
}
