package platforms

import (
	"time"

	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
)

func NewMockPlatform(
	name string,
	handle string,
) Platform {
	return &MockPlatform{
		id:     uuid.New(),
		name:   name,
		handle: handle,
	}
}

func NewMockPlatformRepository() *MockPlatformRepository {
	return &MockPlatformRepository{
		Items: map[uuid.UUID]Platform{},
	}
}

type MockPlatform struct {
	id      uuid.UUID
	name    string
	handle  string
	baseURL string
}

func (platform *MockPlatform) Identifier() uuid.UUID {
	return platform.id
}

func (platform *MockPlatform) ParticipationKind() participatables.Kind {
	return participatables.PlatformKind
}

func (platform *MockPlatform) Name() string {
	return platform.name
}

func (platform *MockPlatform) Handle() string {
	return platform.handle
}

func (platform *MockPlatform) BaseURL() string {
	return platform.baseURL
}

func (platform *MockPlatform) CreatedOn() time.Time {
	return time.Time{}
}

type MockPlatformRepository struct {
	SaveCalls int
	SaveErr   error

	Items map[uuid.UUID]Platform

	FindByIDCalls int
	FindByIDErr   error

	FindByHandleCalls int
	FindByHandleErr   error

	FindByNameCalls int
	FindByNameErr   error

	FindAllCalls int
	FindAllErr   error
	FindAllValue []Platform
}

func (repository *MockPlatformRepository) Save(
	platform Platform,
) error {
	repository.SaveCalls++

	return repository.SaveErr
}

func (repository *MockPlatformRepository) FindByID(
	id uuid.UUID,
) (Platform, error) {
	repository.FindByIDCalls++

	if repository.FindByIDErr != nil {
		return nil, repository.FindByIDErr
	}

	return repository.Items[id], nil
}

func (repository *MockPlatformRepository) FindByHandle(
	handle string,
) (Platform, error) {
	repository.FindByHandleCalls++

	if repository.FindByHandleErr != nil {
		return nil, repository.FindByHandleErr
	}

	for _, platform := range repository.Items {
		if platform.Handle() == handle {
			return platform, nil
		}
	}

	return nil, nil
}

func (repository *MockPlatformRepository) FindByName(
	name string,
) (Platform, error) {
	repository.FindByNameCalls++

	if repository.FindByNameErr != nil {
		return nil, repository.FindByNameErr
	}

	for _, platform := range repository.Items {
		if platform.Name() == name {
			return platform, nil
		}
	}

	return nil, nil
}

func (repository *MockPlatformRepository) FindAll() ([]Platform, error) {
	repository.FindAllCalls++

	if repository.FindAllErr != nil {
		return nil, repository.FindAllErr
	}

	if repository.FindAllValue != nil {
		return repository.FindAllValue, nil
	}

	out := make([]Platform, 0, len(repository.Items))

	for _, platform := range repository.Items {
		out = append(out, platform)
	}

	return out, nil
}
