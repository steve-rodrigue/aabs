package posts

import (
	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/communities"
	"github.com/steve-rodrigue/aabs/services/saas/domain/platforms"
	domain_posts "github.com/steve-rodrigue/aabs/services/saas/domain/posts"
	"github.com/steve-rodrigue/aabs/services/saas/domain/users"
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
}

func (application *MockPostsApplication) Save(
	post domain_posts.Post,
) error {
	application.SaveCalls++
	application.LastPost = post

	return application.SaveErr
}

func (application *MockPostsApplication) FindByID(
	id uuid.UUID,
) (domain_posts.Post, error) {
	application.FindByIDCalls++

	return application.FindByIDValue, application.FindByIDErr
}

func (application *MockPostsApplication) Find(
	index int,
	amount int,
) ([]domain_posts.Post, error) {
	application.FindCalls++

	return application.FindValue, application.FindErr
}

func (application *MockPostsApplication) FindAfter(
	cursor uuid.UUID,
	amount int,
) ([]domain_posts.Post, error) {
	application.FindAfterCalls++

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

func (application *MockPostsApplication) Count() (
	int64,
	error,
) {
	application.CountCalls++

	return application.CountValue, application.CountErr
}

func (application *MockPostsApplication) FindByUser(
	user users.User,
) ([]domain_posts.Post, error) {
	application.FindByUserCalls++

	return application.FindByUserValue, application.FindByUserErr
}

func (application *MockPostsApplication) FindByCommunity(
	community communities.Community,
) ([]domain_posts.Post, error) {
	application.FindByCommunityCalls++

	return application.FindByCommunityValue, application.FindByCommunityErr
}

func (application *MockPostsApplication) FindByPlatform(
	platform platforms.Platform,
) ([]domain_posts.Post, error) {
	application.FindByPlatformCalls++

	return application.FindByPlatformValue, application.FindByPlatformErr
}
