package postgresql

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/clusterables"
	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/participatables"
	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/relatables"
	domain_communities "github.com/steve-rodrigue/aabs/services/saas/domain/entities/communities"
	domain_platforms "github.com/steve-rodrigue/aabs/services/saas/domain/entities/platforms"
	domain_posts "github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts"
	domain_users "github.com/steve-rodrigue/aabs/services/saas/domain/entities/users"
	domain_campaigns "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/campaigns"
	domain_clusters "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters"
	domain_narratives "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/narratives"
	domain_assignments "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/narratives/assignments"
	domain_participations "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations"
	domain_dirty_participation "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/dirty"
	domain_evidences "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/evidences"
	domain_topics "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/topics"
	domain_relationships "github.com/steve-rodrigue/aabs/services/saas/domain/relationships"
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

// NewGroupingsParticipationRepository creates a new postgresql participation repository
func NewGroupingsParticipationRepository(
	pool *pgxpool.Pool,
	adapter domain_participations.Adapter,
) domain_participations.Repository {
	return &groupingsParticipationRepository{
		pool:    pool,
		adapter: adapter,
	}
}

// NewGroupingsNarrativeRepository creates a new postgresql narrative repository
func NewGroupingsNarrativeRepository(
	pool *pgxpool.Pool,
	adapter domain_narratives.Adapter,
	clusters domain_clusters.Repository,
) domain_narratives.Repository {
	return &groupingsNarrativeRepository{
		pool:     pool,
		adapter:  adapter,
		clusters: clusters,
	}
}

// NewGroupingsCampaignRepository creates a new postgresql campaign repository
func NewGroupingsCampaignRepository(
	pool *pgxpool.Pool,
	adapter domain_campaigns.Adapter,
	clusters domain_clusters.Repository,
) domain_campaigns.Repository {
	return &groupingsCampaignRepository{
		pool:     pool,
		adapter:  adapter,
		clusters: clusters,
	}
}

// NewGroupingsAssignmentRepository creates a new postgresql assignment repository
func NewGroupingsAssignmentRepository(
	pool *pgxpool.Pool,
	adapter domain_assignments.Adapter,
	narratives domain_narratives.Repository,
	campaigns domain_campaigns.Repository,
) domain_assignments.Repository {
	return &groupingsAssignmentRepository{
		pool:       pool,
		adapter:    adapter,
		narratives: narratives,
		campaigns:  campaigns,
	}
}

// NewGroupingsParticipationsDirtyRepository creates a new postgresql dirty participation repository
func NewGroupingsParticipationsDirtyRepository(
	pool *pgxpool.Pool,
	adapter domain_dirty_participation.Adapter,
) domain_dirty_participation.Repository {
	return &groupingsParticipationsDirtyRepository{
		pool:    pool,
		adapter: adapter,
	}
}

// NewConceptParticipatableCounter creates a new postgresql participatable counter
func NewConceptParticipatableCounter(
	pool *pgxpool.Pool,
) participatables.Counter {
	return &conceptParticipatableCounter{
		pool: pool,
	}
}
