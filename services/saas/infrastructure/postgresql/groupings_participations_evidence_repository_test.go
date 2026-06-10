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
	domain_posts "github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts"
	domain_participations "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations"
	domain_evidences "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/evidences"
)

func TestNewGroupingsParticipationsEvidenceRepository(t *testing.T) {
	fixture := newGroupingsParticipationsEvidenceRepositoryFixture(t)

	repository := NewGroupingsParticipationsEvidenceRepository(
		fixture.pool,
		fixture.adapter,
		fixture.participations,
		fixture.posts,
	)

	if repository == nil {
		t.Fatalf("expected repository")
	}
}

func TestGroupingsParticipationsEvidenceRepositorySaveAndFindByID(t *testing.T) {
	fixture := newGroupingsParticipationsEvidenceRepositoryFixture(t)

	evidence := newTestEvidence(t, fixture, nil, nil)

	if err := fixture.repository.Save(fixture.ctx, evidence); err != nil {
		t.Fatal(err)
	}

	result, err := fixture.repository.FindByID(
		fixture.ctx,
		evidence.Identifier(),
	)
	if err != nil {
		t.Fatal(err)
	}

	assertEvidence(t, result, evidence)
}

func TestGroupingsParticipationsEvidenceRepositoryFindByIDReturnsNilWhenNotFound(t *testing.T) {
	fixture := newGroupingsParticipationsEvidenceRepositoryFixture(t)

	result, err := fixture.repository.FindByID(
		fixture.ctx,
		uuid.New(),
	)

	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Fatalf("expected nil evidence")
	}
}

func TestGroupingsParticipationsEvidenceRepositorySaveUpdatesExistingEvidence(t *testing.T) {
	fixture := newGroupingsParticipationsEvidenceRepositoryFixture(t)

	evidence := newTestEvidence(t, fixture, nil, nil)

	if err := fixture.repository.Save(fixture.ctx, evidence); err != nil {
		t.Fatal(err)
	}

	updated := &domain_evidences.MockEvidence{
		ID:                 evidence.Identifier(),
		ParticipationValue: evidence.Participation(),
		ParticipantValue:   evidence.Participant(),
		TargetValue:        evidence.Target(),
		PostValue:          evidence.Post(),
		ScoreValue:         0.42,
		DetectedOnValue:    evidence.DetectedOn(),
	}

	if err := fixture.repository.Save(fixture.ctx, updated); err != nil {
		t.Fatal(err)
	}

	result, err := fixture.repository.FindByID(
		fixture.ctx,
		evidence.Identifier(),
	)
	if err != nil {
		t.Fatal(err)
	}

	assertEvidence(t, result, updated)
}

func TestGroupingsParticipationsEvidenceRepositoryFindByParticipation(t *testing.T) {
	fixture := newGroupingsParticipationsEvidenceRepositoryFixture(t)

	participation := newTestParticipation(
		fixture,
		nil,
		nil,
	)

	first := newTestEvidence(t, fixture, participation, nil)
	second := newTestEvidence(t, fixture, participation, nil)
	other := newTestEvidence(t, fixture, nil, nil)

	saveEvidences(t, fixture, first, second, other)

	result, err := fixture.repository.FindByParticipation(
		fixture.ctx,
		participation.Identifier(),
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 evidences, got %d", len(result))
	}

	for _, evidence := range result {
		if evidence.Participation().Identifier() != participation.Identifier() {
			t.Fatalf("expected participation")
		}
	}
}

func TestGroupingsParticipationsEvidenceRepositoryFindByPost(t *testing.T) {
	fixture := newGroupingsParticipationsEvidenceRepositoryFixture(t)

	post := domain_posts.NewMockPost("shared")

	first := newTestEvidence(t, fixture, nil, post)
	second := newTestEvidence(t, fixture, nil, post)
	other := newTestEvidence(t, fixture, nil, nil)

	saveEvidences(t, fixture, first, second, other)

	result, err := fixture.repository.FindByPost(
		fixture.ctx,
		post.Identifier(),
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 evidences, got %d", len(result))
	}

	for _, evidence := range result {
		if evidence.Post().Identifier() != post.Identifier() {
			t.Fatalf("expected post")
		}
	}
}

func TestGroupingsParticipationsEvidenceRepositoryFindByParticipant(t *testing.T) {
	fixture := newGroupingsParticipationsEvidenceRepositoryFixture(t)

	participant := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.UserKind,
	)

	firstParticipation := newTestParticipation(
		fixture,
		participant,
		nil,
	)

	secondParticipation := newTestParticipation(
		fixture,
		participant,
		nil,
	)

	otherParticipation := newTestParticipation(
		fixture,
		nil,
		nil,
	)

	first := newTestEvidence(t, fixture, firstParticipation, nil)
	second := newTestEvidence(t, fixture, secondParticipation, nil)
	other := newTestEvidence(t, fixture, otherParticipation, nil)

	saveEvidences(t, fixture, first, second, other)

	result, err := fixture.repository.FindByParticipant(
		fixture.ctx,
		participant,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 evidences, got %d", len(result))
	}

	for _, evidence := range result {
		if evidence.Participant().Identifier() != participant.Identifier() {
			t.Fatalf("expected participant")
		}
	}
}

func TestGroupingsParticipationsEvidenceRepositoryFindByTarget(t *testing.T) {
	fixture := newGroupingsParticipationsEvidenceRepositoryFixture(t)

	target := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.TopicKind,
	)

	firstParticipation := newTestParticipation(
		fixture,
		nil,
		target,
	)

	secondParticipation := newTestParticipation(
		fixture,
		nil,
		target,
	)

	otherParticipation := newTestParticipation(
		fixture,
		nil,
		nil,
	)

	first := newTestEvidence(t, fixture, firstParticipation, nil)
	second := newTestEvidence(t, fixture, secondParticipation, nil)
	other := newTestEvidence(t, fixture, otherParticipation, nil)

	saveEvidences(t, fixture, first, second, other)

	result, err := fixture.repository.FindByTarget(
		fixture.ctx,
		target,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 evidences, got %d", len(result))
	}

	for _, evidence := range result {
		if evidence.Target().Identifier() != target.Identifier() {
			t.Fatalf("expected target")
		}
	}
}

type groupingsParticipationsEvidenceRepositoryFixture struct {
	ctx context.Context

	pool *pgxpool.Pool

	adapter        domain_evidences.Adapter
	participations *domain_participations.MockParticipationRepository
	posts          *domain_posts.MockPostRepository

	repository domain_evidences.Repository
}

func newGroupingsParticipationsEvidenceRepositoryFixture(
	t *testing.T,
) *groupingsParticipationsEvidenceRepositoryFixture {
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

	fixture := &groupingsParticipationsEvidenceRepositoryFixture{
		ctx:            ctx,
		pool:           pool,
		adapter:        domain_evidences.NewAdapter(),
		participations: domain_participations.NewMockParticipationRepository(),
		posts:          domain_posts.NewMockPostRepository(),
	}

	createGroupingsParticipationsEvidenceTables(t, fixture)
	truncateGroupingsParticipationsEvidenceTables(t, fixture)

	fixture.repository = NewGroupingsParticipationsEvidenceRepository(
		pool,
		fixture.adapter,
		fixture.participations,
		fixture.posts,
	)

	t.Cleanup(func() {
		truncateGroupingsParticipationsEvidenceTables(t, fixture)
		pool.Close()
	})

	return fixture
}

func createGroupingsParticipationsEvidenceTables(
	t *testing.T,
	fixture *groupingsParticipationsEvidenceRepositoryFixture,
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
		CREATE TABLE IF NOT EXISTS groupings_participations_evidences (
			identifier UUID PRIMARY KEY,
			participation_id UUID NOT NULL,
			post_id UUID NOT NULL,
			score DOUBLE PRECISION NOT NULL,
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
		CREATE INDEX IF NOT EXISTS groupings_participations_evidences_participation_idx
			ON groupings_participations_evidences (participation_id)
		`,
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = fixture.pool.Exec(
		fixture.ctx,
		`
		CREATE INDEX IF NOT EXISTS groupings_participations_evidences_post_idx
			ON groupings_participations_evidences (post_id)
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
}

func truncateGroupingsParticipationsEvidenceTables(
	t *testing.T,
	fixture *groupingsParticipationsEvidenceRepositoryFixture,
) {
	t.Helper()

	_, err := fixture.pool.Exec(
		fixture.ctx,
		`
		TRUNCATE TABLE
			groupings_participations_evidences,
			groupings_participations
		`,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func saveEvidences(
	t *testing.T,
	fixture *groupingsParticipationsEvidenceRepositoryFixture,
	evidences ...domain_evidences.Evidence,
) {
	t.Helper()

	for _, evidence := range evidences {
		if err := fixture.repository.Save(fixture.ctx, evidence); err != nil {
			t.Fatal(err)
		}
	}
}

func newTestParticipation(
	fixture *groupingsParticipationsEvidenceRepositoryFixture,
	participant participatables.Participatable,
	target participatables.Participatable,
) domain_participations.Participation {
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

	participation := domain_participations.NewMockParticipationWithParticipantAndTarget(
		participant,
		target,
	)

	fixture.participations.Items[participation.Identifier()] = participation

	_, _ = fixture.pool.Exec(
		fixture.ctx,
		`
		INSERT INTO groupings_participations (
			identifier,
			participant_id,
			participant_kind,
			target_id,
			target_kind,
			post_count,
			total_post_count,
			percentage,
			detected_on
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		ON CONFLICT (identifier)
		DO NOTHING
		`,
		participation.Identifier(),
		participant.Identifier(),
		string(participant.ParticipationKind()),
		target.Identifier(),
		string(target.ParticipationKind()),
		1,
		1,
		1,
		time.Now().UTC(),
	)

	return participation
}

func newTestEvidence(
	t *testing.T,
	fixture *groupingsParticipationsEvidenceRepositoryFixture,
	participation domain_participations.Participation,
	post domain_posts.Post,
) domain_evidences.Evidence {
	t.Helper()

	if participation == nil {
		participation = newTestParticipation(
			fixture,
			nil,
			nil,
		)
	}

	if post == nil {
		post = domain_posts.NewMockPost("hello")
	}

	fixture.posts.Items[post.Identifier()] = post

	evidence, err := fixture.adapter.ToDomain(
		domain_evidences.EvidenceInput{
			Identifier:    uuid.New(),
			Participation: participation,
			Participant:   participation.Participant(),
			Target:        participation.Target(),
			Post:          post,
			Score:         0.9,
			DetectedOn:    time.Now().UTC().Truncate(time.Microsecond),
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	return evidence
}

func assertEvidence(
	t *testing.T,
	result domain_evidences.Evidence,
	expected domain_evidences.Evidence,
) {
	t.Helper()

	if result == nil {
		t.Fatalf("expected evidence")
	}

	if result.Identifier() != expected.Identifier() {
		t.Fatalf(
			"expected id %s, got %s",
			expected.Identifier(),
			result.Identifier(),
		)
	}

	if result.Participation().Identifier() != expected.Participation().Identifier() {
		t.Fatalf("expected participation")
	}

	if result.Participant().Identifier() != expected.Participant().Identifier() {
		t.Fatalf("expected participant")
	}

	if result.Target().Identifier() != expected.Target().Identifier() {
		t.Fatalf("expected target")
	}

	if result.Post().Identifier() != expected.Post().Identifier() {
		t.Fatalf("expected post")
	}

	if result.Score() != expected.Score() {
		t.Fatalf(
			"expected score %f, got %f",
			expected.Score(),
			result.Score(),
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

func sortEvidences(
	evidences []domain_evidences.Evidence,
) {
	sort.Slice(evidences, func(left int, right int) bool {
		return evidences[left].Identifier().String() <
			evidences[right].Identifier().String()
	})
}
