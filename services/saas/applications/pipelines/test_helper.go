package pipelines

import (
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts"
)

func NewMockPipelineApplication() *MockPipelineApplication {
	return &MockPipelineApplication{}
}

type MockPipelineApplication struct {
	ProcessPostCalls int
	ProcessPostErr   error
	LastPost         posts.Post

	ProcessPostsCalls int
	ProcessPostsErr   error
	LastPosts         []posts.Post

	RebuildCalls int
	RebuildErr   error
}

func (application *MockPipelineApplication) ProcessPost(
	post posts.Post,
) error {
	application.ProcessPostCalls++
	application.LastPost = post

	return application.ProcessPostErr
}

func (application *MockPipelineApplication) ProcessPosts(
	posts []posts.Post,
) error {
	application.ProcessPostsCalls++
	application.LastPosts = posts

	return application.ProcessPostsErr
}

func (application *MockPipelineApplication) Rebuild() error {
	application.RebuildCalls++

	return application.RebuildErr
}
