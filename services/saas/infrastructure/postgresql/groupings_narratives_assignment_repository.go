package postgresql

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	domain_campaigns "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/campaigns"
	domain_narratives "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/narratives"
	domain_assignments "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/narratives/assignments"
)

type groupingsAssignmentRepository struct {
	pool *pgxpool.Pool

	adapter    domain_assignments.Adapter
	narratives domain_narratives.Repository
	campaigns  domain_campaigns.Repository
}

func (repository *groupingsAssignmentRepository) Save(
	ctx context.Context,
	assignment domain_assignments.Assignment,
) error {
	_, err := repository.pool.Exec(
		ctx,
		`
		INSERT INTO groupings_assignments (
			identifier,
			narrative_id,
			campaign_id,
			confidence,
			assigned_on
		)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (identifier)
		DO UPDATE SET
			narrative_id = EXCLUDED.narrative_id,
			campaign_id = EXCLUDED.campaign_id,
			confidence = EXCLUDED.confidence,
			assigned_on = EXCLUDED.assigned_on
		`,
		assignment.Identifier(),
		assignment.Narrative().Identifier(),
		assignment.Campaign().Identifier(),
		assignment.Confidence(),
		assignment.AssignedOn(),
	)

	return err
}

func (repository *groupingsAssignmentRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (domain_assignments.Assignment, error) {
	return repository.findOne(
		ctx,
		`
		SELECT
			identifier,
			narrative_id,
			campaign_id,
			confidence,
			assigned_on
		FROM groupings_assignments
		WHERE identifier = $1
		`,
		id,
	)
}

func (repository *groupingsAssignmentRepository) FindByNarrative(
	ctx context.Context,
	narrative uuid.UUID,
) ([]domain_assignments.Assignment, error) {
	rows, err := repository.pool.Query(
		ctx,
		`
		SELECT
			identifier,
			narrative_id,
			campaign_id,
			confidence,
			assigned_on
		FROM groupings_assignments
		WHERE narrative_id = $1
		ORDER BY identifier
		`,
		narrative,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return repository.scanMany(ctx, rows)
}

func (repository *groupingsAssignmentRepository) FindByCampaign(
	ctx context.Context,
	campaign uuid.UUID,
) ([]domain_assignments.Assignment, error) {
	rows, err := repository.pool.Query(
		ctx,
		`
		SELECT
			identifier,
			narrative_id,
			campaign_id,
			confidence,
			assigned_on
		FROM groupings_assignments
		WHERE campaign_id = $1
		ORDER BY identifier
		`,
		campaign,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return repository.scanMany(ctx, rows)
}

func (repository *groupingsAssignmentRepository) FindBetween(
	ctx context.Context,
	narrative uuid.UUID,
	campaign uuid.UUID,
) (domain_assignments.Assignment, error) {
	return repository.findOne(
		ctx,
		`
		SELECT
			identifier,
			narrative_id,
			campaign_id,
			confidence,
			assigned_on
		FROM groupings_assignments
		WHERE narrative_id = $1
		AND campaign_id = $2
		`,
		narrative,
		campaign,
	)
}

func (repository *groupingsAssignmentRepository) Find(
	ctx context.Context,
	index int,
	amount int,
) ([]domain_assignments.Assignment, error) {
	rows, err := repository.pool.Query(
		ctx,
		`
		SELECT
			identifier,
			narrative_id,
			campaign_id,
			confidence,
			assigned_on
		FROM groupings_assignments
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

func (repository *groupingsAssignmentRepository) FindAfter(
	ctx context.Context,
	cursor uuid.UUID,
	amount int,
) ([]domain_assignments.Assignment, error) {
	if cursor == uuid.Nil {
		return repository.Find(ctx, 0, amount)
	}

	rows, err := repository.pool.Query(
		ctx,
		`
		SELECT
			identifier,
			narrative_id,
			campaign_id,
			confidence,
			assigned_on
		FROM groupings_assignments
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

func (repository *groupingsAssignmentRepository) Count(
	ctx context.Context,
) (int64, error) {
	var count int64

	err := repository.pool.QueryRow(
		ctx,
		`
		SELECT COUNT(*)
		FROM groupings_assignments
		`,
	).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (repository *groupingsAssignmentRepository) findOne(
	ctx context.Context,
	query string,
	args ...any,
) (domain_assignments.Assignment, error) {
	row := repository.pool.QueryRow(
		ctx,
		query,
		args...,
	)

	assignment, err := repository.scanOne(ctx, row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return assignment, nil
}

func (repository *groupingsAssignmentRepository) scanMany(
	ctx context.Context,
	rows pgx.Rows,
) ([]domain_assignments.Assignment, error) {
	out := []domain_assignments.Assignment{}

	for rows.Next() {
		assignment, err := repository.scanOne(ctx, rows)
		if err != nil {
			return nil, err
		}

		out = append(out, assignment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (repository *groupingsAssignmentRepository) scanOne(
	ctx context.Context,
	row pgx.Row,
) (domain_assignments.Assignment, error) {
	var input domain_assignments.AssignmentInput

	var narrativeID uuid.UUID
	var campaignID uuid.UUID

	err := row.Scan(
		&input.Identifier,
		&narrativeID,
		&campaignID,
		&input.Confidence,
		&input.AssignedOn,
	)
	if err != nil {
		return nil, err
	}

	narrative, err := repository.narratives.FindByID(
		ctx,
		narrativeID,
	)
	if err != nil {
		return nil, err
	}

	campaign, err := repository.campaigns.FindByID(
		ctx,
		campaignID,
	)
	if err != nil {
		return nil, err
	}

	input.Narrative = narrative
	input.Campaign = campaign

	return repository.adapter.ToDomain(input)
}
