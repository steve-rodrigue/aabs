package postgresql

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/participatables"
	domain_participations "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations"
)

type groupingsParticipationRepository struct {
	pool    *pgxpool.Pool
	adapter domain_participations.Adapter
}

func (repository *groupingsParticipationRepository) Save(
	ctx context.Context,
	participation domain_participations.Participation,
) error {
	_, err := repository.pool.Exec(
		ctx,
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
		DO UPDATE SET
			participant_id = EXCLUDED.participant_id,
			participant_kind = EXCLUDED.participant_kind,
			target_id = EXCLUDED.target_id,
			target_kind = EXCLUDED.target_kind,
			post_count = EXCLUDED.post_count,
			total_post_count = EXCLUDED.total_post_count,
			percentage = EXCLUDED.percentage,
			detected_on = EXCLUDED.detected_on
		`,
		participation.Identifier(),
		participation.Participant().Identifier(),
		string(participation.Participant().ParticipationKind()),
		participation.Target().Identifier(),
		string(participation.Target().ParticipationKind()),
		participation.PostCount(),
		participation.TotalPostCount(),
		participation.Percentage(),
		participation.DetectedOn(),
	)

	return err
}

func (repository *groupingsParticipationRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (domain_participations.Participation, error) {
	return repository.findOne(
		ctx,
		`
		SELECT
			identifier,
			participant_id,
			participant_kind,
			target_id,
			target_kind,
			post_count,
			total_post_count,
			percentage,
			detected_on
		FROM groupings_participations
		WHERE identifier = $1
		`,
		id,
	)
}

func (repository *groupingsParticipationRepository) FindByParticipant(
	ctx context.Context,
	participant participatables.Participatable,
) ([]domain_participations.Participation, error) {
	rows, err := repository.pool.Query(
		ctx,
		`
		SELECT
			identifier,
			participant_id,
			participant_kind,
			target_id,
			target_kind,
			post_count,
			total_post_count,
			percentage,
			detected_on
		FROM groupings_participations
		WHERE participant_id = $1
		AND participant_kind = $2
		ORDER BY identifier
		`,
		participant.Identifier(),
		string(participant.ParticipationKind()),
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return repository.scanMany(rows)
}

func (repository *groupingsParticipationRepository) FindByTarget(
	ctx context.Context,
	target participatables.Participatable,
) ([]domain_participations.Participation, error) {
	rows, err := repository.pool.Query(
		ctx,
		`
		SELECT
			identifier,
			participant_id,
			participant_kind,
			target_id,
			target_kind,
			post_count,
			total_post_count,
			percentage,
			detected_on
		FROM groupings_participations
		WHERE target_id = $1
		AND target_kind = $2
		ORDER BY identifier
		`,
		target.Identifier(),
		string(target.ParticipationKind()),
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return repository.scanMany(rows)
}

func (repository *groupingsParticipationRepository) FindBetween(
	ctx context.Context,
	participant participatables.Participatable,
	target participatables.Participatable,
) (domain_participations.Participation, error) {
	return repository.findOne(
		ctx,
		`
		SELECT
			identifier,
			participant_id,
			participant_kind,
			target_id,
			target_kind,
			post_count,
			total_post_count,
			percentage,
			detected_on
		FROM groupings_participations
		WHERE participant_id = $1
		AND participant_kind = $2
		AND target_id = $3
		AND target_kind = $4
		`,
		participant.Identifier(),
		string(participant.ParticipationKind()),
		target.Identifier(),
		string(target.ParticipationKind()),
	)
}

func (repository *groupingsParticipationRepository) findOne(
	ctx context.Context,
	query string,
	args ...any,
) (domain_participations.Participation, error) {
	row := repository.pool.QueryRow(ctx, query, args...)

	participation, err := repository.scanOne(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return participation, nil
}

func (repository *groupingsParticipationRepository) scanMany(
	rows pgx.Rows,
) ([]domain_participations.Participation, error) {
	out := []domain_participations.Participation{}

	for rows.Next() {
		participation, err := repository.scanOne(rows)
		if err != nil {
			return nil, err
		}

		out = append(out, participation)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (repository *groupingsParticipationRepository) scanOne(
	row pgx.Row,
) (domain_participations.Participation, error) {
	var input domain_participations.ParticipationInput

	var participantID uuid.UUID
	var participantKind string
	var targetID uuid.UUID
	var targetKind string

	err := row.Scan(
		&input.Identifier,
		&participantID,
		&participantKind,
		&targetID,
		&targetKind,
		&input.PostCount,
		&input.TotalPostCount,
		&input.Percentage,
		&input.DetectedOn,
	)
	if err != nil {
		return nil, err
	}

	input.Participant = participatables.NewMockParticipatable(
		participantID,
		participatables.Kind(participantKind),
	)

	input.Target = participatables.NewMockParticipatable(
		targetID,
		participatables.Kind(targetKind),
	)

	return repository.adapter.ToDomain(input)
}
