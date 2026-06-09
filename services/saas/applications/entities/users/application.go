package users

import (
	"context"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/platforms"
	domain_users "github.com/steve-rodrigue/aabs/services/saas/domain/entities/users"
)

type application struct {
	repository domain_users.Repository
}

func createApplication(
	repository domain_users.Repository,
) Application {
	return &application{
		repository: repository,
	}
}

// Save saves a user
func (app *application) Save(
	ctx context.Context,
	user domain_users.User,
) error {
	return app.repository.Save(ctx, user)
}

// FindByID finds a user by id
func (app *application) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (domain_users.User, error) {
	return app.repository.FindByID(ctx, id)
}

// FindByExternalID finds a user by platform and external id
func (app *application) FindByExternalID(
	ctx context.Context,
	platform platforms.Platform,
	externalID string,
) (domain_users.User, error) {
	return app.repository.FindByPlatformAndExternalID(ctx, platform, externalID)
}

// FindByHandle finds a user by platform and handle
func (app *application) FindByHandle(
	ctx context.Context,
	platform platforms.Platform,
	handle string,
) (domain_users.User, error) {
	return app.repository.FindByPlatformAndHandle(ctx, platform, handle)
}

// Find finds users using offset pagination
func (app *application) Find(
	ctx context.Context,
	index int,
	amount int,
) ([]domain_users.User, error) {
	return app.repository.Find(ctx, index, amount)
}

// FindAfter finds users using cursor pagination
func (app *application) FindAfter(
	ctx context.Context,
	cursor uuid.UUID,
	amount int,
) ([]domain_users.User, error) {
	return app.repository.FindAfter(ctx, cursor, amount)
}

// Count counts users
func (app *application) Count(
	ctx context.Context,
) (int64, error) {
	return app.repository.Count(ctx)
}
