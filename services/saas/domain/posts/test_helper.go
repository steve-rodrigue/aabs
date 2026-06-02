package posts

import (
	"time"

	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/communities"
	"github.com/steve-rodrigue/aabs/services/saas/domain/platforms"
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

func NewMockPostWithUser(
	text string,
	creator users.User,
) Post {
	return &MockPost{
		id:      uuid.New(),
		creator: creator,
		content: &contents.MockContent{
			TextValue: text,
		},
	}
}

func NewMockPostWithCommunities(
	text string,
	communityIDs []uuid.UUID,
) Post {
	return &MockPost{
		id:           uuid.New(),
		communityIDs: communityIDs,
		content: &contents.MockContent{
			TextValue: text,
		},
	}
}

func NewMockPostRepository() *MockPostRepository {
	return &MockPostRepository{
		Items: map[uuid.UUID]Post{},
	}
}

type MockPost struct {
	id uuid.UUID

	communityIDs []uuid.UUID
	creator      users.User
	content      contents.Content
}

func (post *MockPost) Identifier() uuid.UUID {
	return post.id
}

func (post *MockPost) CommunityIDs() []uuid.UUID {
	return post.communityIDs
}

func (post *MockPost) Creator() users.User {
	return post.creator
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

	Items map[uuid.UUID]Post

	FindByIDCalls int
	FindByIDErr   error

	FindAllCalls int
	FindAllErr   error
	FindAllValue []Post

	FindByUserCalls int
	FindByUserErr   error
	FindByUserValue []Post

	FindByCommunityCalls int
	FindByCommunityErr   error
	FindByCommunityValue []Post

	FindByPlatformCalls int
	FindByPlatformErr   error
	FindByPlatformValue []Post
}

func (repository *MockPostRepository) Save(post Post) error {
	repository.SaveCalls++

	return repository.SaveErr
}

func (repository *MockPostRepository) FindByID(id uuid.UUID) (Post, error) {
	repository.FindByIDCalls++

	if repository.FindByIDErr != nil {
		return nil, repository.FindByIDErr
	}

	if repository.Items == nil {
		return nil, nil
	}

	return repository.Items[id], nil
}

func (repository *MockPostRepository) FindAll() ([]Post, error) {
	repository.FindAllCalls++

	if repository.FindAllErr != nil {
		return nil, repository.FindAllErr
	}

	if repository.FindAllValue != nil {
		return repository.FindAllValue, nil
	}

	out := make([]Post, 0, len(repository.Items))

	for _, post := range repository.Items {
		out = append(out, post)
	}

	return out, nil
}

func (repository *MockPostRepository) FindByUser(
	user users.User,
) ([]Post, error) {
	repository.FindByUserCalls++

	if repository.FindByUserErr != nil {
		return nil, repository.FindByUserErr
	}

	if repository.FindByUserValue != nil {
		return repository.FindByUserValue, nil
	}

	out := []Post{}

	for _, post := range repository.Items {
		if post.Creator() == nil {
			continue
		}

		if post.Creator().Identifier() == user.Identifier() {
			out = append(out, post)
		}
	}

	return out, nil
}

func (repository *MockPostRepository) FindByCommunity(
	community communities.Community,
) ([]Post, error) {
	repository.FindByCommunityCalls++

	if repository.FindByCommunityErr != nil {
		return nil, repository.FindByCommunityErr
	}

	if repository.FindByCommunityValue != nil {
		return repository.FindByCommunityValue, nil
	}

	out := []Post{}

	for _, post := range repository.Items {
		for _, communityID := range post.CommunityIDs() {
			if communityID == community.Identifier() {
				out = append(out, post)
				break
			}
		}
	}

	return out, nil
}

func (repository *MockPostRepository) FindByPlatform(
	platform platforms.Platform,
) ([]Post, error) {
	repository.FindByPlatformCalls++

	if repository.FindByPlatformErr != nil {
		return nil, repository.FindByPlatformErr
	}

	if repository.FindByPlatformValue != nil {
		return repository.FindByPlatformValue, nil
	}

	out := []Post{}

	for _, post := range repository.Items {
		if post.Creator() == nil || post.Creator().Platform() == nil {
			continue
		}

		if post.Creator().Platform().Identifier() == platform.Identifier() {
			out = append(out, post)
		}
	}

	return out, nil
}
