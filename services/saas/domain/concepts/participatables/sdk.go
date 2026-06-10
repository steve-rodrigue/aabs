package participatables

import (
	"context"

	"github.com/google/uuid"
)

type Kind string

const (
	UserKind      Kind = "user"
	CommunityKind Kind = "community"
	PlatformKind  Kind = "platform"

	PostKind Kind = "post"

	CampaignKind  Kind = "campaign"
	TopicKind     Kind = "topic"
	NarrativeKind Kind = "narrative"
)

// Participatable represents an entity that can participate in another entity.
type Participatable interface {
	Identifier() uuid.UUID
	ParticipationKind() Kind
}

// Counter counts posts involved in participations
type Counter interface {
	// CountByParticipantAndTarget returns how many posts created by
	// participant belong to target
	CountByParticipantAndTarget(
		ctx context.Context,
		participant Participatable,
		target Participatable,
	) (int, error)
	// CountByTarget returns the total number of posts belonging
	// to the target
	CountByTarget(
		ctx context.Context,
		target Participatable,
	) (int, error)
}

// Repository represents a participatable repository.
type Repository interface {
	CountParticipants(ctx context.Context) (int64, error)
	FindParticipantsAfter(
		ctx context.Context,
		cursor uuid.UUID,
		amount int,
	) ([]Participatable, error)
	CountTargets(ctx context.Context) (int64, error)
	FindTargetsAfter(
		ctx context.Context,
		cursor uuid.UUID,
		amount int,
	) ([]Participatable, error)
}
