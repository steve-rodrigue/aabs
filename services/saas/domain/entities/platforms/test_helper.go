package platforms

import (
	"context"
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

func NewMockPlatformAdapter() *MockPlatformAdapter {
	return &MockPlatformAdapter{}
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

	LastContext context.Context
	LastID      uuid.UUID
	LastHandle  string
	LastName    string
	LastIndex   int
	LastAmount  int
	LastCursor  uuid.UUID
	LastSaved   Platform
}

func (repository *MockPlatformRepository) Save(
	ctx context.Context,
	platform Platform,
) error {
	repository.SaveCalls++
	repository.LastContext = ctx
	repository.LastSaved = platform

	return repository.SaveErr
}

func (repository *MockPlatformRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (Platform, error) {
	repository.FindByIDCalls++
	repository.LastContext = ctx
	repository.LastID = id

	if repository.FindByIDErr != nil {
		return nil, repository.FindByIDErr
	}

	return repository.Items[id], nil
}

func (repository *MockPlatformRepository) FindByHandle(
	ctx context.Context,
	handle string,
) (Platform, error) {
	repository.FindByHandleCalls++
	repository.LastContext = ctx
	repository.LastHandle = handle

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
	ctx context.Context,
	name string,
) (Platform, error) {
	repository.FindByNameCalls++
	repository.LastContext = ctx
	repository.LastName = name

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
	ctx context.Context,
	index int,
	amount int,
) ([]Platform, error) {
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
	ctx context.Context,
	cursor uuid.UUID,
	amount int,
) ([]Platform, error) {
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

func (repository *MockPlatformRepository) Count(
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

type MockPlatformAdapter struct {
	ToDomainCalls int
	ToDomainErr   error
	ToDomainValue Platform

	LastInput PlatformInput
}

func (adapter *MockPlatformAdapter) ToDomain(
	input PlatformInput,
) (Platform, error) {
	adapter.ToDomainCalls++
	adapter.LastInput = input

	if adapter.ToDomainErr != nil {
		return nil, adapter.ToDomainErr
	}

	if adapter.ToDomainValue != nil {
		return adapter.ToDomainValue, nil
	}

	return &MockPlatform{
		id:      input.Identifier,
		name:    input.Name,
		handle:  input.Handle,
		baseURL: input.BaseURL,
	}, nil
}
