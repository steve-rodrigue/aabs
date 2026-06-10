package postgresql

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/relatables"
	domain_relationships "github.com/steve-rodrigue/aabs/services/saas/domain/relationships"
)

type relationshipRepository struct {
	pool       *pgxpool.Pool
	adapter    domain_relationships.Adapter
	relatables relatables.Adapter
}

func (repository *relationshipRepository) Save(
	ctx context.Context,
	relationship domain_relationships.Relationship,
) error {
	_, err := repository.pool.Exec(
		ctx,
		`
		INSERT INTO relationships (
			identifier,
			source_id,
			source_kind,
			target_id,
			target_kind,
			similarity,
			created_on
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (identifier)
		DO UPDATE SET
			source_id = EXCLUDED.source_id,
			source_kind = EXCLUDED.source_kind,
			target_id = EXCLUDED.target_id,
			target_kind = EXCLUDED.target_kind,
			similarity = EXCLUDED.similarity
		`,
		relationship.Identifier(),
		relationship.Source().Identifier(),
		string(relationship.Source().RelationshipKind()),
		relationship.Target().Identifier(),
		string(relationship.Target().RelationshipKind()),
		relationship.Similarity(),
		relationship.CreatedOn(),
	)

	return err
}

func (repository *relationshipRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (domain_relationships.Relationship, error) {
	return repository.findOne(
		ctx,
		`
		SELECT
			identifier,
			source_id,
			source_kind,
			target_id,
			target_kind,
			similarity,
			created_on
		FROM relationships
		WHERE identifier = $1
		`,
		id,
	)
}

func (repository *relationshipRepository) Find(
	ctx context.Context,
	index int,
	amount int,
) ([]domain_relationships.Relationship, error) {
	rows, err := repository.pool.Query(
		ctx,
		`
		SELECT
			identifier,
			source_id,
			source_kind,
			target_id,
			target_kind,
			similarity,
			created_on
		FROM relationships
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

func (repository *relationshipRepository) FindAfter(
	ctx context.Context,
	cursor uuid.UUID,
	amount int,
) ([]domain_relationships.Relationship, error) {
	if cursor == uuid.Nil {
		return repository.Find(ctx, 0, amount)
	}

	rows, err := repository.pool.Query(
		ctx,
		`
		SELECT
			identifier,
			source_id,
			source_kind,
			target_id,
			target_kind,
			similarity,
			created_on
		FROM relationships
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

func (repository *relationshipRepository) Count(
	ctx context.Context,
) (int64, error) {
	var count int64

	err := repository.pool.QueryRow(
		ctx,
		`
		SELECT COUNT(*)
		FROM relationships
		`,
	).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (repository *relationshipRepository) FindBySourceID(
	ctx context.Context,
	source uuid.UUID,
) ([]domain_relationships.Relationship, error) {
	rows, err := repository.pool.Query(
		ctx,
		`
		SELECT
			identifier,
			source_id,
			source_kind,
			target_id,
			target_kind,
			similarity,
			created_on
		FROM relationships
		WHERE source_id = $1
		ORDER BY identifier
		`,
		source,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return repository.scanMany(rows)
}

func (repository *relationshipRepository) FindByTargetID(
	ctx context.Context,
	target uuid.UUID,
) ([]domain_relationships.Relationship, error) {
	rows, err := repository.pool.Query(
		ctx,
		`
		SELECT
			identifier,
			source_id,
			source_kind,
			target_id,
			target_kind,
			similarity,
			created_on
		FROM relationships
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

func (repository *relationshipRepository) FindBySource(
	ctx context.Context,
	source relatables.Relatable,
) ([]domain_relationships.Relationship, error) {
	rows, err := repository.pool.Query(
		ctx,
		`
		SELECT
			identifier,
			source_id,
			source_kind,
			target_id,
			target_kind,
			similarity,
			created_on
		FROM relationships
		WHERE source_id = $1
		AND source_kind = $2
		ORDER BY identifier
		`,
		source.Identifier(),
		string(source.RelationshipKind()),
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return repository.scanMany(rows)
}

func (repository *relationshipRepository) FindByTarget(
	ctx context.Context,
	target relatables.Relatable,
) ([]domain_relationships.Relationship, error) {
	rows, err := repository.pool.Query(
		ctx,
		`
		SELECT
			identifier,
			source_id,
			source_kind,
			target_id,
			target_kind,
			similarity,
			created_on
		FROM relationships
		WHERE target_id = $1
		AND target_kind = $2
		ORDER BY identifier
		`,
		target.Identifier(),
		string(target.RelationshipKind()),
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return repository.scanMany(rows)
}

func (repository *relationshipRepository) FindBetween(
	ctx context.Context,
	source relatables.Relatable,
	target relatables.Relatable,
) (domain_relationships.Relationship, error) {
	return repository.findOne(
		ctx,
		`
		SELECT
			identifier,
			source_id,
			source_kind,
			target_id,
			target_kind,
			similarity,
			created_on
		FROM relationships
		WHERE source_id = $1
		AND source_kind = $2
		AND target_id = $3
		AND target_kind = $4
		`,
		source.Identifier(),
		string(source.RelationshipKind()),
		target.Identifier(),
		string(target.RelationshipKind()),
	)
}

func (repository *relationshipRepository) findOne(
	ctx context.Context,
	query string,
	args ...any,
) (domain_relationships.Relationship, error) {
	row := repository.pool.QueryRow(ctx, query, args...)

	relationship, err := repository.scanOne(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return relationship, nil
}

func (repository *relationshipRepository) scanMany(
	rows pgx.Rows,
) ([]domain_relationships.Relationship, error) {
	out := []domain_relationships.Relationship{}

	for rows.Next() {
		relationship, err := repository.scanOne(rows)
		if err != nil {
			return nil, err
		}

		out = append(out, relationship)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (repository *relationshipRepository) scanOne(
	row pgx.Row,
) (domain_relationships.Relationship, error) {
	var input domain_relationships.RelationshipInput

	var sourceID uuid.UUID
	var sourceKind string

	var targetID uuid.UUID
	var targetKind string

	err := row.Scan(
		&input.Identifier,
		&sourceID,
		&sourceKind,
		&targetID,
		&targetKind,
		&input.Similarity,
		&input.CreatedOn,
	)
	if err != nil {
		return nil, err
	}

	source, err := repository.relatables.ToDomain(
		relatables.RelatableInput{
			Identifier:       sourceID,
			RelationshipKind: relatables.Kind(sourceKind),
		},
	)
	if err != nil {
		return nil, err
	}

	target, err := repository.relatables.ToDomain(
		relatables.RelatableInput{
			Identifier:       targetID,
			RelationshipKind: relatables.Kind(targetKind),
		},
	)
	if err != nil {
		return nil, err
	}

	input.Source = source
	input.Target = target

	return repository.adapter.ToDomain(input)
}
