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

// FindAll finds all platforms
func (app *application) FindAll() ([]domain_platforms.Platform, error) {
	return app.repository.FindAll()
}
