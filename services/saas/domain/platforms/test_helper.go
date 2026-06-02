package platforms

import (
	"sort"
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
	FindByHandleValue Platform

	FindByNameCalls int
	FindByNameErr   error
	FindByNameValue Platform

	FindCalls int
	FindErr   error
	FindValue []Platform

	FindAfterCalls int
	FindAfterErr   error
	FindAfterValue []Platform

	CountCalls int
	CountErr   error
	CountValue int64
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

	if repository.FindByHandleValue != nil {
		return repository.FindByHandleValue, nil
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

	if repository.FindByNameValue != nil {
		return repository.FindByNameValue, nil
	}

	for _, platform := range repository.Items {
		if platform.Name() == name {
			return platform, nil
		}
	}

	return nil, nil
}

func (repository *MockPlatformRepository) Find(
	index int,
	amount int,
) ([]Platform, error) {
	repository.FindCalls++

	if repository.FindErr != nil {
		return nil, repository.FindErr
	}

	if repository.FindValue != nil {
		return repository.FindValue, nil
	}

	items := repository.sortedPlatforms()

	if index >= len(items) {
		return []Platform{}, nil
	}

	end := index + amount
	if end > len(items) {
		end = len(items)
	}

	return items[index:end], nil
}

func (repository *MockPlatformRepository) FindAfter(
	cursor uuid.UUID,
	amount int,
) ([]Platform, error) {
	repository.FindAfterCalls++

	if repository.FindAfterErr != nil {
		return nil, repository.FindAfterErr
	}

	if repository.FindAfterValue != nil {
		return repository.FindAfterValue, nil
	}

	items := repository.sortedPlatforms()

	start := 0

	if cursor != uuid.Nil {
		for index, platform := range items {
			if platform.Identifier() == cursor {
				start = index + 1
				break
			}
		}
	}

	if start >= len(items) {
		return []Platform{}, nil
	}

	end := start + amount
	if end > len(items) {
		end = len(items)
	}

	return items[start:end], nil
}

func (repository *MockPlatformRepository) Count() (int64, error) {
	repository.CountCalls++

	if repository.CountErr != nil {
		return 0, repository.CountErr
	}

	if repository.CountValue != 0 {
		return repository.CountValue, nil
	}

	return int64(len(repository.Items)), nil
}

func (repository *MockPlatformRepository) sortedPlatforms() []Platform {
	out := make([]Platform, 0, len(repository.Items))

	for _, platform := range repository.Items {
		out = append(out, platform)
	}

	sort.Slice(out, func(left int, right int) bool {
		return out[left].Identifier().String() <
			out[right].Identifier().String()
	})

	return out
}
