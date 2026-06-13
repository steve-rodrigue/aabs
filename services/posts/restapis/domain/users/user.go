package users

import (
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/platforms"
)

type user struct {
	identifier uuid.UUID

	platform platforms.Platform

	externalID  string
	handle      string
	displayName string
	profileURL  string

	createdOn time.Time
}

func (user *user) Identifier() uuid.UUID {
	return user.identifier
}

func (user *user) Platform() platforms.Platform {
	return user.platform
}

func (user *user) ExternalID() string {
	return user.externalID
}

func (user *user) Handle() string {
	return user.handle
}

func (user *user) DisplayName() string {
	return user.displayName
}

func (user *user) ProfileURL() string {
	return user.profileURL
}

func (user *user) CreatedOn() time.Time {
	return user.createdOn
}
