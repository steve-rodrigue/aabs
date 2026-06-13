package postgresql

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	domain_posts "github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts/contents"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts/contents/replies"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts/contents/threads"
)

type postService struct {
	pool *pgxpool.Pool
}

func (service *postService) Save(
	ctx context.Context,
	post domain_posts.Post,
) error {
	tx, err := service.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	_, err = tx.Exec(
		ctx,
		`
		INSERT INTO posts (
			identifier,
			creator_id,
			content_id,
			created_on
		)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (identifier)
		DO UPDATE SET
			creator_id = EXCLUDED.creator_id,
			content_id = EXCLUDED.content_id
		`,
		post.Identifier(),
		post.Creator().Identifier(),
		post.Content().Identifier(),
		post.CreatedOn(),
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		ctx,
		`
		DELETE FROM post_communities
		WHERE post_id = $1
		`,
		post.Identifier(),
	)
	if err != nil {
		return err
	}

	for _, communityID := range post.CommunityIDs() {
		_, err = tx.Exec(
			ctx,
			`
			INSERT INTO post_communities (
				post_id,
				community_id
			)
			VALUES ($1, $2)
			ON CONFLICT (post_id, community_id)
			DO NOTHING
			`,
			post.Identifier(),
			communityID,
		)
		if err != nil {
			return err
		}
	}

	if err := service.saveContent(ctx, tx, post.Content()); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (service *postService) saveContent(
	ctx context.Context,
	tx pgx.Tx,
	content contents.Content,
) error {
	_, err := tx.Exec(
		ctx,
		`
		INSERT INTO post_contents (
			identifier,
			kind,
			created_at
		)
		VALUES ($1, $2, $3)
		ON CONFLICT (identifier)
		DO UPDATE SET
			kind = EXCLUDED.kind,
			created_at = EXCLUDED.created_at
		`,
		content.Identifier(),
		service.contentKind(content),
		content.CreatedAt(),
	)
	if err != nil {
		return err
	}

	if content.IsThread() {
		return service.saveThread(ctx, tx, content.Identifier(), content.Thread())
	}

	return service.saveReply(ctx, tx, content.Identifier(), content.Reply())
}

func (service *postService) saveThread(
	ctx context.Context,
	tx pgx.Tx,
	contentID uuid.UUID,
	thread threads.Thread,
) error {
	_, err := tx.Exec(
		ctx,
		`
		INSERT INTO post_content_threads (
			content_id,
			identifier,
			creator_id,
			title,
			text
		)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (content_id)
		DO UPDATE SET
			identifier = EXCLUDED.identifier,
			creator_id = EXCLUDED.creator_id,
			title = EXCLUDED.title,
			text = EXCLUDED.text
		`,
		contentID,
		thread.Identifier(),
		thread.Creator().Identifier(),
		thread.Title(),
		thread.Text(),
	)

	return err
}

func (service *postService) saveReply(
	ctx context.Context,
	tx pgx.Tx,
	contentID uuid.UUID,
	reply replies.Reply,
) error {
	var targetReplyID *uuid.UUID
	var targetThreadID *uuid.UUID

	if reply.Target().IsReply() {
		id := reply.Target().Reply().Identifier()
		targetReplyID = &id
	}

	if reply.Target().IsThread() {
		id := reply.Target().Thread().Identifier()
		targetThreadID = &id
	}

	_, err := tx.Exec(
		ctx,
		`
		INSERT INTO post_content_replies (
			content_id,
			identifier,
			target_reply_id,
			target_thread_id,
			text
		)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (content_id)
		DO UPDATE SET
			identifier = EXCLUDED.identifier,
			target_reply_id = EXCLUDED.target_reply_id,
			target_thread_id = EXCLUDED.target_thread_id,
			text = EXCLUDED.text
		`,
		contentID,
		reply.Identifier(),
		targetReplyID,
		targetThreadID,
		reply.Text(),
	)

	return err
}

func (service *postService) contentKind(
	content contents.Content,
) string {
	if content.IsThread() {
		return "thread"
	}

	return "reply"
}
