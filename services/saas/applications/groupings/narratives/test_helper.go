package narratives

import (
	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/communities"
	domain_narratives "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/narratives"
	"github.com/steve-rodrigue/aabs/services/saas/domain/users"
)

func NewMockNarrativesApplication() *MockNarrativesApplication {
	return &MockNarrativesApplication{}
}

type MockNarrativesApplication struct {
	FindByIDCalls int
	FindByIDErr   error
	FindByIDValue domain_narratives.Narrative

	FindAllCalls int
	FindAllErr   error
	FindAllValue []domain_narratives.Narrative

	FindNarrativesByUserCalls int
	FindNarrativesByUserErr   error
	FindNarrativesByUserValue []domain_narratives.Narrative

	FindNarrativesByCommunityCalls int
	FindNarrativesByCommunityErr   error
	FindNarrativesByCommunityValue []domain_narratives.Narrative

	RebuildNarrativesCalls int
	RebuildNarrativesErr   error
}

func (application *MockNarrativesApplication) FindByID(
	id uuid.UUID,
) (domain_narratives.Narrative, error) {
	application.FindByIDCalls++

	return application.FindByIDValue, application.FindByIDErr
}

func (application *MockNarrativesApplication) FindAll() (
	[]domain_narratives.Narrative,
	error,
) {
	application.FindAllCalls++

	return application.FindAllValue, application.FindAllErr
}

func (application *MockNarrativesApplication) FindNarrativesByUser(
	user users.User,
) ([]domain_narratives.Narrative, error) {
	application.FindNarrativesByUserCalls++

	return application.FindNarrativesByUserValue, application.FindNarrativesByUserErr
}

func (application *MockNarrativesApplication) FindNarrativesByCommunity(
	community communities.Community,
) ([]domain_narratives.Narrative, error) {
	application.FindNarrativesByCommunityCalls++

	return application.FindNarrativesByCommunityValue, application.FindNarrativesByCommunityErr
}

func (application *MockNarrativesApplication) RebuildNarratives() error {
	application.RebuildNarrativesCalls++

	return application.RebuildNarrativesErr
}
