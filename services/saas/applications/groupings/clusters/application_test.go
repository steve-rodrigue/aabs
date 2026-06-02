package clusters

import (
	"errors"
	"testing"

	"github.com/google/uuid"

	domain_clusters "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters/clusterables"
)

var errTest = errors.New("test error")

func TestBuildForTarget(t *testing.T) {
	fixture := newApplicationFixture()

	target := clusterables.NewMockClusterable(clusterables.PostKind)
	member := clusterables.NewMockClusterable(clusterables.PostKind)

	cluster := domain_clusters.NewMockCluster(
		target,
		clusterables.PostKind,
		[]uuid.UUID{member.Identifier()},
	)

	fixture.detector.DetectValue = []domain_clusters.Cluster{cluster}

	result, err := fixture.application.BuildForTarget(
		target,
		[]clusterables.Clusterable{member},
	)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.detector.DetectCalls != 1 {
		t.Fatalf("expected 1 detect call")
	}

	if len(result) != 1 || result[0] != cluster {
		t.Fatalf("expected cluster result")
	}
}

func TestBuildForTargetReturnsDetectorError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.detector.DetectErr = errTest

	_, err := fixture.application.BuildForTarget(
		clusterables.NewMockClusterable(clusterables.PostKind),
		nil,
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected detector error, got %v", err)
	}
}

func TestFindByID(t *testing.T) {
	fixture := newApplicationFixture()

	target := clusterables.NewMockClusterable(clusterables.PostKind)
	cluster := domain_clusters.NewMockCluster(target, clusterables.PostKind, nil)

	fixture.repository.Items[cluster.Identifier()] = cluster

	result, err := fixture.application.FindByID(cluster.Identifier())

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindByIDCalls != 1 {
		t.Fatalf("expected 1 find by id call")
	}

	if result != cluster {
		t.Fatalf("expected cluster result")
	}
}

func TestFindByIDReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.repository.FindByIDErr = errTest

	_, err := fixture.application.FindByID(uuid.New())

	if !errors.Is(err, errTest) {
		t.Fatalf("expected find by id error, got %v", err)
	}
}

func TestFindByTarget(t *testing.T) {
	fixture := newApplicationFixture()

	target := clusterables.NewMockClusterable(clusterables.PostKind)
	cluster := domain_clusters.NewMockCluster(target, clusterables.PostKind, nil)

	fixture.repository.FindByTargetValue = []domain_clusters.Cluster{cluster}

	result, err := fixture.application.FindByTarget(target)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindByTargetCalls != 1 {
		t.Fatalf("expected 1 find by target call")
	}

	if len(result) != 1 || result[0] != cluster {
		t.Fatalf("expected cluster result")
	}
}

func TestFindByTargetReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.repository.FindByTargetErr = errTest

	_, err := fixture.application.FindByTarget(
		clusterables.NewMockClusterable(clusterables.PostKind),
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected find by target error, got %v", err)
	}
}

func TestFindByMember(t *testing.T) {
	fixture := newApplicationFixture()

	target := clusterables.NewMockClusterable(clusterables.PostKind)
	member := clusterables.NewMockClusterable(clusterables.PostKind)

	cluster := domain_clusters.NewMockCluster(
		target,
		clusterables.PostKind,
		[]uuid.UUID{member.Identifier()},
	)

	fixture.repository.FindByMemberValue = []domain_clusters.Cluster{cluster}

	result, err := fixture.application.FindByMember(member)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindByMemberCalls != 1 {
		t.Fatalf("expected 1 find by member call")
	}

	if len(result) != 1 || result[0] != cluster {
		t.Fatalf("expected cluster result")
	}
}

func TestFindByMemberReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.repository.FindByMemberErr = errTest

	_, err := fixture.application.FindByMember(
		clusterables.NewMockClusterable(clusterables.PostKind),
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected find by member error, got %v", err)
	}
}

func TestRebuildPostClusters(t *testing.T) {
	fixture := newApplicationFixture()

	target := clusterables.NewMockClusterable(clusterables.PostKind)
	candidate := clusterables.NewMockClusterable(clusterables.PostKind)

	cluster := domain_clusters.NewMockCluster(
		target,
		clusterables.PostKind,
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

	err := fixture.application.RebuildPostClusters()

	if err != nil {
		t.Fatal(err)
	}

	if fixture.clusterables.FindByKindAfterCalls != 2 {
		t.Fatalf("expected 2 find by kind after calls, got %d", fixture.clusterables.FindByKindAfterCalls)
	}

	if fixture.candidates.FindCandidatesCalls != 1 {
		t.Fatalf("expected 1 find candidates call")
	}

	if fixture.detector.DetectCalls != 1 {
		t.Fatalf("expected 1 detect call")
	}

	if fixture.repository.SaveCalls != 1 {
		t.Fatalf("expected 1 save call")
	}
}

func TestRebuildAll(t *testing.T) {
	fixture := newApplicationFixture()

	err := fixture.application.RebuildAll()

	if err != nil {
		t.Fatal(err)
	}

	if fixture.clusterables.FindByKindAfterCalls != 7 {
		t.Fatalf("expected 7 rebuild calls, got %d", fixture.clusterables.FindByKindAfterCalls)
	}
}

func TestRebuildAllReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.clusterables.FindByKindAfterErr = errTest

	err := fixture.application.RebuildAll()

	if !errors.Is(err, errTest) {
		t.Fatalf("expected rebuild all error, got %v", err)
	}
}

func TestRebuildAllReturnsUserClustersError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.clusterables.FailOnCall = 2
	fixture.clusterables.FindByKindAfterErr = errTest

	err := fixture.application.RebuildAll()

	if !errors.Is(err, errTest) {
		t.Fatalf("expected user clusters error, got %v", err)
	}
}

func TestRebuildAllReturnsCommunityClustersError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.clusterables.FailOnCall = 3
	fixture.clusterables.FindByKindAfterErr = errTest

	err := fixture.application.RebuildAll()

	if !errors.Is(err, errTest) {
		t.Fatalf("expected community clusters error, got %v", err)
	}
}

func TestRebuildAllReturnsPlatformClustersError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.clusterables.FailOnCall = 4
	fixture.clusterables.FindByKindAfterErr = errTest

	err := fixture.application.RebuildAll()

	if !errors.Is(err, errTest) {
		t.Fatalf("expected platform clusters error, got %v", err)
	}
}

func TestRebuildAllReturnsCampaignClustersError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.clusterables.FailOnCall = 5
	fixture.clusterables.FindByKindAfterErr = errTest

	err := fixture.application.RebuildAll()

	if !errors.Is(err, errTest) {
		t.Fatalf("expected campaign clusters error, got %v", err)
	}
}

func TestRebuildAllReturnsTopicClustersError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.clusterables.FailOnCall = 6
	fixture.clusterables.FindByKindAfterErr = errTest

	err := fixture.application.RebuildAll()

	if !errors.Is(err, errTest) {
		t.Fatalf("expected topic clusters error, got %v", err)
	}
}

func TestRebuildAllReturnsNarrativeClustersError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.clusterables.FailOnCall = 7
	fixture.clusterables.FindByKindAfterErr = errTest

	err := fixture.application.RebuildAll()

	if !errors.Is(err, errTest) {
		t.Fatalf("expected narrative clusters error, got %v", err)
	}
}

func TestRebuildUserClusters(t *testing.T) {
	assertRebuildKind(t, clusterables.UserKind, func(app Application) error {
		return app.RebuildUserClusters()
	})
}

func TestRebuildCommunityClusters(t *testing.T) {
	assertRebuildKind(t, clusterables.CommunityKind, func(app Application) error {
		return app.RebuildCommunityClusters()
	})
}

func TestRebuildPlatformClusters(t *testing.T) {
	assertRebuildKind(t, clusterables.PlatformKind, func(app Application) error {
		return app.RebuildPlatformClusters()
	})
}

func TestRebuildCampaignClusters(t *testing.T) {
	assertRebuildKind(t, clusterables.CampaignKind, func(app Application) error {
		return app.RebuildCampaignClusters()
	})
}

func TestRebuildTopicClusters(t *testing.T) {
	assertRebuildKind(t, clusterables.TopicKind, func(app Application) error {
		return app.RebuildTopicClusters()
	})
}

func TestRebuildNarrativeClusters(t *testing.T) {
	assertRebuildKind(t, clusterables.NarrativeKind, func(app Application) error {
		return app.RebuildNarrativeClusters()
	})
}

func TestRebuildReturnsClusterableRepositoryError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.clusterables.FindByKindAfterErr = errTest

	err := fixture.application.RebuildPostClusters()

	if !errors.Is(err, errTest) {
		t.Fatalf("expected clusterable repository error, got %v", err)
	}
}

func TestRebuildReturnsCandidateRepositoryError(t *testing.T) {
	fixture := newApplicationFixture()

	fixture.clusterables.FindByKindAfterValue = []clusterables.Clusterable{
		clusterables.NewMockClusterable(clusterables.PostKind),
	}
	fixture.candidates.FindCandidatesErr = errTest

	err := fixture.application.RebuildPostClusters()

	if !errors.Is(err, errTest) {
		t.Fatalf("expected candidate repository error, got %v", err)
	}
}

func TestRebuildReturnsDetectorError(t *testing.T) {
	fixture := newApplicationFixture()

	fixture.clusterables.FindByKindAfterValue = []clusterables.Clusterable{
		clusterables.NewMockClusterable(clusterables.PostKind),
	}
	fixture.detector.DetectErr = errTest

	err := fixture.application.RebuildPostClusters()

	if !errors.Is(err, errTest) {
		t.Fatalf("expected detector error, got %v", err)
	}
}

func TestRebuildReturnsSaveError(t *testing.T) {
	fixture := newApplicationFixture()

	target := clusterables.NewMockClusterable(clusterables.PostKind)
	cluster := domain_clusters.NewMockCluster(target, clusterables.PostKind, nil)

	fixture.clusterables.FindByKindAfterValue = []clusterables.Clusterable{
		target,
	}
	fixture.detector.DetectValue = []domain_clusters.Cluster{
		cluster,
	}
	fixture.repository.SaveErr = errTest

	err := fixture.application.RebuildPostClusters()

	if !errors.Is(err, errTest) {
		t.Fatalf("expected save error, got %v", err)
	}
}

func assertRebuildKind(
	t *testing.T,
	kind clusterables.Kind,
	rebuild func(app Application) error,
) {
	t.Helper()

	fixture := newApplicationFixture()

	err := rebuild(fixture.application)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.clusterables.FindByKindAfterCalls != 1 {
		t.Fatalf("expected 1 find by kind after call")
	}

	if fixture.clusterables.LastKind != kind {
		t.Fatalf("expected kind %s, got %s", kind, fixture.clusterables.LastKind)
	}
}
