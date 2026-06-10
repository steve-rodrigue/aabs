package campaigns

import (
	"context"
	"errors"
	"testing"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/clusterables"
	"github.com/steve-rodrigue/aabs/services/saas/domain/embeddings"
	domain_posts "github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts"
)

var errTest = errors.New("test error")

func TestClassifierClassify(t *testing.T) {
	ctx := context.Background()

	post := domain_posts.NewMockPost("hello world")
	campaign := NewMockCampaign("Campaign A", "Description A")

	campaigns := NewMockCampaignRepository()
	campaigns.Items[campaign.Identifier()] = campaign

	comparables := clusterables.NewMockComparableRepository()
	comparables.FindNearestValue = []clusterables.Comparable{
		clusterables.NewMockComparableWithID(
			campaign.Identifier(),
			clusterables.CampaignKind,
			[]float32{1, 0},
		),
	}

	comparableAdapter := clusterables.NewMockComparableAdapter()

	embedder := embeddings.NewMockEmbedder()
	embedder.Vector = embeddings.Vector{1, 0}

	classifier := &classifier{
		campaigns:         campaigns,
		comparables:       comparables,
		comparableAdapter: comparableAdapter,
		embedder:          embedder,
		threshold:         0.7,
	}

	result, confidence, err := classifier.Classify(ctx, post)
	if err != nil {
		t.Fatal(err)
	}

	if result != campaign {
		t.Fatalf("expected campaign")
	}

	if confidence != 1 {
		t.Fatalf("expected confidence 1, got %f", confidence)
	}

	if embedder.EmbedCalls != 1 {
		t.Fatalf("expected 1 embed call")
	}

	if embedder.LastContext != ctx {
		t.Fatalf("expected context")
	}

	if embedder.LastText != "hello world" {
		t.Fatalf("expected text hello world, got %s", embedder.LastText)
	}

	if comparableAdapter.ToDomainCalls != 1 {
		t.Fatalf("expected 1 comparable adapter call")
	}

	if comparableAdapter.LastInput.Clusterable.Identifier != post.Identifier() {
		t.Fatalf("expected post identifier")
	}

	if comparableAdapter.LastInput.Clusterable.ClusterKind != clusterables.PostKind {
		t.Fatalf("expected post kind")
	}

	if comparables.FindNearestCalls != 1 {
		t.Fatalf("expected 1 find nearest call")
	}

	if comparables.LastContext != ctx {
		t.Fatalf("expected context")
	}

	if comparables.LastKind != clusterables.CampaignKind {
		t.Fatalf("expected campaign kind")
	}

	if comparables.LastAmount != 1 {
		t.Fatalf("expected amount 1")
	}

	if campaigns.FindByIDCalls != 1 {
		t.Fatalf("expected 1 find by id call")
	}

	if campaigns.LastContext != ctx {
		t.Fatalf("expected context")
	}

	if campaigns.LastID != campaign.Identifier() {
		t.Fatalf("expected campaign id")
	}
}

func TestClassifierClassifyReturnsInvalidPostErrorWhenPostIsNil(t *testing.T) {
	classifier := &classifier{}

	_, _, err := classifier.Classify(context.Background(), nil)

	if !errors.Is(err, ErrInvalidCampaignClassifierPost) {
		t.Fatalf("expected invalid post error, got %v", err)
	}
}

func TestClassifierClassifyReturnsInvalidTextError(t *testing.T) {
	classifier := &classifier{}

	_, _, err := classifier.Classify(
		context.Background(),
		domain_posts.NewMockPost("   "),
	)

	if !errors.Is(err, ErrInvalidCampaignClassifierText) {
		t.Fatalf("expected invalid text error, got %v", err)
	}
}

func TestClassifierClassifyReturnsEmbedderError(t *testing.T) {
	embedder := embeddings.NewMockEmbedder()
	embedder.EmbedErr = errTest

	classifier := &classifier{
		embedder: embedder,
	}

	_, _, err := classifier.Classify(
		context.Background(),
		domain_posts.NewMockPost("hello"),
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected embedder error, got %v", err)
	}
}

func TestClassifierClassifyReturnsInvalidVectorErrorWhenEmbeddingIsEmpty(t *testing.T) {
	embedder := embeddings.NewMockEmbedder()
	embedder.Vector = embeddings.Vector{}

	classifier := &classifier{
		embedder: embedder,
	}

	_, _, err := classifier.Classify(
		context.Background(),
		domain_posts.NewMockPost("hello"),
	)

	if !errors.Is(err, ErrInvalidCampaignClassifierVector) {
		t.Fatalf("expected invalid vector error, got %v", err)
	}
}

func TestClassifierClassifyReturnsComparableAdapterError(t *testing.T) {
	embedder := embeddings.NewMockEmbedder()
	embedder.Vector = embeddings.Vector{1, 0}

	comparableAdapter := clusterables.NewMockComparableAdapter()
	comparableAdapter.ToDomainErr = errTest

	classifier := &classifier{
		comparableAdapter: comparableAdapter,
		embedder:          embedder,
	}

	_, _, err := classifier.Classify(
		context.Background(),
		domain_posts.NewMockPost("hello"),
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected comparable adapter error, got %v", err)
	}
}

func TestClassifierClassifyReturnsComparableRepositoryError(t *testing.T) {
	embedder := embeddings.NewMockEmbedder()
	embedder.Vector = embeddings.Vector{1, 0}

	comparables := clusterables.NewMockComparableRepository()
	comparables.FindNearestErr = errTest

	classifier := &classifier{
		comparables:       comparables,
		comparableAdapter: clusterables.NewMockComparableAdapter(),
		embedder:          embedder,
	}

	_, _, err := classifier.Classify(
		context.Background(),
		domain_posts.NewMockPost("hello"),
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected comparable repository error, got %v", err)
	}
}

func TestClassifierClassifyReturnsNilWhenNoNearestCampaignExists(t *testing.T) {
	embedder := embeddings.NewMockEmbedder()
	embedder.Vector = embeddings.Vector{1, 0}

	classifier := &classifier{
		comparables:       clusterables.NewMockComparableRepository(),
		comparableAdapter: clusterables.NewMockComparableAdapter(),
		embedder:          embedder,
		threshold:         0.7,
	}

	result, confidence, err := classifier.Classify(
		context.Background(),
		domain_posts.NewMockPost("hello"),
	)

	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Fatalf("expected nil campaign")
	}

	if confidence != 0 {
		t.Fatalf("expected confidence 0, got %f", confidence)
	}
}

func TestClassifierClassifyReturnsInvalidComparableErrorWhenNearestIsNil(t *testing.T) {
	embedder := embeddings.NewMockEmbedder()
	embedder.Vector = embeddings.Vector{1, 0}

	comparables := clusterables.NewMockComparableRepository()
	comparables.FindNearestValue = []clusterables.Comparable{nil}

	classifier := &classifier{
		comparables:       comparables,
		comparableAdapter: clusterables.NewMockComparableAdapter(),
		embedder:          embedder,
		threshold:         0.7,
	}

	_, _, err := classifier.Classify(
		context.Background(),
		domain_posts.NewMockPost("hello"),
	)

	if !errors.Is(err, ErrInvalidCampaignClassifierComparable) {
		t.Fatalf("expected invalid comparable error, got %v", err)
	}
}

func TestClassifierClassifyReturnsInvalidVectorErrorWhenVectorsMismatch(t *testing.T) {
	embedder := embeddings.NewMockEmbedder()
	embedder.Vector = embeddings.Vector{1, 0}

	comparables := clusterables.NewMockComparableRepository()
	comparables.FindNearestValue = []clusterables.Comparable{
		clusterables.NewMockComparable(
			clusterables.CampaignKind,
			[]float32{1, 0, 0},
		),
	}

	classifier := &classifier{
		comparables:       comparables,
		comparableAdapter: clusterables.NewMockComparableAdapter(),
		embedder:          embedder,
		threshold:         0.7,
	}

	_, _, err := classifier.Classify(
		context.Background(),
		domain_posts.NewMockPost("hello"),
	)

	if !errors.Is(err, ErrInvalidCampaignClassifierVector) {
		t.Fatalf("expected invalid vector error, got %v", err)
	}
}

func TestClassifierClassifyReturnsNilWhenConfidenceIsBelowThreshold(t *testing.T) {
	embedder := embeddings.NewMockEmbedder()
	embedder.Vector = embeddings.Vector{1, 0}

	comparables := clusterables.NewMockComparableRepository()
	comparables.FindNearestValue = []clusterables.Comparable{
		clusterables.NewMockComparable(
			clusterables.CampaignKind,
			[]float32{0, 1},
		),
	}

	campaigns := NewMockCampaignRepository()

	classifier := &classifier{
		campaigns:         campaigns,
		comparables:       comparables,
		comparableAdapter: clusterables.NewMockComparableAdapter(),
		embedder:          embedder,
		threshold:         0.7,
	}

	result, confidence, err := classifier.Classify(
		context.Background(),
		domain_posts.NewMockPost("hello"),
	)

	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Fatalf("expected nil campaign")
	}

	if confidence != 0 {
		t.Fatalf("expected confidence 0, got %f", confidence)
	}

	if campaigns.FindByIDCalls != 0 {
		t.Fatalf("expected no campaign lookup")
	}
}

func TestClassifierClassifyReturnsCampaignRepositoryError(t *testing.T) {
	embedder := embeddings.NewMockEmbedder()
	embedder.Vector = embeddings.Vector{1, 0}

	campaign := NewMockCampaign("Campaign A", "Description A")

	campaigns := NewMockCampaignRepository()
	campaigns.FindByIDErr = errTest

	comparables := clusterables.NewMockComparableRepository()
	comparables.FindNearestValue = []clusterables.Comparable{
		clusterables.NewMockComparableWithID(
			campaign.Identifier(),
			clusterables.CampaignKind,
			[]float32{1, 0},
		),
	}

	classifier := &classifier{
		campaigns:         campaigns,
		comparables:       comparables,
		comparableAdapter: clusterables.NewMockComparableAdapter(),
		embedder:          embedder,
		threshold:         0.7,
	}

	_, _, err := classifier.Classify(
		context.Background(),
		domain_posts.NewMockPost("hello"),
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected campaign repository error, got %v", err)
	}
}

func TestCampaignConfidence(t *testing.T) {
	confidence, err := campaignConfidence(
		[]float32{1, 0},
		[]float32{1, 0},
	)

	if err != nil {
		t.Fatal(err)
	}

	if confidence != 1 {
		t.Fatalf("expected confidence 1, got %f", confidence)
	}
}

func TestCampaignConfidenceClampsNegativeToZero(t *testing.T) {
	confidence, err := campaignConfidence(
		[]float32{-1, 0},
		[]float32{1, 0},
	)

	if err != nil {
		t.Fatal(err)
	}

	if confidence != 0 {
		t.Fatalf("expected confidence 0, got %f", confidence)
	}
}

func TestCampaignConfidenceReturnsInvalidVectorErrorWhenEmpty(t *testing.T) {
	_, err := campaignConfidence(
		[]float32{},
		[]float32{1, 0},
	)

	if !errors.Is(err, ErrInvalidCampaignClassifierVector) {
		t.Fatalf("expected invalid vector error, got %v", err)
	}
}

func TestCampaignConfidenceReturnsInvalidVectorErrorWhenMismatch(t *testing.T) {
	_, err := campaignConfidence(
		[]float32{1, 0},
		[]float32{1, 0, 0},
	)

	if !errors.Is(err, ErrInvalidCampaignClassifierVector) {
		t.Fatalf("expected invalid vector error, got %v", err)
	}
}

func TestCampaignCosineSimilarityReturnsZeroWhenSourceMagnitudeIsZero(t *testing.T) {
	confidence := campaignCosineSimilarity(
		[]float32{0, 0},
		[]float32{1, 0},
	)

	if confidence != 0 {
		t.Fatalf("expected confidence 0, got %f", confidence)
	}
}

func TestCampaignCosineSimilarityReturnsZeroWhenTargetMagnitudeIsZero(t *testing.T) {
	confidence := campaignCosineSimilarity(
		[]float32{1, 0},
		[]float32{0, 0},
	)

	if confidence != 0 {
		t.Fatalf("expected confidence 0, got %f", confidence)
	}
}
