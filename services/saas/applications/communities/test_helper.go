package communities

import (
	"github.com/google/uuid"

	domain_communities "github.com/steve-rodrigue/aabs/services/saas/domain/communities"
	"github.com/steve-rodrigue/aabs/services/saas/domain/platforms"
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
}

func (application *MockCommunitiesApplication) Save(
	community domain_communities.Community,
) error {
	application.SaveCalls++

	return application.SaveErr
}

func (application *MockCommunitiesApplication) FindByID(
	id uuid.UUID,
) (domain_communities.Community, error) {
	application.FindByIDCalls++

	return application.FindByIDValue,
		application.FindByIDErr
}

func (application *MockCommunitiesApplication) FindByHandle(
	platform platforms.Platform,
	handle string,
) (domain_communities.Community, error) {
	application.FindByHandleCalls++

	return application.FindByHandleValue,
		application.FindByHandleErr
}

func (application *MockCommunitiesApplication) Find(
	index int,
	amount int,
) ([]domain_communities.Community, error) {
	application.FindCalls++

	return application.FindValue,
		application.FindErr
}

func (application *MockCommunitiesApplication) FindAfter(
	cursor uuid.UUID,
	amount int,
) ([]domain_communities.Community, error) {
	application.FindAfterCalls++

	return application.FindAfterValue,
		application.FindAfterErr
}

func (application *MockCommunitiesApplication) FindByPlatform(
	platform platforms.Platform,
) ([]domain_communities.Community, error) {
	application.FindByPlatformCalls++

	return application.FindByPlatformValue,
		application.FindByPlatformErr
}

func (application *MockCommunitiesApplication) Count() (int64, error) {
	application.CountCalls++

	return application.CountValue,
		application.CountErr
}
