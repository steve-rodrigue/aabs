package clusters

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	domain_clusters "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters/clusterables"
)

var errTest = errors.New("test error")

func TestBuildForTarget(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	target := clusterables.NewMockClusterable(clusterables.PostKind)
	member := clusterables.NewMockClusterable(clusterables.PostKind)
	cluster := domain_clusters.NewMockCluster(
		target,
		clusterables.PostKind,
		[]uuid.UUID{member.Identifier()},
	)

	fixture.detector.DetectValue = []domain_clusters.Cluster{
		cluster,
	}

	result, err := fixture.application.BuildForTarget(
		ctx,
		target,
		[]clusterables.Clusterable{member},
	)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.detector.DetectCalls != 1 {
		t.Fatalf("expected 1 detect call")
	}

	if fixture.detector.LastContext != ctx {
		t.Fatalf("expected context to be passed")
	}

	if fixture.detector.LastTarget != target {
		t.Fatalf("expected target to be passed")
	}

	if len(fixture.detector.LastMembers) != 1 ||
		fixture.detector.LastMembers[0] != member {
		t.Fatalf("expected members to be passed")
	}

	if len(result) != 1 || result[0] != cluster {
		t.Fatalf("expected cluster result")
	}
}

func TestBuildForTargetReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.detector.DetectErr = errTest

	_, err := fixture.application.BuildForTarget(
		context.Background(),
		clusterables.NewMockClusterable(clusterables.PostKind),
		nil,
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestFindByID(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	cluster := domain_clusters.NewMockCluster(
		clusterables.NewMockClusterable(clusterables.PostKind),
		clusterables.PostKind,
		[]uuid.UUID{uuid.New()},
	)

	fixture.repository.Items[cluster.Identifier()] = cluster

	result, err := fixture.application.FindByID(
		ctx,
		cluster.Identifier(),
	)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindByIDCalls != 1 {
		t.Fatalf("expected 1 find by id call")
	}

	if fixture.repository.LastContext != ctx {
		t.Fatalf("expected context to be passed")
	}

	if fixture.repository.LastID != cluster.Identifier() {
		t.Fatalf("expected id to be passed")
	}

	if result != cluster {
		t.Fatalf("expected cluster result")
	}
}

func TestFindByIDReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.repository.FindByIDErr = errTest

	_, err := fixture.application.FindByID(
		context.Background(),
		uuid.New(),
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestFindByTarget(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	target := clusterables.NewMockClusterable(clusterables.PostKind)

	cluster := domain_clusters.NewMockCluster(
		target,
		clusterables.PostKind,
		[]uuid.UUID{uuid.New()},
	)

	fixture.repository.FindByTargetValue = []domain_clusters.Cluster{
		cluster,
	}

	result, err := fixture.application.FindByTarget(
		ctx,
		target,
	)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindByTargetCalls != 1 {
		t.Fatalf("expected 1 find by target call")
	}

	if fixture.repository.LastContext != ctx {
		t.Fatalf("expected context to be passed")
	}

	if fixture.repository.LastTarget != target.Identifier() {
		t.Fatalf("expected target id to be passed")
	}

	if len(result) != 1 || result[0] != cluster {
		t.Fatalf("expected cluster result")
	}
}

func TestFindByTargetReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.repository.FindByTargetErr = errTest

	_, err := fixture.application.FindByTarget(
		context.Background(),
		clusterables.NewMockClusterable(clusterables.PostKind),
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestFindByMember(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	member := clusterables.NewMockClusterable(clusterables.PostKind)

	cluster := domain_clusters.NewMockCluster(
		clusterables.NewMockClusterable(clusterables.PostKind),
		clusterables.PostKind,
		[]uuid.UUID{member.Identifier()},
	)

	fixture.repository.FindByMemberValue = []domain_clusters.Cluster{
		cluster,
	}

	result, err := fixture.application.FindByMember(
		ctx,
		member,
	)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindByMemberCalls != 1 {
		t.Fatalf("expected 1 find by member call")
	}

	if fixture.repository.LastContext != ctx {
		t.Fatalf("expected context to be passed")
	}

	if fixture.repository.LastMember != member.Identifier() {
		t.Fatalf("expected member id to be passed")
	}

	if len(result) != 1 || result[0] != cluster {
		t.Fatalf("expected cluster result")
	}
}

func TestFindByMemberReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.repository.FindByMemberErr = errTest

	_, err := fixture.application.FindByMember(
		context.Background(),
		clusterables.NewMockClusterable(clusterables.PostKind),
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestRebuildPostClusters(t *testing.T) {
	assertRebuildKind(t, clusterables.PostKind, func(
		ctx context.Context,
		app Application,
	) error {
		return app.RebuildPostClusters(ctx)
	})
}

func TestRebuildUserClusters(t *testing.T) {
	assertRebuildKind(t, clusterables.UserKind, func(
		ctx context.Context,
		app Application,
	) error {
		return app.RebuildUserClusters(ctx)
	})
}

func TestRebuildCommunityClusters(t *testing.T) {
	assertRebuildKind(t, clusterables.CommunityKind, func(
		ctx context.Context,
		app Application,
	) error {
		return app.RebuildCommunityClusters(ctx)
	})
}

func TestRebuildPlatformClusters(t *testing.T) {
	assertRebuildKind(t, clusterables.PlatformKind, func(
		ctx context.Context,
		app Application,
	) error {
		return app.RebuildPlatformClusters(ctx)
	})
}

func TestRebuildCampaignClusters(t *testing.T) {
	assertRebuildKind(t, clusterables.CampaignKind, func(
		ctx context.Context,
		app Application,
	) error {
		return app.RebuildCampaignClusters(ctx)
	})
}

func TestRebuildTopicClusters(t *testing.T) {
	assertRebuildKind(t, clusterables.TopicKind, func(
		ctx context.Context,
		app Application,
	) error {
		return app.RebuildTopicClusters(ctx)
	})
}

func TestRebuildNarrativeClusters(t *testing.T) {
	assertRebuildKind(t, clusterables.NarrativeKind, func(
		ctx context.Context,
		app Application,
	) error {
		return app.RebuildNarrativeClusters(ctx)
	})
}

func TestRebuildAll(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	err := fixture.application.RebuildAll(ctx)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.clusterables.FindByKindAfterCalls != 7 {
		t.Fatalf(
			"expected 7 find by kind after calls, got %d",
			fixture.clusterables.FindByKindAfterCalls,
		)
	}
}

func TestRebuildAllReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.clusterables.FindByKindAfterErr = errTest

	err := fixture.application.RebuildAll(
		context.Background(),
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestRebuildClustersReturnsClusterableRepositoryError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.clusterables.FindByKindAfterErr = errTest

	err := fixture.application.RebuildPostClusters(
		context.Background(),
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestRebuildClustersReturnsCandidateRepositoryError(t *testing.T) {
	fixture := newApplicationFixture()

	target := clusterables.NewMockClusterable(clusterables.PostKind)

	fixture.clusterables.FindByKindAfterValue = []clusterables.Clusterable{
		target,
	}

	fixture.candidates.FindCandidatesErr = errTest

	err := fixture.application.RebuildPostClusters(
		context.Background(),
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestRebuildClustersReturnsDetectorError(t *testing.T) {
	fixture := newApplicationFixture()

	target := clusterables.NewMockClusterable(clusterables.PostKind)

	fixture.clusterables.FindByKindAfterValue = []clusterables.Clusterable{
		target,
	}

	fixture.detector.DetectErr = errTest

	err := fixture.application.RebuildPostClusters(
		context.Background(),
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestRebuildClustersReturnsRepositorySaveError(t *testing.T) {
	fixture := newApplicationFixture()

	target := clusterables.NewMockClusterable(clusterables.PostKind)

	cluster := domain_clusters.NewMockCluster(
		target,
		clusterables.PostKind,
		[]uuid.UUID{uuid.New()},
	)

	fixture.clusterables.FindByKindAfterValue = []clusterables.Clusterable{
		target,
	}

	fixture.detector.DetectValue = []domain_clusters.Cluster{
		cluster,
	}

	fixture.repository.SaveErr = errTest

	err := fixture.application.RebuildPostClusters(
		context.Background(),
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected error, got %v", err)
	}
}

func assertRebuildKind(
	t *testing.T,
	kind clusterables.Kind,
	rebuild func(ctx context.Context, app Application) error,
) {
	t.Helper()

	fixture := newApplicationFixture()
	ctx := context.Background()

	target := clusterables.NewMockClusterable(kind)
	candidate := clusterables.NewMockClusterable(kind)

	cluster := domain_clusters.NewMockCluster(
		target,
		kind,
		[]uuid.UUID{candidate.Identifier()},
	)

	fixture.clusterables.FindByKindAfterValue = []clusterables.Clusterable{
		target,
	}

	fixture.candidates.FindCandidatesValue = []clusterables.Clusterable{
		candidate,
	}

	fixture.detector.DetectValue = []domain_clusters.Cluster{
		cluster,
	}

	err := rebuild(ctx, fixture.application)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.clusterables.FindByKindAfterCalls != 2 {
		t.Fatalf(
			"expected 2 find by kind after calls, got %d",
			fixture.clusterables.FindByKindAfterCalls,
		)
	}

	if fixture.clusterables.LastContext != ctx {
		t.Fatalf("expected context to be passed to clusterables")
	}

	if fixture.clusterables.LastKind != kind {
		t.Fatalf("expected kind %s, got %s", kind, fixture.clusterables.LastKind)
	}

	if fixture.clusterables.LastAmount != testRebuildBatchSize {
		t.Fatalf(
			"expected rebuild batch size %d, got %d",
			testRebuildBatchSize,
			fixture.clusterables.LastAmount,
		)
	}

	if fixture.candidates.FindCandidatesCalls != 1 {
		t.Fatalf("expected 1 find candidates call")
	}

	if fixture.candidates.LastContext != ctx {
		t.Fatalf("expected context to be passed to candidates")
	}

	if fixture.candidates.LastTarget != target {
		t.Fatalf("expected target to be passed to candidates")
	}

	if fixture.candidates.LastKind != kind {
		t.Fatalf("expected kind %s, got %s", kind, fixture.candidates.LastKind)
	}

	if fixture.candidates.LastAmount != testCandidateAmount {
		t.Fatalf(
			"expected candidate amount %d, got %d",
			testCandidateAmount,
			fixture.candidates.LastAmount,
		)
	}

	if fixture.detector.DetectCalls != 1 {
		t.Fatalf("expected 1 detect call")
	}

	if fixture.detector.LastContext != ctx {
		t.Fatalf("expected context to be passed to detector")
	}

	if fixture.detector.LastTarget != target {
		t.Fatalf("expected target to be passed to detector")
	}

	if len(fixture.detector.LastMembers) != 1 ||
		fixture.detector.LastMembers[0] != candidate {
		t.Fatalf("expected candidates to be passed to detector")
	}

	if fixture.repository.SaveCalls != 1 {
		t.Fatalf("expected 1 save call")
	}

	if fixture.repository.LastContext != ctx {
		t.Fatalf("expected context to be passed to repository")
	}

	if fixture.repository.LastCluster != cluster {
		t.Fatalf("expected cluster to be saved")
	}
}
