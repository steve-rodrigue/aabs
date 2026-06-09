package posts

import (
	"context"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/communities"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/platforms"
	domain_posts "github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/users"
)

func NewMockPostsApplication() *MockPostsApplication {
	return &MockPostsApplication{}
}

type MockPostsApplication struct {
	SaveCalls int
	SaveErr   error
	LastPost  domain_posts.Post

	FindByIDCalls int
	FindByIDErr   error
	FindByIDValue domain_posts.Post

	FindCalls int
	FindErr   error
	FindValue []domain_posts.Post

	FindAfterCalls int
	FindAfterErr   error
	FindAfterValue []domain_posts.Post

	CountCalls int
	CountErr   error
	CountValue int64

	FindByUserCalls int
	FindByUserErr   error
	FindByUserValue []domain_posts.Post

	FindByCommunityCalls int
	FindByCommunityErr   error
	FindByCommunityValue []domain_posts.Post

	FindByPlatformCalls int
	FindByPlatformErr   error
	FindByPlatformValue []domain_posts.Post

	LastContext   context.Context
	LastID        uuid.UUID
	LastIndex     int
	LastAmount    int
	LastCursor    uuid.UUID
	LastUser      users.User
	LastCommunity communities.Community
	LastPlatform  platforms.Platform
}

func (application *MockPostsApplication) Save(
	ctx context.Context,
	post domain_posts.Post,
) error {
	application.SaveCalls++
	application.LastContext = ctx
	application.LastPost = post

	return application.SaveErr
}

func (application *MockPostsApplication) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (domain_posts.Post, error) {
	application.FindByIDCalls++
	application.LastContext = ctx
	application.LastID = id

	return application.FindByIDValue, application.FindByIDErr
}

func (application *MockPostsApplication) Find(
	ctx context.Context,
	index int,
	amount int,
) ([]domain_posts.Post, error) {
	application.FindCalls++
	application.LastContext = ctx
	application.LastIndex = index
	application.LastAmount = amount

	return application.FindValue, application.FindErr
}

func (application *MockPostsApplication) FindAfter(
	ctx context.Context,
	cursor uuid.UUID,
	amount int,
) ([]domain_posts.Post, error) {
	application.FindAfterCalls++
	application.LastContext = ctx
	application.LastCursor = cursor
	application.LastAmount = amount

	if application.FindAfterErr != nil {
		return nil, application.FindAfterErr
	}

	if application.FindAfterValue != nil {
		if application.FindAfterCalls == 1 {
			return application.FindAfterValue, nil
		}

		return []domain_posts.Post{}, nil
	}

	return []domain_posts.Post{}, nil
}

func (application *MockPostsApplication) Count(
	ctx context.Context,
) (
	int64,
	error,
) {
	application.CountCalls++
	application.LastContext = ctx

	return application.CountValue, application.CountErr
}

func (application *MockPostsApplication) FindByUser(
	ctx context.Context,
	user users.User,
) ([]domain_posts.Post, error) {
	application.FindByUserCalls++
	application.LastContext = ctx
	application.LastUser = user

	return application.FindByUserValue, application.FindByUserErr
}

func (application *MockPostsApplication) FindByCommunity(
	ctx context.Context,
	community communities.Community,
) ([]domain_posts.Post, error) {
	application.FindByCommunityCalls++
	application.LastContext = ctx
	application.LastCommunity = community

	return application.FindByCommunityValue, application.FindByCommunityErr
}

func (application *MockPostsApplication) FindByPlatform(
	ctx context.Context,
	platform platforms.Platform,
) ([]domain_posts.Post, error) {
	application.FindByPlatformCalls++
	application.LastContext = ctx
	application.LastPlatform = platform

	return application.FindByPlatformValue, application.FindByPlatformErr
}
