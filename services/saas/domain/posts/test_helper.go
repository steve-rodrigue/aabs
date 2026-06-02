package posts

import (
	"sort"
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

	FindCalls int
	FindErr   error
	FindValue []Post

	FindAfterCalls int
	FindAfterErr   error
	FindAfterValue []Post

	CountCalls int
	CountErr   error
	CountValue int64

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

func (repository *MockPostRepository) Find(
	index int,
	amount int,
) ([]Post, error) {
	repository.FindCalls++

	if repository.FindErr != nil {
		return nil, repository.FindErr
	}

	if repository.FindValue != nil {
		return repository.FindValue, nil
	}

	items := repository.sortedPosts()

	if index >= len(items) {
		return []Post{}, nil
	}

	end := index + amount
	if end > len(items) {
		end = len(items)
	}

	return items[index:end], nil
}

func (repository *MockPostRepository) FindAfter(
	cursor uuid.UUID,
	amount int,
) ([]Post, error) {
	repository.FindAfterCalls++

	if repository.FindAfterErr != nil {
		return nil, repository.FindAfterErr
	}

	if repository.FindAfterValue != nil {
		if repository.FindAfterCalls == 1 {
			return repository.FindAfterValue, nil
		}

		return []Post{}, nil
	}

	items := repository.sortedPosts()

	start := 0

	if cursor != uuid.Nil {
		for index, post := range items {
			if post.Identifier() == cursor {
				start = index + 1
				break
			}
		}
	}

	if start >= len(items) {
		return []Post{}, nil
	}

	end := start + amount
	if end > len(items) {
		end = len(items)
	}

	return items[start:end], nil
}

func (repository *MockPostRepository) Count() (int64, error) {
	repository.CountCalls++

	if repository.CountErr != nil {
		return 0, repository.CountErr
	}

	if repository.CountValue != 0 {
		return repository.CountValue, nil
	}

	return int64(len(repository.Items)), nil
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

func (repository *MockPostRepository) sortedPosts() []Post {
	out := make([]Post, 0, len(repository.Items))

	for _, post := range repository.Items {
		out = append(out, post)
	}

	sort.Slice(out, func(left int, right int) bool {
		return out[left].Identifier().String() <
			out[right].Identifier().String()
	})

	return out
}
