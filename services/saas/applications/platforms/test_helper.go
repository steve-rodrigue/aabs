package platforms

import (
	"github.com/google/uuid"

	domain_platforms "github.com/steve-rodrigue/aabs/services/saas/domain/platforms"
)

func NewMockPlatformsApplication() *MockPlatformsApplication {
	return &MockPlatformsApplication{}
}

type MockPlatformsApplication struct {
	SaveCalls int
	SaveErr   error

	FindByIDCalls int
	FindByIDErr   error
	FindByIDValue domain_platforms.Platform

	FindByHandleCalls int
	FindByHandleErr   error
	FindByHandleValue domain_platforms.Platform

	FindAllCalls int
	FindAllErr   error
	FindAllValue []domain_platforms.Platform
}

func (application *MockPlatformsApplication) Save(
	platform domain_platforms.Platform,
) error {
	application.SaveCalls++

	return application.SaveErr
}

func (application *MockPlatformsApplication) FindByID(
	id uuid.UUID,
) (domain_platforms.Platform, error) {
	application.FindByIDCalls++

	return application.FindByIDValue, application.FindByIDErr
}

func (application *MockPlatformsApplication) FindByHandle(
	handle string,
) (domain_platforms.Platform, error) {
	application.FindByHandleCalls++

	return application.FindByHandleValue, application.FindByHandleErr
}

func (application *MockPlatformsApplication) FindAll() ([]domain_platforms.Platform, error) {
	application.FindAllCalls++

	return application.FindAllValue, application.FindAllErr
}
