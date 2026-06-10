package postgresql

import (
	"context"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/clusterables"
	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/participatables"
	domain_clusters "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters"
	domain_topics "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/topics"
)

func TestNewGroupingsTopicRepository(t *testing.T) {
	fixture := newGroupingsTopicRepositoryFixture(t)

	repository := NewGroupingsTopicRepository(
		fixture.pool,
		fixture.adapter,
		fixture.clusters,
	)

	if repository == nil {
		t.Fatalf("expected repository")
	}
}

func TestGroupingsTopicRepositorySaveAndFindByID(t *testing.T) {
	fixture := newGroupingsTopicRepositoryFixture(t)

	topic := newTestTopic(
		t,
		fixture,
		"Electric Vehicles",
		"Posts about electric vehicles",
		nil,
	)

	if err := fixture.repository.Save(fixture.ctx, topic); err != nil {
		t.Fatal(err)
	}

	result, err := fixture.repository.FindByID(
		fixture.ctx,
		topic.Identifier(),
	)
	if err != nil {
		t.Fatal(err)
	}

	assertTopic(t, result, topic)
}

func TestGroupingsTopicRepositoryFindByIDReturnsNilWhenNotFound(t *testing.T) {
	fixture := newGroupingsTopicRepositoryFixture(t)

	result, err := fixture.repository.FindByID(
		fixture.ctx,
		uuid.New(),
	)

	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Fatalf("expected nil topic")
	}
}

func TestGroupingsTopicRepositorySaveUpdatesExistingTopic(t *testing.T) {
	fixture := newGroupingsTopicRepositoryFixture(t)

	topic := newTestTopic(
		t,
		fixture,
		"Electric Vehicles",
		"Old description",
		nil,
	)

	if err := fixture.repository.Save(fixture.ctx, topic); err != nil {
		t.Fatal(err)
	}

	updated := &testTopic{
		id:          topic.Identifier(),
		cluster:     topic.Cluster(),
		name:        "Electric Vehicles Updated",
		description: "Updated description",
		createdOn:   topic.CreatedOn(),
	}

	if err := fixture.repository.Save(fixture.ctx, updated); err != nil {
		t.Fatal(err)
	}

	result, err := fixture.repository.FindByID(
		fixture.ctx,
		topic.Identifier(),
	)
	if err != nil {
		t.Fatal(err)
	}

	assertTopic(t, result, updated)
}

func TestGroupingsTopicRepositoryFindByName(t *testing.T) {
	fixture := newGroupingsTopicRepositoryFixture(t)

	topic := newTestTopic(
		t,
		fixture,
		"Electric Vehicles",
		"Posts about electric vehicles",
		nil,
	)

	if err := fixture.repository.Save(fixture.ctx, topic); err != nil {
		t.Fatal(err)
	}

	result, err := fixture.repository.FindByName(
		fixture.ctx,
		"Electric Vehicles",
	)
	if err != nil {
		t.Fatal(err)
	}

	assertTopic(t, result, topic)
}

func TestGroupingsTopicRepositoryFindByNameReturnsNilWhenNotFound(t *testing.T) {
	fixture := newGroupingsTopicRepositoryFixture(t)

	result, err := fixture.repository.FindByName(
		fixture.ctx,
		"Unknown",
	)

	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Fatalf("expected nil topic")
	}
}

func TestGroupingsTopicRepositoryFind(t *testing.T) {
	fixture := newGroupingsTopicRepositoryFixture(t)

	topics := []domain_topics.Topic{
		newTestTopic(t, fixture, "A", "A description", nil),
		newTestTopic(t, fixture, "B", "B description", nil),
		newTestTopic(t, fixture, "C", "C description", nil),
	}

	sortTopics(topics)
	saveTopics(t, fixture, topics...)

	result, err := fixture.repository.Find(
		fixture.ctx,
		1,
		1,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 topic, got %d", len(result))
	}

	assertTopic(t, result[0], topics[1])
}

func TestGroupingsTopicRepositoryFindAfterWithNilCursor(t *testing.T) {
	fixture := newGroupingsTopicRepositoryFixture(t)

	topic := newTestTopic(
		t,
		fixture,
		"Electric Vehicles",
		"Posts about electric vehicles",
		nil,
	)

	if err := fixture.repository.Save(fixture.ctx, topic); err != nil {
		t.Fatal(err)
	}

	result, err := fixture.repository.FindAfter(
		fixture.ctx,
		uuid.Nil,
		10,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 topic, got %d", len(result))
	}

	assertTopic(t, result[0], topic)
}

func TestGroupingsTopicRepositoryFindAfter(t *testing.T) {
	fixture := newGroupingsTopicRepositoryFixture(t)

	topics := []domain_topics.Topic{
		newTestTopic(t, fixture, "A", "A description", nil),
		newTestTopic(t, fixture, "B", "B description", nil),
		newTestTopic(t, fixture, "C", "C description", nil),
	}

	sortTopics(topics)
	saveTopics(t, fixture, topics...)

	result, err := fixture.repository.FindAfter(
		fixture.ctx,
		topics[0].Identifier(),
		10,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 topics, got %d", len(result))
	}

	assertTopic(t, result[0], topics[1])
	assertTopic(t, result[1], topics[2])
}

func TestGroupingsTopicRepositoryCount(t *testing.T) {
	fixture := newGroupingsTopicRepositoryFixture(t)

	saveTopics(
		t,
		fixture,
		newTestTopic(t, fixture, "A", "A description", nil),
		newTestTopic(t, fixture, "B", "B description", nil),
	)

	count, err := fixture.repository.Count(fixture.ctx)
	if err != nil {
		t.Fatal(err)
	}

	if count != 2 {
		t.Fatalf("expected count 2, got %d", count)
	}
}

func TestGroupingsTopicRepositoryFindChildren(t *testing.T) {
	fixture := newGroupingsTopicRepositoryFixture(t)

	parent := newTestTopic(
		t,
		fixture,
		"Technology",
		"Technology posts",
		nil,
	)

	first := newTestTopic(
		t,
		fixture,
		"Electric Vehicles",
		"EV posts",
		parent,
	)

	second := newTestTopic(
		t,
		fixture,
		"Artificial Intelligence",
		"AI posts",
		parent,
	)

	other := newTestTopic(
		t,
		fixture,
		"Politics",
		"Politics posts",
		nil,
	)

	saveTopics(t, fixture, parent, first, second, other)

	result, err := fixture.repository.FindChildren(
		fixture.ctx,
		parent.Identifier(),
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 children, got %d", len(result))
	}

	for _, topic := range result {
		if !topic.HasParent() {
			t.Fatalf("expected child to have parent")
		}

		if topic.Parent().Identifier() != parent.Identifier() {
			t.Fatalf(
				"expected parent %s, got %s",
				parent.Identifier(),
				topic.Parent().Identifier(),
			)
		}
	}
}

func TestGroupingsTopicRepositoryFindRoots(t *testing.T) {
	fixture := newGroupingsTopicRepositoryFixture(t)

	parent := newTestTopic(
		t,
		fixture,
		"Technology",
		"Technology posts",
		nil,
	)

	child := newTestTopic(
		t,
		fixture,
		"Electric Vehicles",
		"EV posts",
		parent,
	)

	root := newTestTopic(
		t,
		fixture,
		"Politics",
		"Politics posts",
		nil,
	)

	saveTopics(t, fixture, parent, child, root)

	result, err := fixture.repository.FindRoots(
		fixture.ctx,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 root topics, got %d", len(result))
	}

	for _, topic := range result {
		if topic.HasParent() {
			t.Fatalf("expected root topic")
		}
	}
}

type groupingsTopicRepositoryFixture struct {
	ctx        context.Context
	pool       *pgxpool.Pool
	adapter    domain_topics.Adapter
	clusters   *domain_clusters.MockClusterRepository
	repository domain_topics.Repository
}

func newGroupingsTopicRepositoryFixture(
	t *testing.T,
) *groupingsTopicRepositoryFixture {
	t.Helper()

	dsn := os.Getenv("POSTGRES_TEST_DSN")
	if dsn == "" {
		t.Skip("POSTGRES_TEST_DSN is not set")
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}

	adapter := domain_topics.NewAdapter()
	clusterRepository := domain_clusters.NewMockClusterRepository()

	fixture := &groupingsTopicRepositoryFixture{
		ctx:      ctx,
		pool:     pool,
		adapter:  adapter,
		clusters: clusterRepository,
	}

	createGroupingsTopicsTable(t, fixture)
	truncateGroupingsTopicsTable(t, fixture)

	fixture.repository = NewGroupingsTopicRepository(
		pool,
		adapter,
		clusterRepository,
	)

	t.Cleanup(func() {
		truncateGroupingsTopicsTable(t, fixture)
		pool.Close()
	})

	return fixture
}

func createGroupingsTopicsTable(
	t *testing.T,
	fixture *groupingsTopicRepositoryFixture,
) {
	t.Helper()

	_, err := fixture.pool.Exec(
		fixture.ctx,
		`
		CREATE TABLE IF NOT EXISTS groupings_topics (
			identifier UUID PRIMARY KEY,
			cluster_id UUID NOT NULL,
			name TEXT NOT NULL,
			description TEXT NOT NULL DEFAULT '',
			parent_id UUID NULL,
			created_on TIMESTAMPTZ NOT NULL
		)
		`,
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = fixture.pool.Exec(
		fixture.ctx,
		`
		CREATE UNIQUE INDEX IF NOT EXISTS groupings_topics_name_idx
			ON groupings_topics (name)
		`,
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = fixture.pool.Exec(
		fixture.ctx,
		`
		CREATE INDEX IF NOT EXISTS groupings_topics_parent_id_idx
			ON groupings_topics (parent_id)
		`,
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = fixture.pool.Exec(
		fixture.ctx,
		`
		CREATE INDEX IF NOT EXISTS groupings_topics_cluster_id_idx
			ON groupings_topics (cluster_id)
		`,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func truncateGroupingsTopicsTable(
	t *testing.T,
	fixture *groupingsTopicRepositoryFixture,
) {
	t.Helper()

	_, err := fixture.pool.Exec(
		fixture.ctx,
		`
		TRUNCATE TABLE groupings_topics
		`,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func saveTopics(
	t *testing.T,
	fixture *groupingsTopicRepositoryFixture,
	topics ...domain_topics.Topic,
) {
	t.Helper()

	for _, topic := range topics {
		if err := fixture.repository.Save(fixture.ctx, topic); err != nil {
			t.Fatal(err)
		}
	}
}

type testTopic struct {
	id uuid.UUID

	cluster domain_clusters.Cluster

	name        string
	description string

	parent domain_topics.Topic

	createdOn time.Time
}

func newTestTopic(
	t *testing.T,
	fixture *groupingsTopicRepositoryFixture,
	name string,
	description string,
	parent domain_topics.Topic,
) domain_topics.Topic {
	t.Helper()

	cluster := domain_clusters.NewMockCluster(
		clusterables.NewMockClusterable(clusterables.TopicKind),
		clusterables.PostKind,
		[]uuid.UUID{uuid.New()},
	)

	fixture.clusters.Items[cluster.Identifier()] = cluster

	return &testTopic{
		id:          uuid.New(),
		cluster:     cluster,
		name:        name,
		description: description,
		parent:      parent,
		createdOn:   time.Now().UTC().Truncate(time.Microsecond),
	}
}

func (topic *testTopic) Identifier() uuid.UUID {
	return topic.id
}

func (topic *testTopic) ParticipationKind() participatables.Kind {
	return participatables.TopicKind
}

func (topic *testTopic) Cluster() domain_clusters.Cluster {
	return topic.cluster
}

func (topic *testTopic) Name() string {
	return topic.name
}

func (topic *testTopic) Description() string {
	return topic.description
}

func (topic *testTopic) CreatedOn() time.Time {
	return topic.createdOn
}

func (topic *testTopic) HasParent() bool {
	return topic.parent != nil
}

func (topic *testTopic) Parent() domain_topics.Topic {
	return topic.parent
}

func assertTopic(
	t *testing.T,
	result domain_topics.Topic,
	expected domain_topics.Topic,
) {
	t.Helper()

	if result == nil {
		t.Fatalf("expected topic")
	}

	if result.Identifier() != expected.Identifier() {
		t.Fatalf(
			"expected id %s, got %s",
			expected.Identifier(),
			result.Identifier(),
		)
	}

	if result.Cluster().Identifier() != expected.Cluster().Identifier() {
		t.Fatalf(
			"expected cluster %s, got %s",
			expected.Cluster().Identifier(),
			result.Cluster().Identifier(),
		)
	}

	if result.Name() != expected.Name() {
		t.Fatalf(
			"expected name %s, got %s",
			expected.Name(),
			result.Name(),
		)
	}

	if result.Description() != expected.Description() {
		t.Fatalf(
			"expected description %s, got %s",
			expected.Description(),
			result.Description(),
		)
	}

	if result.HasParent() != expected.HasParent() {
		t.Fatalf(
			"expected has parent %v, got %v",
			expected.HasParent(),
			result.HasParent(),
		)
	}

	if expected.HasParent() {
		if result.Parent() == nil {
			t.Fatalf("expected parent")
		}

		if result.Parent().Identifier() != expected.Parent().Identifier() {
			t.Fatalf(
				"expected parent %s, got %s",
				expected.Parent().Identifier(),
				result.Parent().Identifier(),
			)
		}
	}

	if !result.CreatedOn().Equal(expected.CreatedOn()) {
		t.Fatalf(
			"expected created on %s, got %s",
			expected.CreatedOn(),
			result.CreatedOn(),
		)
	}
}

func sortTopics(
	topics []domain_topics.Topic,
) {
	sort.Slice(topics, func(left int, right int) bool {
		return topics[left].Identifier().String() <
			topics[right].Identifier().String()
	})
}
