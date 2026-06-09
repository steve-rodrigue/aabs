package communities

import (
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/platforms"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/users"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
)

type community struct {
	identifier uuid.UUID

	platform platforms.Platform

	handle string
	title  string
	text   string

	createdOn time.Time

	moderators []users.User
}

func (community *community) Identifier() uuid.UUID {
	return community.identifier
}

func (community *community) ParticipationKind() participatables.Kind {
	return participatables.CommunityKind
}

func (community *community) Platform() platforms.Platform {
	return community.platform
}

func (community *community) Handle() string {
	return community.handle
}

func (community *community) Title() string {
	return community.title
}

func (community *community) Text() string {
	return community.text
}

func (community *community) CreatedOn() time.Time {
	return community.createdOn
}

func (community *community) HasModerators() bool {
	return len(community.moderators) > 0
}

func (community *community) Moderators() []users.User {
	out := make([]users.User, len(community.moderators))
	copy(out, community.moderators)

	return out
}
