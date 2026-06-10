package postgresql

import (
	"context"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/participatables"
	domain_participations "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations"
)

func TestNewGroupingsParticipationRepository(t *testing.T) {
	fixture := newGroupingsParticipationRepositoryFixture(t)

	repository := NewGroupingsParticipationRepository(
		fixture.pool,
		fixture.adapter,
	)

	if repository == nil {
		t.Fatalf("expected repository")
	}
}

func TestGroupingsParticipationRepositorySaveAndFindByID(t *testing.T) {
	fixture := newGroupingsParticipationRepositoryFixture(t)

	participation := newTestParticipationForRepository(
		t,
		nil,
		nil,
		4,
		10,
		0.4,
	)

	if err := fixture.repository.Save(fixture.ctx, participation); err != nil {
		t.Fatal(err)
	}

	result, err := fixture.repository.FindByID(
		fixture.ctx,
		participation.Identifier(),
	)
	if err != nil {
		t.Fatal(err)
	}

	assertParticipation(t, result, participation)
}

func TestGroupingsParticipationRepositoryFindByIDReturnsNilWhenNotFound(t *testing.T) {
	fixture := newGroupingsParticipationRepositoryFixture(t)

	result, err := fixture.repository.FindByID(
		fixture.ctx,
		uuid.New(),
	)

	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Fatalf("expected nil participation")
	}
}

func TestGroupingsParticipationRepositorySaveUpdatesExistingParticipation(t *testing.T) {
	fixture := newGroupingsParticipationRepositoryFixture(t)

	participation := newTestParticipationForRepository(
		t,
		nil,
		nil,
		4,
		10,
		0.4,
	)

	if err := fixture.repository.Save(fixture.ctx, participation); err != nil {
		t.Fatal(err)
	}

	updated := &domain_participations.MockParticipation{
		ID:                  participation.Identifier(),
		ParticipantValue:    participation.Participant(),
		TargetValue:         participation.Target(),
		PostCountValue:      8,
		TotalPostCountValue: 10,
		PercentageValue:     0.8,
		DetectedOnValue:     participation.DetectedOn(),
	}

	if err := fixture.repository.Save(fixture.ctx, updated); err != nil {
		t.Fatal(err)
	}

	result, err := fixture.repository.FindByID(
		fixture.ctx,
		participation.Identifier(),
	)
	if err != nil {
		t.Fatal(err)
	}

	assertParticipation(t, result, updated)
}

func TestGroupingsParticipationRepositoryFindByParticipant(t *testing.T) {
	fixture := newGroupingsParticipationRepositoryFixture(t)

	participant := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.UserKind,
	)

	first := newTestParticipationForRepository(
		t,
		participant,
		nil,
		2,
		10,
		0.2,
	)

	second := newTestParticipationForRepository(
		t,
		participant,
		nil,
		3,
		10,
		0.3,
	)

	other := newTestParticipationForRepository(
		t,
		nil,
		nil,
		4,
		10,
		0.4,
	)

	saveParticipations(t, fixture, first, second, other)

	result, err := fixture.repository.FindByParticipant(
		fixture.ctx,
		participant,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 participations, got %d", len(result))
	}

	for _, participation := range result {
		if participation.Participant().Identifier() != participant.Identifier() {
			t.Fatalf("expected participant")
		}

		if participation.Participant().ParticipationKind() != participant.ParticipationKind() {
			t.Fatalf("expected participant kind")
		}
	}
}

func TestGroupingsParticipationRepositoryFindByTarget(t *testing.T) {
	fixture := newGroupingsParticipationRepositoryFixture(t)

	target := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.TopicKind,
	)

	first := newTestParticipationForRepository(
		t,
		nil,
		target,
		2,
		10,
		0.2,
	)

	second := newTestParticipationForRepository(
		t,
		nil,
		target,
		3,
		10,
		0.3,
	)

	other := newTestParticipationForRepository(
		t,
		nil,
		nil,
		4,
		10,
		0.4,
	)

	saveParticipations(t, fixture, first, second, other)

	result, err := fixture.repository.FindByTarget(
		fixture.ctx,
		target,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 participations, got %d", len(result))
	}

	for _, participation := range result {
		if participation.Target().Identifier() != target.Identifier() {
			t.Fatalf("expected target")
		}

		if participation.Target().ParticipationKind() != target.ParticipationKind() {
			t.Fatalf("expected target kind")
		}
	}
}

func TestGroupingsParticipationRepositoryFindBetween(t *testing.T) {
	fixture := newGroupingsParticipationRepositoryFixture(t)

	participant := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.UserKind,
	)

	target := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.TopicKind,
	)

	expected := newTestParticipationForRepository(
		t,
		participant,
		target,
		6,
		10,
		0.6,
	)

	other := newTestParticipationForRepository(
		t,
		participant,
		nil,
		3,
		10,
		0.3,
	)

	saveParticipations(t, fixture, expected, other)

	result, err := fixture.repository.FindBetween(
		fixture.ctx,
		participant,
		target,
	)
	if err != nil {
		t.Fatal(err)
	}

	assertParticipation(t, result, expected)
}

func TestGroupingsParticipationRepositoryFindBetweenReturnsNilWhenNotFound(t *testing.T) {
	fixture := newGroupingsParticipationRepositoryFixture(t)

	result, err := fixture.repository.FindBetween(
		fixture.ctx,
		participatables.NewMockParticipatable(
			uuid.New(),
			participatables.UserKind,
		),
		participatables.NewMockParticipatable(
			uuid.New(),
			participatables.TopicKind,
		),
	)

	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Fatalf("expected nil participation")
	}
}

type groupingsParticipationRepositoryFixture struct {
	ctx context.Context

	pool *pgxpool.Pool

	adapter    domain_participations.Adapter
	repository domain_participations.Repository
}

func newGroupingsParticipationRepositoryFixture(
	t *testing.T,
) *groupingsParticipationRepositoryFixture {
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

	fixture := &groupingsParticipationRepositoryFixture{
		ctx:     ctx,
		pool:    pool,
		adapter: domain_participations.NewAdapter(),
	}

	createGroupingsParticipationTable(t, fixture)
	truncateGroupingsParticipationTable(t, fixture)

	fixture.repository = NewGroupingsParticipationRepository(
		pool,
		fixture.adapter,
	)

	t.Cleanup(func() {
		truncateGroupingsParticipationTable(t, fixture)
		pool.Close()
	})

	return fixture
}

func createGroupingsParticipationTable(
	t *testing.T,
	fixture *groupingsParticipationRepositoryFixture,
) {
	t.Helper()

	_, err := fixture.pool.Exec(
		fixture.ctx,
		`
		CREATE TABLE IF NOT EXISTS groupings_participations (
			identifier UUID PRIMARY KEY,
			participant_id UUID NOT NULL,
			participant_kind TEXT NOT NULL,
			target_id UUID NOT NULL,
			target_kind TEXT NOT NULL,
			post_count INTEGER NOT NULL,
			total_post_count INTEGER NOT NULL,
			percentage DOUBLE PRECISION NOT NULL,
			detected_on TIMESTAMPTZ NOT NULL
		)
		`,
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = fixture.pool.Exec(
		fixture.ctx,
		`
		CREATE INDEX IF NOT EXISTS groupings_participations_participant_idx
			ON groupings_participations (participant_id, participant_kind)
		`,
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = fixture.pool.Exec(
		fixture.ctx,
		`
		CREATE INDEX IF NOT EXISTS groupings_participations_target_idx
			ON groupings_participations (target_id, target_kind)
		`,
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = fixture.pool.Exec(
		fixture.ctx,
		`
		CREATE INDEX IF NOT EXISTS groupings_participations_between_idx
			ON groupings_participations (
				participant_id,
				participant_kind,
				target_id,
				target_kind
			)
		`,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func truncateGroupingsParticipationTable(
	t *testing.T,
	fixture *groupingsParticipationRepositoryFixture,
) {
	t.Helper()

	_, err := fixture.pool.Exec(
		fixture.ctx,
		`
		TRUNCATE TABLE groupings_participations
		`,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func saveParticipations(
	t *testing.T,
	fixture *groupingsParticipationRepositoryFixture,
	participations ...domain_participations.Participation,
) {
	t.Helper()

	for _, participation := range participations {
		if err := fixture.repository.Save(fixture.ctx, participation); err != nil {
			t.Fatal(err)
		}
	}
}

func newTestParticipationForRepository(
	t *testing.T,
	participant participatables.Participatable,
	target participatables.Participatable,
	postCount int,
	totalPostCount int,
	percentage float64,
) domain_participations.Participation {
	t.Helper()

	if participant == nil {
		participant = participatables.NewMockParticipatable(
			uuid.New(),
			participatables.UserKind,
		)
	}

	if target == nil {
		target = participatables.NewMockParticipatable(
			uuid.New(),
			participatables.TopicKind,
		)
	}

	participation, err := domain_participations.NewAdapter().ToDomain(
		domain_participations.ParticipationInput{
			Identifier:     uuid.New(),
			Participant:    participant,
			Target:         target,
			PostCount:      postCount,
			TotalPostCount: totalPostCount,
			Percentage:     percentage,
			DetectedOn:     time.Now().UTC().Truncate(time.Microsecond),
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	return participation
}

func assertParticipation(
	t *testing.T,
	result domain_participations.Participation,
	expected domain_participations.Participation,
) {
	t.Helper()

	if result == nil {
		t.Fatalf("expected participation")
	}

	if result.Identifier() != expected.Identifier() {
		t.Fatalf(
			"expected id %s, got %s",
			expected.Identifier(),
			result.Identifier(),
		)
	}

	if result.Participant().Identifier() != expected.Participant().Identifier() {
		t.Fatalf("expected participant")
	}

	if result.Participant().ParticipationKind() != expected.Participant().ParticipationKind() {
		t.Fatalf("expected participant kind")
	}

	if result.Target().Identifier() != expected.Target().Identifier() {
		t.Fatalf("expected target")
	}

	if result.Target().ParticipationKind() != expected.Target().ParticipationKind() {
		t.Fatalf("expected target kind")
	}

	if result.PostCount() != expected.PostCount() {
		t.Fatalf(
			"expected post count %d, got %d",
			expected.PostCount(),
			result.PostCount(),
		)
	}

	if result.TotalPostCount() != expected.TotalPostCount() {
		t.Fatalf(
			"expected total post count %d, got %d",
			expected.TotalPostCount(),
			result.TotalPostCount(),
		)
	}

	if result.Percentage() != expected.Percentage() {
		t.Fatalf(
			"expected percentage %f, got %f",
			expected.Percentage(),
			result.Percentage(),
		)
	}

	if !result.DetectedOn().Equal(expected.DetectedOn()) {
		t.Fatalf(
			"expected detected on %s, got %s",
			expected.DetectedOn(),
			result.DetectedOn(),
		)
	}
}

func sortParticipations(
	participations []domain_participations.Participation,
) {
	sort.Slice(participations, func(left int, right int) bool {
		return participations[left].Identifier().String() <
			participations[right].Identifier().String()
	})
}
