package redis

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/relatables"
)

type relationshipRelatableRepository struct {
	client  redis.UniversalClient
	adapter relatables.Adapter
}

func (repository *relationshipRelatableRepository) Save(
	ctx context.Context,
	relatable relatables.Relatable,
) error {
	member := repository.encode(relatable)

	pipe := repository.client.TxPipeline()

	pipe.ZAdd(ctx, relationshipRelatablesKey, redis.Z{
		Score:  0,
		Member: member,
	})

	pipe.ZAdd(ctx, repository.kindKey(relatable.RelationshipKind()), redis.Z{
		Score:  0,
		Member: member,
	})

	pipe.HSet(
		ctx,
		relationshipRelatableKindsKey,
		relatable.Identifier().String(),
		string(relatable.RelationshipKind()),
	)

	_, err := pipe.Exec(ctx)

	return err
}

func (repository *relationshipRelatableRepository) Delete(
	ctx context.Context,
	relatable relatables.Relatable,
) error {
	member := repository.encode(relatable)

	pipe := repository.client.TxPipeline()

	pipe.ZRem(ctx, relationshipRelatablesKey, member)
	pipe.ZRem(ctx, repository.kindKey(relatable.RelationshipKind()), member)
	pipe.HDel(ctx, relationshipRelatableKindsKey, relatable.Identifier().String())

	_, err := pipe.Exec(ctx)

	return err
}

func (repository *relationshipRelatableRepository) DeleteByID(
	ctx context.Context,
	id uuid.UUID,
) error {
	kind, err := repository.client.HGet(
		ctx,
		relationshipRelatableKindsKey,
		id.String(),
	).Result()
	if err == redis.Nil {
		return nil
	}

	if err != nil {
		return err
	}

	relatable, err := repository.adapter.ToDomain(
		relatables.RelatableInput{
			Identifier:       id,
			RelationshipKind: relatables.Kind(kind),
		},
	)
	if err != nil {
		return err
	}

	return repository.Delete(ctx, relatable)
}

func (repository *relationshipRelatableRepository) Find(
	ctx context.Context,
	index int,
	amount int,
) ([]relatables.Relatable, error) {
	members, err := repository.client.ZRange(
		ctx,
		relationshipRelatablesKey,
		int64(index),
		int64(index+amount-1),
	).Result()
	if err != nil {
		return nil, err
	}

	return repository.decodeMany(members)
}

func (repository *relationshipRelatableRepository) FindAfter(
	ctx context.Context,
	cursor uuid.UUID,
	amount int,
) ([]relatables.Relatable, error) {
	if cursor == uuid.Nil {
		return repository.Find(ctx, 0, amount)
	}

	kind, err := repository.client.HGet(
		ctx,
		relationshipRelatableKindsKey,
		cursor.String(),
	).Result()
	if err == redis.Nil {
		return []relatables.Relatable{}, nil
	}

	if err != nil {
		return nil, err
	}

	cursorMember := string(kind) + ":" + cursor.String()

	members, err := repository.client.ZRangeByLex(
		ctx,
		relationshipRelatablesKey,
		&redis.ZRangeBy{
			Min:    "(" + cursorMember,
			Max:    "+",
			Offset: 0,
			Count:  int64(amount),
		},
	).Result()
	if err != nil {
		return nil, err
	}

	return repository.decodeMany(members)
}

func (repository *relationshipRelatableRepository) Count(
	ctx context.Context,
) (int64, error) {
	return repository.client.ZCard(ctx, relationshipRelatablesKey).Result()
}

func (repository *relationshipRelatableRepository) FindByKind(
	ctx context.Context,
	kind relatables.Kind,
	index int,
	amount int,
) ([]relatables.Relatable, error) {
	members, err := repository.client.ZRange(
		ctx,
		repository.kindKey(kind),
		int64(index),
		int64(index+amount-1),
	).Result()
	if err != nil {
		return nil, err
	}

	return repository.decodeMany(members)
}

func (repository *relationshipRelatableRepository) CountByKind(
	ctx context.Context,
	kind relatables.Kind,
) (int64, error) {
	return repository.client.ZCard(ctx, repository.kindKey(kind)).Result()
}

func (repository *relationshipRelatableRepository) encode(
	relatable relatables.Relatable,
) string {
	return string(relatable.RelationshipKind()) + ":" + relatable.Identifier().String()
}

func (repository *relationshipRelatableRepository) decode(
	value string,
) (relatables.Relatable, error) {
	parts := strings.SplitN(value, ":", 2)
	if len(parts) != 2 {
		return nil, relatables.ErrInvalidRelatableRelationshipKind
	}

	id, err := uuid.Parse(parts[1])
	if err != nil {
		return nil, relatables.ErrInvalidRelatableIdentifier
	}

	return repository.adapter.ToDomain(
		relatables.RelatableInput{
			Identifier:       id,
			RelationshipKind: relatables.Kind(parts[0]),
		},
	)
}

func (repository *relationshipRelatableRepository) decodeMany(
	values []string,
) ([]relatables.Relatable, error) {
	out := make([]relatables.Relatable, 0, len(values))

	for _, value := range values {
		relatable, err := repository.decode(value)
		if err != nil {
			return nil, err
		}

		out = append(out, relatable)
	}

	return out, nil
}

func (repository *relationshipRelatableRepository) kindKey(
	kind relatables.Kind,
) string {
	return "relationship:relatables:kind:" + string(kind)
}
