package campaigns

import (
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
)

type campaign struct {
	identifier uuid.UUID

	name        string
	description string

	cluster clusters.Cluster

	postCount  int
	confidence float64

	createdOn time.Time
}

func (campaign *campaign) Identifier() uuid.UUID {
	return campaign.identifier
}

func (campaign *campaign) ParticipationKind() participatables.Kind {
	return participatables.CampaignKind
}

func (campaign *campaign) Name() string {
	return campaign.name
}

func (campaign *campaign) Description() string {
	return campaign.description
}

func (campaign *campaign) Cluster() clusters.Cluster {
	return campaign.cluster
}

func (campaign *campaign) PostCount() int {
	return campaign.postCount
}

func (campaign *campaign) Confidence() float64 {
	return campaign.confidence
}

func (campaign *campaign) CreatedOn() time.Time {
	return campaign.createdOn
}
