package posts

import (
	"time"

	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/posts/contents"
	"github.com/steve-rodrigue/aabs/services/saas/domain/users"
)

func NewMockPost(text string) Post {
	return &MockPost{
		id: uuid.New(),
		content: &contents.MockContent{
			TextValue: text,
		},
	}
}

type MockPost struct {
	id      uuid.UUID
	content contents.Content
}

func (post *MockPost) Identifier() uuid.UUID {
	return post.id
}

func (post *MockPost) CommunityIDs() []uuid.UUID {
	return nil
}

func (post *MockPost) Creator() users.User {
	return nil
}

func (post *MockPost) Content() contents.Content {
	return post.content
}

func (post *MockPost) CreatedOn() time.Time {
	return time.Time{}
}

type MockPostRepository struct {
	SaveCalls int
	SaveErr   error
}

func (repository *MockPostRepository) Save(post Post) error {
	repository.SaveCalls++

	return repository.SaveErr
}

func (repository *MockPostRepository) FindByID(id uuid.UUID) (Post, error) {
	return nil, nil
}
