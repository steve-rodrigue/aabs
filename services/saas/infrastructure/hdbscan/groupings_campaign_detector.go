package hdbscan

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	domain_campaigns "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/campaigns"
	domain_clusters "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters/clusterables"
)

type groupingsCampaignDetector struct {
	endpoint string
	client   *http.Client

	campaignAdapter domain_campaigns.Adapter
	clusterAdapter  domain_clusters.Adapter
	comparables     clusterables.ComparableRepository

	minClusterSize int
	minSamples     *int
}

type clusterRequest struct {
	Embeddings     [][]float32 `json:"embeddings"`
	MinClusterSize int         `json:"min_cluster_size"`
	MinSamples     *int        `json:"min_samples,omitempty"`
}

type clusterResult struct {
	Index       int     `json:"index"`
	ClusterID   int     `json:"cluster_id"`
	Probability float64 `json:"probability"`
	IsNoise     bool    `json:"is_noise"`
}

type clusterSummary struct {
	ClusterID     int       `json:"cluster_id"`
	Size          int       `json:"size"`
	Centroid      []float32 `json:"centroid"`
	MemberIndexes []int     `json:"member_indexes"`
}

type clusterResponse struct {
	TotalEmbeddings int              `json:"total_embeddings"`
	TotalClusters   int              `json:"total_clusters"`
	NoiseCount      int              `json:"noise_count"`
	Results         []clusterResult  `json:"results"`
	Clusters        []clusterSummary `json:"clusters"`
}

func (detector *groupingsCampaignDetector) Detect(
	ctx context.Context,
	candidates []clusterables.Clusterable,
) ([]domain_campaigns.Campaign, error) {
	if detector.endpoint == "" {
		return nil, ErrInvalidGroupingsCampaignDetectorEndpoint
	}

	if len(candidates) == 0 {
		return []domain_campaigns.Campaign{}, nil
	}

	minClusterSize := detector.minClusterSize
	if minClusterSize <= 0 {
		minClusterSize = 5
	}

	comparables := make([]clusterables.Comparable, 0, len(candidates))
	embeddings := make([][]float32, 0, len(candidates))
	memberIDs := make([]uuid.UUID, 0, len(candidates))

	var memberKind clusterables.Kind

	for index, candidate := range candidates {
		if candidate == nil {
			return nil, ErrInvalidGroupingsCampaignDetectorCandidates
		}

		if index == 0 {
			memberKind = candidate.ClusterKind()
		}

		if candidate.ClusterKind() != memberKind {
			return nil, ErrInvalidGroupingsCampaignDetectorKind
		}

		comparable, err := detector.comparables.FindByID(
			ctx,
			candidate.Identifier(),
		)
		if err != nil {
			return nil, err
		}

		if comparable == nil {
			return nil, ErrInvalidGroupingsCampaignDetectorComparable
		}

		vector := comparable.Vector()
		if len(vector) == 0 {
			return nil, ErrInvalidGroupingsCampaignDetectorVector
		}

		comparables = append(comparables, comparable)
		embeddings = append(embeddings, vector)
		memberIDs = append(memberIDs, candidate.Identifier())
	}

	if len(embeddings) < minClusterSize {
		return []domain_campaigns.Campaign{}, nil
	}

	response, err := detector.cluster(
		ctx,
		clusterRequest{
			Embeddings:     embeddings,
			MinClusterSize: minClusterSize,
			MinSamples:     detector.minSamples,
		},
	)
	if err != nil {
		return nil, err
	}

	out := []domain_campaigns.Campaign{}

	for _, summary := range response.Clusters {
		campaignID := uuid.New()

		cluster, err := detector.clusterAdapter.ToDomain(
			domain_clusters.ClusterInput{
				Identifier: uuid.New(),
				Target: clusterables.ClusterableInput{
					Identifier:  campaignID,
					ClusterKind: clusterables.CampaignKind,
				},
				MemberIDs:       memberIDsForCluster(summary, memberIDs),
				MemberKind:      memberKind,
				ConfidenceScore: averageProbability(summary.ClusterID, response.Results),
				Centroid:        summary.Centroid,
				CreatedOn:       time.Now().UTC(),
			},
		)
		if err != nil {
			return nil, err
		}

		campaign, err := detector.campaignAdapter.ToDomain(
			domain_campaigns.CampaignInput{
				Identifier:  campaignID,
				Name:        fmt.Sprintf("Campaign %d", summary.ClusterID),
				Description: fmt.Sprintf("Detected campaign with %d members", summary.Size),
				Cluster:     cluster,
				PostCount:   summary.Size,
				Confidence:  averageProbability(summary.ClusterID, response.Results),
				CreatedOn:   time.Now().UTC(),
			},
		)
		if err != nil {
			return nil, err
		}

		out = append(out, campaign)
	}

	return out, nil
}

func (detector *groupingsCampaignDetector) cluster(
	ctx context.Context,
	payload clusterRequest,
) (clusterResponse, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return clusterResponse{}, err
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		detector.endpoint+"/cluster",
		bytes.NewReader(body),
	)
	if err != nil {
		return clusterResponse{}, err
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := detector.client.Do(request)
	if err != nil {
		return clusterResponse{}, err
	}

	defer response.Body.Close()

	if response.StatusCode < http.StatusOK ||
		response.StatusCode >= http.StatusMultipleChoices {
		return clusterResponse{}, fmt.Errorf(
			"%w: status %d",
			ErrInvalidGroupingsCampaignDetectorResponse,
			response.StatusCode,
		)
	}

	var result clusterResponse

	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return clusterResponse{}, err
	}

	return result, nil
}

func memberIDsForCluster(
	cluster clusterSummary,
	memberIDs []uuid.UUID,
) []uuid.UUID {
	out := make([]uuid.UUID, 0, len(cluster.MemberIndexes))

	for _, index := range cluster.MemberIndexes {
		if index < 0 || index >= len(memberIDs) {
			continue
		}

		out = append(out, memberIDs[index])
	}

	return out
}

func averageProbability(
	clusterID int,
	results []clusterResult,
) float64 {
	var total float64
	var count int

	for _, result := range results {
		if result.ClusterID != clusterID || result.IsNoise {
			continue
		}

		total += result.Probability
		count++
	}

	if count == 0 {
		return 0
	}

	return total / float64(count)
}
