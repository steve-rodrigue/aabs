package assignments

import (
	"context"
	"math"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/campaigns"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters/clusterables"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/narratives"
)

type assigner struct {
	adapter Adapter

	comparables clusterables.ComparableRepository

	threshold float64
}

func (assigner *assigner) Assign(
	ctx context.Context,
	campaign campaigns.Campaign,
	narratives []narratives.Narrative,
) ([]Assignment, error) {
	if campaign == nil {
		return nil, ErrInvalidAssignmentAssignerCampaign
	}

	if len(narratives) == 0 {
		return []Assignment{}, nil
	}

	campaignComparable, err := assigner.comparables.FindByID(
		ctx,
		campaign.Identifier(),
	)
	if err != nil {
		return nil, err
	}

	if campaignComparable == nil {
		return nil, ErrInvalidAssignmentAssignerComparable
	}

	campaignVector := campaignComparable.Vector()
	if len(campaignVector) == 0 {
		return nil, ErrInvalidAssignmentAssignerVector
	}

	out := []Assignment{}

	for _, narrative := range narratives {
		if narrative == nil {
			return nil, ErrInvalidAssignmentAssignerNarrative
		}

		narrativeComparable, err := assigner.comparables.FindByID(
			ctx,
			narrative.Identifier(),
		)
		if err != nil {
			return nil, err
		}

		if narrativeComparable == nil {
			return nil, ErrInvalidAssignmentAssignerComparable
		}

		confidence, err := assignmentConfidence(
			campaignVector,
			narrativeComparable.Vector(),
		)
		if err != nil {
			return nil, err
		}

		if confidence < assigner.threshold {
			continue
		}

		assignment, err := assigner.adapter.ToDomain(
			AssignmentInput{
				Identifier: uuid.New(),
				Narrative:  narrative,
				Campaign:   campaign,
				Confidence: confidence,
				AssignedOn: time.Now().UTC(),
			},
		)
		if err != nil {
			return nil, err
		}

		out = append(out, assignment)
	}

	return out, nil
}

func assignmentConfidence(
	source []float32,
	target []float32,
) (float64, error) {
	if len(source) == 0 ||
		len(target) == 0 ||
		len(source) != len(target) {
		return 0, ErrInvalidAssignmentAssignerVector
	}

	score := assignmentCosineSimilarity(source, target)

	if score < 0 {
		return 0, nil
	}

	if score > 1 {
		return 1, nil
	}

	return score, nil
}

func assignmentCosineSimilarity(
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
