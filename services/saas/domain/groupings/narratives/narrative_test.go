package narratives

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/clusterables"
	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/participatables"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters"
)

func TestNarrative(t *testing.T) {
	id := uuid.New()

	cluster := clusters.NewMockCluster(
		clusterables.NewMockClusterable(clusterables.NarrativeKind),
		clusterables.PostKind,
		[]uuid.UUID{uuid.New()},
	)

	createdOn := time.Now().UTC()

	narrative := &narrative{
		identifier:        id,
		participationKind: participatables.NarrativeKind,
		cluster:           cluster,
		name:              "Election Integrity",
		description:       "Discussion about election integrity.",
		createdOn:         createdOn,
	}

	if narrative.Identifier() != id {
		t.Fatalf("expected identifier %s, got %s", id, narrative.Identifier())
	}

	if narrative.ParticipationKind() != participatables.NarrativeKind {
		t.Fatalf(
			"expected participation kind %s, got %s",
			participatables.NarrativeKind,
			narrative.ParticipationKind(),
		)
	}

	if narrative.Cluster() != cluster {
		t.Fatalf("expected cluster")
	}

	if narrative.Name() != "Election Integrity" {
		t.Fatalf("expected name Election Integrity, got %s", narrative.Name())
	}

	if narrative.Description() != "Discussion about election integrity." {
		t.Fatalf(
			"expected description, got %s",
			narrative.Description(),
		)
	}

	if !narrative.CreatedOn().Equal(createdOn) {
		t.Fatalf(
			"expected created on %s, got %s",
			createdOn,
			narrative.CreatedOn(),
		)
	}
}
