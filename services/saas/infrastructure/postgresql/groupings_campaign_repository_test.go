package postgresql

import (
	"context"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	domain_campaigns "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/campaigns"
	domain_clusters "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters/clusterables"
)

func TestNewGroupingsCampaignRepository(t *testing.T) {
	fixture := newGroupingsCampaignRepositoryFixture(t)

	repository := NewGroupingsCampaignRepository(
		fixture.pool,
		fixture.adapter,
		fixture.clusters,
	)

	if repository == nil {
		t.Fatalf("expected repository")
	}
}

func TestGroupingsCampaignRepositorySaveAndFindByID(t *testing.T) {
	fixture := newGroupingsCampaignRepositoryFixture(t)

	campaign := newTestCampaign(t, fixture, "Campaign A", "Description A")

	if err := fixture.repository.Save(fixture.ctx, campaign); err != nil {
		t.Fatal(err)
	}

	result, err := fixture.repository.FindByID(
		fixture.ctx,
		campaign.Identifier(),
	)
	if err != nil {
		t.Fatal(err)
	}

	assertCampaign(t, result, campaign)
}

func TestGroupingsCampaignRepositoryFindByIDReturnsNilWhenNotFound(t *testing.T) {
	fixture := newGroupingsCampaignRepositoryFixture(t)

	result, err := fixture.repository.FindByID(
		fixture.ctx,
		uuid.New(),
	)
	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Fatalf("expected nil campaign")
	}
}

func TestGroupingsCampaignRepositorySaveUpdatesExistingCampaign(t *testing.T) {
	fixture := newGroupingsCampaignRepositoryFixture(t)

	campaign := newTestCampaign(t, fixture, "Campaign A", "Description A")

	if err := fixture.repository.Save(fixture.ctx, campaign); err != nil {
		t.Fatal(err)
	}

	updated := &domain_campaigns.MockCampaign{
		ID:                     campaign.Identifier(),
		ParticipationKindValue: campaign.ParticipationKind(),
		NameValue:              "Campaign B",
		DescriptionValue:       "Description B",
		ClusterValue:           campaign.Cluster(),
		PostCountValue:         12,
		ConfidenceValue:        0.84,
		CreatedOnValue:         campaign.CreatedOn(),
	}

	if err := fixture.repository.Save(fixture.ctx, updated); err != nil {
		t.Fatal(err)
	}

	result, err := fixture.repository.FindByID(
		fixture.ctx,
		campaign.Identifier(),
	)
	if err != nil {
		t.Fatal(err)
	}

	assertCampaign(t, result, updated)
}

func TestGroupingsCampaignRepositoryFindByName(t *testing.T) {
	fixture := newGroupingsCampaignRepositoryFixture(t)

	expected := newTestCampaign(t, fixture, "Campaign A", "Description A")
	other := newTestCampaign(t, fixture, "Campaign B", "Description B")

	saveCampaigns(t, fixture, expected, other)

	result, err := fixture.repository.FindByName(
		fixture.ctx,
		expected.Name(),
	)
	if err != nil {
		t.Fatal(err)
	}

	assertCampaign(t, result, expected)
}

func TestGroupingsCampaignRepositoryFindByNameReturnsNilWhenNotFound(t *testing.T) {
	fixture := newGroupingsCampaignRepositoryFixture(t)

	result, err := fixture.repository.FindByName(
		fixture.ctx,
		"missing",
	)
	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Fatalf("expected nil campaign")
	}
}

func TestGroupingsCampaignRepositoryFind(t *testing.T) {
	fixture := newGroupingsCampaignRepositoryFixture(t)

	first := newTestCampaign(t, fixture, "Campaign A", "Description A")
	second := newTestCampaign(t, fixture, "Campaign B", "Description B")
	third := newTestCampaign(t, fixture, "Campaign C", "Description C")

	saveCampaigns(t, fixture, first, second, third)

	result, err := fixture.repository.Find(
		fixture.ctx,
		0,
		2,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 campaigns, got %d", len(result))
	}
}

func TestGroupingsCampaignRepositoryFindAfter(t *testing.T) {
	fixture := newGroupingsCampaignRepositoryFixture(t)

	first := newTestCampaign(t, fixture, "Campaign A", "Description A")
	second := newTestCampaign(t, fixture, "Campaign B", "Description B")
	third := newTestCampaign(t, fixture, "Campaign C", "Description C")

	campaigns := []domain_campaigns.Campaign{
		first,
		second,
		third,
	}

	sortCampaignsByID(campaigns)

	saveCampaigns(t, fixture, campaigns...)

	result, err := fixture.repository.FindAfter(
		fixture.ctx,
		campaigns[0].Identifier(),
		10,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 campaigns, got %d", len(result))
	}

	if result[0].Identifier() != campaigns[1].Identifier() {
		t.Fatalf("expected second campaign first")
	}
}

func TestGroupingsCampaignRepositoryFindAfterWithNilCursor(t *testing.T) {
	fixture := newGroupingsCampaignRepositoryFixture(t)

	first := newTestCampaign(t, fixture, "Campaign A", "Description A")
	second := newTestCampaign(t, fixture, "Campaign B", "Description B")

	saveCampaigns(t, fixture, first, second)

	result, err := fixture.repository.FindAfter(
		fixture.ctx,
		uuid.Nil,
		10,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 campaigns, got %d", len(result))
	}
}

func TestGroupingsCampaignRepositoryCount(t *testing.T) {
	fixture := newGroupingsCampaignRepositoryFixture(t)

	first := newTestCampaign(t, fixture, "Campaign A", "Description A")
	second := newTestCampaign(t, fixture, "Campaign B", "Description B")

	saveCampaigns(t, fixture, first, second)

	count, err := fixture.repository.Count(fixture.ctx)
	if err != nil {
		t.Fatal(err)
	}

	if count != 2 {
		t.Fatalf("expected count 2, got %d", count)
	}
}

type groupingsCampaignRepositoryFixture struct {
	ctx context.Context

	pool *pgxpool.Pool

	adapter    domain_campaigns.Adapter
	clusters   *domain_clusters.MockClusterRepository
	repository domain_campaigns.Repository
}

func newGroupingsCampaignRepositoryFixture(
	t *testing.T,
) *groupingsCampaignRepositoryFixture {
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

	fixture := &groupingsCampaignRepositoryFixture{
		ctx:      ctx,
		pool:     pool,
		adapter:  domain_campaigns.NewAdapter(),
		clusters: domain_clusters.NewMockClusterRepository(),
	}

	createGroupingsCampaignTable(t, fixture)
	truncateGroupingsCampaignTable(t, fixture)

	fixture.repository = NewGroupingsCampaignRepository(
		pool,
		fixture.adapter,
		fixture.clusters,
	)

	t.Cleanup(func() {
		truncateGroupingsCampaignTable(t, fixture)
		pool.Close()
	})

	return fixture
}

func createGroupingsCampaignTable(
	t *testing.T,
	fixture *groupingsCampaignRepositoryFixture,
) {
	t.Helper()

	_, err := fixture.pool.Exec(
		fixture.ctx,
		`
		CREATE TABLE IF NOT EXISTS groupings_campaigns (
			identifier UUID PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT NOT NULL,
			cluster_id UUID NOT NULL,
			post_count INTEGER NOT NULL,
			confidence DOUBLE PRECISION NOT NULL,
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
		CREATE INDEX IF NOT EXISTS groupings_campaigns_name_idx
			ON groupings_campaigns (name)
		`,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func truncateGroupingsCampaignTable(
	t *testing.T,
	fixture *groupingsCampaignRepositoryFixture,
) {
	t.Helper()

	_, err := fixture.pool.Exec(
		fixture.ctx,
		`
		TRUNCATE TABLE groupings_campaigns
		`,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func saveCampaigns(
	t *testing.T,
	fixture *groupingsCampaignRepositoryFixture,
	campaigns ...domain_campaigns.Campaign,
) {
	t.Helper()

	for _, campaign := range campaigns {
		if err := fixture.repository.Save(fixture.ctx, campaign); err != nil {
			t.Fatal(err)
		}
	}
}

func newTestCampaign(
	t *testing.T,
	fixture *groupingsCampaignRepositoryFixture,
	name string,
	description string,
) domain_campaigns.Campaign {
	t.Helper()

	cluster := domain_clusters.NewMockCluster(
		clusterables.NewMockClusterable(clusterables.CampaignKind),
		clusterables.PostKind,
		[]uuid.UUID{uuid.New()},
	)

	fixture.clusters.Items[cluster.Identifier()] = cluster

	campaign, err := fixture.adapter.ToDomain(
		domain_campaigns.CampaignInput{
			Identifier:  uuid.New(),
			Name:        name,
			Description: description,
			Cluster:     cluster,
			PostCount:   5,
			Confidence:  0.8,
			CreatedOn:   time.Now().UTC().Truncate(time.Microsecond),
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	return campaign
}

func assertCampaign(
	t *testing.T,
	result domain_campaigns.Campaign,
	expected domain_campaigns.Campaign,
) {
	t.Helper()

	if result == nil {
		t.Fatalf("expected campaign")
	}

	if result.Identifier() != expected.Identifier() {
		t.Fatalf(
			"expected id %s, got %s",
			expected.Identifier(),
			result.Identifier(),
		)
	}

	if result.ParticipationKind() != expected.ParticipationKind() {
		t.Fatalf(
			"expected participation kind %s, got %s",
			expected.ParticipationKind(),
			result.ParticipationKind(),
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

	if result.Cluster().Identifier() != expected.Cluster().Identifier() {
		t.Fatalf("expected cluster")
	}

	if result.PostCount() != expected.PostCount() {
		t.Fatalf(
			"expected post count %d, got %d",
			expected.PostCount(),
			result.PostCount(),
		)
	}

	if result.Confidence() != expected.Confidence() {
		t.Fatalf(
			"expected confidence %f, got %f",
			expected.Confidence(),
			result.Confidence(),
		)
	}

	if !result.CreatedOn().Equal(expected.CreatedOn()) {
		t.Fatalf(
			"expected created on %s, got %s",
			expected.CreatedOn(),
			result.CreatedOn(),
		)
	}
}

func sortCampaignsByID(
	campaigns []domain_campaigns.Campaign,
) {
	sort.Slice(campaigns, func(left int, right int) bool {
		return campaigns[left].Identifier().String() <
			campaigns[right].Identifier().String()
	})
}
