package users

import (
	"context"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/platforms"
	domain_users "github.com/steve-rodrigue/aabs/services/posts/restapis/domain/users"
)

func NewMockUsersApplication() *MockUsersApplication {
	return &MockUsersApplication{}
}

type MockUsersApplication struct {
	SaveCalls int
	SaveErr   error
	LastUser  domain_users.User

	FindByIDCalls int
	FindByIDErr   error
	FindByIDValue domain_users.User

	FindByExternalIDCalls int
	FindByExternalIDErr   error
	FindByExternalIDValue domain_users.User

	FindByHandleCalls int
	FindByHandleErr   error
	FindByHandleValue domain_users.User

	FindCalls int
	FindErr   error
	FindValue []domain_users.User

	FindAfterCalls int
	FindAfterErr   error
	FindAfterValue []domain_users.User

	CountCalls int
	CountErr   error
	CountValue int64

	LastContext    context.Context
	LastID         uuid.UUID
	LastPlatform   platforms.Platform
	LastExternalID string
	LastHandle     string
	LastIndex      int
	LastAmount     int
	LastCursor     uuid.UUID
}

func (application *MockUsersApplication) Save(
	ctx context.Context,
	user domain_users.User,
) error {
	application.SaveCalls++
	application.LastContext = ctx
	application.LastUser = user

	return application.SaveErr
}

func (application *MockUsersApplication) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (domain_users.User, error) {
	application.FindByIDCalls++
	application.LastContext = ctx
	application.LastID = id

	return application.FindByIDValue, application.FindByIDErr
}

func (application *MockUsersApplication) FindByExternalID(
	ctx context.Context,
	platform platforms.Platform,
	externalID string,
) (domain_users.User, error) {
	application.FindByExternalIDCalls++
	application.LastContext = ctx
	application.LastPlatform = platform
	application.LastExternalID = externalID

	return application.FindByExternalIDValue,
		application.FindByExternalIDErr
}

func (application *MockUsersApplication) FindByHandle(
	ctx context.Context,
	platform platforms.Platform,
	handle string,
) (domain_users.User, error) {
	application.FindByHandleCalls++
	application.LastContext = ctx
	application.LastPlatform = platform
	application.LastHandle = handle

	return application.FindByHandleValue,
		application.FindByHandleErr
}

func (application *MockUsersApplication) Find(
	ctx context.Context,
	index int,
	amount int,
) ([]domain_users.User, error) {
	application.FindCalls++
	application.LastContext = ctx
	application.LastIndex = index
	application.LastAmount = amount

	return application.FindValue,
		application.FindErr
}

func (application *MockUsersApplication) FindAfter(
	ctx context.Context,
	cursor uuid.UUID,
	amount int,
) ([]domain_users.User, error) {
	application.FindAfterCalls++
	application.LastContext = ctx
	application.LastCursor = cursor
	application.LastAmount = amount

	return application.FindAfterValue,
		application.FindAfterErr
}

func (application *MockUsersApplication) Count(
	ctx context.Context,
) (int64, error) {
	application.CountCalls++
	application.LastContext = ctx

	return application.CountValue,
		application.CountErr
}
