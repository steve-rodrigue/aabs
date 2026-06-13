package applications

import (
	"context"

	hatchetsdk "github.com/hatchet-dev/hatchet/sdks/go"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/steve-rodrigue/aabs/services/posts/restapis/applications"

	application_communities "github.com/steve-rodrigue/aabs/services/posts/restapis/applications/communities"
	application_platforms "github.com/steve-rodrigue/aabs/services/posts/restapis/applications/platforms"
	application_posts "github.com/steve-rodrigue/aabs/services/posts/restapis/applications/posts"
	application_users "github.com/steve-rodrigue/aabs/services/posts/restapis/applications/users"

	domain_communities "github.com/steve-rodrigue/aabs/services/posts/restapis/domain/communities"
	domain_platforms "github.com/steve-rodrigue/aabs/services/posts/restapis/domain/platforms"
	domain_posts "github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts/contents"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts/contents/replies"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts/contents/threads"
	domain_users "github.com/steve-rodrigue/aabs/services/posts/restapis/domain/users"

	"github.com/steve-rodrigue/aabs/services/posts/restapis/infrastructure/hatchet"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/infrastructure/postgresql"
)

// New creates a new entities application
func New(
	pool *pgxpool.Pool,
	hatchetClient *hatchetsdk.Client,
) applications.Application {
	platformRepository := postgresql.NewPlatformRepository(
		pool,
		domain_platforms.NewAdapter(),
	)

	userRepository := postgresql.NewUserRepository(
		pool,
		domain_users.NewAdapter(),
		platformRepository,
	)

	communityRepository := postgresql.NewCommunityRepository(
		pool,
		domain_communities.NewAdapter(),
		platformRepository,
		userRepository,
	)

	contentAdapter := contents.NewAdapter(
		replies.NewAdapter(),
		threads.NewAdapter(),
	)

	postRepository := postgresql.NewPostRepository(
		pool,
		domain_posts.NewAdapter(contentAdapter),
		userRepository,
	)

	postServices := []domain_posts.Service{
		postgresql.NewPostService(pool),
	}

	if hatchetClient != nil {
		postServices = append(
			postServices,
			hatchet.NewPostService(hatchetClient),
		)
	}

	postService := domain_posts.NewService(postServices...)

	return applications.New(
		application_posts.New(
			postRepository,
			postService,
		),
		application_users.New(
			userRepository,
		),
		application_communities.New(
			communityRepository,
		),
		application_platforms.New(
			platformRepository,
		),
	)
}

// Install installs the application
func Install(
	ctx context.Context,
	pool *pgxpool.Pool,
) error {
	_, err := pool.Exec(
		ctx,
		`
		CREATE TABLE IF NOT EXISTS platforms (
			identifier UUID PRIMARY KEY,
			name TEXT NOT NULL,
			handle TEXT NOT NULL UNIQUE,
			base_url TEXT NOT NULL,
			created_on TIMESTAMPTZ NOT NULL
		);
		
		CREATE TABLE IF NOT EXISTS users (
			identifier UUID PRIMARY KEY,
			platform_id UUID NOT NULL
				REFERENCES platforms(identifier)
				ON DELETE CASCADE,
			external_id TEXT NOT NULL,
			handle TEXT NOT NULL,
			display_name TEXT NOT NULL,
			profile_url TEXT NOT NULL,
			created_on TIMESTAMPTZ NOT NULL,
		
			UNIQUE (platform_id, external_id),
			UNIQUE (platform_id, handle)
		);
		
		CREATE TABLE IF NOT EXISTS communities (
			identifier UUID PRIMARY KEY,
			platform_id UUID NOT NULL
				REFERENCES platforms(identifier)
				ON DELETE CASCADE,
			handle TEXT NOT NULL,
			title TEXT NOT NULL,
			text TEXT NOT NULL,
			created_on TIMESTAMPTZ NOT NULL,
		
			UNIQUE (platform_id, handle)
		);
		
		CREATE TABLE IF NOT EXISTS community_moderators (
			community_id UUID NOT NULL
				REFERENCES communities(identifier)
				ON DELETE CASCADE,
			user_id UUID NOT NULL
				REFERENCES users(identifier)
				ON DELETE CASCADE,
		
			PRIMARY KEY (community_id, user_id)
		);
		
		CREATE TABLE IF NOT EXISTS post_contents (
			identifier UUID PRIMARY KEY,
			kind TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL
		);
		
		CREATE TABLE IF NOT EXISTS posts (
			identifier UUID PRIMARY KEY,
			creator_id UUID NOT NULL
				REFERENCES users(identifier)
				ON DELETE CASCADE,
			content_id UUID NOT NULL UNIQUE
				REFERENCES post_contents(identifier)
				ON DELETE CASCADE,
			created_on TIMESTAMPTZ NOT NULL
		);
		
		CREATE TABLE IF NOT EXISTS post_communities (
			post_id UUID NOT NULL
				REFERENCES posts(identifier)
				ON DELETE CASCADE,
			community_id UUID NOT NULL
				REFERENCES communities(identifier)
				ON DELETE CASCADE,
		
			PRIMARY KEY (post_id, community_id)
		);
		
		CREATE TABLE IF NOT EXISTS post_content_threads (
			content_id UUID PRIMARY KEY
				REFERENCES post_contents(identifier)
				ON DELETE CASCADE,
			identifier UUID NOT NULL UNIQUE,
			creator_id UUID NOT NULL
				REFERENCES users(identifier)
				ON DELETE CASCADE,
			title TEXT NOT NULL,
			text TEXT NOT NULL
		);
		
		CREATE TABLE IF NOT EXISTS post_content_replies (
			content_id UUID PRIMARY KEY
				REFERENCES post_contents(identifier)
				ON DELETE CASCADE,
			identifier UUID NOT NULL UNIQUE,
		
			target_reply_id UUID NULL
				REFERENCES post_content_replies(identifier)
				ON DELETE CASCADE,
		
			target_thread_id UUID NULL
				REFERENCES post_content_threads(identifier)
				ON DELETE CASCADE,
		
			text TEXT NOT NULL
		);
		`,
	)

	return err
}
