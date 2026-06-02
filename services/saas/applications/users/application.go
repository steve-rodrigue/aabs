package users

import (
	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/platforms"
	domain_users "github.com/steve-rodrigue/aabs/services/saas/domain/users"
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
	user domain_users.User,
) error {
	return app.repository.Save(user)
}

// FindByID finds a user by id
func (app *application) FindByID(
	id uuid.UUID,
) (domain_users.User, error) {
	return app.repository.FindByID(id)
}

// FindByExternalID finds a user by platform and external id
func (app *application) FindByExternalID(
	platform platforms.Platform,
	externalID string,
) (domain_users.User, error) {
	return app.repository.FindByPlatformAndExternalID(platform, externalID)
}

// Find finds users using offset pagination
func (app *application) Find(
	index int,
	amount int,
) ([]domain_users.User, error) {
	return app.repository.Find(index, amount)
}

// FindAfter finds users using cursor pagination
func (app *application) FindAfter(
	cursor uuid.UUID,
	amount int,
) ([]domain_users.User, error) {
	return app.repository.FindAfter(cursor, amount)
}

// Count counts users
func (app *application) Count() (int64, error) {
	return app.repository.Count()
}
