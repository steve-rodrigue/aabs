package topics

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	domain_posts "github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters/clusterables"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/namers"
)

var errTest = errors.New("test error")

func TestNewBuilder(t *testing.T) {
	builder := NewBuilder(
		NewMockTopicAdapter(),
		namers.NewMockNamer(),
	)

	if builder == nil {
		t.Fatalf("expected builder")
	}
}

func TestBuilderBuild(t *testing.T) {
	ctx := context.Background()

	adapter := NewMockTopicAdapter()
	namer := namers.NewMockNamer()
	namer.NameValue = "Electric Vehicles"

	builder := NewBuilder(
		adapter,
		namer,
	)

	cluster := clusters.NewMockCluster(
		clusterables.NewMockClusterable(clusterables.TopicKind),
		clusterables.PostKind,
		[]uuid.UUID{uuid.New()},
	)

	first := newMockClusteredPost(
		"EV subsidies are changing the auto market",
		cluster,
	)

	second := newMockClusteredPost(
		"Electric cars are getting cheaper",
		cluster,
	)

	topic := NewMockTopic(
		"Electric Vehicles",
		"",
	)

	adapter.ToDomainValue = topic

	result, err := builder.Build(
		ctx,
		[]domain_posts.Post{
			first,
			second,
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	if namer.NameCalls != 1 {
		t.Fatalf("expected 1 namer call")
	}

	if namer.LastContext != ctx {
		t.Fatalf("expected context to be passed to namer")
	}

	if len(namer.LastPosts) != 2 {
		t.Fatalf("expected 2 posts to be passed to namer")
	}

	if namer.LastPosts[0] != first {
		t.Fatalf("expected first post to be passed to namer")
	}

	if adapter.ToDomainCalls != 1 {
		t.Fatalf("expected 1 adapter call")
	}

	if adapter.LastInput.Cluster != cluster {
		t.Fatalf("expected cluster to be passed")
	}

	if adapter.LastInput.Name != "Electric Vehicles" {
		t.Fatalf(
			"expected name Electric Vehicles, got %s",
			adapter.LastInput.Name,
		)
	}

	if adapter.LastInput.Identifier == uuid.Nil {
		t.Fatalf("expected generated identifier")
	}

	if adapter.LastInput.CreatedOn.IsZero() {
		t.Fatalf("expected created on")
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 topic, got %d", len(result))
	}

	if result[0] != topic {
		t.Fatalf("expected topic result")
	}
}

func TestBuilderBuildReturnsEmptyWhenPostsAreEmpty(t *testing.T) {
	adapter := NewMockTopicAdapter()
	namer := namers.NewMockNamer()

	builder := NewBuilder(
		adapter,
		namer,
	)

	result, err := builder.Build(
		context.Background(),
		[]domain_posts.Post{},
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 0 {
		t.Fatalf("expected empty result")
	}

	if namer.NameCalls != 0 {
		t.Fatalf("expected no namer call")
	}

	if adapter.ToDomainCalls != 0 {
		t.Fatalf("expected no adapter call")
	}
}

func TestBuilderBuildReturnsInvalidClusterErrorWhenPostsAreNilOnly(t *testing.T) {
	builder := NewBuilder(
		NewMockTopicAdapter(),
		namers.NewMockNamer(),
	)

	_, err := builder.Build(
		context.Background(),
		[]domain_posts.Post{
			nil,
		},
	)

	if !errors.Is(err, ErrInvalidTopicBuilderCluster) {
		t.Fatalf(
			"expected invalid topic builder cluster error, got %v",
			err,
		)
	}
}

func TestBuilderBuildReturnsInvalidClusterErrorWhenPostHasNoCluster(t *testing.T) {
	builder := NewBuilder(
		NewMockTopicAdapter(),
		namers.NewMockNamer(),
	)

	_, err := builder.Build(
		context.Background(),
		[]domain_posts.Post{
			domain_posts.NewMockPost("hello"),
		},
	)

	if !errors.Is(err, ErrInvalidTopicBuilderCluster) {
		t.Fatalf(
			"expected invalid topic builder cluster error, got %v",
			err,
		)
	}
}

func TestBuilderBuildUsesFirstPostCluster(t *testing.T) {
	ctx := context.Background()

	adapter := NewMockTopicAdapter()
	namer := namers.NewMockNamer()
	namer.NameValue = "Electric Vehicles"

	builder := NewBuilder(
		adapter,
		namer,
	)

	firstCluster := clusters.NewMockCluster(
		clusterables.NewMockClusterable(clusterables.TopicKind),
		clusterables.PostKind,
		[]uuid.UUID{uuid.New()},
	)

	secondCluster := clusters.NewMockCluster(
		clusterables.NewMockClusterable(clusterables.TopicKind),
		clusterables.PostKind,
		[]uuid.UUID{uuid.New()},
	)

	first := newMockClusteredPost(
		"first",
		firstCluster,
	)

	second := newMockClusteredPost(
		"second",
		secondCluster,
	)

	_, err := builder.Build(
		ctx,
		[]domain_posts.Post{
			first,
			second,
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	if adapter.LastInput.Cluster != firstCluster {
		t.Fatalf("expected first post cluster")
	}
}

func TestBuilderBuildSkipsNilPostsUntilClusterIsFound(t *testing.T) {
	ctx := context.Background()

	adapter := NewMockTopicAdapter()
	namer := namers.NewMockNamer()
	namer.NameValue = "Electric Vehicles"

	builder := NewBuilder(
		adapter,
		namer,
	)

	cluster := clusters.NewMockCluster(
		clusterables.NewMockClusterable(clusterables.TopicKind),
		clusterables.PostKind,
		[]uuid.UUID{uuid.New()},
	)

	post := newMockClusteredPost(
		"hello",
		cluster,
	)

	_, err := builder.Build(
		ctx,
		[]domain_posts.Post{
			nil,
			post,
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	if adapter.LastInput.Cluster != cluster {
		t.Fatalf("expected cluster")
	}
}

func TestBuilderBuildReturnsNamerError(t *testing.T) {
	adapter := NewMockTopicAdapter()
	namer := namers.NewMockNamer()
	namer.NameErr = errTest

	builder := NewBuilder(
		adapter,
		namer,
	)

	_, err := builder.Build(
		context.Background(),
		[]domain_posts.Post{
			newMockClusteredPost(
				"hello",
				clusters.NewMockCluster(
					clusterables.NewMockClusterable(clusterables.TopicKind),
					clusterables.PostKind,
					[]uuid.UUID{uuid.New()},
				),
			),
		},
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected namer error, got %v", err)
	}

	if adapter.ToDomainCalls != 0 {
		t.Fatalf("expected no adapter call")
	}
}

func TestBuilderBuildReturnsAdapterError(t *testing.T) {
	adapter := NewMockTopicAdapter()
	adapter.ToDomainErr = errTest

	namer := namers.NewMockNamer()
	namer.NameValue = "Electric Vehicles"

	builder := NewBuilder(
		adapter,
		namer,
	)

	_, err := builder.Build(
		context.Background(),
		[]domain_posts.Post{
			newMockClusteredPost(
				"hello",
				clusters.NewMockCluster(
					clusterables.NewMockClusterable(clusterables.TopicKind),
					clusterables.PostKind,
					[]uuid.UUID{uuid.New()},
				),
			),
		},
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected adapter error, got %v", err)
	}
}

func TestFirstPostCluster(t *testing.T) {
	cluster := clusters.NewMockCluster(
		clusterables.NewMockClusterable(clusterables.TopicKind),
		clusterables.PostKind,
		[]uuid.UUID{uuid.New()},
	)

	result := firstPostCluster(
		[]domain_posts.Post{
			newMockClusteredPost(
				"hello",
				cluster,
			),
		},
	)

	if result != cluster {
		t.Fatalf("expected cluster")
	}
}

func TestFirstPostClusterReturnsNilWhenNoClusterIsFound(t *testing.T) {
	result := firstPostCluster(
		[]domain_posts.Post{
			nil,
			domain_posts.NewMockPost("hello"),
		},
	)

	if result != nil {
		t.Fatalf("expected nil cluster")
	}
}

type mockClusteredPost struct {
	domain_posts.Post
	cluster clusters.Cluster
}

func newMockClusteredPost(
	text string,
	cluster clusters.Cluster,
) domain_posts.Post {
	return &mockClusteredPost{
		Post:    domain_posts.NewMockPost(text),
		cluster: cluster,
	}
}

func (post *mockClusteredPost) Cluster() clusters.Cluster {
	return post.cluster
}
