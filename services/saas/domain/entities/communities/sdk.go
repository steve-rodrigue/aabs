package communities

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/platforms"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/users"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
)

var (
	ErrInvalidCommunityIdentifier = errors.New("invalid community identifier")
	ErrInvalidCommunityPlatform   = errors.New("invalid community platform")
	ErrInvalidCommunityHandle     = errors.New("invalid community handle")
	ErrInvalidCommunityTitle      = errors.New("invalid community title")
	ErrInvalidCommunityText       = errors.New("invalid community text")
	ErrInvalidCommunityCreatedOn  = errors.New("invalid community created on")
	ErrInvalidCommunityModerator  = errors.New("invalid community moderator")
)

// NewAdapter creates a new community adapter
func NewAdapter() Adapter {
	return &adapter{}
}

// CommunityInput represents a community input
type CommunityInput struct {
	Identifier uuid.UUID
	Platform   platforms.Platform
	Handle     string
	Title      string
	Text       string
	CreatedOn  time.Time
	Moderators []users.User
}

// Adapter represents a community adapter
type Adapter interface {
	ToDomain(input CommunityInput) (Community, error)
}

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
	Save(
		ctx context.Context,
		community Community,
	) error

	FindByID(
		ctx context.Context,
		id uuid.UUID,
	) (Community, error)

	FindByHandle(
		ctx context.Context,
		platform platforms.Platform,
		handle string,
	) (Community, error)

	Find(
		ctx context.Context,
		index int,
		amount int,
	) ([]Community, error)

	FindAfter(
		ctx context.Context,
		cursor uuid.UUID,
		amount int,
	) ([]Community, error)

	FindByPlatform(
		ctx context.Context,
		platform platforms.Platform,
	) ([]Community, error)

	Count(
		ctx context.Context,
	) (int64, error)
}
