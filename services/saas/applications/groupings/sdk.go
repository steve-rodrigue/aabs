package groupings

import (
	"github.com/steve-rodrigue/aabs/services/saas/applications/groupings/campaigns"
	"github.com/steve-rodrigue/aabs/services/saas/applications/groupings/clusters"
	"github.com/steve-rodrigue/aabs/services/saas/applications/groupings/narratives"
	"github.com/steve-rodrigue/aabs/services/saas/applications/groupings/participations"
	"github.com/steve-rodrigue/aabs/services/saas/applications/groupings/topics"
)

// Application represents the grouping application
type Application interface {
	Campaigns() campaigns.Application
	Topics() topics.Application
	Narratives() narratives.Application
	Participations() participations.Application
	Clusters() clusters.Application
}
