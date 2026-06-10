package clusters

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/clusterables"
)

func TestCluster(t *testing.T) {
	id := uuid.New()
	target := clusterables.NewMockClusterable(clusterables.PostKind)
	memberIDs := []uuid.UUID{
		uuid.New(),
		uuid.New(),
	}
	centroid := []float32{
		0.1,
		0.2,
	}
	createdOn := time.Now().UTC()

	cluster := &cluster{
		identifier:      id,
		target:          target,
		memberIDs:       memberIDs,
		memberKind:      clusterables.PostKind,
		confidenceScore: 0.75,
		centroid:        centroid,
		createdOn:       createdOn,
	}

	if cluster.Identifier() != id {
		t.Fatalf("expected identifier %s, got %s", id, cluster.Identifier())
	}

	if cluster.Target() != target {
		t.Fatalf("expected target")
	}

	resultMemberIDs := cluster.MemberIDs()

	if len(resultMemberIDs) != len(memberIDs) {
		t.Fatalf("expected %d member ids, got %d", len(memberIDs), len(resultMemberIDs))
	}

	for index := range memberIDs {
		if resultMemberIDs[index] != memberIDs[index] {
			t.Fatalf(
				"expected member id[%d] %s, got %s",
				index,
				memberIDs[index],
				resultMemberIDs[index],
			)
		}
	}

	if cluster.MemberKind() != clusterables.PostKind {
		t.Fatalf(
			"expected member kind %s, got %s",
			clusterables.PostKind,
			cluster.MemberKind(),
		)
	}

	if cluster.ConfidenceScore() != 0.75 {
		t.Fatalf("expected confidence score 0.75, got %f", cluster.ConfidenceScore())
	}

	resultCentroid := cluster.Centroid()

	if len(resultCentroid) != len(centroid) {
		t.Fatalf("expected centroid length %d, got %d", len(centroid), len(resultCentroid))
	}

	for index := range centroid {
		if resultCentroid[index] != centroid[index] {
			t.Fatalf(
				"expected centroid[%d] %f, got %f",
				index,
				centroid[index],
				resultCentroid[index],
			)
		}
	}

	if !cluster.CreatedOn().Equal(createdOn) {
		t.Fatalf("expected created on %s, got %s", createdOn, cluster.CreatedOn())
	}
}

func TestClusterMemberIDsReturnsCopy(t *testing.T) {
	cluster := &cluster{
		memberIDs: []uuid.UUID{
			uuid.New(),
		},
	}

	memberIDs := cluster.MemberIDs()
	memberIDs[0] = uuid.New()

	again := cluster.MemberIDs()

	if again[0] == memberIDs[0] {
		t.Fatalf("expected member ids copy")
	}
}

func TestClusterCentroidReturnsCopy(t *testing.T) {
	cluster := &cluster{
		centroid: []float32{
			0.1,
		},
	}

	centroid := cluster.Centroid()
	centroid[0] = 99

	again := cluster.Centroid()

	if again[0] == 99 {
		t.Fatalf("expected centroid copy")
	}
}
