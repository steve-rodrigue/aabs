package postgresql

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	domain_campaigns "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/campaigns"
	domain_clusters "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters"
)

type groupingsCampaignRepository struct {
	pool *pgxpool.Pool

	adapter  domain_campaigns.Adapter
	clusters domain_clusters.Repository
}

func (repository *groupingsCampaignRepository) Save(
	ctx context.Context,
	campaign domain_campaigns.Campaign,
) error {
	_, err := repository.pool.Exec(
		ctx,
		`
		INSERT INTO groupings_campaigns (
			identifier,
			name,
			description,
			cluster_id,
			post_count,
			confidence,
			created_on
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (identifier)
		DO UPDATE SET
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			cluster_id = EXCLUDED.cluster_id,
			post_count = EXCLUDED.post_count,
			confidence = EXCLUDED.confidence
		`,
		campaign.Identifier(),
		campaign.Name(),
		campaign.Description(),
		campaign.Cluster().Identifier(),
		campaign.PostCount(),
		campaign.Confidence(),
		campaign.CreatedOn(),
	)

	return err
}

func (repository *groupingsCampaignRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (domain_campaigns.Campaign, error) {
	return repository.findOne(
		ctx,
		`
		SELECT
			identifier,
			name,
			description,
			cluster_id,
			post_count,
			confidence,
			created_on
		FROM groupings_campaigns
		WHERE identifier = $1
		`,
		id,
	)
}

func (repository *groupingsCampaignRepository) FindByName(
	ctx context.Context,
	name string,
) (domain_campaigns.Campaign, error) {
	return repository.findOne(
		ctx,
		`
		SELECT
			identifier,
			name,
			description,
			cluster_id,
			post_count,
			confidence,
			created_on
		FROM groupings_campaigns
		WHERE name = $1
		`,
		name,
	)
}

func (repository *groupingsCampaignRepository) Find(
	ctx context.Context,
	index int,
	amount int,
) ([]domain_campaigns.Campaign, error) {
	rows, err := repository.pool.Query(
		ctx,
		`
		SELECT
			identifier,
			name,
			description,
			cluster_id,
			post_count,
			confidence,
			created_on
		FROM groupings_campaigns
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

func (repository *groupingsCampaignRepository) FindAfter(
	ctx context.Context,
	cursor uuid.UUID,
	amount int,
) ([]domain_campaigns.Campaign, error) {
	if cursor == uuid.Nil {
		return repository.Find(ctx, 0, amount)
	}

	rows, err := repository.pool.Query(
		ctx,
		`
		SELECT
			identifier,
			name,
			description,
			cluster_id,
			post_count,
			confidence,
			created_on
		FROM groupings_campaigns
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

func (repository *groupingsCampaignRepository) Count(
	ctx context.Context,
) (int64, error) {
	var count int64

	err := repository.pool.QueryRow(
		ctx,
		`
		SELECT COUNT(*)
		FROM groupings_campaigns
		`,
	).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (repository *groupingsCampaignRepository) findOne(
	ctx context.Context,
	query string,
	args ...any,
) (domain_campaigns.Campaign, error) {
	row := repository.pool.QueryRow(ctx, query, args...)

	campaign, err := repository.scanOne(ctx, row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return campaign, nil
}

func (repository *groupingsCampaignRepository) scanMany(
	ctx context.Context,
	rows pgx.Rows,
) ([]domain_campaigns.Campaign, error) {
	out := []domain_campaigns.Campaign{}

	for rows.Next() {
		campaign, err := repository.scanOne(ctx, rows)
		if err != nil {
			return nil, err
		}

		out = append(out, campaign)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (repository *groupingsCampaignRepository) scanOne(
	ctx context.Context,
	row pgx.Row,
) (domain_campaigns.Campaign, error) {
	var input domain_campaigns.CampaignInput

	var clusterID uuid.UUID

	err := row.Scan(
		&input.Identifier,
		&input.Name,
		&input.Description,
		&clusterID,
		&input.PostCount,
		&input.Confidence,
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

	input.Cluster = cluster

	return repository.adapter.ToDomain(input)
}
