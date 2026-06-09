package communities

import (
	"context"

	"github.com/google/uuid"

	domain_communities "github.com/steve-rodrigue/aabs/services/saas/domain/entities/communities"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/platforms"
)

type application struct {
	repository domain_communities.Repository
}

func createApplication(
	repository domain_communities.Repository,
) Application {
	return &application{
		repository: repository,
	}
}

// Save saves a community
func (app *application) Save(
	ctx context.Context,
	community domain_communities.Community,
) error {
	return app.repository.Save(ctx, community)
}

// FindByID finds a community by id
func (app *application) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (domain_communities.Community, error) {
	return app.repository.FindByID(ctx, id)
}

// FindByHandle finds a community by platform and handle
func (app *application) FindByHandle(
	ctx context.Context,
	platform platforms.Platform,
	handle string,
) (domain_communities.Community, error) {
	return app.repository.FindByHandle(ctx, platform, handle)
}

// Find finds communities using offset pagination
func (app *application) Find(
	ctx context.Context,
	index int,
	amount int,
) ([]domain_communities.Community, error) {
	return app.repository.Find(ctx, index, amount)
}

// FindAfter finds communities using cursor pagination
func (app *application) FindAfter(
	ctx context.Context,
	cursor uuid.UUID,
	amount int,
) ([]domain_communities.Community, error) {
	return app.repository.FindAfter(ctx, cursor, amount)
}

// FindByPlatform finds communities by platform
func (app *application) FindByPlatform(
	ctx context.Context,
	platform platforms.Platform,
) ([]domain_communities.Community, error) {
	return app.repository.FindByPlatform(ctx, platform)
}

// Count counts communities
func (app *application) Count(
	ctx context.Context,
) (int64, error) {
	return app.repository.Count(ctx)
}
