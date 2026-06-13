package threads

import (
	"strings"

	"github.com/google/uuid"
)

type adapter struct{}

func (adapter *adapter) ToDomain(
	input ThreadInput,
) (Thread, error) {
	if input.Identifier == uuid.Nil {
		return nil, ErrInvalidThreadIdentifier
	}

	if input.Creator == nil {
		return nil, ErrInvalidThreadCreator
	}

	input.Title = strings.TrimSpace(input.Title)
	if input.Title == "" {
		return nil, ErrInvalidThreadTitle
	}

	input.Text = strings.TrimSpace(input.Text)
	if input.Text == "" {
		return nil, ErrInvalidThreadText
	}

	return &thread{
		identifier: input.Identifier,
		creator:    input.Creator,
		title:      input.Title,
		text:       input.Text,
	}, nil
}
