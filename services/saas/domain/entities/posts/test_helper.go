package posts

import (
	"context"
	"sort"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts/contents"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/users"
)

func NewMockPost(text string) Post {
	return &MockPost{
		ID: uuid.New(),
		ContentValue: &contents.MockContent{
			TextValue: text,
		},
		CreatedOnValue: time.Now().UTC(),
	}
}

func NewMockPostWithUser(
	text string,
	creator users.User,
) Post {
	return &MockPost{
		ID:           uuid.New(),
		CreatorValue: creator,
		ContentValue: &contents.MockContent{
			TextValue: text,
		},
		CreatedOnValue: time.Now().UTC(),
	}
}

func NewMockPostWithCommunities(
	text string,
	communityIDs []uuid.UUID,
) Post {
	ids := make([]uuid.UUID, len(communityIDs))
	copy(ids, communityIDs)

	return &MockPost{
		ID:                uuid.New(),
		CommunityIDsValue: ids,
		ContentValue: &contents.MockContent{
			TextValue: text,
		},
		CreatedOnValue: time.Now().UTC(),
	}
}

func NewMockPostRepository() *MockPostRepository {
	return &MockPostRepository{
		Items: map[uuid.UUID]Post{},
	}
}

func NewMockPostAdapter() *MockPostAdapter {
	return &MockPostAdapter{}
}

type MockPost struct {
	ID uuid.UUID

	CommunityIDsValue []uuid.UUID
	CreatorValue      users.User
	ContentValue      contents.Content
	CreatedOnValue    time.Time
}

func (post *MockPost) Identifier() uuid.UUID {
	return post.ID
}

func (post *MockPost) CommunityIDs() []uuid.UUID {
	out := make([]uuid.UUID, len(post.CommunityIDsValue))
	copy(out, post.CommunityIDsValue)

	return out
}

func (post *MockPost) Creator() users.User {
	return post.CreatorValue
}

func (post *MockPost) Content() contents.Content {
	return post.ContentValue
}

func (post *MockPost) CreatedOn() time.Time {
	return post.CreatedOnValue
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
		ID:                input.Identifier,
		CommunityIDsValue: communityIDs,
		CreatorValue:      input.Creator,
		ContentValue: &contents.MockContent{
			TextValue: input.Content.Thread.Text,
		},
		CreatedOnValue: input.CreatedOn,
	}, nil
}

type MockPostRepository struct {
	SaveCalls int
	SaveErr   error

	Items map[uuid.UUID]Post

	FindByIDCalls int
	FindByIDErr   error
	FindByIDValue Post

	FindCalls int
	FindErr   error
	FindValue []Post

	FindAfterCalls int
	FindAfterErr   error
	FindAfterValue []Post

	FindByCriteriaCalls int
	FindByCriteriaErr   error
	FindByCriteriaValue []Post

	FindByCriteriaAfterCalls int
	FindByCriteriaAfterErr   error
	FindByCriteriaAfterValue []Post

	CountCalls int
	CountErr   error
	CountValue int64

	CountByCriteriaCalls int
	CountByCriteriaErr   error
	CountByCriteriaValue int64

	LastContext  context.Context
	LastSaved    Post
	LastID       uuid.UUID
	LastIndex    int
	LastAmount   int
	LastCursor   uuid.UUID
	LastCriteria Criteria
}

func (repository *MockPostRepository) Save(
	ctx context.Context,
	post Post,
) error {
	repository.SaveCalls++
	repository.LastContext = ctx
	repository.LastSaved = post

	if repository.Items != nil && post != nil {
		repository.Items[post.Identifier()] = post
	}

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

	if repository.FindByIDValue != nil {
		return repository.FindByIDValue, nil
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

	return paginatePosts(
		repository.sortedPosts(),
		index,
		amount,
	), nil
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

	return paginatePostsAfter(
		repository.sortedPosts(),
		cursor,
		amount,
	), nil
}

func (repository *MockPostRepository) FindByCriteria(
	ctx context.Context,
	criteria Criteria,
	index int,
	amount int,
) ([]Post, error) {
	repository.FindByCriteriaCalls++
	repository.LastContext = ctx
	repository.LastCriteria = criteria
	repository.LastIndex = index
	repository.LastAmount = amount

	if repository.FindByCriteriaErr != nil {
		return nil, repository.FindByCriteriaErr
	}

	if repository.FindByCriteriaValue != nil {
		return repository.FindByCriteriaValue, nil
	}

	return paginatePosts(
		repository.filterByCriteria(criteria),
		index,
		amount,
	), nil
}

func (repository *MockPostRepository) FindByCriteriaAfter(
	ctx context.Context,
	criteria Criteria,
	cursor uuid.UUID,
	amount int,
) ([]Post, error) {
	repository.FindByCriteriaAfterCalls++
	repository.LastContext = ctx
	repository.LastCriteria = criteria
	repository.LastCursor = cursor
	repository.LastAmount = amount

	if repository.FindByCriteriaAfterErr != nil {
		return nil, repository.FindByCriteriaAfterErr
	}

	if repository.FindByCriteriaAfterValue != nil {
		if repository.FindByCriteriaAfterCalls == 1 {
			return repository.FindByCriteriaAfterValue, nil
		}

		return []Post{}, nil
	}

	return paginatePostsAfter(
		repository.filterByCriteria(criteria),
		cursor,
		amount,
	), nil
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

func (repository *MockPostRepository) CountByCriteria(
	ctx context.Context,
	criteria Criteria,
) (int64, error) {
	repository.CountByCriteriaCalls++
	repository.LastContext = ctx
	repository.LastCriteria = criteria

	if repository.CountByCriteriaErr != nil {
		return 0, repository.CountByCriteriaErr
	}

	if repository.CountByCriteriaValue != 0 {
		return repository.CountByCriteriaValue, nil
	}

	return int64(len(repository.filterByCriteria(criteria))), nil
}

func (repository *MockPostRepository) filterByCriteria(
	criteria Criteria,
) []Post {
	out := []Post{}

	for _, post := range repository.sortedPosts() {
		if !postMatchesCriteria(post, criteria) {
			continue
		}

		out = append(out, post)
	}

	return out
}

func postMatchesCriteria(
	post Post,
	criteria Criteria,
) bool {
	if len(criteria.UserIDs) > 0 {
		if post.Creator() == nil ||
			!uuidIn(criteria.UserIDs, post.Creator().Identifier()) {
			return false
		}
	}

	if len(criteria.PlatformIDs) > 0 {
		if post.Creator() == nil ||
			post.Creator().Platform() == nil ||
			!uuidIn(criteria.PlatformIDs, post.Creator().Platform().Identifier()) {
			return false
		}
	}

	if len(criteria.CommunityIDs) > 0 {
		matched := false

		for _, communityID := range post.CommunityIDs() {
			if uuidIn(criteria.CommunityIDs, communityID) {
				matched = true
				break
			}
		}

		if !matched {
			return false
		}
	}

	return true
}

func paginatePosts(
	items []Post,
	index int,
	amount int,
) []Post {
	if index >= len(items) {
		return []Post{}
	}

	end := index + amount
	if end > len(items) {
		end = len(items)
	}

	return items[index:end]
}

func paginatePostsAfter(
	items []Post,
	cursor uuid.UUID,
	amount int,
) []Post {
	start := 0

	if cursor != uuid.Nil {
		for index, post := range items {
			if post.Identifier() == cursor {
				start = index + 1
				break
			}
		}
	}

	return paginatePosts(items, start, amount)
}

func uuidIn(
	ids []uuid.UUID,
	id uuid.UUID,
) bool {
	for _, current := range ids {
		if current == id {
			return true
		}
	}

	return false
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
