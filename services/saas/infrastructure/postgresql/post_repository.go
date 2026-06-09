package postgresql

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	domain_communities "github.com/steve-rodrigue/aabs/services/saas/domain/entities/communities"
	domain_platforms "github.com/steve-rodrigue/aabs/services/saas/domain/entities/platforms"
	domain_posts "github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts/contents"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts/contents/replies"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts/contents/threads"
	domain_users "github.com/steve-rodrigue/aabs/services/saas/domain/entities/users"
)

type postRepository struct {
	pool    *pgxpool.Pool
	adapter domain_posts.Adapter
	users   domain_users.Repository
}

func (repository *postRepository) Save(
	ctx context.Context,
	post domain_posts.Post,
) error {
	tx, err := repository.pool.Begin(ctx)
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

	if err := repository.saveContent(ctx, tx, post.Content()); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (repository *postRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (domain_posts.Post, error) {
	return repository.findOne(
		ctx,
		`
		SELECT
			identifier,
			creator_id,
			content_id,
			created_on
		FROM posts
		WHERE identifier = $1
		`,
		id,
	)
}

func (repository *postRepository) Find(
	ctx context.Context,
	index int,
	amount int,
) ([]domain_posts.Post, error) {
	rows, err := repository.pool.Query(
		ctx,
		`
		SELECT
			identifier,
			creator_id,
			content_id,
			created_on
		FROM posts
		ORDER BY identifier
		OFFSET $1
		LIMIT $2
		`,
		index,
		amount,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return repository.scanMany(ctx, rows)
}

func (repository *postRepository) FindAfter(
	ctx context.Context,
	cursor uuid.UUID,
	amount int,
) ([]domain_posts.Post, error) {
	if cursor == uuid.Nil {
		rows, err := repository.pool.Query(
			ctx,
			`
			SELECT
				identifier,
				creator_id,
				content_id,
				created_on
			FROM posts
			ORDER BY identifier
			LIMIT $1
			`,
			amount,
		)
		if err != nil {
			return nil, err
		}

		defer rows.Close()

		return repository.scanMany(ctx, rows)
	}

	rows, err := repository.pool.Query(
		ctx,
		`
		SELECT
			identifier,
			creator_id,
			content_id,
			created_on
		FROM posts
		WHERE identifier > $1
		ORDER BY identifier
		LIMIT $2
		`,
		cursor,
		amount,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return repository.scanMany(ctx, rows)
}

func (repository *postRepository) Count(
	ctx context.Context,
) (int64, error) {
	var count int64

	err := repository.pool.QueryRow(
		ctx,
		`
		SELECT COUNT(*)
		FROM posts
		`,
	).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (repository *postRepository) FindByUser(
	ctx context.Context,
	user domain_users.User,
) ([]domain_posts.Post, error) {
	rows, err := repository.pool.Query(
		ctx,
		`
		SELECT
			identifier,
			creator_id,
			content_id,
			created_on
		FROM posts
		WHERE creator_id = $1
		ORDER BY identifier
		`,
		user.Identifier(),
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return repository.scanMany(ctx, rows)
}

func (repository *postRepository) FindByCommunity(
	ctx context.Context,
	community domain_communities.Community,
) ([]domain_posts.Post, error) {
	rows, err := repository.pool.Query(
		ctx,
		`
		SELECT
			posts.identifier,
			posts.creator_id,
			posts.content_id,
			posts.created_on
		FROM posts
		INNER JOIN post_communities
			ON post_communities.post_id = posts.identifier
		WHERE post_communities.community_id = $1
		ORDER BY posts.identifier
		`,
		community.Identifier(),
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return repository.scanMany(ctx, rows)
}

func (repository *postRepository) FindByPlatform(
	ctx context.Context,
	platform domain_platforms.Platform,
) ([]domain_posts.Post, error) {
	rows, err := repository.pool.Query(
		ctx,
		`
		SELECT
			posts.identifier,
			posts.creator_id,
			posts.content_id,
			posts.created_on
		FROM posts
		INNER JOIN users
			ON users.identifier = posts.creator_id
		WHERE users.platform_id = $1
		ORDER BY posts.identifier
		`,
		platform.Identifier(),
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return repository.scanMany(ctx, rows)
}

func (repository *postRepository) findOne(
	ctx context.Context,
	query string,
	args ...any,
) (domain_posts.Post, error) {
	row := repository.pool.QueryRow(ctx, query, args...)

	post, err := repository.scanOne(ctx, row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return post, nil
}

func (repository *postRepository) scanMany(
	ctx context.Context,
	rows pgx.Rows,
) ([]domain_posts.Post, error) {
	out := []domain_posts.Post{}

	for rows.Next() {
		post, err := repository.scanOne(ctx, rows)
		if err != nil {
			return nil, err
		}

		out = append(out, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (repository *postRepository) scanOne(
	ctx context.Context,
	row pgx.Row,
) (domain_posts.Post, error) {
	var input domain_posts.PostInput
	var creatorID uuid.UUID
	var contentID uuid.UUID

	err := row.Scan(
		&input.Identifier,
		&creatorID,
		&contentID,
		&input.CreatedOn,
	)
	if err != nil {
		return nil, err
	}

	creator, err := repository.users.FindByID(ctx, creatorID)
	if err != nil {
		return nil, err
	}

	input.Creator = creator

	communityIDs, err := repository.findCommunityIDs(ctx, input.Identifier)
	if err != nil {
		return nil, err
	}

	input.CommunityIDs = communityIDs

	content, err := repository.findContentInput(ctx, contentID)
	if err != nil {
		return nil, err
	}

	input.Content = content

	return repository.adapter.ToDomain(input)
}

func (repository *postRepository) findCommunityIDs(
	ctx context.Context,
	postID uuid.UUID,
) ([]uuid.UUID, error) {
	rows, err := repository.pool.Query(
		ctx,
		`
		SELECT community_id
		FROM post_communities
		WHERE post_id = $1
		ORDER BY community_id
		`,
		postID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	out := []uuid.UUID{}

	for rows.Next() {
		var id uuid.UUID

		if err := rows.Scan(&id); err != nil {
			return nil, err
		}

		out = append(out, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (repository *postRepository) saveContent(
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
		repository.contentKind(content),
		content.CreatedAt(),
	)
	if err != nil {
		return err
	}

	if content.IsThread() {
		return repository.saveThread(ctx, tx, content.Identifier(), content.Thread())
	}

	return repository.saveReply(ctx, tx, content.Identifier(), content.Reply())
}

func (repository *postRepository) saveThread(
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

func (repository *postRepository) saveReply(
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

func (repository *postRepository) findContentInput(
	ctx context.Context,
	contentID uuid.UUID,
) (contents.ContentInput, error) {
	var input contents.ContentInput
	var kind string

	err := repository.pool.QueryRow(
		ctx,
		`
		SELECT
			identifier,
			kind,
			created_at
		FROM post_contents
		WHERE identifier = $1
		`,
		contentID,
	).Scan(
		&input.Identifier,
		&kind,
		&input.CreatedAt,
	)
	if err != nil {
		return contents.ContentInput{}, err
	}

	if kind == "thread" {
		thread, err := repository.findThreadInput(ctx, contentID)
		if err != nil {
			return contents.ContentInput{}, err
		}

		input.Thread = &thread
		return input, nil
	}

	reply, err := repository.findReplyInput(ctx, contentID)
	if err != nil {
		return contents.ContentInput{}, err
	}

	input.Reply = &reply

	return input, nil
}

func (repository *postRepository) findThreadInput(
	ctx context.Context,
	contentID uuid.UUID,
) (threads.ThreadInput, error) {
	var input threads.ThreadInput
	var creatorID uuid.UUID

	err := repository.pool.QueryRow(
		ctx,
		`
		SELECT
			identifier,
			creator_id,
			title,
			text
		FROM post_content_threads
		WHERE content_id = $1
		`,
		contentID,
	).Scan(
		&input.Identifier,
		&creatorID,
		&input.Title,
		&input.Text,
	)
	if err != nil {
		return threads.ThreadInput{}, err
	}

	creator, err := repository.users.FindByID(ctx, creatorID)
	if err != nil {
		return threads.ThreadInput{}, err
	}

	input.Creator = creator

	return input, nil
}

func (repository *postRepository) findReplyInput(
	ctx context.Context,
	contentID uuid.UUID,
) (replies.ReplyInput, error) {
	var input replies.ReplyInput
	var targetReplyID *uuid.UUID
	var targetThreadID *uuid.UUID

	err := repository.pool.QueryRow(
		ctx,
		`
		SELECT
			identifier,
			target_reply_id,
			target_thread_id,
			text
		FROM post_content_replies
		WHERE content_id = $1
		`,
		contentID,
	).Scan(
		&input.Identifier,
		&targetReplyID,
		&targetThreadID,
		&input.Text,
	)
	if err != nil {
		return replies.ReplyInput{}, err
	}

	target, err := repository.findReplyTarget(ctx, targetReplyID, targetThreadID)
	if err != nil {
		return replies.ReplyInput{}, err
	}

	input.Target = target

	return input, nil
}

func (repository *postRepository) findReplyTarget(
	ctx context.Context,
	targetReplyID *uuid.UUID,
	targetThreadID *uuid.UUID,
) (replies.TargetInput, error) {
	if targetThreadID != nil {
		thread, err := repository.findThreadByID(ctx, *targetThreadID)
		if err != nil {
			return replies.TargetInput{}, err
		}

		return replies.TargetInput{
			Thread: thread,
		}, nil
	}

	if targetReplyID != nil {
		reply, err := repository.findReplyByID(ctx, *targetReplyID)
		if err != nil {
			return replies.TargetInput{}, err
		}

		return replies.TargetInput{
			Reply: reply,
		}, nil
	}

	return replies.TargetInput{}, nil
}

func (repository *postRepository) findThreadByID(
	ctx context.Context,
	id uuid.UUID,
) (threads.Thread, error) {
	var input threads.ThreadInput
	var creatorID uuid.UUID

	err := repository.pool.QueryRow(
		ctx,
		`
		SELECT
			identifier,
			creator_id,
			title,
			text
		FROM post_content_threads
		WHERE identifier = $1
		`,
		id,
	).Scan(
		&input.Identifier,
		&creatorID,
		&input.Title,
		&input.Text,
	)
	if err != nil {
		return nil, err
	}

	creator, err := repository.users.FindByID(ctx, creatorID)
	if err != nil {
		return nil, err
	}

	input.Creator = creator

	return threads.NewAdapter().ToDomain(input)
}

func (repository *postRepository) findReplyByID(
	ctx context.Context,
	id uuid.UUID,
) (replies.Reply, error) {
	var input replies.ReplyInput
	var targetReplyID *uuid.UUID
	var targetThreadID *uuid.UUID

	err := repository.pool.QueryRow(
		ctx,
		`
		SELECT
			identifier,
			target_reply_id,
			target_thread_id,
			text
		FROM post_content_replies
		WHERE identifier = $1
		`,
		id,
	).Scan(
		&input.Identifier,
		&targetReplyID,
		&targetThreadID,
		&input.Text,
	)
	if err != nil {
		return nil, err
	}

	target, err := repository.findReplyTarget(ctx, targetReplyID, targetThreadID)
	if err != nil {
		return nil, err
	}

	input.Target = target

	return replies.NewAdapter().ToDomain(input)
}

func (repository *postRepository) contentKind(
	content contents.Content,
) string {
	if content.IsThread() {
		return "thread"
	}

	return "reply"
}
