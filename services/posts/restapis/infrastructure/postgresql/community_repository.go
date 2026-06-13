package postgresql

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	domain_communities "github.com/steve-rodrigue/aabs/services/posts/restapis/domain/communities"
	domain_platforms "github.com/steve-rodrigue/aabs/services/posts/restapis/domain/platforms"
	domain_users "github.com/steve-rodrigue/aabs/services/posts/restapis/domain/users"
)

type communityRepository struct {
	pool      *pgxpool.Pool
	adapter   domain_communities.Adapter
	platforms domain_platforms.Repository
	users     domain_users.Repository
}

func (repository *communityRepository) Save(
	ctx context.Context,
	community domain_communities.Community,
) error {
	tx, err := repository.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	_, err = tx.Exec(
		ctx,
		`
		INSERT INTO communities (
			identifier,
			platform_id,
			handle,
			title,
			text,
			created_on
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (identifier)
		DO UPDATE SET
			platform_id = EXCLUDED.platform_id,
			handle = EXCLUDED.handle,
			title = EXCLUDED.title,
			text = EXCLUDED.text
		`,
		community.Identifier(),
		community.Platform().Identifier(),
		community.Handle(),
		community.Title(),
		community.Text(),
		community.CreatedOn(),
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		ctx,
		`
		DELETE FROM community_moderators
		WHERE community_id = $1
		`,
		community.Identifier(),
	)
	if err != nil {
		return err
	}

	for _, moderator := range community.Moderators() {
		_, err = tx.Exec(
			ctx,
			`
			INSERT INTO community_moderators (
				community_id,
				user_id
			)
			VALUES ($1, $2)
			ON CONFLICT (community_id, user_id)
			DO NOTHING
			`,
			community.Identifier(),
			moderator.Identifier(),
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (repository *communityRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (domain_communities.Community, error) {
	return repository.findOne(
		ctx,
		`
		SELECT
			identifier,
			platform_id,
			handle,
			title,
			text,
			created_on
		FROM communities
		WHERE identifier = $1
		`,
		id,
	)
}

func (repository *communityRepository) FindByHandle(
	ctx context.Context,
	platform domain_platforms.Platform,
	handle string,
) (domain_communities.Community, error) {
	return repository.findOne(
		ctx,
		`
		SELECT
			identifier,
			platform_id,
			handle,
			title,
			text,
			created_on
		FROM communities
		WHERE platform_id = $1
		AND handle = $2
		`,
		platform.Identifier(),
		handle,
	)
}

func (repository *communityRepository) Find(
	ctx context.Context,
	index int,
	amount int,
) ([]domain_communities.Community, error) {
	rows, err := repository.pool.Query(
		ctx,
		`
		SELECT
			identifier,
			platform_id,
			handle,
			title,
			text,
			created_on
		FROM communities
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

func (repository *communityRepository) FindAfter(
	ctx context.Context,
	cursor uuid.UUID,
	amount int,
) ([]domain_communities.Community, error) {
	if cursor == uuid.Nil {
		rows, err := repository.pool.Query(
			ctx,
			`
			SELECT
				identifier,
				platform_id,
				handle,
				title,
				text,
				created_on
			FROM communities
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
			handle,
			title,
			text,
			created_on
		FROM communities
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

func (repository *communityRepository) FindByPlatform(
	ctx context.Context,
	platform domain_platforms.Platform,
) ([]domain_communities.Community, error) {
	rows, err := repository.pool.Query(
		ctx,
		`
		SELECT
			identifier,
			platform_id,
			handle,
			title,
			text,
			created_on
		FROM communities
		WHERE platform_id = $1
		ORDER BY identifier
		`,
		platform.Identifier(),
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return repository.scanMany(ctx, rows)
}

func (repository *communityRepository) Count(
	ctx context.Context,
) (int64, error) {
	var count int64

	err := repository.pool.QueryRow(
		ctx,
		`
		SELECT COUNT(*)
		FROM communities
		`,
	).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (repository *communityRepository) findOne(
	ctx context.Context,
	query string,
	args ...any,
) (domain_communities.Community, error) {
	row := repository.pool.QueryRow(ctx, query, args...)

	community, err := repository.scanOne(ctx, row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return community, nil
}

func (repository *communityRepository) scanMany(
	ctx context.Context,
	rows pgx.Rows,
) ([]domain_communities.Community, error) {
	out := []domain_communities.Community{}

	for rows.Next() {
		community, err := repository.scanOne(ctx, rows)
		if err != nil {
			return nil, err
		}

		out = append(out, community)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (repository *communityRepository) scanOne(
	ctx context.Context,
	row pgx.Row,
) (domain_communities.Community, error) {
	var input domain_communities.CommunityInput
	var platformID uuid.UUID

	err := row.Scan(
		&input.Identifier,
		&platformID,
		&input.Handle,
		&input.Title,
		&input.Text,
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

	moderators, err := repository.findModerators(ctx, input.Identifier)
	if err != nil {
		return nil, err
	}

	input.Moderators = moderators

	return repository.adapter.ToDomain(input)
}

func (repository *communityRepository) findModerators(
	ctx context.Context,
	communityID uuid.UUID,
) ([]domain_users.User, error) {
	rows, err := repository.pool.Query(
		ctx,
		`
		SELECT user_id
		FROM community_moderators
		WHERE community_id = $1
		ORDER BY user_id
		`,
		communityID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	out := []domain_users.User{}

	for rows.Next() {
		var userID uuid.UUID

		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}

		user, err := repository.users.FindByID(ctx, userID)
		if err != nil {
			return nil, err
		}

		if user != nil {
			out = append(out, user)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}
