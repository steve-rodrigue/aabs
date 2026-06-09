package evidences

import (
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
)

type evidence struct {
	identifier uuid.UUID

	participation participations.Participation

	participant participatables.Participatable
	target      participatables.Participatable

	post  posts.Post
	score float64

	detectedOn time.Time
}

func (evidence *evidence) Identifier() uuid.UUID {
	return evidence.identifier
}

func (evidence *evidence) Participation() participations.Participation {
	return evidence.participation
}

func (evidence *evidence) Participant() participatables.Participatable {
	return evidence.participant
}

func (evidence *evidence) Target() participatables.Participatable {
	return evidence.target
}

func (evidence *evidence) Post() posts.Post {
	return evidence.post
}

func (evidence *evidence) Score() float64 {
	return evidence.score
}

func (evidence *evidence) DetectedOn() time.Time {
	return evidence.detectedOn
}
