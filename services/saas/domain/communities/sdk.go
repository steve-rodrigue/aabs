package communities

import (
	"time"

	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
	"github.com/steve-rodrigue/aabs/services/saas/domain/platforms"
	"github.com/steve-rodrigue/aabs/services/saas/domain/users"
)

// Community represents a community
type Community interface {
	Identifier() uuid.UUID
	ParticipationKind() participatables.Kind
	Platform() platforms.Platform
	Handle() string
	Title() string
	Text() string
	CreatedOn() time.Time
	HasModerators() bool
	Moderators() []users.User
}

// Repository represents a community repository
type Repository interface {
	Save(community Community) error
	FindByID(id uuid.UUID) (Community, error)
	FindByHandle(handle string) (Community, error)
}
