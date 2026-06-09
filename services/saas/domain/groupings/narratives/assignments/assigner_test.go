package assignments

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/campaigns"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters/clusterables"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/narratives"
)

var errTest = errors.New("test error")

func TestNewAssigner(t *testing.T) {
	assigner := NewAssigner(
		NewMockAssignmentAdapter(),
		clusterables.NewMockComparableRepository(),
		0.7,
	)

	if assigner == nil {
		t.Fatalf("expected assigner")
	}
}

func TestAssignerAssign(t *testing.T) {
	ctx := context.Background()

	adapter := NewMockAssignmentAdapter()
	comparables := clusterables.NewMockComparableRepository()

	campaign := campaigns.NewMockCampaign(
		"Campaign A",
		"Description A",
	)

	narrative := narratives.NewMockNarrative(
		"Narrative A",
		"Description A",
	)

	expected := NewMockAssignment(
		narrative,
		campaign,
	)

	adapter.ToDomainValue = expected

	comparables.Items[campaign.Identifier()] =
		clusterables.NewMockComparableWithID(
			campaign.Identifier(),
			clusterables.CampaignKind,
			[]float32{1, 0},
		)

	comparables.Items[narrative.Identifier()] =
		clusterables.NewMockComparableWithID(
			narrative.Identifier(),
			clusterables.NarrativeKind,
			[]float32{1, 0},
		)

	assigner := NewAssigner(
		adapter,
		comparables,
		0.7,
	)

	result, err := assigner.Assign(
		ctx,
		campaign,
		[]narratives.Narrative{
			narrative,
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 assignment, got %d", len(result))
	}

	if result[0] != expected {
		t.Fatalf("expected assignment")
	}

	if comparables.FindByIDCalls != 2 {
		t.Fatalf(
			"expected 2 comparable lookups, got %d",
			comparables.FindByIDCalls,
		)
	}

	if adapter.ToDomainCalls != 1 {
		t.Fatalf("expected 1 adapter call")
	}

	if adapter.LastInput.Identifier == uuid.Nil {
		t.Fatalf("expected generated identifier")
	}

	if adapter.LastInput.Narrative != narrative {
		t.Fatalf("expected narrative")
	}

	if adapter.LastInput.Campaign != campaign {
		t.Fatalf("expected campaign")
	}

	if adapter.LastInput.Confidence != 1 {
		t.Fatalf(
			"expected confidence 1, got %f",
			adapter.LastInput.Confidence,
		)
	}

	if adapter.LastInput.AssignedOn.IsZero() {
		t.Fatalf("expected assigned on")
	}
}

func TestAssignerAssignReturnsEmptyWhenNoNarratives(t *testing.T) {
	assigner := NewAssigner(
		NewMockAssignmentAdapter(),
		clusterables.NewMockComparableRepository(),
		0.7,
	)

	result, err := assigner.Assign(
		context.Background(),
		campaigns.NewMockCampaign("Campaign A", "Description A"),
		[]narratives.Narrative{},
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 0 {
		t.Fatalf("expected empty result")
	}
}

func TestAssignerAssignSkipsNarrativesBelowThreshold(t *testing.T) {
	adapter := NewMockAssignmentAdapter()
	comparables := clusterables.NewMockComparableRepository()

	campaign := campaigns.NewMockCampaign(
		"Campaign A",
		"Description A",
	)

	narrative := narratives.NewMockNarrative(
		"Narrative A",
		"Description A",
	)

	comparables.Items[campaign.Identifier()] =
		clusterables.NewMockComparableWithID(
			campaign.Identifier(),
			clusterables.CampaignKind,
			[]float32{1, 0},
		)

	comparables.Items[narrative.Identifier()] =
		clusterables.NewMockComparableWithID(
			narrative.Identifier(),
			clusterables.NarrativeKind,
			[]float32{0, 1},
		)

	assigner := NewAssigner(
		adapter,
		comparables,
		0.7,
	)

	result, err := assigner.Assign(
		context.Background(),
		campaign,
		[]narratives.Narrative{
			narrative,
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 0 {
		t.Fatalf("expected empty result")
	}

	if adapter.ToDomainCalls != 0 {
		t.Fatalf("expected no adapter call")
	}
}

func TestAssignerAssignReturnsInvalidCampaignError(t *testing.T) {
	assigner := NewAssigner(
		NewMockAssignmentAdapter(),
		clusterables.NewMockComparableRepository(),
		0.7,
	)

	_, err := assigner.Assign(
		context.Background(),
		nil,
		[]narratives.Narrative{},
	)

	if !errors.Is(err, ErrInvalidAssignmentAssignerCampaign) {
		t.Fatalf("expected invalid campaign error, got %v", err)
	}
}

func TestAssignerAssignReturnsInvalidNarrativeError(t *testing.T) {
	campaign := campaigns.NewMockCampaign(
		"Campaign A",
		"Description A",
	)

	comparables := clusterables.NewMockComparableRepository()
	comparables.Items[campaign.Identifier()] =
		clusterables.NewMockComparableWithID(
			campaign.Identifier(),
			clusterables.CampaignKind,
			[]float32{1, 0},
		)

	assigner := NewAssigner(
		NewMockAssignmentAdapter(),
		comparables,
		0.7,
	)

	_, err := assigner.Assign(
		context.Background(),
		campaign,
		[]narratives.Narrative{
			nil,
		},
	)

	if !errors.Is(err, ErrInvalidAssignmentAssignerNarrative) {
		t.Fatalf("expected invalid narrative error, got %v", err)
	}
}

func TestAssignerAssignReturnsComparableRepositoryErrorForCampaign(t *testing.T) {
	comparables := clusterables.NewMockComparableRepository()
	comparables.FindByIDErr = errTest

	assigner := NewAssigner(
		NewMockAssignmentAdapter(),
		comparables,
		0.7,
	)

	_, err := assigner.Assign(
		context.Background(),
		campaigns.NewMockCampaign("Campaign A", "Description A"),
		[]narratives.Narrative{
			narratives.NewMockNarrative("Narrative A", "Description A"),
		},
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected comparable repository error, got %v", err)
	}
}

func TestAssignerAssignReturnsInvalidComparableErrorForCampaign(t *testing.T) {
	assigner := NewAssigner(
		NewMockAssignmentAdapter(),
		clusterables.NewMockComparableRepository(),
		0.7,
	)

	_, err := assigner.Assign(
		context.Background(),
		campaigns.NewMockCampaign("Campaign A", "Description A"),
		[]narratives.Narrative{
			narratives.NewMockNarrative("Narrative A", "Description A"),
		},
	)

	if !errors.Is(err, ErrInvalidAssignmentAssignerComparable) {
		t.Fatalf("expected invalid comparable error, got %v", err)
	}
}

func TestAssignerAssignReturnsInvalidVectorErrorForCampaign(t *testing.T) {
	campaign := campaigns.NewMockCampaign(
		"Campaign A",
		"Description A",
	)

	comparables := clusterables.NewMockComparableRepository()
	comparables.Items[campaign.Identifier()] =
		clusterables.NewMockComparableWithID(
			campaign.Identifier(),
			clusterables.CampaignKind,
			[]float32{},
		)

	assigner := NewAssigner(
		NewMockAssignmentAdapter(),
		comparables,
		0.7,
	)

	_, err := assigner.Assign(
		context.Background(),
		campaign,
		[]narratives.Narrative{
			narratives.NewMockNarrative("Narrative A", "Description A"),
		},
	)

	if !errors.Is(err, ErrInvalidAssignmentAssignerVector) {
		t.Fatalf("expected invalid vector error, got %v", err)
	}
}

func TestAssignerAssignReturnsComparableRepositoryErrorForNarrative(t *testing.T) {
	campaign := campaigns.NewMockCampaign(
		"Campaign A",
		"Description A",
	)

	narrative := narratives.NewMockNarrative(
		"Narrative A",
		"Description A",
	)

	comparables := clusterables.NewMockComparableRepository()
	comparables.Items[campaign.Identifier()] =
		clusterables.NewMockComparableWithID(
			campaign.Identifier(),
			clusterables.CampaignKind,
			[]float32{1, 0},
		)

	comparables.FindByIDErr = errTest

	assigner := NewAssigner(
		NewMockAssignmentAdapter(),
		comparables,
		0.7,
	)

	_, err := assigner.Assign(
		context.Background(),
		campaign,
		[]narratives.Narrative{
			narrative,
		},
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected comparable repository error, got %v", err)
	}
}

func TestAssignerAssignReturnsInvalidComparableErrorForNarrative(t *testing.T) {
	campaign := campaigns.NewMockCampaign(
		"Campaign A",
		"Description A",
	)

	narrative := narratives.NewMockNarrative(
		"Narrative A",
		"Description A",
	)

	comparables := clusterables.NewMockComparableRepository()
	comparables.Items[campaign.Identifier()] =
		clusterables.NewMockComparableWithID(
			campaign.Identifier(),
			clusterables.CampaignKind,
			[]float32{1, 0},
		)

	assigner := NewAssigner(
		NewMockAssignmentAdapter(),
		comparables,
		0.7,
	)

	_, err := assigner.Assign(
		context.Background(),
		campaign,
		[]narratives.Narrative{
			narrative,
		},
	)

	if !errors.Is(err, ErrInvalidAssignmentAssignerComparable) {
		t.Fatalf("expected invalid comparable error, got %v", err)
	}
}

func TestAssignerAssignReturnsVectorErrorWhenVectorsMismatch(t *testing.T) {
	campaign := campaigns.NewMockCampaign(
		"Campaign A",
		"Description A",
	)

	narrative := narratives.NewMockNarrative(
		"Narrative A",
		"Description A",
	)

	comparables := clusterables.NewMockComparableRepository()
	comparables.Items[campaign.Identifier()] =
		clusterables.NewMockComparableWithID(
			campaign.Identifier(),
			clusterables.CampaignKind,
			[]float32{1, 0},
		)

	comparables.Items[narrative.Identifier()] =
		clusterables.NewMockComparableWithID(
			narrative.Identifier(),
			clusterables.NarrativeKind,
			[]float32{1, 0, 0},
		)

	assigner := NewAssigner(
		NewMockAssignmentAdapter(),
		comparables,
		0.7,
	)

	_, err := assigner.Assign(
		context.Background(),
		campaign,
		[]narratives.Narrative{
			narrative,
		},
	)

	if !errors.Is(err, ErrInvalidAssignmentAssignerVector) {
		t.Fatalf("expected vector error, got %v", err)
	}
}

func TestAssignerAssignReturnsAdapterError(t *testing.T) {
	adapter := NewMockAssignmentAdapter()
	adapter.ToDomainErr = errTest

	comparables := clusterables.NewMockComparableRepository()

	campaign := campaigns.NewMockCampaign(
		"Campaign A",
		"Description A",
	)

	narrative := narratives.NewMockNarrative(
		"Narrative A",
		"Description A",
	)

	comparables.Items[campaign.Identifier()] =
		clusterables.NewMockComparableWithID(
			campaign.Identifier(),
			clusterables.CampaignKind,
			[]float32{1, 0},
		)

	comparables.Items[narrative.Identifier()] =
		clusterables.NewMockComparableWithID(
			narrative.Identifier(),
			clusterables.NarrativeKind,
			[]float32{1, 0},
		)

	assigner := NewAssigner(
		adapter,
		comparables,
		0.7,
	)

	_, err := assigner.Assign(
		context.Background(),
		campaign,
		[]narratives.Narrative{
			narrative,
		},
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected adapter error, got %v", err)
	}
}

func TestAssignmentConfidence(t *testing.T) {
	confidence, err := assignmentConfidence(
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

func TestAssignmentConfidenceClampsNegativeToZero(t *testing.T) {
	confidence, err := assignmentConfidence(
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

func TestAssignmentConfidenceReturnsVectorErrorWhenEmpty(t *testing.T) {
	_, err := assignmentConfidence(
		[]float32{},
		[]float32{1, 0},
	)

	if !errors.Is(err, ErrInvalidAssignmentAssignerVector) {
		t.Fatalf("expected vector error, got %v", err)
	}
}

func TestAssignmentConfidenceReturnsVectorErrorWhenMismatch(t *testing.T) {
	_, err := assignmentConfidence(
		[]float32{1, 0},
		[]float32{1, 0, 0},
	)

	if !errors.Is(err, ErrInvalidAssignmentAssignerVector) {
		t.Fatalf("expected vector error, got %v", err)
	}
}

func TestAssignmentCosineSimilarityReturnsZeroWhenSourceMagnitudeIsZero(t *testing.T) {
	confidence := assignmentCosineSimilarity(
		[]float32{0, 0},
		[]float32{1, 0},
	)

	if confidence != 0 {
		t.Fatalf("expected confidence 0, got %f", confidence)
	}
}

func TestAssignmentCosineSimilarityReturnsZeroWhenTargetMagnitudeIsZero(t *testing.T) {
	confidence := assignmentCosineSimilarity(
		[]float32{1, 0},
		[]float32{0, 0},
	)

	if confidence != 0 {
		t.Fatalf("expected confidence 0, got %f", confidence)
	}
}
