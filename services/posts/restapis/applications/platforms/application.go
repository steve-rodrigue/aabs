package platforms

import (
	"context"

	"github.com/google/uuid"

	domain_platforms "github.com/steve-rodrigue/aabs/services/posts/restapis/domain/platforms"
)

type application struct {
	repository domain_platforms.Repository
}

func createApplication(
	repository domain_platforms.Repository,
) Application {
	return &application{
		repository: repository,
	}
}

// Save saves a platform
func (app *application) Save(
	ctx context.Context,
	platform domain_platforms.Platform,
) error {
	return app.repository.Save(ctx, platform)
}

// FindByID finds a platform by id
func (app *application) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (domain_platforms.Platform, error) {
	return app.repository.FindByID(ctx, id)
}

// FindByHandle finds a platform by handle
func (app *application) FindByHandle(
	ctx context.Context,
	handle string,
) (domain_platforms.Platform, error) {
	return app.repository.FindByHandle(ctx, handle)
}

// Find finds platforms using offset pagination
func (app *application) Find(
	ctx context.Context,
	index int,
	amount int,
) ([]domain_platforms.Platform, error) {
	return app.repository.Find(ctx, index, amount)
}

// FindAfter finds platforms using cursor pagination
func (app *application) FindAfter(
	ctx context.Context,
	cursor uuid.UUID,
	amount int,
) ([]domain_platforms.Platform, error) {
	return app.repository.FindAfter(ctx, cursor, amount)
}

// Count counts platforms
func (app *application) Count(
	ctx context.Context,
) (int64, error) {
	return app.repository.Count(ctx)
}
