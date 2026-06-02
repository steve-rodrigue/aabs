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

	FindAllCalls int
	FindAllErr   error
	FindAllValue []domain_posts.Post

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

func (application *MockPostsApplication) FindAll() ([]domain_posts.Post, error) {
	application.FindAllCalls++

	return application.FindAllValue, application.FindAllErr
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
