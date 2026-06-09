package postgresql

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	domain_clusters "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters/clusterables"
)

type groupingsClusterRepository struct {
	pool    *pgxpool.Pool
	adapter domain_clusters.Adapter
}

func (repository *groupingsClusterRepository) Save(
	ctx context.Context,
	cluster domain_clusters.Cluster,
) error {
	_, err := repository.pool.Exec(
		ctx,
		`
		INSERT INTO groupings_clusters (
			identifier,
			target_id,
			target_kind,
			member_ids,
			member_kind,
			confidence_score,
			centroid,
			created_on
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (identifier)
		DO UPDATE SET
			target_id = EXCLUDED.target_id,
			target_kind = EXCLUDED.target_kind,
			member_ids = EXCLUDED.member_ids,
			member_kind = EXCLUDED.member_kind,
			confidence_score = EXCLUDED.confidence_score,
			centroid = EXCLUDED.centroid
		`,
		cluster.Identifier(),
		cluster.Target().Identifier(),
		string(cluster.Target().ClusterKind()),
		cluster.MemberIDs(),
		string(cluster.MemberKind()),
		cluster.ConfidenceScore(),
		cluster.Centroid(),
		cluster.CreatedOn(),
	)

	return err
}

func (repository *groupingsClusterRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (domain_clusters.Cluster, error) {
	return repository.findOne(
		ctx,
		`
		SELECT
			identifier,
			target_id,
			target_kind,
			member_ids,
			member_kind,
			confidence_score,
			centroid,
			created_on
		FROM groupings_clusters
		WHERE identifier = $1
		`,
		id,
	)
}

func (repository *groupingsClusterRepository) FindByTarget(
	ctx context.Context,
	target uuid.UUID,
) ([]domain_clusters.Cluster, error) {
	rows, err := repository.pool.Query(
		ctx,
		`
		SELECT
			identifier,
			target_id,
			target_kind,
			member_ids,
			member_kind,
			confidence_score,
			centroid,
			created_on
		FROM groupings_clusters
		WHERE target_id = $1
		ORDER BY identifier
		`,
		target,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return repository.scanMany(rows)
}

func (repository *groupingsClusterRepository) FindByMember(
	ctx context.Context,
	member uuid.UUID,
) ([]domain_clusters.Cluster, error) {
	rows, err := repository.pool.Query(
		ctx,
		`
		SELECT
			identifier,
			target_id,
			target_kind,
			member_ids,
			member_kind,
			confidence_score,
			centroid,
			created_on
		FROM groupings_clusters
		WHERE $1 = ANY(member_ids)
		ORDER BY identifier
		`,
		member,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return repository.scanMany(rows)
}

func (repository *groupingsClusterRepository) Find(
	ctx context.Context,
	index int,
	amount int,
) ([]domain_clusters.Cluster, error) {
	rows, err := repository.pool.Query(
		ctx,
		`
		SELECT
			identifier,
			target_id,
			target_kind,
			member_ids,
			member_kind,
			confidence_score,
			centroid,
			created_on
		FROM groupings_clusters
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

	return repository.scanMany(rows)
}

func (repository *groupingsClusterRepository) FindAfter(
	ctx context.Context,
	cursor uuid.UUID,
	amount int,
) ([]domain_clusters.Cluster, error) {
	if cursor == uuid.Nil {
		return repository.Find(ctx, 0, amount)
	}

	rows, err := repository.pool.Query(
		ctx,
		`
		SELECT
			identifier,
			target_id,
			target_kind,
			member_ids,
			member_kind,
			confidence_score,
			centroid,
			created_on
		FROM groupings_clusters
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

	return repository.scanMany(rows)
}

func (repository *groupingsClusterRepository) Count(
	ctx context.Context,
) (int64, error) {
	var count int64

	err := repository.pool.QueryRow(
		ctx,
		`
		SELECT COUNT(*)
		FROM groupings_clusters
		`,
	).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (repository *groupingsClusterRepository) findOne(
	ctx context.Context,
	query string,
	args ...any,
) (domain_clusters.Cluster, error) {
	row := repository.pool.QueryRow(ctx, query, args...)

	cluster, err := repository.scanOne(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return cluster, nil
}

func (repository *groupingsClusterRepository) scanMany(
	rows pgx.Rows,
) ([]domain_clusters.Cluster, error) {
	out := []domain_clusters.Cluster{}

	for rows.Next() {
		cluster, err := repository.scanOne(rows)
		if err != nil {
			return nil, err
		}

		out = append(out, cluster)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (repository *groupingsClusterRepository) scanOne(
	row pgx.Row,
) (domain_clusters.Cluster, error) {
	var input domain_clusters.ClusterInput

	var targetID uuid.UUID
	var targetKind string
	var memberKind string

	err := row.Scan(
		&input.Identifier,
		&targetID,
		&targetKind,
		&input.MemberIDs,
		&memberKind,
		&input.ConfidenceScore,
		&input.Centroid,
		&input.CreatedOn,
	)
	if err != nil {
		return nil, err
	}

	input.Target = clusterables.ClusterableInput{
		Identifier:  targetID,
		ClusterKind: clusterables.Kind(targetKind),
	}

	input.MemberKind = clusterables.Kind(memberKind)

	return repository.adapter.ToDomain(input)
}
