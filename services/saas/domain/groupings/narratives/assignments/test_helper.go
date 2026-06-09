package assignments

import (
	"context"
	"sort"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/campaigns"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/narratives"
)

func NewMockAssignment(
	narrative narratives.Narrative,
	campaign campaigns.Campaign,
) Assignment {
	return &MockAssignment{
		ID:              uuid.New(),
		NarrativeValue:  narrative,
		CampaignValue:   campaign,
		ConfidenceValue: 1.0,
		AssignedOnValue: time.Now().UTC(),
	}
}

func NewMockAssignmentWithID(
	id uuid.UUID,
	narrative narratives.Narrative,
	campaign campaigns.Campaign,
) Assignment {
	return &MockAssignment{
		ID:              id,
		NarrativeValue:  narrative,
		CampaignValue:   campaign,
		ConfidenceValue: 1.0,
		AssignedOnValue: time.Now().UTC(),
	}
}

func NewMockAssignmentRepository() *MockAssignmentRepository {
	return &MockAssignmentRepository{
		Items: map[uuid.UUID]Assignment{},
	}
}

func NewMockAssignmentAdapter() *MockAssignmentAdapter {
	return &MockAssignmentAdapter{}
}

func NewMockAssigner() *MockAssigner {
	return &MockAssigner{}
}

type MockAssignment struct {
	ID uuid.UUID

	NarrativeValue narratives.Narrative
	CampaignValue  campaigns.Campaign

	ConfidenceValue float64
	AssignedOnValue time.Time
}

func (assignment *MockAssignment) Identifier() uuid.UUID {
	return assignment.ID
}

func (assignment *MockAssignment) Narrative() narratives.Narrative {
	return assignment.NarrativeValue
}

func (assignment *MockAssignment) Campaign() campaigns.Campaign {
	return assignment.CampaignValue
}

func (assignment *MockAssignment) Confidence() float64 {
	return assignment.ConfidenceValue
}

func (assignment *MockAssignment) AssignedOn() time.Time {
	return assignment.AssignedOnValue
}

type MockAssignmentAdapter struct {
	ToDomainCalls int
	ToDomainErr   error
	ToDomainValue Assignment

	LastInput AssignmentInput
}

func (adapter *MockAssignmentAdapter) ToDomain(
	input AssignmentInput,
) (Assignment, error) {
	adapter.ToDomainCalls++
	adapter.LastInput = input

	if adapter.ToDomainErr != nil {
		return nil, adapter.ToDomainErr
	}

	if adapter.ToDomainValue != nil {
		return adapter.ToDomainValue, nil
	}

	return &MockAssignment{
		ID:              input.Identifier,
		NarrativeValue:  input.Narrative,
		CampaignValue:   input.Campaign,
		ConfidenceValue: input.Confidence,
		AssignedOnValue: input.AssignedOn,
	}, nil
}

type MockAssignmentRepository struct {
	Items map[uuid.UUID]Assignment

	SaveCalls int
	SaveErr   error

	FindByIDCalls int
	FindByIDErr   error
	FindByIDValue Assignment

	FindByNarrativeCalls int
	FindByNarrativeErr   error
	FindByNarrativeValue []Assignment

	FindByCampaignCalls int
	FindByCampaignErr   error
	FindByCampaignValue []Assignment

	FindBetweenCalls int
	FindBetweenErr   error
	FindBetweenValue Assignment

	FindCalls int
	FindErr   error
	FindValue []Assignment

	FindAfterCalls int
	FindAfterErr   error
	FindAfterValue []Assignment

	CountCalls int
	CountErr   error
	CountValue int64

	LastContext     context.Context
	LastSaved       Assignment
	LastID          uuid.UUID
	LastNarrativeID uuid.UUID
	LastCampaignID  uuid.UUID
	LastIndex       int
	LastAmount      int
	LastCursor      uuid.UUID
}

func (repository *MockAssignmentRepository) Save(
	ctx context.Context,
	assignment Assignment,
) error {
	repository.SaveCalls++
	repository.LastContext = ctx
	repository.LastSaved = assignment

	if repository.Items != nil && assignment != nil {
		repository.Items[assignment.Identifier()] = assignment
	}

	return repository.SaveErr
}

func (repository *MockAssignmentRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (Assignment, error) {
	repository.FindByIDCalls++
	repository.LastContext = ctx
	repository.LastID = id

	if repository.FindByIDErr != nil {
		return nil, repository.FindByIDErr
	}

	if repository.FindByIDValue != nil {
		return repository.FindByIDValue, nil
	}

	return repository.Items[id], nil
}

func (repository *MockAssignmentRepository) FindByNarrative(
	ctx context.Context,
	narrative uuid.UUID,
) ([]Assignment, error) {
	repository.FindByNarrativeCalls++
	repository.LastContext = ctx
	repository.LastNarrativeID = narrative

	if repository.FindByNarrativeErr != nil {
		return nil, repository.FindByNarrativeErr
	}

	if repository.FindByNarrativeValue != nil {
		return repository.FindByNarrativeValue, nil
	}

	out := []Assignment{}

	for _, assignment := range repository.Items {
		if assignment.Narrative() != nil &&
			assignment.Narrative().Identifier() == narrative {
			out = append(out, assignment)
		}
	}

	return out, nil
}

func (repository *MockAssignmentRepository) FindByCampaign(
	ctx context.Context,
	campaign uuid.UUID,
) ([]Assignment, error) {
	repository.FindByCampaignCalls++
	repository.LastContext = ctx
	repository.LastCampaignID = campaign

	if repository.FindByCampaignErr != nil {
		return nil, repository.FindByCampaignErr
	}

	if repository.FindByCampaignValue != nil {
		return repository.FindByCampaignValue, nil
	}

	out := []Assignment{}

	for _, assignment := range repository.Items {
		if assignment.Campaign() != nil &&
			assignment.Campaign().Identifier() == campaign {
			out = append(out, assignment)
		}
	}

	return out, nil
}

func (repository *MockAssignmentRepository) FindBetween(
	ctx context.Context,
	narrative uuid.UUID,
	campaign uuid.UUID,
) (Assignment, error) {
	repository.FindBetweenCalls++
	repository.LastContext = ctx
	repository.LastNarrativeID = narrative
	repository.LastCampaignID = campaign

	if repository.FindBetweenErr != nil {
		return nil, repository.FindBetweenErr
	}

	if repository.FindBetweenValue != nil {
		return repository.FindBetweenValue, nil
	}

	for _, assignment := range repository.Items {
		if assignment.Narrative() != nil &&
			assignment.Campaign() != nil &&
			assignment.Narrative().Identifier() == narrative &&
			assignment.Campaign().Identifier() == campaign {
			return assignment, nil
		}
	}

	return nil, nil
}

func (repository *MockAssignmentRepository) Find(
	ctx context.Context,
	index int,
	amount int,
) ([]Assignment, error) {
	repository.FindCalls++
	repository.LastContext = ctx
	repository.LastIndex = index
	repository.LastAmount = amount

	if repository.FindErr != nil {
		return nil, repository.FindErr
	}

	if repository.FindValue != nil {
		return repository.FindValue, nil
	}

	items := repository.sortedAssignments()

	if index >= len(items) {
		return []Assignment{}, nil
	}

	end := index + amount
	if end > len(items) {
		end = len(items)
	}

	return items[index:end], nil
}

func (repository *MockAssignmentRepository) FindAfter(
	ctx context.Context,
	cursor uuid.UUID,
	amount int,
) ([]Assignment, error) {
	repository.FindAfterCalls++
	repository.LastContext = ctx
	repository.LastCursor = cursor
	repository.LastAmount = amount

	if repository.FindAfterErr != nil {
		return nil, repository.FindAfterErr
	}

	if repository.FindAfterValue != nil {
		return repository.FindAfterValue, nil
	}

	items := repository.sortedAssignments()

	start := 0

	if cursor != uuid.Nil {
		for i, assignment := range items {
			if assignment.Identifier() == cursor {
				start = i + 1
				break
			}
		}
	}

	if start >= len(items) {
		return []Assignment{}, nil
	}

	end := start + amount
	if end > len(items) {
		end = len(items)
	}

	return items[start:end], nil
}

func (repository *MockAssignmentRepository) Count(
	ctx context.Context,
) (int64, error) {
	repository.CountCalls++
	repository.LastContext = ctx

	if repository.CountErr != nil {
		return 0, repository.CountErr
	}

	if repository.CountValue != 0 {
		return repository.CountValue, nil
	}

	return int64(len(repository.Items)), nil
}

func (repository *MockAssignmentRepository) sortedAssignments() []Assignment {
	out := make([]Assignment, 0, len(repository.Items))

	for _, assignment := range repository.Items {
		out = append(out, assignment)
	}

	sort.Slice(out, func(left int, right int) bool {
		return out[left].Identifier().String() <
			out[right].Identifier().String()
	})

	return out
}

type MockAssigner struct {
	AssignCalls int
	AssignErr   error
	AssignValue []Assignment

	LastContext    context.Context
	LastCampaign   campaigns.Campaign
	LastNarratives []narratives.Narrative
}

func (assigner *MockAssigner) Assign(
	ctx context.Context,
	campaign campaigns.Campaign,
	narratives []narratives.Narrative,
) ([]Assignment, error) {
	assigner.AssignCalls++
	assigner.LastContext = ctx
	assigner.LastCampaign = campaign
	assigner.LastNarratives = narratives

	return assigner.AssignValue, assigner.AssignErr
}
