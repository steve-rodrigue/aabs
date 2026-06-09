package posts

import (
	"context"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/communities"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/platforms"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts/contents"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/users"
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
	ids := make([]uuid.UUID, len(communityIDs))
	copy(ids, communityIDs)

	return &MockPost{
		id:           uuid.New(),
		communityIDs: ids,
		content: &contents.MockContent{
			TextValue: text,
		},
	}
}

type MockPost struct {
	id uuid.UUID

	communityIDs []uuid.UUID
	creator      users.User
	content      contents.Content
	createdOn    time.Time
}

func (post *MockPost) Identifier() uuid.UUID {
	return post.id
}

func (post *MockPost) CommunityIDs() []uuid.UUID {
	out := make([]uuid.UUID, len(post.communityIDs))
	copy(out, post.communityIDs)

	return out
}

func (post *MockPost) Creator() users.User {
	return post.creator
}

func (post *MockPost) Content() contents.Content {
	return post.content
}

func (post *MockPost) CreatedOn() time.Time {
	return post.createdOn
}

func NewMockPostAdapter() *MockPostAdapter {
	return &MockPostAdapter{}
}

type MockPostAdapter struct {
	ToDomainCalls int
	ToDomainErr   error
	ToDomainValue Post

	LastInput PostInput
}

func (adapter *MockPostAdapter) ToDomain(
	input PostInput,
) (Post, error) {
	adapter.ToDomainCalls++
	adapter.LastInput = input

	if adapter.ToDomainErr != nil {
		return nil, adapter.ToDomainErr
	}

	if adapter.ToDomainValue != nil {
		return adapter.ToDomainValue, nil
	}

	communityIDs := make([]uuid.UUID, len(input.CommunityIDs))
	copy(communityIDs, input.CommunityIDs)

	return &MockPost{
		id:           input.Identifier,
		communityIDs: communityIDs,
		creator:      input.Creator,
		content: &contents.MockContent{
			TextValue: input.Content.Thread.Text,
		},
		createdOn: input.CreatedOn,
	}, nil
}

func NewMockPostRepository() *MockPostRepository {
	return &MockPostRepository{
		Items: map[uuid.UUID]Post{},
	}
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

	LastContext   context.Context
	LastSaved     Post
	LastID        uuid.UUID
	LastIndex     int
	LastAmount    int
	LastCursor    uuid.UUID
	LastUser      users.User
	LastCommunity communities.Community
	LastPlatform  platforms.Platform
}

func (repository *MockPostRepository) Save(
	ctx context.Context,
	post Post,
) error {
	repository.SaveCalls++
	repository.LastContext = ctx
	repository.LastSaved = post

	return repository.SaveErr
}

func (repository *MockPostRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (Post, error) {
	repository.FindByIDCalls++
	repository.LastContext = ctx
	repository.LastID = id

	if repository.FindByIDErr != nil {
		return nil, repository.FindByIDErr
	}

	if repository.Items == nil {
		return nil, nil
	}

	return repository.Items[id], nil
}

func (repository *MockPostRepository) Find(
	ctx context.Context,
	index int,
	amount int,
) ([]Post, error) {
	repository.FindCalls++
	repository.LastContext = ctx
	repository.LastIndex = index
	repository.LastAmount = amount

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
	ctx context.Context,
	cursor uuid.UUID,
	amount int,
) ([]Post, error) {
	repository.FindAfterCalls++
	repository.LastContext = ctx
	repository.LastCursor = cursor
	repository.LastAmount = amount

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

func (repository *MockPostRepository) Count(
	ctx context.Context,
) (int64, error) {
	repository.CountCalls++
	repository.LastContext = ctx

	if repository.CountErr != nil {
		return 0, repository.CountErr
	}

	if repository.CountValue != 0 {
		return repository.CountValue, nil
	}

	return int64(len(repository.Items)), nil
}

func (repository *MockPostRepository) FindByUser(
	ctx context.Context,
	user users.User,
) ([]Post, error) {
	repository.FindByUserCalls++
	repository.LastContext = ctx
	repository.LastUser = user

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
	ctx context.Context,
	community communities.Community,
) ([]Post, error) {
	repository.FindByCommunityCalls++
	repository.LastContext = ctx
	repository.LastCommunity = community

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
	ctx context.Context,
	platform platforms.Platform,
) ([]Post, error) {
	repository.FindByPlatformCalls++
	repository.LastContext = ctx
	repository.LastPlatform = platform

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
