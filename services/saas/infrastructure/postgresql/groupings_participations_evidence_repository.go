package postgresql

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/participatables"
	domain_posts "github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts"
	domain_participations "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations"
	domain_evidences "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/evidences"
)

type groupingsParticipationsEvidenceRepository struct {
	pool *pgxpool.Pool

	adapter domain_evidences.Adapter

	participations domain_participations.Repository
	posts          domain_posts.Repository
}

func (repository *groupingsParticipationsEvidenceRepository) Save(
	ctx context.Context,
	evidence domain_evidences.Evidence,
) error {
	_, err := repository.pool.Exec(
		ctx,
		`
		INSERT INTO groupings_participations_evidences (
			identifier,
			participation_id,
			post_id,
			score,
			detected_on
		)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (identifier)
		DO UPDATE SET
			participation_id = EXCLUDED.participation_id,
			post_id = EXCLUDED.post_id,
			score = EXCLUDED.score,
			detected_on = EXCLUDED.detected_on
		`,
		evidence.Identifier(),
		evidence.Participation().Identifier(),
		evidence.Post().Identifier(),
		evidence.Score(),
		evidence.DetectedOn(),
	)

	return err
}

func (repository *groupingsParticipationsEvidenceRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (domain_evidences.Evidence, error) {
	return repository.findOne(
		ctx,
		`
		SELECT
			identifier,
			participation_id,
			post_id,
			score,
			detected_on
		FROM groupings_participations_evidences
		WHERE identifier = $1
		`,
		id,
	)
}

func (repository *groupingsParticipationsEvidenceRepository) FindByParticipation(
	ctx context.Context,
	participation uuid.UUID,
) ([]domain_evidences.Evidence, error) {
	rows, err := repository.pool.Query(
		ctx,
		`
		SELECT
			identifier,
			participation_id,
			post_id,
			score,
			detected_on
		FROM groupings_participations_evidences
		WHERE participation_id = $1
		ORDER BY identifier
		`,
		participation,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return repository.scanMany(ctx, rows)
}

func (repository *groupingsParticipationsEvidenceRepository) FindByPost(
	ctx context.Context,
	post uuid.UUID,
) ([]domain_evidences.Evidence, error) {
	rows, err := repository.pool.Query(
		ctx,
		`
		SELECT
			identifier,
			participation_id,
			post_id,
			score,
			detected_on
		FROM groupings_participations_evidences
		WHERE post_id = $1
		ORDER BY identifier
		`,
		post,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return repository.scanMany(ctx, rows)
}

func (repository *groupingsParticipationsEvidenceRepository) FindByParticipant(
	ctx context.Context,
	participant participatables.Participatable,
) ([]domain_evidences.Evidence, error) {
	rows, err := repository.pool.Query(
		ctx,
		`
		SELECT
			e.identifier,
			e.participation_id,
			e.post_id,
			e.score,
			e.detected_on
		FROM groupings_participations_evidences e
		INNER JOIN groupings_participations p
			ON p.identifier = e.participation_id
		WHERE p.participant_id = $1
		AND p.participant_kind = $2
		ORDER BY e.identifier
		`,
		participant.Identifier(),
		string(participant.ParticipationKind()),
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return repository.scanMany(ctx, rows)
}

func (repository *groupingsParticipationsEvidenceRepository) FindByTarget(
	ctx context.Context,
	target participatables.Participatable,
) ([]domain_evidences.Evidence, error) {
	rows, err := repository.pool.Query(
		ctx,
		`
		SELECT
			e.identifier,
			e.participation_id,
			e.post_id,
			e.score,
			e.detected_on
		FROM groupings_participations_evidences e
		INNER JOIN groupings_participations p
			ON p.identifier = e.participation_id
		WHERE p.target_id = $1
		AND p.target_kind = $2
		ORDER BY e.identifier
		`,
		target.Identifier(),
		string(target.ParticipationKind()),
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return repository.scanMany(ctx, rows)
}

func (repository *groupingsParticipationsEvidenceRepository) findOne(
	ctx context.Context,
	query string,
	args ...any,
) (domain_evidences.Evidence, error) {
	row := repository.pool.QueryRow(
		ctx,
		query,
		args...,
	)

	evidence, err := repository.scanOne(ctx, row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return evidence, nil
}

func (repository *groupingsParticipationsEvidenceRepository) scanMany(
	ctx context.Context,
	rows pgx.Rows,
) ([]domain_evidences.Evidence, error) {
	out := []domain_evidences.Evidence{}

	for rows.Next() {
		evidence, err := repository.scanOne(ctx, rows)
		if err != nil {
			return nil, err
		}

		out = append(out, evidence)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (repository *groupingsParticipationsEvidenceRepository) scanOne(
	ctx context.Context,
	row pgx.Row,
) (domain_evidences.Evidence, error) {
	var input domain_evidences.EvidenceInput

	var participationID uuid.UUID
	var postID uuid.UUID

	err := row.Scan(
		&input.Identifier,
		&participationID,
		&postID,
		&input.Score,
		&input.DetectedOn,
	)
	if err != nil {
		return nil, err
	}

	participation, err := repository.participations.FindByID(
		ctx,
		participationID,
	)
	if err != nil {
		return nil, err
	}

	post, err := repository.posts.FindByID(
		ctx,
		postID,
	)
	if err != nil {
		return nil, err
	}

	input.Participation = participation
	input.Post = post

	if participation != nil {
		input.Participant = participation.Participant()
		input.Target = participation.Target()
	}

	return repository.adapter.ToDomain(input)
}
