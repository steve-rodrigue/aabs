package clusterables

import "github.com/google/uuid"

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

// ClusterableInput represents a clusterable input
type ClusterableInput struct {
	Identifier  uuid.UUID
	ClusterKind Kind
}

// Adapter represents a clusterable adapter
type Adapter interface {
	ToDomain(input ClusterableInput) (Clusterable, error)
}

// Clusterable represents a clusterable entity
type Clusterable interface {
	Identifier() uuid.UUID
	ClusterKind() Kind
}

// Repository represents a clusterable repository
type Repository interface {
	FindByKind(kind Kind, index int, amount int) ([]Clusterable, error)
	FindByKindAfter(kind Kind, cursor uuid.UUID, amount int) ([]Clusterable, error)
	CountByKind(kind Kind) (int64, error)
}

// CandidateRepository represents a repository that finds likely clustering candidates.
type CandidateRepository interface {
	FindCandidates(target Clusterable, kind Kind, amount int) ([]Clusterable, error)
}
