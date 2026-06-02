package groupings

import (
	"github.com/steve-rodrigue/aabs/services/saas/applications/groupings/campaigns"
	"github.com/steve-rodrigue/aabs/services/saas/applications/groupings/clusters"
	"github.com/steve-rodrigue/aabs/services/saas/applications/groupings/narratives"
	"github.com/steve-rodrigue/aabs/services/saas/applications/groupings/participations"
	"github.com/steve-rodrigue/aabs/services/saas/applications/groupings/topics"
)

type application struct {
	campaigns      campaigns.Application
	topics         topics.Application
	narratives     narratives.Application
	participations participations.Application
	clusters       clusters.Application
}

func createApplication(
	campaigns campaigns.Application,
	topics topics.Application,
	narratives narratives.Application,
	participations participations.Application,
	clusters clusters.Application,
) Application {
	return &application{
		campaigns:      campaigns,
		topics:         topics,
		narratives:     narratives,
		participations: participations,
		clusters:       clusters,
	}
}

func (app *application) Campaigns() campaigns.Application {
	return app.campaigns
}

func (app *application) Topics() topics.Application {
	return app.topics
}

func (app *application) Narratives() narratives.Application {
	return app.narratives
}

func (app *application) Participations() participations.Application {
	return app.participations
}

func (app *application) Clusters() clusters.Application {
	return app.clusters
}
