package relatables

import "github.com/google/uuid"

type Kind string

const (
	CampaignKind  Kind = "campaign"
	TopicKind     Kind = "topic"
	UserKind      Kind = "user"
	PostKind      Kind = "post"
	NarrativeKind Kind = "narrative"
)

type Relatable interface {
	Identifier() uuid.UUID
	RelationshipKind() Kind
}

type Repository interface {
	Find(index int, amount int) ([]Relatable, error)
	FindAfter(cursor uuid.UUID, amount int) ([]Relatable, error)
	Count() (int64, error)
}

type CandidateRepository interface {
	FindCandidates(source Relatable, amount int) ([]Relatable, error)
}
