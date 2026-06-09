package users

import (
	"context"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/platforms"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
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

type MockUser struct {
	id                uuid.UUID
	participationKind participatables.Kind

	platform    platforms.Platform
	externalID  string
	handle      string
	displayName string
	profileURL  string
	createdOn   time.Time
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
	return user.createdOn
}

func NewMockUserAdapter() *MockUserAdapter {
	return &MockUserAdapter{}
}

type MockUserAdapter struct {
	ToDomainCalls int
	ToDomainErr   error
	ToDomainValue User

	LastInput UserInput
}

func (adapter *MockUserAdapter) ToDomain(
	input UserInput,
) (User, error) {
	adapter.ToDomainCalls++
	adapter.LastInput = input

	if adapter.ToDomainErr != nil {
		return nil, adapter.ToDomainErr
	}

	if adapter.ToDomainValue != nil {
		return adapter.ToDomainValue, nil
	}

	return &MockUser{
		id:                input.Identifier,
		participationKind: participatables.UserKind,
		platform:          input.Platform,
		externalID:        input.ExternalID,
		handle:            input.Handle,
		displayName:       input.DisplayName,
		profileURL:        input.ProfileURL,
		createdOn:         input.CreatedOn,
	}, nil
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		Items: map[uuid.UUID]User{},
	}
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
	FindByPlatformAndHandleValue User
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

	LastContext    context.Context
	LastSaved      User
	LastID         uuid.UUID
	LastPlatform   platforms.Platform
	LastExternalID string
	LastHandle     string
	LastIndex      int
	LastAmount     int
	LastCursor     uuid.UUID
}

func (repository *MockUserRepository) Save(
	ctx context.Context,
	user User,
) error {
	repository.SaveCalls++
	repository.LastContext = ctx
	repository.LastSaved = user

	return repository.SaveErr
}

func (repository *MockUserRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (User, error) {
	repository.FindByIDCalls++
	repository.LastContext = ctx
	repository.LastID = id

	if repository.FindByIDErr != nil {
		return nil, repository.FindByIDErr
	}

	return repository.Items[id], nil
}

func (repository *MockUserRepository) FindByPlatformAndExternalID(
	ctx context.Context,
	platform platforms.Platform,
	externalID string,
) (User, error) {
	repository.FindByPlatformAndExternalIDCalls++
	repository.LastContext = ctx
	repository.LastPlatform = platform
	repository.LastExternalID = externalID

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
	ctx context.Context,
	platform platforms.Platform,
	handle string,
) (User, error) {
	repository.FindByPlatformAndHandleCalls++
	repository.LastContext = ctx
	repository.LastPlatform = platform
	repository.LastHandle = handle

	if repository.FindByPlatformAndHandleErr != nil {
		return nil, repository.FindByPlatformAndHandleErr
	}

	if repository.FindByPlatformAndHandleValue != nil {
		return repository.FindByPlatformAndHandleValue, nil
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
	ctx context.Context,
	index int,
	amount int,
) ([]User, error) {
	repository.FindCalls++
	repository.LastContext = ctx
	repository.LastIndex = index
	repository.LastAmount = amount

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
	ctx context.Context,
	cursor uuid.UUID,
	amount int,
) ([]User, error) {
	repository.FindAfterCalls++
	repository.LastContext = ctx
	repository.LastCursor = cursor
	repository.LastAmount = amount

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

func (repository *MockUserRepository) Count(
	ctx context.Context,
) (int64, error) {
	repository.CountCalls++
	repository.LastContext = ctx

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
