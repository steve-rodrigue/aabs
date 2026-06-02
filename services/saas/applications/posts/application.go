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

// Save saves a post
func (app *application) Save(
	post domain_posts.Post,
) error {
	return app.repository.Save(post)
}

// FindByID finds a post by id
func (app *application) FindByID(
	id uuid.UUID,
) (domain_posts.Post, error) {
	return app.repository.FindByID(id)
}

// FindAll finds all posts
func (app *application) FindAll() ([]domain_posts.Post, error) {
	return app.repository.FindAll()
}

// FindByUser finds posts by user
func (app *application) FindByUser(
	user users.User,
) ([]domain_posts.Post, error) {
	return app.repository.FindByUser(user)
}

// FindByCommunity finds posts by community
func (app *application) FindByCommunity(
	community communities.Community,
) ([]domain_posts.Post, error) {
	return app.repository.FindByCommunity(community)
}

// FindByPlatform finds posts by platform
func (app *application) FindByPlatform(
	platform platforms.Platform,
) ([]domain_posts.Post, error) {
	return app.repository.FindByPlatform(platform)
}
