package narratives

import (
	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/communities"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/narratives"
	"github.com/steve-rodrigue/aabs/services/saas/domain/users"
)

type MockNarrativesApplication struct {
	RebuildNarrativesCalls int
	RebuildNarrativesErr   error
}

func (application *MockNarrativesApplication) FindByID(id uuid.UUID) (narratives.Narrative, error) {
	return nil, nil
}

func (application *MockNarrativesApplication) FindAll() ([]narratives.Narrative, error) {
	return nil, nil
}

func (application *MockNarrativesApplication) FindNarrativesByUser(
	user users.User,
) ([]narratives.Narrative, error) {
	return nil, nil
}

func (application *MockNarrativesApplication) FindNarrativesByCommunity(
	community communities.Community,
) ([]narratives.Narrative, error) {
	return nil, nil
}

func (application *MockNarrativesApplication) RebuildNarratives() error {
	application.RebuildNarrativesCalls++

	return application.RebuildNarrativesErr
}
