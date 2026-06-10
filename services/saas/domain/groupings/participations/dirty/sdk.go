package dirty

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/participatables"
)

var (
	// adapter
	ErrInvalidDirtyIdentifier  = errors.New("invalid dirty participation identifier")
	ErrInvalidDirtyParticipant = errors.New("invalid dirty participation participant")
	ErrInvalidDirtyTarget      = errors.New("invalid dirty participation target")
	ErrInvalidDirtyMarkedOn    = errors.New("invalid dirty participation marked on")
)

// NewAdapter creates a new dirty participation adapter
func NewAdapter() Adapter {
	return &adapter{}
}

// DirtyInput represents a dirty participation input
type DirtyInput struct {
	Identifier uuid.UUID

	Participant participatables.Participatable
	Target      participatables.Participatable

	MarkedOn time.Time
}

// Adapter represents a dirty participation adapter
type Adapter interface {
	ToDomain(
		input DirtyInput,
	) (Dirty, error)
}

// Dirty represents a participation that must be rebuilt
type Dirty interface {
	Identifier() uuid.UUID

	Participant() participatables.Participatable
	Target() participatables.Participatable

	MarkedOn() time.Time
}

// Repository represents a dirty participation repository
type Repository interface {
	Save(
		ctx context.Context,
		dirty Dirty,
	) error

	Delete(
		ctx context.Context,
		id uuid.UUID,
	) error

	FindByID(
		ctx context.Context,
		id uuid.UUID,
	) (Dirty, error)

	FindBetween(
		ctx context.Context,
		participant participatables.Participatable,
		target participatables.Participatable,
	) (Dirty, error)

	Find(
		ctx context.Context,
		index int,
		amount int,
	) ([]Dirty, error)

	FindAfter(
		ctx context.Context,
		cursor uuid.UUID,
		amount int,
	) ([]Dirty, error)

	Count(
		ctx context.Context,
	) (int64, error)
}
