package replies

import (
	"strings"

	"github.com/google/uuid"
)

type adapter struct{}

func (adapter *adapter) ToDomain(
	input ReplyInput,
) (Reply, error) {
	if input.Identifier == uuid.Nil {
		return nil, ErrInvalidReplyIdentifier
	}

	target, err := adapter.toTarget(input.Target)
	if err != nil {
		return nil, err
	}

	input.Text = strings.TrimSpace(input.Text)
	if input.Text == "" {
		return nil, ErrInvalidReplyText
	}

	return &reply{
		identifier: input.Identifier,
		target:     target,
		text:       input.Text,
	}, nil
}

func (adapter *adapter) toTarget(
	input TargetInput,
) (Target, error) {
	hasReply := input.Reply != nil
	hasThread := input.Thread != nil

	if hasReply == hasThread {
		return nil, ErrInvalidReplyTarget
	}

	return &target{
		reply:  input.Reply,
		thread: input.Thread,
	}, nil
}
