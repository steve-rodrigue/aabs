package platforms

import (
	"time"

	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
)

// Platform represents a website or app where users publish content
type Platform interface {
	Identifier() uuid.UUID
	ParticipationKind() participatables.Kind
	Name() string
	Handle() string
	BaseURL() string
	CreatedOn() time.Time
}

// Repository represents a platform repository
type Repository interface {
	Save(platform Platform) error
	FindByID(id uuid.UUID) (Platform, error)
	FindByHandle(handle string) (Platform, error)
	FindByName(name string) (Platform, error)
}
