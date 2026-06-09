package redis

import (
	"github.com/redis/go-redis/v9"

	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships/relatables"
)

const (
	relationshipRelatablesKey     = "relationship:relatables"
	relationshipRelatableKindsKey = "relationship:relatable:kinds"
)

// NewRelationshipRelatableRepository creates a new redis relationship relatable repository
func NewRelationshipRelatableRepository(
	client redis.UniversalClient,
	adapter relatables.Adapter,
) relatables.Repository {
	return &relationshipRelatableRepository{
		client:  client,
		adapter: adapter,
	}
}

// NewRelationshipRelatableCandidateRepository creates a new redis relationship relatable candidate repository
func NewRelationshipRelatableCandidateRepository(
	client redis.UniversalClient,
	adapter relatables.Adapter,
) relatables.CandidateRepository {
	return &relationshipRelatableCandidateRepository{
		client:  client,
		adapter: adapter,
	}
}
