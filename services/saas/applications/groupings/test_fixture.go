package groupings

import (
	"github.com/steve-rodrigue/aabs/services/saas/applications/groupings/campaigns"
	"github.com/steve-rodrigue/aabs/services/saas/applications/groupings/clusters"
	"github.com/steve-rodrigue/aabs/services/saas/applications/groupings/narratives"
	"github.com/steve-rodrigue/aabs/services/saas/applications/groupings/participations"
	"github.com/steve-rodrigue/aabs/services/saas/applications/groupings/topics"
)

type applicationFixture struct {
	application Application

	campaigns      *campaigns.MockCampaignsApplication
	topics         *topics.MockTopicsApplication
	narratives     *narratives.MockNarrativesApplication
	participations *participations.MockParticipationsApplication
	clusters       *clusters.MockClustersApplication
}

func newApplicationFixture() *applicationFixture {
	campaigns := campaigns.NewMockCampaignsApplication()
	topics := topics.NewMockTopicsApplication()
	narratives := narratives.NewMockNarrativesApplication()
	participations := participations.NewMockParticipationsApplication()
	clusters := clusters.NewMockClustersApplication()

	application := New(
		campaigns,
		topics,
		narratives,
		participations,
		clusters,
	)

	return &applicationFixture{
		application:    application,
		campaigns:      campaigns,
		topics:         topics,
		narratives:     narratives,
		participations: participations,
		clusters:       clusters,
	}
}
