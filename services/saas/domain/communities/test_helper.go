package communities

import (
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
		id:    uuid.New(),
		title: title,
		text:  text,
	}
}

type MockCommunity struct {
	id         uuid.UUID
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
	return participatables.CommunityKind
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
}

func NewMockCommunityRepository() *MockCommunityRepository {
	return &MockCommunityRepository{
		Items: map[uuid.UUID]Community{},
	}
}

func (repository *MockCommunityRepository) Save(community Community) error {
	repository.SaveCalls++

	return repository.SaveErr
}

func (repository *MockCommunityRepository) FindByID(id uuid.UUID) (Community, error) {
	repository.FindByIDCalls++

	if repository.FindByIDErr != nil {
		return nil, repository.FindByIDErr
	}

	return repository.Items[id], nil
}

func (repository *MockCommunityRepository) FindByHandle(handle string) (Community, error) {
	repository.FindByHandleCalls++

	if repository.FindByHandleErr != nil {
		return nil, repository.FindByHandleErr
	}

	for _, community := range repository.Items {
		if community.Handle() == handle {
			return community, nil
		}
	}

	return nil, nil
}
