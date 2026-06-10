package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/participatables"
)

var ErrInvalidConceptParticipatableCounterKind = errors.New(
	"invalid concept participatable counter kind",
)

type conceptParticipatableCounter struct {
	pool *pgxpool.Pool
}

func (counter *conceptParticipatableCounter) CountByParticipantAndTarget(
	ctx context.Context,
	participant participatables.Participatable,
	target participatables.Participatable,
) (int, error) {
	switch participant.ParticipationKind() {
	case participatables.UserKind:
		return counter.countUserByTarget(ctx, participant, target)

	case participatables.CommunityKind:
		return counter.countCommunityByTarget(ctx, participant, target)

	case participatables.PlatformKind:
		return counter.countPlatformByTarget(ctx, participant, target)

	case participatables.PostKind:
		return counter.countPostByTarget(ctx, participant, target)

	case participatables.CampaignKind,
		participatables.TopicKind,
		participatables.NarrativeKind:
		return counter.countGroupingByTarget(ctx, participant, target)

	default:
		return 0, ErrInvalidConceptParticipatableCounterKind
	}
}

func (counter *conceptParticipatableCounter) CountByTarget(
	ctx context.Context,
	target participatables.Participatable,
) (int, error) {
	switch target.ParticipationKind() {
	case participatables.UserKind:
		return counter.countByUser(ctx, target)

	case participatables.CommunityKind:
		return counter.countByCommunity(ctx, target)

	case participatables.PlatformKind:
		return counter.countByPlatform(ctx, target)

	case participatables.PostKind:
		return counter.countByPost(ctx, target)

	case participatables.CampaignKind,
		participatables.TopicKind,
		participatables.NarrativeKind:
		return counter.countByGrouping(ctx, target)

	default:
		return 0, ErrInvalidConceptParticipatableCounterKind
	}
}

func (counter *conceptParticipatableCounter) countUserByTarget(
	ctx context.Context,
	participant participatables.Participatable,
	target participatables.Participatable,
) (int, error) {
	return counter.countParticipantPostsByTarget(
		ctx,
		`
		SELECT p.identifier
		FROM posts p
		WHERE p.user_id = $1
		`,
		[]any{participant.Identifier()},
		target,
	)
}

func (counter *conceptParticipatableCounter) countCommunityByTarget(
	ctx context.Context,
	participant participatables.Participatable,
	target participatables.Participatable,
) (int, error) {
	return counter.countParticipantPostsByTarget(
		ctx,
		`
		SELECT p.identifier
		FROM posts p
		WHERE $1 = ANY(p.community_ids)
		`,
		[]any{participant.Identifier()},
		target,
	)
}

func (counter *conceptParticipatableCounter) countPlatformByTarget(
	ctx context.Context,
	participant participatables.Participatable,
	target participatables.Participatable,
) (int, error) {
	return counter.countParticipantPostsByTarget(
		ctx,
		`
		SELECT p.identifier
		FROM posts p
		INNER JOIN users u
			ON u.identifier = p.user_id
		WHERE u.platform_id = $1
		`,
		[]any{participant.Identifier()},
		target,
	)
}

func (counter *conceptParticipatableCounter) countPostByTarget(
	ctx context.Context,
	participant participatables.Participatable,
	target participatables.Participatable,
) (int, error) {
	return counter.countParticipantPostsByTarget(
		ctx,
		`
		SELECT p.identifier
		FROM posts p
		WHERE p.identifier = $1
		`,
		[]any{participant.Identifier()},
		target,
	)
}

func (counter *conceptParticipatableCounter) countGroupingByTarget(
	ctx context.Context,
	participant participatables.Participatable,
	target participatables.Participatable,
) (int, error) {
	return counter.countParticipantPostsByTarget(
		ctx,
		`
		SELECT DISTINCT m.member_id AS identifier
		FROM groupings_clusters c
		INNER JOIN groupings_clusters_members m
			ON m.cluster_id = c.identifier
		WHERE c.target_id = $1
		AND c.target_kind = $2
		AND m.member_kind = $3
		`,
		[]any{
			participant.Identifier(),
			string(participant.ParticipationKind()),
			string(participatables.PostKind),
		},
		target,
	)
}

func (counter *conceptParticipatableCounter) countByUser(
	ctx context.Context,
	target participatables.Participatable,
) (int, error) {
	return counter.count(
		ctx,
		`
		SELECT COUNT(*)
		FROM posts
		WHERE user_id = $1
		`,
		target.Identifier(),
	)
}

func (counter *conceptParticipatableCounter) countByCommunity(
	ctx context.Context,
	target participatables.Participatable,
) (int, error) {
	return counter.count(
		ctx,
		`
		SELECT COUNT(*)
		FROM posts
		WHERE $1 = ANY(community_ids)
		`,
		target.Identifier(),
	)
}

func (counter *conceptParticipatableCounter) countByPlatform(
	ctx context.Context,
	target participatables.Participatable,
) (int, error) {
	return counter.count(
		ctx,
		`
		SELECT COUNT(*)
		FROM posts p
		INNER JOIN users u
			ON u.identifier = p.user_id
		WHERE u.platform_id = $1
		`,
		target.Identifier(),
	)
}

func (counter *conceptParticipatableCounter) countByPost(
	ctx context.Context,
	target participatables.Participatable,
) (int, error) {
	return counter.count(
		ctx,
		`
		SELECT COUNT(*)
		FROM posts
		WHERE identifier = $1
		`,
		target.Identifier(),
	)
}

func (counter *conceptParticipatableCounter) countByGrouping(
	ctx context.Context,
	target participatables.Participatable,
) (int, error) {
	return counter.count(
		ctx,
		`
		SELECT COUNT(DISTINCT m.member_id)
		FROM groupings_clusters c
		INNER JOIN groupings_clusters_members m
			ON m.cluster_id = c.identifier
		WHERE c.target_id = $1
		AND c.target_kind = $2
		AND m.member_kind = $3
		`,
		target.Identifier(),
		string(target.ParticipationKind()),
		string(participatables.PostKind),
	)
}

func (counter *conceptParticipatableCounter) countParticipantPostsByTarget(
	ctx context.Context,
	participantPostQuery string,
	args []any,
	target participatables.Participatable,
) (int, error) {
	switch target.ParticipationKind() {
	case participatables.UserKind:
		return counter.countWithParticipantPosts(
			ctx,
			participantPostQuery,
			args,
			`
			INNER JOIN posts p
				ON p.identifier = participant_posts.identifier
			WHERE p.user_id = $%d
			`,
			target.Identifier(),
		)

	case participatables.CommunityKind:
		return counter.countWithParticipantPosts(
			ctx,
			participantPostQuery,
			args,
			`
			INNER JOIN posts p
				ON p.identifier = participant_posts.identifier
			WHERE $%d = ANY(p.community_ids)
			`,
			target.Identifier(),
		)

	case participatables.PlatformKind:
		return counter.countWithParticipantPosts(
			ctx,
			participantPostQuery,
			args,
			`
			INNER JOIN posts p
				ON p.identifier = participant_posts.identifier
			INNER JOIN users u
				ON u.identifier = p.user_id
			WHERE u.platform_id = $%d
			`,
			target.Identifier(),
		)

	case participatables.PostKind:
		return counter.countWithParticipantPosts(
			ctx,
			participantPostQuery,
			args,
			`
			WHERE participant_posts.identifier = $%d
			`,
			target.Identifier(),
		)

	case participatables.CampaignKind,
		participatables.TopicKind,
		participatables.NarrativeKind:
		return counter.countWithParticipantPosts(
			ctx,
			participantPostQuery,
			args,
			`
			INNER JOIN groupings_clusters c
				ON c.target_id = $%d
				AND c.target_kind = $%d
			INNER JOIN groupings_clusters_members m
				ON m.cluster_id = c.identifier
				AND m.member_id = participant_posts.identifier
				AND m.member_kind = $%d
			`,
			target.Identifier(),
			string(target.ParticipationKind()),
			string(participatables.PostKind),
		)

	default:
		return 0, ErrInvalidConceptParticipatableCounterKind
	}
}

func (counter *conceptParticipatableCounter) countWithParticipantPosts(
	ctx context.Context,
	participantPostQuery string,
	args []any,
	targetJoinFormat string,
	targetArgs ...any,
) (int, error) {
	firstTargetArgPosition := len(args) + 1

	for _, arg := range targetArgs {
		args = append(args, arg)
	}

	placeholders := make([]any, len(targetArgs))
	for index := range targetArgs {
		placeholders[index] = firstTargetArgPosition + index
	}

	query := `
		WITH participant_posts AS (
	` + participantPostQuery + `
		)
		SELECT COUNT(DISTINCT participant_posts.identifier)
		FROM participant_posts
	` + fmt.Sprintf(targetJoinFormat, placeholders...)

	return counter.count(ctx, query, args...)
}

func (counter *conceptParticipatableCounter) count(
	ctx context.Context,
	query string,
	args ...any,
) (int, error) {
	var count int

	err := counter.pool.QueryRow(
		ctx,
		query,
		args...,
	).Scan(&count)

	return count, err
}
