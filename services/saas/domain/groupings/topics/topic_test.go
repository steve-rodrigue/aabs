package topics

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters/clusterables"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
)

func TestTopic(t *testing.T) {
	id := uuid.New()
	cluster := clusters.NewMockCluster(
		clusterables.NewMockClusterable(clusterables.TopicKind),
		clusterables.PostKind,
		[]uuid.UUID{uuid.New()},
	)
	parent := &topic{
		identifier:  uuid.New(),
		cluster:     cluster,
		name:        "Parent Topic",
		description: "Parent description",
		createdOn:   time.Now().UTC(),
	}
	createdOn := time.Now().UTC()

	topic := &topic{
		identifier:  id,
		cluster:     cluster,
		name:        "AI Spam",
		description: "Posts about AI spam",
		parent:      parent,
		createdOn:   createdOn,
	}

	if topic.Identifier() != id {
		t.Fatalf("expected identifier %s, got %s", id, topic.Identifier())
	}

	if topic.ParticipationKind() != participatables.TopicKind {
		t.Fatalf(
			"expected participation kind %s, got %s",
			participatables.TopicKind,
			topic.ParticipationKind(),
		)
	}

	if topic.Cluster() != cluster {
		t.Fatalf("expected cluster")
	}

	if topic.Name() != "AI Spam" {
		t.Fatalf("expected name AI Spam, got %s", topic.Name())
	}

	if topic.Description() != "Posts about AI spam" {
		t.Fatalf("expected description, got %s", topic.Description())
	}

	if !topic.CreatedOn().Equal(createdOn) {
		t.Fatalf("expected created on %s, got %s", createdOn, topic.CreatedOn())
	}

	if !topic.HasParent() {
		t.Fatalf("expected topic to have parent")
	}

	if topic.Parent() != parent {
		t.Fatalf("expected parent")
	}
}

func TestTopicWithoutParent(t *testing.T) {
	topic := &topic{
		identifier: uuid.New(),
		cluster: clusters.NewMockCluster(
			clusterables.NewMockClusterable(clusterables.TopicKind),
			clusterables.PostKind,
			[]uuid.UUID{uuid.New()},
		),
		name:      "Root Topic",
		createdOn: time.Now().UTC(),
	}

	if topic.HasParent() {
		t.Fatalf("expected topic to not have parent")
	}

	if topic.Parent() != nil {
		t.Fatalf("expected nil parent")
	}
}
