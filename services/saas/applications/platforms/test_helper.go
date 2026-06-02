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

	FindCalls int
	FindErr   error
	FindValue []domain_platforms.Platform

	FindAfterCalls int
	FindAfterErr   error
	FindAfterValue []domain_platforms.Platform

	CountCalls int
	CountErr   error
	CountValue int64
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

func (application *MockPlatformsApplication) Find(
	index int,
	amount int,
) ([]domain_platforms.Platform, error) {
	application.FindCalls++

	return application.FindValue, application.FindErr
}

func (application *MockPlatformsApplication) FindAfter(
	cursor uuid.UUID,
	amount int,
) ([]domain_platforms.Platform, error) {
	application.FindAfterCalls++

	return application.FindAfterValue, application.FindAfterErr
}

func (application *MockPlatformsApplication) Count() (int64, error) {
	application.CountCalls++

	return application.CountValue, application.CountErr
}
