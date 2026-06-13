package platforms

import (
	"context"

	"github.com/google/uuid"

	domain_platforms "github.com/steve-rodrigue/aabs/services/posts/restapis/domain/platforms"
)

func NewMockPlatformsApplication() *MockPlatformsApplication {
	return &MockPlatformsApplication{}
}

type MockPlatformsApplication struct {
	SaveCalls int
	SaveErr   error
	LastSaved domain_platforms.Platform

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

	LastContext context.Context
	LastID      uuid.UUID
	LastHandle  string
	LastIndex   int
	LastAmount  int
	LastCursor  uuid.UUID
}

func (application *MockPlatformsApplication) Save(
	ctx context.Context,
	platform domain_platforms.Platform,
) error {
	application.SaveCalls++
	application.LastContext = ctx
	application.LastSaved = platform

	return application.SaveErr
}

func (application *MockPlatformsApplication) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (domain_platforms.Platform, error) {
	application.FindByIDCalls++
	application.LastContext = ctx
	application.LastID = id

	return application.FindByIDValue, application.FindByIDErr
}

func (application *MockPlatformsApplication) FindByHandle(
	ctx context.Context,
	handle string,
) (domain_platforms.Platform, error) {
	application.FindByHandleCalls++
	application.LastContext = ctx
	application.LastHandle = handle

	return application.FindByHandleValue, application.FindByHandleErr
}

func (application *MockPlatformsApplication) Find(
	ctx context.Context,
	index int,
	amount int,
) ([]domain_platforms.Platform, error) {
	application.FindCalls++
	application.LastContext = ctx
	application.LastIndex = index
	application.LastAmount = amount

	return application.FindValue, application.FindErr
}

func (application *MockPlatformsApplication) FindAfter(
	ctx context.Context,
	cursor uuid.UUID,
	amount int,
) ([]domain_platforms.Platform, error) {
	application.FindAfterCalls++
	application.LastContext = ctx
	application.LastCursor = cursor
	application.LastAmount = amount

	return application.FindAfterValue, application.FindAfterErr
}

func (application *MockPlatformsApplication) Count(
	ctx context.Context,
) (int64, error) {
	application.CountCalls++
	application.LastContext = ctx

	return application.CountValue, application.CountErr
}
