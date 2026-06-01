package groupings

import (
	"github.com/steve-rodrigue/aabs/services/saas/applications/groupings/campaigns"
	"github.com/steve-rodrigue/aabs/services/saas/applications/groupings/clusters"
	"github.com/steve-rodrigue/aabs/services/saas/applications/groupings/narratives"
	"github.com/steve-rodrigue/aabs/services/saas/applications/groupings/participations"
	"github.com/steve-rodrigue/aabs/services/saas/applications/groupings/topics"
)

type MockGroupingsApplication struct {
	ClustersIns       clusters.Application
	CampaignsIns      campaigns.Application
	TopicsIns         topics.Application
	NarrativesIns     narratives.Application
	ParticipationsIns participations.Application
}

func (application *MockGroupingsApplication) Campaigns() campaigns.Application {
	return application.CampaignsIns
}

func (application *MockGroupingsApplication) Topics() topics.Application {
	return application.TopicsIns
}

func (application *MockGroupingsApplication) Narratives() narratives.Application {
	return application.NarrativesIns
}

func (application *MockGroupingsApplication) Participations() participations.Application {
	return application.ParticipationsIns
}

func (application *MockGroupingsApplication) Clusters() clusters.Application {
	return application.ClustersIns
}
