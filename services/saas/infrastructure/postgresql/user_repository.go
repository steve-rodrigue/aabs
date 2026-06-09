package postgresql

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	domain_platforms "github.com/steve-rodrigue/aabs/services/saas/domain/entities/platforms"
	domain_users "github.com/steve-rodrigue/aabs/services/saas/domain/entities/users"
)

type userRepository struct {
	pool      *pgxpool.Pool
	adapter   domain_users.Adapter
	platforms domain_platforms.Repository
}

func (repository *userRepository) Save(
	ctx context.Context,
	user domain_users.User,
) error {
	_, err := repository.pool.Exec(
		ctx,
		`
		INSERT INTO users (
			identifier,
			platform_id,
			external_id,
			handle,
			display_name,
			profile_url,
			created_on
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (identifier)
		DO UPDATE SET
			platform_id = EXCLUDED.platform_id,
			external_id = EXCLUDED.external_id,
			handle = EXCLUDED.handle,
			display_name = EXCLUDED.display_name,
			profile_url = EXCLUDED.profile_url
		`,
		user.Identifier(),
		user.Platform().Identifier(),
		user.ExternalID(),
		user.Handle(),
		user.DisplayName(),
		user.ProfileURL(),
		user.CreatedOn(),
	)

	return err
}

func (repository *userRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (domain_users.User, error) {
	return repository.findOne(
		ctx,
		`
		SELECT
			identifier,
			platform_id,
			external_id,
			handle,
			display_name,
			profile_url,
			created_on
		FROM users
		WHERE identifier = $1
		`,
		id,
	)
}

func (repository *userRepository) FindByPlatformAndExternalID(
	ctx context.Context,
	platform domain_platforms.Platform,
	externalID string,
) (domain_users.User, error) {
	return repository.findOne(
		ctx,
		`
		SELECT
			identifier,
			platform_id,
			external_id,
			handle,
			display_name,
			profile_url,
			created_on
		FROM users
		WHERE platform_id = $1
		AND external_id = $2
		`,
		platform.Identifier(),
		externalID,
	)
}

func (repository *userRepository) FindByPlatformAndHandle(
	ctx context.Context,
	platform domain_platforms.Platform,
	handle string,
) (domain_users.User, error) {
	return repository.findOne(
		ctx,
		`
		SELECT
			identifier,
			platform_id,
			external_id,
			handle,
			display_name,
			profile_url,
			created_on
		FROM users
		WHERE platform_id = $1
		AND handle = $2
		`,
		platform.Identifier(),
		handle,
	)
}

func (repository *userRepository) Find(
	ctx context.Context,
	index int,
	amount int,
) ([]domain_users.User, error) {
	rows, err := repository.pool.Query(
		ctx,
		`
		SELECT
			identifier,
			platform_id,
			external_id,
			handle,
			display_name,
			profile_url,
			created_on
		FROM users
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

func (repository *userRepository) FindAfter(
	ctx context.Context,
	cursor uuid.UUID,
	amount int,
) ([]domain_users.User, error) {
	if cursor == uuid.Nil {
		rows, err := repository.pool.Query(
			ctx,
			`
			SELECT
				identifier,
				platform_id,
				external_id,
				handle,
				display_name,
				profile_url,
				created_on
			FROM users
			ORDER BY identifier
			LIMIT $1
			`,
			amount,
		)
		if err != nil {
			return nil, err
		}

		defer rows.Close()

		return repository.scanMany(ctx, rows)
	}

	rows, err := repository.pool.Query(
		ctx,
		`
		SELECT
			identifier,
			platform_id,
			external_id,
			handle,
			display_name,
			profile_url,
			created_on
		FROM users
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

func (repository *userRepository) Count(
	ctx context.Context,
) (int64, error) {
	var count int64

	err := repository.pool.QueryRow(
		ctx,
		`
		SELECT COUNT(*)
		FROM users
		`,
	).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (repository *userRepository) findOne(
	ctx context.Context,
	query string,
	args ...any,
) (domain_users.User, error) {
	row := repository.pool.QueryRow(ctx, query, args...)

	user, err := repository.scanOne(ctx, row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repository *userRepository) scanMany(
	ctx context.Context,
	rows pgx.Rows,
) ([]domain_users.User, error) {
	out := []domain_users.User{}

	for rows.Next() {
		user, err := repository.scanOne(ctx, rows)
		if err != nil {
			return nil, err
		}

		out = append(out, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (repository *userRepository) scanOne(
	ctx context.Context,
	row pgx.Row,
) (domain_users.User, error) {
	var input domain_users.UserInput
	var platformID uuid.UUID

	err := row.Scan(
		&input.Identifier,
		&platformID,
		&input.ExternalID,
		&input.Handle,
		&input.DisplayName,
		&input.ProfileURL,
		&input.CreatedOn,
	)
	if err != nil {
		return nil, err
	}

	platform, err := repository.platforms.FindByID(ctx, platformID)
	if err != nil {
		return nil, err
	}

	input.Platform = platform

	return repository.adapter.ToDomain(input)
}
