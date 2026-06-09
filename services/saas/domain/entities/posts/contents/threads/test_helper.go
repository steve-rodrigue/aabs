package threads

import (
	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/users"
)

func NewMockThread(
	title string,
	text string,
) Thread {
	return &MockThread{
		id:      uuid.New(),
		title:   title,
		text:    text,
		creator: users.NewMockUser("@user", "User"),
	}
}

func NewMockThreadAdapter() *MockThreadAdapter {
	return &MockThreadAdapter{}
}

type MockThread struct {
	id      uuid.UUID
	creator users.User
	title   string
	text    string
}

func (thread *MockThread) Identifier() uuid.UUID {
	return thread.id
}

func (thread *MockThread) Creator() users.User {
	return thread.creator
}

func (thread *MockThread) Title() string {
	return thread.title
}

func (thread *MockThread) Text() string {
	return thread.text
}

type MockThreadAdapter struct {
	ToDomainCalls int
	ToDomainErr   error
	ToDomainValue Thread

	LastInput ThreadInput
}

func (adapter *MockThreadAdapter) ToDomain(
	input ThreadInput,
) (Thread, error) {
	adapter.ToDomainCalls++
	adapter.LastInput = input

	if adapter.ToDomainErr != nil {
		return nil, adapter.ToDomainErr
	}

	if adapter.ToDomainValue != nil {
		return adapter.ToDomainValue, nil
	}

	return &MockThread{
		id:      input.Identifier,
		creator: input.Creator,
		title:   input.Title,
		text:    input.Text,
	}, nil
}
