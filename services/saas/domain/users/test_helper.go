package users

import (
	"time"

	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/platforms"
)

func NewMockUser(
	handle string,
	displayName string,
) User {
	return &MockUser{
		id:          uuid.New(),
		handle:      handle,
		displayName: displayName,
	}
}

type MockUser struct {
	id          uuid.UUID
	platform    platforms.Platform
	externalID  string
	handle      string
	displayName string
	profileURL  string
}

func (user *MockUser) Identifier() uuid.UUID {
	return user.id
}

func (user *MockUser) Platform() platforms.Platform {
	return user.platform
}

func (user *MockUser) ExternalID() string {
	return user.externalID
}

func (user *MockUser) Handle() string {
	return user.handle
}

func (user *MockUser) DisplayName() string {
	return user.displayName
}

func (user *MockUser) ProfileURL() string {
	return user.profileURL
}

func (user *MockUser) CreatedOn() time.Time {
	return time.Time{}
}

type MockUserRepository struct {
	SaveCalls int
	SaveErr   error

	Items map[uuid.UUID]User

	FindByIDCalls int
	FindByIDErr   error

	FindByPlatformAndExternalIDCalls int
	FindByPlatformAndExternalIDErr   error

	FindByPlatformAndHandleCalls int
	FindByPlatformAndHandleErr   error
}

func (repository *MockUserRepository) Save(user User) error {
	repository.SaveCalls++

	return repository.SaveErr
}

func (repository *MockUserRepository) FindByID(id uuid.UUID) (User, error) {
	repository.FindByIDCalls++

	if repository.FindByIDErr != nil {
		return nil, repository.FindByIDErr
	}

	return repository.Items[id], nil
}

func (repository *MockUserRepository) FindByPlatformAndExternalID(
	platform platforms.Platform,
	externalID string,
) (User, error) {
	repository.FindByPlatformAndExternalIDCalls++

	if repository.FindByPlatformAndExternalIDErr != nil {
		return nil, repository.FindByPlatformAndExternalIDErr
	}

	for _, user := range repository.Items {
		if user.Platform() == platform && user.ExternalID() == externalID {
			return user, nil
		}
	}

	return nil, nil
}

func (repository *MockUserRepository) FindByPlatformAndHandle(
	platform platforms.Platform,
	handle string,
) (User, error) {
	repository.FindByPlatformAndHandleCalls++

	if repository.FindByPlatformAndHandleErr != nil {
		return nil, repository.FindByPlatformAndHandleErr
	}

	for _, user := range repository.Items {
		if user.Platform() == platform && user.Handle() == handle {
			return user, nil
		}
	}

	return nil, nil
}
