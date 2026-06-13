package platforms

import (
	"time"

	"github.com/google/uuid"
)

type platform struct {
	identifier uuid.UUID

	name    string
	handle  string
	baseURL string

	createdOn time.Time
}

func (platform *platform) Identifier() uuid.UUID {
	return platform.identifier
}

func (platform *platform) Name() string {
	return platform.name
}

func (platform *platform) Handle() string {
	return platform.handle
}

func (platform *platform) BaseURL() string {
	return platform.baseURL
}

func (platform *platform) CreatedOn() time.Time {
	return platform.createdOn
}
