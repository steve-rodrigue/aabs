package namers

import (
	"context"

	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts"
)

func NewMockNamer() *MockNamer {
	return &MockNamer{}
}

type MockNamer struct {
	NameCalls int
	NameErr   error
	NameValue string

	LastContext context.Context
	LastPosts   []posts.Post
}

func (namer *MockNamer) Name(
	ctx context.Context,
	posts []posts.Post,
) (string, error) {
	namer.NameCalls++

	namer.LastContext = ctx
	namer.LastPosts = posts

	return namer.NameValue,
		namer.NameErr
}
