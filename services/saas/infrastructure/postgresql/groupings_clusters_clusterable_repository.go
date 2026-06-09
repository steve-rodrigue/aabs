package postgresql

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters/clusterables"
)

type groupingsClustersClusterableRepository struct {
	pool    *pgxpool.Pool
	adapter clusterables.Adapter
}

func (repository *groupingsClustersClusterableRepository) FindByKind(
	ctx context.Context,
	kind clusterables.Kind,
	index int,
	amount int,
) ([]clusterables.Clusterable, error) {
	rows, err := repository.pool.Query(
		ctx,
		`
		SELECT
			identifier,
			kind
		FROM groupings_clusterables
		WHERE kind = $1
		ORDER BY identifier
		OFFSET $2
		LIMIT $3
		`,
		string(kind),
		index,
		amount,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return repository.scanMany(rows)
}

func (repository *groupingsClustersClusterableRepository) FindByKindAfter(
	ctx context.Context,
	kind clusterables.Kind,
	cursor uuid.UUID,
	amount int,
) ([]clusterables.Clusterable, error) {
	if cursor == uuid.Nil {
		return repository.FindByKind(ctx, kind, 0, amount)
	}

	rows, err := repository.pool.Query(
		ctx,
		`
		SELECT
			identifier,
			kind
		FROM groupings_clusterables
		WHERE kind = $1
		AND identifier > $2
		ORDER BY identifier
		LIMIT $3
		`,
		string(kind),
		cursor,
		amount,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return repository.scanMany(rows)
}

func (repository *groupingsClustersClusterableRepository) CountByKind(
	ctx context.Context,
	kind clusterables.Kind,
) (int64, error) {
	var count int64

	err := repository.pool.QueryRow(
		ctx,
		`
		SELECT COUNT(*)
		FROM groupings_clusterables
		WHERE kind = $1
		`,
		string(kind),
	).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (repository *groupingsClustersClusterableRepository) scanMany(
	rows pgx.Rows,
) ([]clusterables.Clusterable, error) {
	out := []clusterables.Clusterable{}

	for rows.Next() {
		clusterable, err := repository.scanOne(rows)
		if err != nil {
			return nil, err
		}

		out = append(out, clusterable)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (repository *groupingsClustersClusterableRepository) scanOne(
	row pgx.Row,
) (clusterables.Clusterable, error) {
	var id uuid.UUID
	var kind string

	if err := row.Scan(&id, &kind); err != nil {
		return nil, err
	}

	return repository.adapter.ToDomain(
		clusterables.ClusterableInput{
			Identifier:  id,
			ClusterKind: clusterables.Kind(kind),
		},
	)
}
