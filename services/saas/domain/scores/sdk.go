package scores

import (
	"time"

	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/scorables"
	"github.com/steve-rodrigue/aabs/services/saas/domain/scores/factors"
)

type Type string

const (
	TrustType              Type = "trust"
	SpamType               Type = "spam"
	BotLikelihoodType      Type = "bot_likelihood"
	CoordinationType       Type = "coordination"
	AuthenticityType       Type = "authenticity"
	NarrativeInfluenceType Type = "narrative_influence"
	PolarizationType       Type = "polarization"
	ManipulationType       Type = "manipulation"
)

// ScoreInput represents a score input
type ScoreInput struct {
	Identifier   uuid.UUID
	Type         Type
	Target       scorables.Scorable
	Value        float64
	Confidence   float64
	Factors      []factors.Factor
	CalculatedOn time.Time
}

// Adapter represents a score adapter
type Adapter interface {
	ToDomain(input ScoreInput) (Score, error)
}

// Score represents the trust score
type Score interface {
	Identifier() uuid.UUID
	Type() Type
	Target() scorables.Scorable
	Value() float64
	Confidence() float64
	Factors() []factors.Factor
	CalculatedOn() time.Time
}

// Repository represents a trust score repository
type Repository interface {
	Save(score Score) error
	FindLatestByTarget(target scorables.Scorable, scoreType Type) (Score, error)
	FindHistoryByTarget(target scorables.Scorable, scoreType Type) ([]Score, error)
	FindLatestByTargetAndTypes(target scorables.Scorable, scoreTypes []Type) ([]Score, error)
}

// Calculator represents a trust score calculator
type Calculator interface {
	Type() Type
	Calculate(target scorables.Scorable) (Score, error)
}
