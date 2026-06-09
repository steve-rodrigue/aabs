package posts

import (
	"context"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/communities"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/platforms"
	domain_posts "github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/users"
)

type application struct {
	repository domain_posts.Repository
}

func createApplication(
	repository domain_posts.Repository,
) Application {
	return &application{
		repository: repository,
	}
}

func (app *application) Save(
	ctx context.Context,
	post domain_posts.Post,
) error {
	return app.repository.Save(ctx, post)
}

func (app *application) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (domain_posts.Post, error) {
	return app.repository.FindByID(ctx, id)
}

func (app *application) Find(
	ctx context.Context,
	index int,
	amount int,
) ([]domain_posts.Post, error) {
	return app.repository.Find(ctx, index, amount)
}

func (app *application) FindAfter(
	ctx context.Context,
	cursor uuid.UUID,
	amount int,
) ([]domain_posts.Post, error) {
	return app.repository.FindAfter(ctx, cursor, amount)
}

func (app *application) Count(
	ctx context.Context,
) (int64, error) {
	return app.repository.Count(ctx)
}

func (app *application) FindByUser(
	ctx context.Context,
	user users.User,
) ([]domain_posts.Post, error) {
	return app.repository.FindByUser(ctx, user)
}

func (app *application) FindByCommunity(
	ctx context.Context,
	community communities.Community,
) ([]domain_posts.Post, error) {
	return app.repository.FindByCommunity(ctx, community)
}

func (app *application) FindByPlatform(
	ctx context.Context,
	platform platforms.Platform,
) ([]domain_posts.Post, error) {
	return app.repository.FindByPlatform(ctx, platform)
}
