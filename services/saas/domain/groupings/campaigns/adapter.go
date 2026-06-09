package campaigns

import (
	"math"
	"strings"

	"github.com/google/uuid"
)

type adapter struct{}

func (adapter *adapter) ToDomain(
	input CampaignInput,
) (Campaign, error) {
	if input.Identifier == uuid.Nil {
		return nil, ErrInvalidCampaignIdentifier
	}

	input.Name = strings.TrimSpace(input.Name)
	if input.Name == "" {
		return nil, ErrInvalidCampaignName
	}

	input.Description = strings.TrimSpace(input.Description)

	if input.Cluster == nil {
		return nil, ErrInvalidCampaignCluster
	}

	if input.PostCount <= 0 {
		return nil, ErrInvalidCampaignPostCount
	}

	if math.IsNaN(input.Confidence) ||
		math.IsInf(input.Confidence, 0) ||
		input.Confidence < 0 ||
		input.Confidence > 1 {
		return nil, ErrInvalidCampaignConfidence
	}

	if input.CreatedOn.IsZero() {
		return nil, ErrInvalidCampaignCreatedOn
	}

	return &campaign{
		identifier:  input.Identifier,
		name:        input.Name,
		description: input.Description,
		cluster:     input.Cluster,
		postCount:   input.PostCount,
		confidence:  input.Confidence,
		createdOn:   input.CreatedOn.UTC(),
	}, nil
}
