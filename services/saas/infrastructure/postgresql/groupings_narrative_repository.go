package postgresql

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/participatables"
	domain_clusters "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters"
	domain_narratives "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/narratives"
)

type groupingsNarrativeRepository struct {
	pool     *pgxpool.Pool
	adapter  domain_narratives.Adapter
	clusters domain_clusters.Repository
}

func (repository *groupingsNarrativeRepository) Save(
	ctx context.Context,
	narrative domain_narratives.Narrative,
) error {
	_, err := repository.pool.Exec(
		ctx,
		`
		INSERT INTO groupings_narratives (
			identifier,
			participation_kind,
			cluster_id,
			name,
			description,
			created_on
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (identifier)
		DO UPDATE SET
			participation_kind = EXCLUDED.participation_kind,
			cluster_id = EXCLUDED.cluster_id,
			name = EXCLUDED.name,
			description = EXCLUDED.description
		`,
		narrative.Identifier(),
		string(narrative.ParticipationKind()),
		narrative.Cluster().Identifier(),
		narrative.Name(),
		narrative.Description(),
		narrative.CreatedOn(),
	)

	return err
}

func (repository *groupingsNarrativeRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (domain_narratives.Narrative, error) {
	return repository.findOne(
		ctx,
		`
		SELECT
			identifier,
			participation_kind,
			cluster_id,
			name,
			description,
			created_on
		FROM groupings_narratives
		WHERE identifier = $1
		`,
		id,
	)
}

func (repository *groupingsNarrativeRepository) FindByName(
	ctx context.Context,
	name string,
) (domain_narratives.Narrative, error) {
	return repository.findOne(
		ctx,
		`
		SELECT
			identifier,
			participation_kind,
			cluster_id,
			name,
			description,
			created_on
		FROM groupings_narratives
		WHERE name = $1
		`,
		name,
	)
}

func (repository *groupingsNarrativeRepository) Find(
	ctx context.Context,
	index int,
	amount int,
) ([]domain_narratives.Narrative, error) {
	rows, err := repository.pool.Query(
		ctx,
		`
		SELECT
			identifier,
			participation_kind,
			cluster_id,
			name,
			description,
			created_on
		FROM groupings_narratives
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

func (repository *groupingsNarrativeRepository) FindAfter(
	ctx context.Context,
	cursor uuid.UUID,
	amount int,
) ([]domain_narratives.Narrative, error) {
	if cursor == uuid.Nil {
		return repository.Find(ctx, 0, amount)
	}

	rows, err := repository.pool.Query(
		ctx,
		`
		SELECT
			identifier,
			participation_kind,
			cluster_id,
			name,
			description,
			created_on
		FROM groupings_narratives
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

func (repository *groupingsNarrativeRepository) Count(
	ctx context.Context,
) (int64, error) {
	var count int64

	err := repository.pool.QueryRow(
		ctx,
		`
		SELECT COUNT(*)
		FROM groupings_narratives
		`,
	).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (repository *groupingsNarrativeRepository) findOne(
	ctx context.Context,
	query string,
	args ...any,
) (domain_narratives.Narrative, error) {
	row := repository.pool.QueryRow(
		ctx,
		query,
		args...,
	)

	narrative, err := repository.scanOne(ctx, row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return narrative, nil
}

func (repository *groupingsNarrativeRepository) scanMany(
	ctx context.Context,
	rows pgx.Rows,
) ([]domain_narratives.Narrative, error) {
	out := []domain_narratives.Narrative{}

	for rows.Next() {
		narrative, err := repository.scanOne(ctx, rows)
		if err != nil {
			return nil, err
		}

		out = append(out, narrative)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (repository *groupingsNarrativeRepository) scanOne(
	ctx context.Context,
	row pgx.Row,
) (domain_narratives.Narrative, error) {
	var input domain_narratives.NarrativeInput

	var participationKind string
	var clusterID uuid.UUID

	err := row.Scan(
		&input.Identifier,
		&participationKind,
		&clusterID,
		&input.Name,
		&input.Description,
		&input.CreatedOn,
	)
	if err != nil {
		return nil, err
	}

	cluster, err := repository.clusters.FindByID(
		ctx,
		clusterID,
	)
	if err != nil {
		return nil, err
	}

	input.ParticipationKind = participatables.Kind(participationKind)
	input.Cluster = cluster

	return repository.adapter.ToDomain(input)
}
