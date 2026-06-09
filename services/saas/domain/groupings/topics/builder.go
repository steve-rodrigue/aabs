package topics

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/namers"
)

type builder struct {
	adapter Adapter
	namer   namers.Namer
}

func (builder *builder) Build(
	ctx context.Context,
	posts []posts.Post,
) ([]Topic, error) {
	if len(posts) == 0 {
		return []Topic{}, nil
	}

	cluster := firstPostCluster(posts)
	if cluster == nil {
		return nil, ErrInvalidTopicBuilderCluster
	}

	name, err := builder.namer.Name(ctx, posts)
	if err != nil {
		return nil, err
	}

	topic, err := builder.adapter.ToDomain(
		TopicInput{
			Identifier: uuid.New(),
			Cluster:    cluster,
			Name:       name,
			CreatedOn:  time.Now().UTC(),
		},
	)
	if err != nil {
		return nil, err
	}

	return []Topic{
		topic,
	}, nil
}

func firstPostCluster(
	posts []posts.Post,
) clusters.Cluster {
	for _, post := range posts {
		if post == nil {
			continue
		}

		clustered, ok := post.(interface {
			Cluster() clusters.Cluster
		})
		if !ok {
			continue
		}

		cluster := clustered.Cluster()
		if cluster != nil {
			return cluster
		}
	}

	return nil
}
