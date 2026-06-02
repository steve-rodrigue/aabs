package users

import (
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
	"github.com/steve-rodrigue/aabs/services/saas/domain/platforms"
)

func NewMockUser(
	handle string,
	displayName string,
) User {
	return &MockUser{
		id:                uuid.New(),
		participationKind: participatables.UserKind,
		handle:            handle,
		displayName:       displayName,
	}
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		Items: map[uuid.UUID]User{},
	}
}

type MockUser struct {
	id                uuid.UUID
	participationKind participatables.Kind

	platform    platforms.Platform
	externalID  string
	handle      string
	displayName string
	profileURL  string
}

func (user *MockUser) Identifier() uuid.UUID {
	return user.id
}

func (user *MockUser) ParticipationKind() participatables.Kind {
	return user.participationKind
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
	FindByPlatformAndExternalIDValue User
	FindByPlatformAndExternalIDErr   error

	FindByPlatformAndHandleCalls int
	FindByPlatformAndHandleErr   error

	FindCalls int
	FindErr   error
	FindValue []User

	FindAfterCalls int
	FindAfterErr   error
	FindAfterValue []User

	CountCalls int
	CountErr   error
	CountValue int64
}

func (repository *MockUserRepository) Save(
	user User,
) error {
	repository.SaveCalls++

	return repository.SaveErr
}

func (repository *MockUserRepository) FindByID(
	id uuid.UUID,
) (User, error) {
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

	if repository.FindByPlatformAndExternalIDValue != nil {
		return repository.FindByPlatformAndExternalIDValue, nil
	}

	for _, user := range repository.Items {
		if user.Platform() == platform &&
			user.ExternalID() == externalID {
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
		if user.Platform() == platform &&
			user.Handle() == handle {
			return user, nil
		}
	}

	return nil, nil
}

func (repository *MockUserRepository) Find(
	index int,
	amount int,
) ([]User, error) {
	repository.FindCalls++

	if repository.FindErr != nil {
		return nil, repository.FindErr
	}

	if repository.FindValue != nil {
		return repository.FindValue, nil
	}

	users := repository.sortedUsers()

	if index >= len(users) {
		return []User{}, nil
	}

	end := index + amount
	if end > len(users) {
		end = len(users)
	}

	return users[index:end], nil
}

func (repository *MockUserRepository) FindAfter(
	cursor uuid.UUID,
	amount int,
) ([]User, error) {
	repository.FindAfterCalls++

	if repository.FindAfterErr != nil {
		return nil, repository.FindAfterErr
	}

	if repository.FindAfterValue != nil {
		return repository.FindAfterValue, nil
	}

	users := repository.sortedUsers()

	start := 0

	if cursor != uuid.Nil {
		for index, user := range users {
			if user.Identifier() == cursor {
				start = index + 1
				break
			}
		}
	}

	if start >= len(users) {
		return []User{}, nil
	}

	end := start + amount
	if end > len(users) {
		end = len(users)
	}

	return users[start:end], nil
}

func (repository *MockUserRepository) Count() (int64, error) {
	repository.CountCalls++

	if repository.CountErr != nil {
		return 0, repository.CountErr
	}

	if repository.CountValue != 0 {
		return repository.CountValue, nil
	}

	return int64(len(repository.Items)), nil
}

func (repository *MockUserRepository) sortedUsers() []User {
	users := make([]User, 0, len(repository.Items))

	for _, user := range repository.Items {
		users = append(users, user)
	}

	sort.Slice(users, func(left int, right int) bool {
		return users[left].Identifier().String() <
			users[right].Identifier().String()
	})

	return users
}
