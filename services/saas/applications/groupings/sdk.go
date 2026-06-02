package groupings

import (
	"github.com/steve-rodrigue/aabs/services/saas/applications/groupings/campaigns"
	"github.com/steve-rodrigue/aabs/services/saas/applications/groupings/clusters"
	"github.com/steve-rodrigue/aabs/services/saas/applications/groupings/narratives"
	"github.com/steve-rodrigue/aabs/services/saas/applications/groupings/participations"
	"github.com/steve-rodrigue/aabs/services/saas/applications/groupings/topics"
)

// New creates a new grouping application
func New(
	campaigns campaigns.Application,
	topics topics.Application,
	narratives narratives.Application,
	participations participations.Application,
	clusters clusters.Application,
) Application {
	return createApplication(
		campaigns,
		topics,
		narratives,
		participations,
		clusters,
	)
}

// Application represents the grouping application
type Application interface {
	Campaigns() campaigns.Application
	Topics() topics.Application
	Narratives() narratives.Application
	Participations() participations.Application
	Clusters() clusters.Application
}
