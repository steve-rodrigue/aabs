package evidences

import (
	"context"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/participatables"
	domain_evidences "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/evidences"
)

// New creates a new participation evidence application
func New(
	repository domain_evidences.Repository,
) Application {
	return createApplication(repository)
}

// Application represents the participation evidence application
type Application interface {
	FindByID(
		ctx context.Context,
		id uuid.UUID,
	) (domain_evidences.Evidence, error)

	FindByParticipation(
		ctx context.Context,
		participation uuid.UUID,
	) ([]domain_evidences.Evidence, error)

	FindByPost(
		ctx context.Context,
		post uuid.UUID,
	) ([]domain_evidences.Evidence, error)

	FindByParticipant(
		ctx context.Context,
		participant participatables.Participatable,
	) ([]domain_evidences.Evidence, error)

	FindByTarget(
		ctx context.Context,
		target participatables.Participatable,
	) ([]domain_evidences.Evidence, error)
}
