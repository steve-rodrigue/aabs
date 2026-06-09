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
	domain_narratives "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/narratives"
	domain_assignments "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/narratives/assignments"
)

func TestNewGroupingsAssignmentRepository(t *testing.T) {
	fixture := newGroupingsNarrativesAssignmentRepositoryFixture(t)

	repository := NewGroupingsAssignmentRepository(
		fixture.pool,
		fixture.adapter,
		fixture.narratives,
		fixture.campaigns,
	)

	if repository == nil {
		t.Fatalf("expected repository")
	}
}

func TestGroupingsNarrativesAssignmentRepositorySaveAndFindByID(t *testing.T) {
	fixture := newGroupingsNarrativesAssignmentRepositoryFixture(t)

	assignment := newTestAssignment(t, fixture)

	if err := fixture.repository.Save(fixture.ctx, assignment); err != nil {
		t.Fatal(err)
	}

	result, err := fixture.repository.FindByID(
		fixture.ctx,
		assignment.Identifier(),
	)
	if err != nil {
		t.Fatal(err)
	}

	assertAssignment(t, result, assignment)
}

func TestGroupingsNarrativesAssignmentRepositoryFindByIDReturnsNilWhenNotFound(t *testing.T) {
	fixture := newGroupingsNarrativesAssignmentRepositoryFixture(t)

	result, err := fixture.repository.FindByID(
		fixture.ctx,
		uuid.New(),
	)
	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Fatalf("expected nil assignment")
	}
}

func TestGroupingsNarrativesAssignmentRepositorySaveUpdatesExistingAssignment(t *testing.T) {
	fixture := newGroupingsNarrativesAssignmentRepositoryFixture(t)

	assignment := newTestAssignment(t, fixture)

	if err := fixture.repository.Save(fixture.ctx, assignment); err != nil {
		t.Fatal(err)
	}

	updated := &domain_assignments.MockAssignment{
		ID:              assignment.Identifier(),
		NarrativeValue:  assignment.Narrative(),
		CampaignValue:   assignment.Campaign(),
		ConfidenceValue: 0.42,
		AssignedOnValue: assignment.AssignedOn(),
	}

	if err := fixture.repository.Save(fixture.ctx, updated); err != nil {
		t.Fatal(err)
	}

	result, err := fixture.repository.FindByID(
		fixture.ctx,
		assignment.Identifier(),
	)
	if err != nil {
		t.Fatal(err)
	}

	assertAssignment(t, result, updated)
}

func TestGroupingsNarrativesAssignmentRepositoryFindByNarrative(t *testing.T) {
	fixture := newGroupingsNarrativesAssignmentRepositoryFixture(t)

	expected := newTestAssignment(t, fixture)
	other := newTestAssignment(t, fixture)

	saveAssignments(t, fixture, expected, other)

	result, err := fixture.repository.FindByNarrative(
		fixture.ctx,
		expected.Narrative().Identifier(),
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 assignment, got %d", len(result))
	}

	assertAssignment(t, result[0], expected)
}

func TestGroupingsNarrativesAssignmentRepositoryFindByCampaign(t *testing.T) {
	fixture := newGroupingsNarrativesAssignmentRepositoryFixture(t)

	expected := newTestAssignment(t, fixture)
	other := newTestAssignment(t, fixture)

	saveAssignments(t, fixture, expected, other)

	result, err := fixture.repository.FindByCampaign(
		fixture.ctx,
		expected.Campaign().Identifier(),
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 assignment, got %d", len(result))
	}

	assertAssignment(t, result[0], expected)
}

func TestGroupingsNarrativesAssignmentRepositoryFindBetween(t *testing.T) {
	fixture := newGroupingsNarrativesAssignmentRepositoryFixture(t)

	expected := newTestAssignment(t, fixture)
	other := newTestAssignment(t, fixture)

	saveAssignments(t, fixture, expected, other)

	result, err := fixture.repository.FindBetween(
		fixture.ctx,
		expected.Narrative().Identifier(),
		expected.Campaign().Identifier(),
	)
	if err != nil {
		t.Fatal(err)
	}

	assertAssignment(t, result, expected)
}

func TestGroupingsNarrativesAssignmentRepositoryFindBetweenReturnsNilWhenNotFound(t *testing.T) {
	fixture := newGroupingsNarrativesAssignmentRepositoryFixture(t)

	result, err := fixture.repository.FindBetween(
		fixture.ctx,
		uuid.New(),
		uuid.New(),
	)
	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Fatalf("expected nil assignment")
	}
}

func TestGroupingsNarrativesAssignmentRepositoryFind(t *testing.T) {
	fixture := newGroupingsNarrativesAssignmentRepositoryFixture(t)

	first := newTestAssignment(t, fixture)
	second := newTestAssignment(t, fixture)
	third := newTestAssignment(t, fixture)

	saveAssignments(t, fixture, first, second, third)

	result, err := fixture.repository.Find(
		fixture.ctx,
		0,
		2,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 assignments, got %d", len(result))
	}
}

func TestGroupingsNarrativesAssignmentRepositoryFindAfter(t *testing.T) {
	fixture := newGroupingsNarrativesAssignmentRepositoryFixture(t)

	first := newTestAssignment(t, fixture)
	second := newTestAssignment(t, fixture)
	third := newTestAssignment(t, fixture)

	assignments := []domain_assignments.Assignment{
		first,
		second,
		third,
	}

	sortAssignmentsByID(assignments)

	saveAssignments(t, fixture, assignments...)

	result, err := fixture.repository.FindAfter(
		fixture.ctx,
		assignments[0].Identifier(),
		10,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 assignments, got %d", len(result))
	}

	if result[0].Identifier() != assignments[1].Identifier() {
		t.Fatalf("expected second assignment first")
	}
}

func TestGroupingsNarrativesAssignmentRepositoryFindAfterWithNilCursor(t *testing.T) {
	fixture := newGroupingsNarrativesAssignmentRepositoryFixture(t)

	first := newTestAssignment(t, fixture)
	second := newTestAssignment(t, fixture)

	saveAssignments(t, fixture, first, second)

	result, err := fixture.repository.FindAfter(
		fixture.ctx,
		uuid.Nil,
		10,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 assignments, got %d", len(result))
	}
}

func TestGroupingsNarrativesAssignmentRepositoryCount(t *testing.T) {
	fixture := newGroupingsNarrativesAssignmentRepositoryFixture(t)

	first := newTestAssignment(t, fixture)
	second := newTestAssignment(t, fixture)

	saveAssignments(t, fixture, first, second)

	count, err := fixture.repository.Count(fixture.ctx)
	if err != nil {
		t.Fatal(err)
	}

	if count != 2 {
		t.Fatalf("expected count 2, got %d", count)
	}
}

type groupingsNarrativesAssignmentRepositoryFixture struct {
	ctx context.Context

	pool *pgxpool.Pool

	adapter    domain_assignments.Adapter
	narratives *domain_narratives.MockNarrativeRepository
	campaigns  *domain_campaigns.MockCampaignRepository

	repository domain_assignments.Repository
}

func newGroupingsNarrativesAssignmentRepositoryFixture(
	t *testing.T,
) *groupingsNarrativesAssignmentRepositoryFixture {
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

	fixture := &groupingsNarrativesAssignmentRepositoryFixture{
		ctx:        ctx,
		pool:       pool,
		adapter:    domain_assignments.NewAdapter(),
		narratives: domain_narratives.NewMockNarrativeRepository(),
		campaigns:  domain_campaigns.NewMockCampaignRepository(),
	}

	createGroupingsNarrativesAssignmentTable(t, fixture)
	truncateGroupingsNarrativesAssignmentTable(t, fixture)

	fixture.repository = NewGroupingsAssignmentRepository(
		pool,
		fixture.adapter,
		fixture.narratives,
		fixture.campaigns,
	)

	t.Cleanup(func() {
		truncateGroupingsNarrativesAssignmentTable(t, fixture)
		pool.Close()
	})

	return fixture
}

func createGroupingsNarrativesAssignmentTable(
	t *testing.T,
	fixture *groupingsNarrativesAssignmentRepositoryFixture,
) {
	t.Helper()

	_, err := fixture.pool.Exec(
		fixture.ctx,
		`
		CREATE TABLE IF NOT EXISTS groupings_assignments (
			identifier UUID PRIMARY KEY,
			narrative_id UUID NOT NULL,
			campaign_id UUID NOT NULL,
			confidence DOUBLE PRECISION NOT NULL,
			assigned_on TIMESTAMPTZ NOT NULL
		)
		`,
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = fixture.pool.Exec(
		fixture.ctx,
		`
		CREATE INDEX IF NOT EXISTS groupings_assignments_narrative_idx
			ON groupings_assignments (narrative_id)
		`,
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = fixture.pool.Exec(
		fixture.ctx,
		`
		CREATE INDEX IF NOT EXISTS groupings_assignments_campaign_idx
			ON groupings_assignments (campaign_id)
		`,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func truncateGroupingsNarrativesAssignmentTable(
	t *testing.T,
	fixture *groupingsNarrativesAssignmentRepositoryFixture,
) {
	t.Helper()

	_, err := fixture.pool.Exec(
		fixture.ctx,
		`
		TRUNCATE TABLE groupings_assignments
		`,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func saveAssignments(
	t *testing.T,
	fixture *groupingsNarrativesAssignmentRepositoryFixture,
	assignments ...domain_assignments.Assignment,
) {
	t.Helper()

	for _, assignment := range assignments {
		if err := fixture.repository.Save(fixture.ctx, assignment); err != nil {
			t.Fatal(err)
		}
	}
}

func newTestAssignment(
	t *testing.T,
	fixture *groupingsNarrativesAssignmentRepositoryFixture,
) domain_assignments.Assignment {
	t.Helper()

	narrative := domain_narratives.NewMockNarrative(
		"Narrative "+uuid.NewString(),
		"Description",
	)

	campaign := domain_campaigns.NewMockCampaign(
		"Campaign "+uuid.NewString(),
		"Description",
	)

	fixture.narratives.Items[narrative.Identifier()] = narrative
	fixture.campaigns.Items[campaign.Identifier()] = campaign

	assignment, err := fixture.adapter.ToDomain(
		domain_assignments.AssignmentInput{
			Identifier: uuid.New(),
			Narrative:  narrative,
			Campaign:   campaign,
			Confidence: 0.8,
			AssignedOn: time.Now().UTC().Truncate(time.Microsecond),
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	return assignment
}

func assertAssignment(
	t *testing.T,
	result domain_assignments.Assignment,
	expected domain_assignments.Assignment,
) {
	t.Helper()

	if result == nil {
		t.Fatalf("expected assignment")
	}

	if result.Identifier() != expected.Identifier() {
		t.Fatalf(
			"expected id %s, got %s",
			expected.Identifier(),
			result.Identifier(),
		)
	}

	if result.Narrative().Identifier() != expected.Narrative().Identifier() {
		t.Fatalf("expected narrative")
	}

	if result.Campaign().Identifier() != expected.Campaign().Identifier() {
		t.Fatalf("expected campaign")
	}

	if result.Confidence() != expected.Confidence() {
		t.Fatalf(
			"expected confidence %f, got %f",
			expected.Confidence(),
			result.Confidence(),
		)
	}

	if !result.AssignedOn().Equal(expected.AssignedOn()) {
		t.Fatalf(
			"expected assigned on %s, got %s",
			expected.AssignedOn(),
			result.AssignedOn(),
		)
	}
}

func sortAssignmentsByID(
	assignments []domain_assignments.Assignment,
) {
	sort.Slice(assignments, func(left int, right int) bool {
		return assignments[left].Identifier().String() <
			assignments[right].Identifier().String()
	})
}
