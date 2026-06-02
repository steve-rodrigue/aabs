package evidences

import (
	"time"

	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
	"github.com/steve-rodrigue/aabs/services/saas/domain/posts"
)

// Evidence represents a post that contributed to a participation score
type Evidence interface {
	Identifier() uuid.UUID

	Participation() participations.Participation

	Participant() participatables.Participatable
	Target() participatables.Participatable

	Post() posts.Post
	Score() float64

	DetectedOn() time.Time
}

// Repository represents a participation evidence repository
type Repository interface {
	Save(evidence Evidence) error
	FindByID(id uuid.UUID) (Evidence, error)
	FindByParticipation(participation uuid.UUID) ([]Evidence, error)
	FindByPost(post uuid.UUID) ([]Evidence, error)
	FindByParticipant(participant participatables.Participatable) ([]Evidence, error)
	FindByTarget(target participatables.Participatable) ([]Evidence, error)
}

// Calculator represents a participation evidence calculator
type Calculator interface {
	Calculate(
		participation participations.Participation,
	) ([]Evidence, error)
}
