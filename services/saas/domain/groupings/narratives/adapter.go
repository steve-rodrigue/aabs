package narratives

import (
	"strings"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
)

type adapter struct{}

func (adapter *adapter) ToDomain(
	input NarrativeInput,
) (Narrative, error) {
	if input.Identifier == uuid.Nil {
		return nil, ErrInvalidNarrativeIdentifier
	}

	if input.ParticipationKind != participatables.NarrativeKind {
		return nil, ErrInvalidNarrativeParticipationKind
	}

	if input.Cluster == nil {
		return nil, ErrInvalidNarrativeCluster
	}

	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, ErrInvalidNarrativeName
	}

	description := strings.TrimSpace(input.Description)
	if description == "" {
		return nil, ErrInvalidNarrativeDescription
	}

	if input.CreatedOn.IsZero() {
		return nil, ErrInvalidNarrativeCreatedOn
	}

	return &narrative{
		identifier:        input.Identifier,
		participationKind: input.ParticipationKind,
		cluster:           input.Cluster,
		name:              name,
		description:       description,
		createdOn:         input.CreatedOn.UTC(),
	}, nil
}
