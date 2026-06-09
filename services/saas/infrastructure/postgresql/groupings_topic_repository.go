package postgresql

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	domain_clusters "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters"
	domain_topics "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/topics"
)

type groupingsTopicRepository struct {
	pool     *pgxpool.Pool
	adapter  domain_topics.Adapter
	clusters domain_clusters.Repository
}

func (repository *groupingsTopicRepository) Save(
	ctx context.Context,
	topic domain_topics.Topic,
) error {
	var parentID *uuid.UUID

	if topic.HasParent() {
		id := topic.Parent().Identifier()
		parentID = &id
	}

	_, err := repository.pool.Exec(
		ctx,
		`
		INSERT INTO groupings_topics (
			identifier,
			cluster_id,
			name,
			description,
			parent_id,
			created_on
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (identifier)
		DO UPDATE SET
			cluster_id = EXCLUDED.cluster_id,
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			parent_id = EXCLUDED.parent_id
		`,
		topic.Identifier(),
		topic.Cluster().Identifier(),
		topic.Name(),
		topic.Description(),
		parentID,
		topic.CreatedOn(),
	)

	return err
}

func (repository *groupingsTopicRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (domain_topics.Topic, error) {
	return repository.findOne(
		ctx,
		`
		SELECT
			identifier,
			cluster_id,
			name,
			description,
			parent_id,
			created_on
		FROM groupings_topics
		WHERE identifier = $1
		`,
		id,
	)
}

func (repository *groupingsTopicRepository) FindByName(
	ctx context.Context,
	name string,
) (domain_topics.Topic, error) {
	return repository.findOne(
		ctx,
		`
		SELECT
			identifier,
			cluster_id,
			name,
			description,
			parent_id,
			created_on
		FROM groupings_topics
		WHERE name = $1
		`,
		name,
	)
}

func (repository *groupingsTopicRepository) Find(
	ctx context.Context,
	index int,
	amount int,
) ([]domain_topics.Topic, error) {
	rows, err := repository.pool.Query(
		ctx,
		`
		SELECT
			identifier,
			cluster_id,
			name,
			description,
			parent_id,
			created_on
		FROM groupings_topics
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

func (repository *groupingsTopicRepository) FindAfter(
	ctx context.Context,
	cursor uuid.UUID,
	amount int,
) ([]domain_topics.Topic, error) {
	if cursor == uuid.Nil {
		return repository.Find(ctx, 0, amount)
	}

	rows, err := repository.pool.Query(
		ctx,
		`
		SELECT
			identifier,
			cluster_id,
			name,
			description,
			parent_id,
			created_on
		FROM groupings_topics
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

func (repository *groupingsTopicRepository) Count(
	ctx context.Context,
) (int64, error) {
	var count int64

	err := repository.pool.QueryRow(
		ctx,
		`
		SELECT COUNT(*)
		FROM groupings_topics
		`,
	).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (repository *groupingsTopicRepository) FindChildren(
	ctx context.Context,
	parent uuid.UUID,
) ([]domain_topics.Topic, error) {
	rows, err := repository.pool.Query(
		ctx,
		`
		SELECT
			identifier,
			cluster_id,
			name,
			description,
			parent_id,
			created_on
		FROM groupings_topics
		WHERE parent_id = $1
		ORDER BY identifier
		`,
		parent,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return repository.scanMany(ctx, rows)
}

func (repository *groupingsTopicRepository) FindRoots(
	ctx context.Context,
) ([]domain_topics.Topic, error) {
	rows, err := repository.pool.Query(
		ctx,
		`
		SELECT
			identifier,
			cluster_id,
			name,
			description,
			parent_id,
			created_on
		FROM groupings_topics
		WHERE parent_id IS NULL
		ORDER BY identifier
		`,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return repository.scanMany(ctx, rows)
}

func (repository *groupingsTopicRepository) findOne(
	ctx context.Context,
	query string,
	args ...any,
) (domain_topics.Topic, error) {
	row := repository.pool.QueryRow(ctx, query, args...)

	topic, err := repository.scanOne(ctx, row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return topic, nil
}

func (repository *groupingsTopicRepository) scanMany(
	ctx context.Context,
	rows pgx.Rows,
) ([]domain_topics.Topic, error) {
	out := []domain_topics.Topic{}

	for rows.Next() {
		topic, err := repository.scanOne(ctx, rows)
		if err != nil {
			return nil, err
		}

		out = append(out, topic)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (repository *groupingsTopicRepository) scanOne(
	ctx context.Context,
	row pgx.Row,
) (domain_topics.Topic, error) {
	var input domain_topics.TopicInput

	var clusterID uuid.UUID
	var parentID pgtype.UUID

	err := row.Scan(
		&input.Identifier,
		&clusterID,
		&input.Name,
		&input.Description,
		&parentID,
		&input.CreatedOn,
	)
	if err != nil {
		return nil, err
	}

	cluster, err := repository.clusters.FindByID(ctx, clusterID)
	if err != nil {
		return nil, err
	}

	input.Cluster = cluster

	if parentID.Valid {
		id, err := uuid.FromBytes(parentID.Bytes[:])
		if err != nil {
			return nil, err
		}

		parent, err := repository.FindByID(ctx, id)
		if err != nil {
			return nil, err
		}

		input.Parent = parent
	}

	return repository.adapter.ToDomain(input)
}
