package users

import (
	"time"

	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
	"github.com/steve-rodrigue/aabs/services/saas/domain/platforms"
)

// User represents a user
type User interface {
	Identifier() uuid.UUID
	ParticipationKind() participatables.Kind
	Platform() platforms.Platform
	ExternalID() string
	Handle() string
	DisplayName() string
	ProfileURL() string
	CreatedOn() time.Time
}

// Repository represents a user repository
type Repository interface {
	Save(user User) error
	FindByID(id uuid.UUID) (User, error)
	FindByPlatformAndExternalID(platform platforms.Platform, externalID string) (User, error)
	FindByPlatformAndHandle(platform platforms.Platform, handle string) (User, error)
}
