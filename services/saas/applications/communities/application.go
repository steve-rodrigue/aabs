package communities

import (
	"github.com/google/uuid"

	domain_communities "github.com/steve-rodrigue/aabs/services/saas/domain/communities"
	"github.com/steve-rodrigue/aabs/services/saas/domain/platforms"
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
	community domain_communities.Community,
) error {
	return app.repository.Save(community)
}

// FindByID finds a community by id
func (app *application) FindByID(
	id uuid.UUID,
) (domain_communities.Community, error) {
	return app.repository.FindByID(id)
}

// FindByHandle finds a community by platform and handle
func (app *application) FindByHandle(
	platform platforms.Platform,
	handle string,
) (domain_communities.Community, error) {
	return app.repository.FindByHandle(platform, handle)
}

// Find finds communities using offset pagination
func (app *application) Find(
	index int,
	amount int,
) ([]domain_communities.Community, error) {
	return app.repository.Find(index, amount)
}

// FindAfter finds communities using cursor pagination
func (app *application) FindAfter(
	cursor uuid.UUID,
	amount int,
) ([]domain_communities.Community, error) {
	return app.repository.FindAfter(cursor, amount)
}

// FindByPlatform finds communities by platform
func (app *application) FindByPlatform(
	platform platforms.Platform,
) ([]domain_communities.Community, error) {
	return app.repository.FindByPlatform(platform)
}

// Count counts communities
func (app *application) Count() (int64, error) {
	return app.repository.Count()
}
