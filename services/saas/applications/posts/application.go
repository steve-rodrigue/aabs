package posts

import (
	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/communities"
	"github.com/steve-rodrigue/aabs/services/saas/domain/platforms"
	domain_posts "github.com/steve-rodrigue/aabs/services/saas/domain/posts"
	"github.com/steve-rodrigue/aabs/services/saas/domain/users"
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
	post domain_posts.Post,
) error {
	return app.repository.Save(post)
}

func (app *application) FindByID(
	id uuid.UUID,
) (domain_posts.Post, error) {
	return app.repository.FindByID(id)
}

func (app *application) Find(
	index int,
	amount int,
) ([]domain_posts.Post, error) {
	return app.repository.Find(index, amount)
}

func (app *application) FindAfter(
	cursor uuid.UUID,
	amount int,
) ([]domain_posts.Post, error) {
	return app.repository.FindAfter(cursor, amount)
}

func (app *application) Count() (int64, error) {
	return app.repository.Count()
}

func (app *application) FindByUser(
	user users.User,
) ([]domain_posts.Post, error) {
	return app.repository.FindByUser(user)
}

func (app *application) FindByCommunity(
	community communities.Community,
) ([]domain_posts.Post, error) {
	return app.repository.FindByCommunity(community)
}

func (app *application) FindByPlatform(
	platform platforms.Platform,
) ([]domain_posts.Post, error) {
	return app.repository.FindByPlatform(platform)
}
