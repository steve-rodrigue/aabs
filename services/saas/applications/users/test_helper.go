package users

import (
	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/platforms"
	domain_users "github.com/steve-rodrigue/aabs/services/saas/domain/users"
)

func NewMockUsersApplication() *MockUsersApplication {
	return &MockUsersApplication{}
}

type MockUsersApplication struct {
	SaveCalls int
	SaveErr   error

	FindByIDCalls int
	FindByIDErr   error
	FindByIDValue domain_users.User

	FindByExternalIDCalls int
	FindByExternalIDErr   error
	FindByExternalIDValue domain_users.User

	FindCalls int
	FindErr   error
	FindValue []domain_users.User

	FindAfterCalls int
	FindAfterErr   error
	FindAfterValue []domain_users.User

	CountCalls int
	CountErr   error
	CountValue int64
}

func (application *MockUsersApplication) Save(
	user domain_users.User,
) error {
	application.SaveCalls++

	return application.SaveErr
}

func (application *MockUsersApplication) FindByID(
	id uuid.UUID,
) (domain_users.User, error) {
	application.FindByIDCalls++

	return application.FindByIDValue, application.FindByIDErr
}

func (application *MockUsersApplication) FindByExternalID(
	platform platforms.Platform,
	externalID string,
) (domain_users.User, error) {
	application.FindByExternalIDCalls++

	return application.FindByExternalIDValue,
		application.FindByExternalIDErr
}

func (application *MockUsersApplication) Find(
	index int,
	amount int,
) ([]domain_users.User, error) {
	application.FindCalls++

	return application.FindValue,
		application.FindErr
}

func (application *MockUsersApplication) FindAfter(
	cursor uuid.UUID,
	amount int,
) ([]domain_users.User, error) {
	application.FindAfterCalls++

	return application.FindAfterValue,
		application.FindAfterErr
}

func (application *MockUsersApplication) Count() (int64, error) {
	application.CountCalls++

	return application.CountValue,
		application.CountErr
}
