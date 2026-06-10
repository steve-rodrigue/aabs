package narratives

import (
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/participatables"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters"
)

type narrative struct {
	identifier uuid.UUID

	participationKind participatables.Kind

	cluster clusters.Cluster

	name        string
	description string

	createdOn time.Time
}

func (narrative *narrative) Identifier() uuid.UUID {
	return narrative.identifier
}

func (narrative *narrative) ParticipationKind() participatables.Kind {
	return narrative.participationKind
}

func (narrative *narrative) Cluster() clusters.Cluster {
	return narrative.cluster
}

func (narrative *narrative) Name() string {
	return narrative.name
}

func (narrative *narrative) Description() string {
	return narrative.description
}

func (narrative *narrative) CreatedOn() time.Time {
	return narrative.createdOn
}
