package relatables

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type Kind string

const (
	CampaignKind  Kind = "campaign"
	TopicKind     Kind = "topic"
	UserKind      Kind = "user"
	PostKind      Kind = "post"
	NarrativeKind Kind = "narrative"
)

var (
	ErrInvalidRelatableIdentifier       = errors.New("invalid relatable identifier")
	ErrInvalidRelatableRelationshipKind = errors.New("invalid relatable relationship kind")
)

// NewAdapter creates a new relatable adapter
func NewAdapter() Adapter {
	return &adapter{}
}

// RelatableInput represents a relatable input
type RelatableInput struct {
	Identifier       uuid.UUID
	RelationshipKind Kind
}

// Adapter represents a relatable adapter
type Adapter interface {
	ToDomain(input RelatableInput) (Relatable, error)
}

// Relatable represents a relatable
type Relatable interface {
	Identifier() uuid.UUID
	RelationshipKind() Kind
}

// Repository represents a relatable repository
type Repository interface {
	Save(ctx context.Context, relatable Relatable) error
	Delete(ctx context.Context, relatable Relatable) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
	Find(ctx context.Context, index int, amount int) ([]Relatable, error)
	FindAfter(ctx context.Context, cursor uuid.UUID, amount int) ([]Relatable, error)
	Count(ctx context.Context) (int64, error)
	FindByKind(ctx context.Context, kind Kind, index int, amount int) ([]Relatable, error)
	CountByKind(ctx context.Context, kind Kind) (int64, error)
}

// CandidateRepository represents a candidate repository
type CandidateRepository interface {
	FindCandidates(
		ctx context.Context,
		source Relatable,
		amount int,
	) ([]Relatable, error)
}
