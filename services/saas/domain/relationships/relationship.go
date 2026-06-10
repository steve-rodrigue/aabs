package relationships

import (
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/relatables"
)

type relationship struct {
	identifier uuid.UUID

	source relatables.Relatable
	target relatables.Relatable

	similarity float64
	createdOn  time.Time
}

func (relationship *relationship) Identifier() uuid.UUID {
	return relationship.identifier
}

func (relationship *relationship) Source() relatables.Relatable {
	return relationship.source
}

func (relationship *relationship) Target() relatables.Relatable {
	return relationship.target
}

func (relationship *relationship) Similarity() float64 {
	return relationship.similarity
}

func (relationship *relationship) CreatedOn() time.Time {
	return relationship.createdOn
}
