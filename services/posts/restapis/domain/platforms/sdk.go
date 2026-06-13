package platforms

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidPlatformIdentifier        = errors.New("invalid platform identifier")
	ErrInvalidPlatformParticipationKind = errors.New("invalid platform participation kind")
	ErrInvalidPlatformName              = errors.New("invalid platform name")
	ErrInvalidPlatformHandle            = errors.New("invalid platform handle")
	ErrInvalidPlatformBaseURL           = errors.New("invalid platform base url")
	ErrInvalidPlatformCreatedOn         = errors.New("invalid platform created on")
)

// NewAdapter creates a new platform adapter
func NewAdapter() Adapter {
	return &adapter{}
}

// PlatformInput represents a platform input
type PlatformInput struct {
	Identifier uuid.UUID
	Name       string
	Handle     string
	BaseURL    string
	CreatedOn  time.Time
}

// Adapter represents a user adapter
type Adapter interface {
	ToDomain(input PlatformInput) (Platform, error)
}

// Platform represents a website or app where users publish content
type Platform interface {
	Identifier() uuid.UUID
	Name() string
	Handle() string
	BaseURL() string
	CreatedOn() time.Time
}

// Repository represents a platform repository
type Repository interface {
	Save(ctx context.Context, platform Platform) error

	FindByID(ctx context.Context, id uuid.UUID) (Platform, error)
	FindByHandle(ctx context.Context, handle string) (Platform, error)
	FindByName(ctx context.Context, name string) (Platform, error)

	Find(ctx context.Context, index int, amount int) ([]Platform, error)
	FindAfter(ctx context.Context, cursor uuid.UUID, amount int) ([]Platform, error)

	Count(ctx context.Context) (int64, error)
}
