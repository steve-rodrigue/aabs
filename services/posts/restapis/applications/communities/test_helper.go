package communities

import (
	"context"

	"github.com/google/uuid"

	domain_communities "github.com/steve-rodrigue/aabs/services/posts/restapis/domain/communities"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/platforms"
)

func NewMockCommunitiesApplication() *MockCommunitiesApplication {
	return &MockCommunitiesApplication{}
}

type MockCommunitiesApplication struct {
	SaveCalls int
	SaveErr   error

	FindByIDCalls int
	FindByIDErr   error
	FindByIDValue domain_communities.Community

	FindByHandleCalls int
	FindByHandleErr   error
	FindByHandleValue domain_communities.Community

	FindCalls int
	FindErr   error
	FindValue []domain_communities.Community

	FindAfterCalls int
	FindAfterErr   error
	FindAfterValue []domain_communities.Community

	FindByPlatformCalls int
	FindByPlatformErr   error
	FindByPlatformValue []domain_communities.Community

	CountCalls int
	CountErr   error
	CountValue int64

	LastContext   context.Context
	LastCommunity domain_communities.Community
	LastID        uuid.UUID
	LastPlatform  platforms.Platform
	LastHandle    string
	LastIndex     int
	LastAmount    int
	LastCursor    uuid.UUID
}

func (application *MockCommunitiesApplication) Save(
	ctx context.Context,
	community domain_communities.Community,
) error {
	application.SaveCalls++
	application.LastContext = ctx
	application.LastCommunity = community

	return application.SaveErr
}

func (application *MockCommunitiesApplication) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (domain_communities.Community, error) {
	application.FindByIDCalls++
	application.LastContext = ctx
	application.LastID = id

	return application.FindByIDValue,
		application.FindByIDErr
}

func (application *MockCommunitiesApplication) FindByHandle(
	ctx context.Context,
	platform platforms.Platform,
	handle string,
) (domain_communities.Community, error) {
	application.FindByHandleCalls++
	application.LastContext = ctx
	application.LastPlatform = platform
	application.LastHandle = handle

	return application.FindByHandleValue,
		application.FindByHandleErr
}

func (application *MockCommunitiesApplication) Find(
	ctx context.Context,
	index int,
	amount int,
) ([]domain_communities.Community, error) {
	application.FindCalls++
	application.LastContext = ctx
	application.LastIndex = index
	application.LastAmount = amount

	return application.FindValue,
		application.FindErr
}

func (application *MockCommunitiesApplication) FindAfter(
	ctx context.Context,
	cursor uuid.UUID,
	amount int,
) ([]domain_communities.Community, error) {
	application.FindAfterCalls++
	application.LastContext = ctx
	application.LastCursor = cursor
	application.LastAmount = amount

	return application.FindAfterValue,
		application.FindAfterErr
}

func (application *MockCommunitiesApplication) FindByPlatform(
	ctx context.Context,
	platform platforms.Platform,
) ([]domain_communities.Community, error) {
	application.FindByPlatformCalls++
	application.LastContext = ctx
	application.LastPlatform = platform

	return application.FindByPlatformValue,
		application.FindByPlatformErr
}

func (application *MockCommunitiesApplication) Count(
	ctx context.Context,
) (int64, error) {
	application.CountCalls++
	application.LastContext = ctx

	return application.CountValue,
		application.CountErr
}
