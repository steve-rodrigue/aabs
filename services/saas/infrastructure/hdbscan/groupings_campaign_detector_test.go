package hdbscan

import (
	"context"
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/clusterables"
	domain_campaigns "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/campaigns"
	domain_clusters "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters"
)

var errTest = errors.New("test error")

func TestNewGroupingsCampaignDetector(t *testing.T) {
	detector := NewGroupingsCampaignDetector(
		"http://localhost:8000",
		domain_campaigns.NewMockCampaignAdapter(),
		domain_clusters.NewAdapter(clusterables.NewAdapter()),
		clusterables.NewMockComparableRepository(),
		2,
		nil,
	)

	if detector == nil {
		t.Fatalf("expected detector")
	}
}

func TestGroupingsCampaignDetectorDetect(t *testing.T) {
	endpoint := os.Getenv("HDBSCAN_TEST_ADDR")
	if endpoint == "" {
		t.Skip("HDBSCAN_TEST_ADDR is not set")
	}

	ctx := context.Background()

	candidateA := clusterables.NewMockClusterable(clusterables.PostKind)
	candidateB := clusterables.NewMockClusterable(clusterables.PostKind)
	candidateC := clusterables.NewMockClusterable(clusterables.PostKind)
	candidateD := clusterables.NewMockClusterable(clusterables.PostKind)

	comparables := clusterables.NewMockComparableRepository()
	comparables.Items[candidateA.Identifier()] = clusterables.NewMockComparableWithID(
		candidateA.Identifier(),
		clusterables.PostKind,
		[]float32{1, 1},
	)
	comparables.Items[candidateB.Identifier()] = clusterables.NewMockComparableWithID(
		candidateB.Identifier(),
		clusterables.PostKind,
		[]float32{1.1, 1},
	)
	comparables.Items[candidateC.Identifier()] = clusterables.NewMockComparableWithID(
		candidateC.Identifier(),
		clusterables.PostKind,
		[]float32{10, 10},
	)
	comparables.Items[candidateD.Identifier()] = clusterables.NewMockComparableWithID(
		candidateD.Identifier(),
		clusterables.PostKind,
		[]float32{10.1, 10},
	)

	campaignAdapter := domain_campaigns.NewMockCampaignAdapter()

	detector := NewGroupingsCampaignDetector(
		endpoint,
		campaignAdapter,
		domain_clusters.NewAdapter(clusterables.NewAdapter()),
		comparables,
		2,
		nil,
	)

	result, err := detector.Detect(
		ctx,
		[]clusterables.Clusterable{
			candidateA,
			candidateB,
			candidateC,
			candidateD,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) == 0 {
		t.Fatalf("expected detected campaigns")
	}

	if campaignAdapter.ToDomainCalls != len(result) {
		t.Fatalf(
			"expected %d campaign adapter calls, got %d",
			len(result),
			campaignAdapter.ToDomainCalls,
		)
	}

	for _, campaign := range result {
		if campaign.Identifier() == uuid.Nil {
			t.Fatalf("expected campaign identifier")
		}

		if campaign.Cluster() == nil {
			t.Fatalf("expected campaign cluster")
		}

		if campaign.PostCount() <= 0 {
			t.Fatalf("expected positive post count")
		}

		if campaign.Confidence() < 0 || campaign.Confidence() > 1 {
			t.Fatalf("expected confidence between 0 and 1")
		}
	}
}

func TestGroupingsCampaignDetectorDetectReturnsInvalidEndpointError(t *testing.T) {
	detector := NewGroupingsCampaignDetector(
		"",
		domain_campaigns.NewMockCampaignAdapter(),
		domain_clusters.NewAdapter(clusterables.NewAdapter()),
		clusterables.NewMockComparableRepository(),
		2,
		nil,
	)

	_, err := detector.Detect(
		context.Background(),
		[]clusterables.Clusterable{
			clusterables.NewMockClusterable(clusterables.PostKind),
		},
	)

	if !errors.Is(err, ErrInvalidGroupingsCampaignDetectorEndpoint) {
		t.Fatalf("expected invalid endpoint error, got %v", err)
	}
}

func TestGroupingsCampaignDetectorDetectReturnsEmptyWhenNoCandidates(t *testing.T) {
	detector := NewGroupingsCampaignDetector(
		"http://localhost:8000",
		domain_campaigns.NewMockCampaignAdapter(),
		domain_clusters.NewAdapter(clusterables.NewAdapter()),
		clusterables.NewMockComparableRepository(),
		2,
		nil,
	)

	result, err := detector.Detect(
		context.Background(),
		[]clusterables.Clusterable{},
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 0 {
		t.Fatalf("expected empty result")
	}
}

func TestGroupingsCampaignDetectorDetectReturnsInvalidCandidateError(t *testing.T) {
	detector := NewGroupingsCampaignDetector(
		"http://localhost:8000",
		domain_campaigns.NewMockCampaignAdapter(),
		domain_clusters.NewAdapter(clusterables.NewAdapter()),
		clusterables.NewMockComparableRepository(),
		2,
		nil,
	)

	_, err := detector.Detect(
		context.Background(),
		[]clusterables.Clusterable{
			nil,
		},
	)

	if !errors.Is(err, ErrInvalidGroupingsCampaignDetectorCandidates) {
		t.Fatalf("expected invalid candidate error, got %v", err)
	}
}

func TestGroupingsCampaignDetectorDetectReturnsInvalidKindError(t *testing.T) {
	detector := NewGroupingsCampaignDetector(
		"http://localhost:8000",
		domain_campaigns.NewMockCampaignAdapter(),
		domain_clusters.NewAdapter(clusterables.NewAdapter()),
		clusterables.NewMockComparableRepository(),
		2,
		nil,
	)

	_, err := detector.Detect(
		context.Background(),
		[]clusterables.Clusterable{
			clusterables.NewMockClusterable(clusterables.PostKind),
			clusterables.NewMockClusterable(clusterables.UserKind),
		},
	)

	if !errors.Is(err, ErrInvalidGroupingsCampaignDetectorKind) {
		t.Fatalf("expected invalid kind error, got %v", err)
	}
}

func TestGroupingsCampaignDetectorDetectReturnsComparableRepositoryError(t *testing.T) {
	comparables := clusterables.NewMockComparableRepository()
	comparables.FindByIDErr = errTest

	detector := NewGroupingsCampaignDetector(
		"http://localhost:8000",
		domain_campaigns.NewMockCampaignAdapter(),
		domain_clusters.NewAdapter(clusterables.NewAdapter()),
		comparables,
		2,
		nil,
	)

	_, err := detector.Detect(
		context.Background(),
		[]clusterables.Clusterable{
			clusterables.NewMockClusterable(clusterables.PostKind),
		},
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected comparable repository error, got %v", err)
	}
}

func TestGroupingsCampaignDetectorDetectReturnsInvalidComparableError(t *testing.T) {
	detector := NewGroupingsCampaignDetector(
		"http://localhost:8000",
		domain_campaigns.NewMockCampaignAdapter(),
		domain_clusters.NewAdapter(clusterables.NewAdapter()),
		clusterables.NewMockComparableRepository(),
		2,
		nil,
	)

	_, err := detector.Detect(
		context.Background(),
		[]clusterables.Clusterable{
			clusterables.NewMockClusterable(clusterables.PostKind),
		},
	)

	if !errors.Is(err, ErrInvalidGroupingsCampaignDetectorComparable) {
		t.Fatalf("expected invalid comparable error, got %v", err)
	}
}

func TestGroupingsCampaignDetectorDetectReturnsInvalidVectorError(t *testing.T) {
	candidate := clusterables.NewMockClusterable(clusterables.PostKind)

	comparables := clusterables.NewMockComparableRepository()
	comparables.Items[candidate.Identifier()] = clusterables.NewMockComparableWithID(
		candidate.Identifier(),
		clusterables.PostKind,
		[]float32{},
	)

	detector := NewGroupingsCampaignDetector(
		"http://localhost:8000",
		domain_campaigns.NewMockCampaignAdapter(),
		domain_clusters.NewAdapter(clusterables.NewAdapter()),
		comparables,
		2,
		nil,
	)

	_, err := detector.Detect(
		context.Background(),
		[]clusterables.Clusterable{
			candidate,
		},
	)

	if !errors.Is(err, ErrInvalidGroupingsCampaignDetectorVector) {
		t.Fatalf("expected invalid vector error, got %v", err)
	}
}

func TestGroupingsCampaignDetectorDetectReturnsEmptyWhenNotEnoughEmbeddings(t *testing.T) {
	candidate := clusterables.NewMockClusterable(clusterables.PostKind)

	comparables := clusterables.NewMockComparableRepository()
	comparables.Items[candidate.Identifier()] = clusterables.NewMockComparableWithID(
		candidate.Identifier(),
		clusterables.PostKind,
		[]float32{1, 0},
	)

	detector := NewGroupingsCampaignDetector(
		"http://localhost:8000",
		domain_campaigns.NewMockCampaignAdapter(),
		domain_clusters.NewAdapter(clusterables.NewAdapter()),
		comparables,
		2,
		nil,
	)

	result, err := detector.Detect(
		context.Background(),
		[]clusterables.Clusterable{
			candidate,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 0 {
		t.Fatalf("expected empty result")
	}
}

func TestGroupingsCampaignDetectorClusterReturnsResponseError(t *testing.T) {
	detector := &groupingsCampaignDetector{
		endpoint: "http://127.0.0.1:1",
		client:   &http.Client{},
	}

	_, err := detector.cluster(
		context.Background(),
		clusterRequest{
			Embeddings: [][]float32{
				{1, 0},
				{0, 1},
			},
			MinClusterSize: 2,
		},
	)

	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestMemberIDsForCluster(t *testing.T) {
	first := uuid.New()
	second := uuid.New()
	third := uuid.New()

	result := memberIDsForCluster(
		clusterSummary{
			MemberIndexes: []int{
				0,
				2,
				99,
				-1,
			},
		},
		[]uuid.UUID{
			first,
			second,
			third,
		},
	)

	if len(result) != 2 {
		t.Fatalf("expected 2 ids, got %d", len(result))
	}

	if result[0] != first {
		t.Fatalf("expected first id")
	}

	if result[1] != third {
		t.Fatalf("expected third id")
	}
}

func TestAverageProbability(t *testing.T) {
	result := averageProbability(
		1,
		[]clusterResult{
			{
				ClusterID:   1,
				Probability: 0.8,
				IsNoise:     false,
			},
			{
				ClusterID:   1,
				Probability: 1,
				IsNoise:     false,
			},
			{
				ClusterID:   2,
				Probability: 0.2,
				IsNoise:     false,
			},
			{
				ClusterID:   1,
				Probability: 0.1,
				IsNoise:     true,
			},
		},
	)

	if result != 0.9 {
		t.Fatalf("expected 0.9, got %f", result)
	}
}

func TestAverageProbabilityReturnsZeroWhenNoMembers(t *testing.T) {
	result := averageProbability(
		1,
		[]clusterResult{},
	)

	if result != 0 {
		t.Fatalf("expected 0, got %f", result)
	}
}
