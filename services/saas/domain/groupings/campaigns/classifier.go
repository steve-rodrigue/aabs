package campaigns

import (
	"context"
	"math"
	"strings"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/clusterables"
	"github.com/steve-rodrigue/aabs/services/saas/domain/embeddings"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts"
)

type classifier struct {
	campaigns Repository

	comparables       clusterables.ComparableRepository
	comparableAdapter clusterables.ComparableAdapter

	embedder embeddings.Embedder

	threshold float64
}

func (classifier *classifier) Classify(
	ctx context.Context,
	post posts.Post,
) (Campaign, float64, error) {
	if post == nil || post.Content() == nil {
		return nil, 0, ErrInvalidCampaignClassifierPost
	}

	text := strings.TrimSpace(post.Content().Text())
	if text == "" {
		return nil, 0, ErrInvalidCampaignClassifierText
	}

	vector, err := classifier.embedder.Embed(ctx, text)
	if err != nil {
		return nil, 0, err
	}

	if len(vector) == 0 {
		return nil, 0, ErrInvalidCampaignClassifierVector
	}

	target, err := classifier.comparableAdapter.ToDomain(
		clusterables.ComparableInput{
			Clusterable: clusterables.ClusterableInput{
				Identifier:  post.Identifier(),
				ClusterKind: clusterables.PostKind,
			},
			Vector: []float32(vector),
		},
	)
	if err != nil {
		return nil, 0, err
	}

	nearest, err := classifier.comparables.FindNearest(
		ctx,
		target,
		clusterables.CampaignKind,
		1,
	)
	if err != nil {
		return nil, 0, err
	}

	if len(nearest) == 0 {
		return nil, 0, nil
	}

	campaignComparable := nearest[0]
	if campaignComparable == nil {
		return nil, 0, ErrInvalidCampaignClassifierComparable
	}

	confidence, err := campaignConfidence(
		target.Vector(),
		campaignComparable.Vector(),
	)
	if err != nil {
		return nil, 0, err
	}

	if confidence < classifier.threshold {
		return nil, confidence, nil
	}

	campaign, err := classifier.campaigns.FindByID(
		ctx,
		campaignComparable.Identifier(),
	)
	if err != nil {
		return nil, 0, err
	}

	return campaign, confidence, nil
}

func campaignConfidence(
	source []float32,
	target []float32,
) (float64, error) {
	if len(source) == 0 ||
		len(target) == 0 ||
		len(source) != len(target) {
		return 0, ErrInvalidCampaignClassifierVector
	}

	score := campaignCosineSimilarity(source, target)

	if score < 0 {
		return 0, nil
	}

	if score > 1 {
		return 1, nil
	}

	return score, nil
}

func campaignCosineSimilarity(
	source []float32,
	target []float32,
) float64 {
	var dot float64
	var sourceMagnitude float64
	var targetMagnitude float64

	for index := range source {
		sourceValue := float64(source[index])
		targetValue := float64(target[index])

		dot += sourceValue * targetValue
		sourceMagnitude += sourceValue * sourceValue
		targetMagnitude += targetValue * targetValue
	}

	if sourceMagnitude == 0 ||
		targetMagnitude == 0 {
		return 0
	}

	return dot / (math.Sqrt(sourceMagnitude) * math.Sqrt(targetMagnitude))
}
