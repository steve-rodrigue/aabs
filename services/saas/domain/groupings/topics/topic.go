package topics

import (
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
)

type topic struct {
	identifier uuid.UUID

	cluster clusters.Cluster

	name        string
	description string

	parent Topic

	createdOn time.Time
}

func (topic *topic) Identifier() uuid.UUID {
	return topic.identifier
}

func (topic *topic) ParticipationKind() participatables.Kind {
	return participatables.TopicKind
}

func (topic *topic) Cluster() clusters.Cluster {
	return topic.cluster
}

func (topic *topic) Name() string {
	return topic.name
}

func (topic *topic) Description() string {
	return topic.description
}

func (topic *topic) CreatedOn() time.Time {
	return topic.createdOn
}

func (topic *topic) HasParent() bool {
	return topic.parent != nil
}

func (topic *topic) Parent() Topic {
	return topic.parent
}
