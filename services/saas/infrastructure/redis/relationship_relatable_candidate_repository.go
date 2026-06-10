package redis

import (
	"context"

	"github.com/redis/go-redis/v9"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/relatables"
)

type relationshipRelatableCandidateRepository struct {
	client  redis.UniversalClient
	adapter relatables.Adapter
}

func (repository *relationshipRelatableCandidateRepository) FindCandidates(
	ctx context.Context,
	source relatables.Relatable,
	amount int,
) ([]relatables.Relatable, error) {
	if amount <= 0 {
		return []relatables.Relatable{}, nil
	}

	key := repository.kindKey(source.RelationshipKind())

	members, err := repository.client.ZRange(
		ctx,
		key,
		0,
		int64(amount),
	).Result()
	if err != nil {
		return nil, err
	}

	out := []relatables.Relatable{}

	relatableRepository := &relationshipRelatableRepository{
		client:  repository.client,
		adapter: repository.adapter,
	}

	for _, member := range members {
		candidate, err := relatableRepository.decode(member)
		if err != nil {
			return nil, err
		}

		if candidate.Identifier() == source.Identifier() {
			continue
		}

		out = append(out, candidate)

		if len(out) >= amount {
			break
		}
	}

	return out, nil
}

func (repository *relationshipRelatableCandidateRepository) kindKey(
	kind relatables.Kind,
) string {
	return "relationship:relatables:kind:" + string(kind)
}
