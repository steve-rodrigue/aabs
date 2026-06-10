package participations

import (
	"context"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/participatables"
	domain_participations "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations"
	domain_evidences "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/evidences"
)

// New creates a new participation application
func New(
	repository domain_participations.Repository,
	calculator domain_participations.Calculator,
	participatableRepository participatables.Repository,
	evidenceRepository domain_evidences.Repository,
	evidenceCalculator domain_evidences.Calculator,
) Application {
	return createApplication(
		repository,
		calculator,
		participatableRepository,
		evidenceRepository,
		evidenceCalculator,
	)
}

// Application represents the participation application
type Application interface {
	FindByID(
		ctx context.Context,
		id uuid.UUID,
	) (domain_participations.Participation, error)

	FindByParticipant(
		ctx context.Context,
		participant participatables.Participatable,
	) ([]domain_participations.Participation, error)

	FindByTarget(
		ctx context.Context,
		target participatables.Participatable,
	) ([]domain_participations.Participation, error)

	FindBetween(
		ctx context.Context,
		participant participatables.Participatable,
		target participatables.Participatable,
	) (domain_participations.Participation, error)

	FindEvidencesByParticipation(
		ctx context.Context,
		participation uuid.UUID,
	) ([]domain_evidences.Evidence, error)

	FindEvidencesByPost(
		ctx context.Context,
		post uuid.UUID,
	) ([]domain_evidences.Evidence, error)

	FindEvidencesByParticipant(
		ctx context.Context,
		participant participatables.Participatable,
	) ([]domain_evidences.Evidence, error)

	FindEvidencesByTarget(
		ctx context.Context,
		target participatables.Participatable,
	) ([]domain_evidences.Evidence, error)

	RebuildParticipations(ctx context.Context) error
}
