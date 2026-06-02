package platforms

import (
	"github.com/google/uuid"

	domain_platforms "github.com/steve-rodrigue/aabs/services/saas/domain/platforms"
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
	platform domain_platforms.Platform,
) error {
	return app.repository.Save(platform)
}

// FindByID finds a platform by id
func (app *application) FindByID(
	id uuid.UUID,
) (domain_platforms.Platform, error) {
	return app.repository.FindByID(id)
}

// FindByHandle finds a platform by handle
func (app *application) FindByHandle(
	handle string,
) (domain_platforms.Platform, error) {
	return app.repository.FindByHandle(handle)
}

// Find finds platforms using offset pagination
func (app *application) Find(
	index int,
	amount int,
) ([]domain_platforms.Platform, error) {
	return app.repository.Find(index, amount)
}

// FindAfter finds platforms using cursor pagination
func (app *application) FindAfter(
	cursor uuid.UUID,
	amount int,
) ([]domain_platforms.Platform, error) {
	return app.repository.FindAfter(cursor, amount)
}

// Count counts platforms
func (app *application) Count() (int64, error) {
	return app.repository.Count()
}
