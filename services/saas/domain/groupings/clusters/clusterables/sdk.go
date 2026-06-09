package clusterables

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type Kind string

const (
	PostKind      Kind = "post"
	UserKind      Kind = "user"
	CommunityKind Kind = "community"
	PlatformKind  Kind = "platform"
	CampaignKind  Kind = "campaign"
	TopicKind     Kind = "topic"
	NarrativeKind Kind = "narrative"
)

var (
	ErrInvalidClusterableIdentifier = errors.New("invalid clusterable identifier")
	ErrInvalidClusterableKind       = errors.New("invalid clusterable kind")
	ErrInvalidComparableVector      = errors.New("invalid comparable vector")
)

// NewAdapter creates a new clusterable adapter
func NewAdapter() Adapter {
	return &adapter{}
}

// NewComparableAdapter creates a new comparable adapter
func NewComparableAdapter(
	clusterableAdapter Adapter,
) ComparableAdapter {
	return &comparableAdapter{
		clusterableAdapter: clusterableAdapter,
	}
}

// NewCandidateRepository creates a new clusterable candidate repository
func NewCandidateRepository(
	comparableRepository ComparableRepository,
) CandidateRepository {
	return &candidateRepository{
		comparableRepository: comparableRepository,
	}
}

// ClusterableInput represents a clusterable input
type ClusterableInput struct {
	Identifier  uuid.UUID
	ClusterKind Kind
}

// ComparableInput represents a comparable input
type ComparableInput struct {
	Clusterable ClusterableInput
	Vector      []float32
}

// Adapter represents a clusterable adapter
type Adapter interface {
	ToDomain(input ClusterableInput) (Clusterable, error)
}

// ComparableAdapter represents a comparable adapter
type ComparableAdapter interface {
	ToDomain(input ComparableInput) (Comparable, error)
}

// Clusterable represents a clusterable entity
type Clusterable interface {
	Identifier() uuid.UUID
	ClusterKind() Kind
}

// Comparable represents a comparable clusterable
type Comparable interface {
	Clusterable
	Vector() []float32
}

// Repository represents a clusterable repository
type Repository interface {
	FindByKind(
		ctx context.Context,
		kind Kind,
		index int,
		amount int,
	) ([]Clusterable, error)

	FindByKindAfter(
		ctx context.Context,
		kind Kind,
		cursor uuid.UUID,
		amount int,
	) ([]Clusterable, error)

	CountByKind(
		ctx context.Context,
		kind Kind,
	) (int64, error)
}

// CandidateRepository represents a repository that finds likely clustering candidates
type CandidateRepository interface {
	FindCandidates(
		ctx context.Context,
		target Clusterable,
		kind Kind,
		amount int,
	) ([]Clusterable, error)
}

// Repository represents a comparable repository
type ComparableRepository interface {
	FindByID(
		ctx context.Context,
		id uuid.UUID,
	) (Comparable, error)

	FindByKind(
		ctx context.Context,
		kind Kind,
		index int,
		amount int,
	) ([]Comparable, error)

	FindNearest(
		ctx context.Context,
		target Comparable,
		kind Kind,
		amount int,
	) ([]Comparable, error)
}
