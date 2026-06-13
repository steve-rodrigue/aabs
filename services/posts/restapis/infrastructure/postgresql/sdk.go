package postgresql

import (
	"github.com/jackc/pgx/v5/pgxpool"

	domain_communities "github.com/steve-rodrigue/aabs/services/posts/restapis/domain/communities"
	domain_platforms "github.com/steve-rodrigue/aabs/services/posts/restapis/domain/platforms"
	domain_posts "github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts"
	domain_users "github.com/steve-rodrigue/aabs/services/posts/restapis/domain/users"
)

// NewPlatformRepository creates a new postgresql platform repository
func NewPlatformRepository(
	pool *pgxpool.Pool,
	adapter domain_platforms.Adapter,
) domain_platforms.Repository {
	return &platformRepository{
		pool:    pool,
		adapter: adapter,
	}
}

// NewUserRepository creates a new postgresql user repository
func NewUserRepository(
	pool *pgxpool.Pool,
	adapter domain_users.Adapter,
	platforms domain_platforms.Repository,
) domain_users.Repository {
	return &userRepository{
		pool:      pool,
		adapter:   adapter,
		platforms: platforms,
	}

}

// NewCommunityRepository creates a new postgresql community repository
func NewCommunityRepository(
	pool *pgxpool.Pool,
	adapter domain_communities.Adapter,
	platforms domain_platforms.Repository,
	users domain_users.Repository,
) domain_communities.Repository {
	return &communityRepository{
		pool:      pool,
		adapter:   adapter,
		platforms: platforms,
		users:     users,
	}
}

// NewPostRepository creates a new postgresql post repository
func NewPostRepository(
	pool *pgxpool.Pool,
	adapter domain_posts.Adapter,
	users domain_users.Repository,
) domain_posts.Repository {
	return &postRepository{
		pool:    pool,
		adapter: adapter,
		users:   users,
	}
}

// NewPostService creates a new postgresql post service
func NewPostService(
	pool *pgxpool.Pool,
) domain_posts.Service {
	return &postService{
		pool: pool,
	}
}
