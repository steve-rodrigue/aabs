package clusterables

import (
	"testing"

	"github.com/google/uuid"
)

func TestClusterable(t *testing.T) {
	id := uuid.New()

	clusterable := &clusterable{
		identifier:  id,
		clusterKind: PostKind,
	}

	if clusterable.Identifier() != id {
		t.Fatalf("expected identifier %s, got %s", id, clusterable.Identifier())
	}

	if clusterable.ClusterKind() != PostKind {
		t.Fatalf(
			"expected cluster kind %s, got %s",
			PostKind,
			clusterable.ClusterKind(),
		)
	}
}
