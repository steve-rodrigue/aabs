package postgresql

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/participatables"
	domain_dirty_participation "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/dirty"
)

type groupingsParticipationsDirtyRepository struct {
	pool    *pgxpool.Pool
	adapter domain_dirty_participation.Adapter
}

func (repository *groupingsParticipationsDirtyRepository) Save(
	ctx context.Context,
	dirty domain_dirty_participation.Dirty,
) error {
	_, err := repository.pool.Exec(
		ctx,
		`
		INSERT INTO groupings_participations_dirty (
			identifier,
			participant_id,
			participant_kind,
			target_id,
			target_kind,
			marked_on
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (
			participant_id,
			participant_kind,
			target_id,
			target_kind
		)
		DO UPDATE SET
			marked_on = EXCLUDED.marked_on
		`,
		dirty.Identifier(),
		dirty.Participant().Identifier(),
		string(dirty.Participant().ParticipationKind()),
		dirty.Target().Identifier(),
		string(dirty.Target().ParticipationKind()),
		dirty.MarkedOn(),
	)

	return err
}

func (repository *groupingsParticipationsDirtyRepository) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	_, err := repository.pool.Exec(
		ctx,
		`
		DELETE FROM groupings_participations_dirty
		WHERE identifier = $1
		`,
		id,
	)

	return err
}

func (repository *groupingsParticipationsDirtyRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (domain_dirty_participation.Dirty, error) {
	return repository.findOne(
		ctx,
		`
		SELECT
			identifier,
			participant_id,
			participant_kind,
			target_id,
			target_kind,
			marked_on
		FROM groupings_participations_dirty
		WHERE identifier = $1
		`,
		id,
	)
}

func (repository *groupingsParticipationsDirtyRepository) FindBetween(
	ctx context.Context,
	participant participatables.Participatable,
	target participatables.Participatable,
) (domain_dirty_participation.Dirty, error) {
	return repository.findOne(
		ctx,
		`
		SELECT
			identifier,
			participant_id,
			participant_kind,
			target_id,
			target_kind,
			marked_on
		FROM groupings_participations_dirty
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

func (repository *groupingsParticipationsDirtyRepository) Find(
	ctx context.Context,
	index int,
	amount int,
) ([]domain_dirty_participation.Dirty, error) {
	rows, err := repository.pool.Query(
		ctx,
		`
		SELECT
			identifier,
			participant_id,
			participant_kind,
			target_id,
			target_kind,
			marked_on
		FROM groupings_participations_dirty
		ORDER BY identifier
		OFFSET $1
		LIMIT $2
		`,
		index,
		amount,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return repository.scanMany(ctx, rows)
}

func (repository *groupingsParticipationsDirtyRepository) FindAfter(
	ctx context.Context,
	cursor uuid.UUID,
	amount int,
) ([]domain_dirty_participation.Dirty, error) {
	if cursor == uuid.Nil {
		return repository.Find(ctx, 0, amount)
	}

	rows, err := repository.pool.Query(
		ctx,
		`
		SELECT
			identifier,
			participant_id,
			participant_kind,
			target_id,
			target_kind,
			marked_on
		FROM groupings_participations_dirty
		WHERE identifier > $1
		ORDER BY identifier
		LIMIT $2
		`,
		cursor,
		amount,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return repository.scanMany(ctx, rows)
}

func (repository *groupingsParticipationsDirtyRepository) Count(
	ctx context.Context,
) (int64, error) {
	var count int64

	err := repository.pool.QueryRow(
		ctx,
		`
		SELECT COUNT(*)
		FROM groupings_participations_dirty
		`,
	).Scan(&count)

	return count, err
}

func (repository *groupingsParticipationsDirtyRepository) findOne(
	ctx context.Context,
	query string,
	args ...any,
) (domain_dirty_participation.Dirty, error) {
	row := repository.pool.QueryRow(ctx, query, args...)

	dirty, err := repository.scanOne(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return dirty, nil
}

func (repository *groupingsParticipationsDirtyRepository) scanMany(
	ctx context.Context,
	rows pgx.Rows,
) ([]domain_dirty_participation.Dirty, error) {
	out := []domain_dirty_participation.Dirty{}

	for rows.Next() {
		dirty, err := repository.scanOne(rows)
		if err != nil {
			return nil, err
		}

		out = append(out, dirty)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (repository *groupingsParticipationsDirtyRepository) scanOne(
	row pgx.Row,
) (domain_dirty_participation.Dirty, error) {
	var input domain_dirty_participation.DirtyInput

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
		&input.MarkedOn,
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
