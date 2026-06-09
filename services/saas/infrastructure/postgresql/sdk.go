package postgresql

import (
	"github.com/jackc/pgx/v5/pgxpool"

	domain_communities "github.com/steve-rodrigue/aabs/services/saas/domain/entities/communities"
	domain_platforms "github.com/steve-rodrigue/aabs/services/saas/domain/entities/platforms"
	domain_posts "github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts"
	domain_users "github.com/steve-rodrigue/aabs/services/saas/domain/entities/users"
	domain_clusters "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters/clusterables"
	domain_participations "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations"
	domain_evidences "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/evidences"
	domain_topics "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/topics"
	domain_relationships "github.com/steve-rodrigue/aabs/services/saas/domain/relationships"
	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships/relatables"
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

// NewRelationshipRepository creates a new postgresql relationship repository
func NewRelationshipRepository(
	pool *pgxpool.Pool,
	adapter domain_relationships.Adapter,
	relatables relatables.Adapter,
) domain_relationships.Repository {
	return &relationshipRepository{
		pool:       pool,
		adapter:    adapter,
		relatables: relatables,
	}
}

// NewGroupingsClustersClusterablesRepository creates a new postgresql clusterable repository
func NewGroupingsClustersClusterablesRepository(
	pool *pgxpool.Pool,
	adapter clusterables.Adapter,
) clusterables.Repository {
	return &groupingsClustersClusterableRepository{
		pool:    pool,
		adapter: adapter,
	}
}

// NewGroupingsClusterRepository creates a new postgresql cluster repository
func NewGroupingsClusterRepository(
	pool *pgxpool.Pool,
	adapter domain_clusters.Adapter,
) domain_clusters.Repository {
	return &groupingsClusterRepository{
		pool:    pool,
		adapter: adapter,
	}
}

// NewGroupingsTopicRepository creates a new postgresql topic repository
func NewGroupingsTopicRepository(
	pool *pgxpool.Pool,
	adapter domain_topics.Adapter,
	clusters domain_clusters.Repository,
) domain_topics.Repository {
	return &groupingsTopicRepository{
		pool:     pool,
		adapter:  adapter,
		clusters: clusters,
	}
}

// NewGroupingsParticipationsEvidenceRepository creates a new postgresql participation evidence repository
func NewGroupingsParticipationsEvidenceRepository(
	pool *pgxpool.Pool,
	adapter domain_evidences.Adapter,
	participations domain_participations.Repository,
	posts domain_posts.Repository,
) domain_evidences.Repository {
	return &groupingsParticipationsEvidenceRepository{
		pool:           pool,
		adapter:        adapter,
		participations: participations,
		posts:          posts,
	}
}
