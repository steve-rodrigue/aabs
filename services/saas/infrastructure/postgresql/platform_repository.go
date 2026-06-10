package postgresql

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/participatables"
	domain_platforms "github.com/steve-rodrigue/aabs/services/saas/domain/entities/platforms"
)

type platformRepository struct {
	pool    *pgxpool.Pool
	adapter domain_platforms.Adapter
}

func (repository *platformRepository) Save(
	ctx context.Context,
	platform domain_platforms.Platform,
) error {
	_, err := repository.pool.Exec(
		ctx,
		`
		INSERT INTO platforms (
			identifier,
			participation_kind,
			name,
			handle,
			base_url,
			created_on
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (identifier)
		DO UPDATE SET
			participation_kind = EXCLUDED.participation_kind,
			name = EXCLUDED.name,
			handle = EXCLUDED.handle,
			base_url = EXCLUDED.base_url
		`,
		platform.Identifier(),
		string(platform.ParticipationKind()),
		platform.Name(),
		platform.Handle(),
		platform.BaseURL(),
		platform.CreatedOn(),
	)

	return err
}

func (repository *platformRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (domain_platforms.Platform, error) {
	return repository.findOne(
		ctx,
		`
		SELECT
			identifier,
			participation_kind,
			name,
			handle,
			base_url,
			created_on
		FROM platforms
		WHERE identifier = $1
		`,
		id,
	)
}

func (repository *platformRepository) FindByHandle(
	ctx context.Context,
	handle string,
) (domain_platforms.Platform, error) {
	return repository.findOne(
		ctx,
		`
		SELECT
			identifier,
			participation_kind,
			name,
			handle,
			base_url,
			created_on
		FROM platforms
		WHERE handle = $1
		`,
		handle,
	)
}

func (repository *platformRepository) FindByName(
	ctx context.Context,
	name string,
) (domain_platforms.Platform, error) {
	return repository.findOne(
		ctx,
		`
		SELECT
			identifier,
			participation_kind,
			name,
			handle,
			base_url,
			created_on
		FROM platforms
		WHERE name = $1
		`,
		name,
	)
}

func (repository *platformRepository) Find(
	ctx context.Context,
	index int,
	amount int,
) ([]domain_platforms.Platform, error) {
	rows, err := repository.pool.Query(
		ctx,
		`
		SELECT
			identifier,
			participation_kind,
			name,
			handle,
			base_url,
			created_on
		FROM platforms
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

func (repository *platformRepository) FindAfter(
	ctx context.Context,
	cursor uuid.UUID,
	amount int,
) ([]domain_platforms.Platform, error) {
	if cursor == uuid.Nil {
		rows, err := repository.pool.Query(
			ctx,
			`
			SELECT
				identifier,
				participation_kind,
				name,
				handle,
				base_url,
				created_on
			FROM platforms
			ORDER BY identifier
			LIMIT $1
			`,
			amount,
		)
		if err != nil {
			return nil, err
		}

		defer rows.Close()

		return repository.scanMany(rows)
	}

	rows, err := repository.pool.Query(
		ctx,
		`
		SELECT
			identifier,
			participation_kind,
			name,
			handle,
			base_url,
			created_on
		FROM platforms
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

func (repository *platformRepository) Count(
	ctx context.Context,
) (int64, error) {
	var count int64

	err := repository.pool.QueryRow(
		ctx,
		`
		SELECT COUNT(*)
		FROM platforms
		`,
	).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (repository *platformRepository) findOne(
	ctx context.Context,
	query string,
	args ...any,
) (domain_platforms.Platform, error) {
	row := repository.pool.QueryRow(
		ctx,
		query,
		args...,
	)

	platform, err := repository.scanOne(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return platform, nil
}

func (repository *platformRepository) scanMany(
	rows pgx.Rows,
) ([]domain_platforms.Platform, error) {
	out := []domain_platforms.Platform{}

	for rows.Next() {
		platform, err := repository.scanOne(rows)
		if err != nil {
			return nil, err
		}

		out = append(out, platform)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (repository *platformRepository) scanOne(
	row pgx.Row,
) (domain_platforms.Platform, error) {
	var input domain_platforms.PlatformInput
	var participationKind string

	err := row.Scan(
		&input.Identifier,
		&participationKind,
		&input.Name,
		&input.Handle,
		&input.BaseURL,
		&input.CreatedOn,
	)
	if err != nil {
		return nil, err
	}

	input.ParticipationKind = participatables.Kind(participationKind)

	return repository.adapter.ToDomain(input)
}
