package users

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/platforms"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
)

var (
	ErrInvalidUserIdentifier  = errors.New("invalid user identifier")
	ErrInvalidUserPlatform    = errors.New("invalid user platform")
	ErrInvalidUserExternalID  = errors.New("invalid user external id")
	ErrInvalidUserHandle      = errors.New("invalid user handle")
	ErrInvalidUserDisplayName = errors.New("invalid user display name")
	ErrInvalidUserProfileURL  = errors.New("invalid user profile url")
	ErrInvalidUserCreatedOn   = errors.New("invalid user created on")
)

// NewAdapter creates a new user adapter
func NewAdapter() Adapter {
	return &adapter{}
}

// UserInput represents a user input
type UserInput struct {
	Identifier  uuid.UUID
	Platform    platforms.Platform
	ExternalID  string
	Handle      string
	DisplayName string
	ProfileURL  string
	CreatedOn   time.Time
}

// Adapter represents a user adapter
type Adapter interface {
	ToDomain(input UserInput) (User, error)
}

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
	Save(ctx context.Context, user User) error
	FindByID(ctx context.Context, id uuid.UUID) (User, error)
	FindByPlatformAndExternalID(
		ctx context.Context,
		platform platforms.Platform,
		externalID string,
	) (User, error)
	FindByPlatformAndHandle(
		ctx context.Context,
		platform platforms.Platform,
		handle string,
	) (User, error)
	Find(ctx context.Context, index int, amount int) ([]User, error)
	FindAfter(ctx context.Context, cursor uuid.UUID, amount int) ([]User, error)
	Count(ctx context.Context) (int64, error)
}
