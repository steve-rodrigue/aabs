package campaigns

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/clusterables"
	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/participatables"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters"
)

func TestCampaign(t *testing.T) {
	id := uuid.New()
	createdOn := time.Now().UTC()

	cluster := clusters.NewMockCluster(
		clusterables.NewMockClusterable(clusterables.PostKind),
		clusterables.PostKind,
		[]uuid.UUID{uuid.New()},
	)

	campaign := &campaign{
		identifier:  id,
		name:        "Campaign",
		description: "Description",
		cluster:     cluster,
		postCount:   10,
		confidence:  0.8,
		createdOn:   createdOn,
	}

	if campaign.Identifier() != id {
		t.Fatalf("expected identifier")
	}

	if campaign.ParticipationKind() != participatables.CampaignKind {
		t.Fatalf("expected campaign participation kind")
	}

	if campaign.Name() != "Campaign" {
		t.Fatalf("expected name")
	}

	if campaign.Description() != "Description" {
		t.Fatalf("expected description")
	}

	if campaign.Cluster() != cluster {
		t.Fatalf("expected cluster")
	}

	if campaign.PostCount() != 10 {
		t.Fatalf("expected post count 10")
	}

	if campaign.Confidence() != 0.8 {
		t.Fatalf("expected confidence 0.8")
	}

	if !campaign.CreatedOn().Equal(createdOn) {
		t.Fatalf("expected created on")
	}
}
