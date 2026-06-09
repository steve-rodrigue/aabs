package narratives

import (
	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/communities"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/users"
	domain_narratives "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/narratives"
)

func NewMockNarrativesApplication() *MockNarrativesApplication {
	return &MockNarrativesApplication{}
}

type MockNarrativesApplication struct {
	FindByIDCalls int
	FindByIDErr   error
	FindByIDValue domain_narratives.Narrative

	FindCalls int
	FindErr   error
	FindValue []domain_narratives.Narrative

	FindAfterCalls int
	FindAfterErr   error
	FindAfterValue []domain_narratives.Narrative

	FindNarrativesByUserCalls int
	FindNarrativesByUserErr   error
	FindNarrativesByUserValue []domain_narratives.Narrative

	FindNarrativesByCommunityCalls int
	FindNarrativesByCommunityErr   error
	FindNarrativesByCommunityValue []domain_narratives.Narrative

	CountCalls int
	CountErr   error
	CountValue int64

	RebuildNarrativesCalls int
	RebuildNarrativesErr   error
}

func (application *MockNarrativesApplication) FindByID(
	id uuid.UUID,
) (domain_narratives.Narrative, error) {
	application.FindByIDCalls++

	return application.FindByIDValue, application.FindByIDErr
}

func (application *MockNarrativesApplication) Find(
	index int,
	amount int,
) ([]domain_narratives.Narrative, error) {
	application.FindCalls++

	return application.FindValue, application.FindErr
}

func (application *MockNarrativesApplication) FindAfter(
	cursor uuid.UUID,
	amount int,
) ([]domain_narratives.Narrative, error) {
	application.FindAfterCalls++

	return application.FindAfterValue, application.FindAfterErr
}

func (application *MockNarrativesApplication) FindNarrativesByUser(
	user users.User,
) ([]domain_narratives.Narrative, error) {
	application.FindNarrativesByUserCalls++

	return application.FindNarrativesByUserValue,
		application.FindNarrativesByUserErr
}

func (application *MockNarrativesApplication) FindNarrativesByCommunity(
	community communities.Community,
) ([]domain_narratives.Narrative, error) {
	application.FindNarrativesByCommunityCalls++

	return application.FindNarrativesByCommunityValue,
		application.FindNarrativesByCommunityErr
}

func (application *MockNarrativesApplication) Count() (int64, error) {
	application.CountCalls++

	return application.CountValue, application.CountErr
}

func (application *MockNarrativesApplication) RebuildNarratives() error {
	application.RebuildNarrativesCalls++

	return application.RebuildNarrativesErr
}
