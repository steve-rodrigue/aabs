package communities

import (
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
	"github.com/steve-rodrigue/aabs/services/saas/domain/platforms"
	"github.com/steve-rodrigue/aabs/services/saas/domain/users"
)

func NewMockCommunity(
	title string,
	text string,
) Community {
	return &MockCommunity{
		id:                uuid.New(),
		participationKind: participatables.CommunityKind,
		title:             title,
		text:              text,
	}
}

func NewMockCommunityRepository() *MockCommunityRepository {
	return &MockCommunityRepository{
		Items: map[uuid.UUID]Community{},
	}
}

type MockCommunity struct {
	id                uuid.UUID
	participationKind participatables.Kind

	platform   platforms.Platform
	handle     string
	title      string
	text       string
	moderators []users.User
}

func (community *MockCommunity) Identifier() uuid.UUID {
	return community.id
}

func (community *MockCommunity) ParticipationKind() participatables.Kind {
	return community.participationKind
}

func (community *MockCommunity) Platform() platforms.Platform {
	return community.platform
}

func (community *MockCommunity) Handle() string {
	return community.handle
}

func (community *MockCommunity) Title() string {
	return community.title
}

func (community *MockCommunity) Text() string {
	return community.text
}

func (community *MockCommunity) CreatedOn() time.Time {
	return time.Time{}
}

func (community *MockCommunity) HasModerators() bool {
	return len(community.moderators) > 0
}

func (community *MockCommunity) Moderators() []users.User {
	return community.moderators
}

type MockCommunityRepository struct {
	SaveCalls int
	SaveErr   error

	Items map[uuid.UUID]Community

	FindByIDCalls int
	FindByIDErr   error

	FindByHandleCalls int
	FindByHandleErr   error
	FindByHandleValue Community

	FindCalls int
	FindErr   error
	FindValue []Community

	FindAfterCalls int
	FindAfterErr   error
	FindAfterValue []Community

	FindByPlatformCalls int
	FindByPlatformErr   error
	FindByPlatformValue []Community

	CountCalls int
	CountErr   error
	CountValue int64
}

func (repository *MockCommunityRepository) Save(
	community Community,
) error {
	repository.SaveCalls++

	return repository.SaveErr
}

func (repository *MockCommunityRepository) FindByID(
	id uuid.UUID,
) (Community, error) {
	repository.FindByIDCalls++

	if repository.FindByIDErr != nil {
		return nil, repository.FindByIDErr
	}

	return repository.Items[id], nil
}

func (repository *MockCommunityRepository) FindByHandle(
	platform platforms.Platform,
	handle string,
) (Community, error) {
	repository.FindByHandleCalls++

	if repository.FindByHandleErr != nil {
		return nil, repository.FindByHandleErr
	}

	if repository.FindByHandleValue != nil {
		return repository.FindByHandleValue, nil
	}

	for _, community := range repository.Items {
		if community.Platform() == platform &&
			community.Handle() == handle {
			return community, nil
		}
	}

	return nil, nil
}

func (repository *MockCommunityRepository) Find(
	index int,
	amount int,
) ([]Community, error) {
	repository.FindCalls++

	if repository.FindErr != nil {
		return nil, repository.FindErr
	}

	if repository.FindValue != nil {
		return repository.FindValue, nil
	}

	communities := repository.sortedCommunities()

	if index >= len(communities) {
		return []Community{}, nil
	}

	end := index + amount
	if end > len(communities) {
		end = len(communities)
	}

	return communities[index:end], nil
}

func (repository *MockCommunityRepository) FindAfter(
	cursor uuid.UUID,
	amount int,
) ([]Community, error) {
	repository.FindAfterCalls++

	if repository.FindAfterErr != nil {
		return nil, repository.FindAfterErr
	}

	if repository.FindAfterValue != nil {
		return repository.FindAfterValue, nil
	}

	communities := repository.sortedCommunities()

	start := 0

	if cursor != uuid.Nil {
		for index, community := range communities {
			if community.Identifier() == cursor {
				start = index + 1
				break
			}
		}
	}

	if start >= len(communities) {
		return []Community{}, nil
	}

	end := start + amount
	if end > len(communities) {
		end = len(communities)
	}

	return communities[start:end], nil
}

func (repository *MockCommunityRepository) FindByPlatform(
	platform platforms.Platform,
) ([]Community, error) {
	repository.FindByPlatformCalls++

	if repository.FindByPlatformErr != nil {
		return nil, repository.FindByPlatformErr
	}

	if repository.FindByPlatformValue != nil {
		return repository.FindByPlatformValue, nil
	}

	out := []Community{}

	for _, community := range repository.Items {
		if community.Platform() == platform {
			out = append(out, community)
		}
	}

	return out, nil
}

func (repository *MockCommunityRepository) Count() (int64, error) {
	repository.CountCalls++

	if repository.CountErr != nil {
		return 0, repository.CountErr
	}

	if repository.CountValue != 0 {
		return repository.CountValue, nil
	}

	return int64(len(repository.Items)), nil
}

func (repository *MockCommunityRepository) sortedCommunities() []Community {
	out := make([]Community, 0, len(repository.Items))

	for _, community := range repository.Items {
		out = append(out, community)
	}

	sort.Slice(out, func(left int, right int) bool {
		return out[left].Identifier().String() <
			out[right].Identifier().String()
	})

	return out
}
