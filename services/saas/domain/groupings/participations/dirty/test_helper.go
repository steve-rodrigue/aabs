package dirty

import (
	"context"
	"sort"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/participatables"
)

func NewMockDirty(
	participant participatables.Participatable,
	target participatables.Participatable,
) Dirty {
	return &MockDirty{
		ID:               uuid.New(),
		ParticipantValue: participant,
		TargetValue:      target,
		MarkedOnValue:    time.Now().UTC(),
	}
}

func NewMockDirtyWithID(
	id uuid.UUID,
	participant participatables.Participatable,
	target participatables.Participatable,
) Dirty {
	return &MockDirty{
		ID:               id,
		ParticipantValue: participant,
		TargetValue:      target,
		MarkedOnValue:    time.Now().UTC(),
	}
}

func NewMockDirtyAdapter() *MockDirtyAdapter {
	return &MockDirtyAdapter{}
}

func NewMockDirtyRepository() *MockDirtyRepository {
	return &MockDirtyRepository{
		Items: map[uuid.UUID]Dirty{},
	}
}

type MockDirty struct {
	ID uuid.UUID

	ParticipantValue participatables.Participatable
	TargetValue      participatables.Participatable

	MarkedOnValue time.Time
}

func (dirty *MockDirty) Identifier() uuid.UUID {
	return dirty.ID
}

func (dirty *MockDirty) Participant() participatables.Participatable {
	return dirty.ParticipantValue
}

func (dirty *MockDirty) Target() participatables.Participatable {
	return dirty.TargetValue
}

func (dirty *MockDirty) MarkedOn() time.Time {
	return dirty.MarkedOnValue
}

type MockDirtyAdapter struct {
	ToDomainCalls int
	ToDomainErr   error
	ToDomainValue Dirty

	LastInput DirtyInput
}

func (adapter *MockDirtyAdapter) ToDomain(
	input DirtyInput,
) (Dirty, error) {
	adapter.ToDomainCalls++
	adapter.LastInput = input

	if adapter.ToDomainErr != nil {
		return nil, adapter.ToDomainErr
	}

	if adapter.ToDomainValue != nil {
		return adapter.ToDomainValue, nil
	}

	return &MockDirty{
		ID:               input.Identifier,
		ParticipantValue: input.Participant,
		TargetValue:      input.Target,
		MarkedOnValue:    input.MarkedOn,
	}, nil
}

type MockDirtyRepository struct {
	Items map[uuid.UUID]Dirty

	SaveCalls int
	SaveErr   error

	DeleteCalls int
	DeleteErr   error

	FindByIDCalls int
	FindByIDErr   error
	FindByIDValue Dirty

	FindBetweenCalls int
	FindBetweenErr   error
	FindBetweenValue Dirty

	FindCalls int
	FindErr   error
	FindValue []Dirty

	FindAfterCalls int
	FindAfterErr   error
	FindAfterValue []Dirty

	CountCalls int
	CountErr   error
	CountValue int64

	LastContext context.Context

	LastSaved     Dirty
	LastDeletedID uuid.UUID

	LastID uuid.UUID

	LastParticipant participatables.Participatable
	LastTarget      participatables.Participatable

	LastIndex  int
	LastAmount int
	LastCursor uuid.UUID
}

func (repository *MockDirtyRepository) Save(
	ctx context.Context,
	dirty Dirty,
) error {
	repository.SaveCalls++
	repository.LastContext = ctx
	repository.LastSaved = dirty

	if repository.Items != nil && dirty != nil {
		repository.Items[dirty.Identifier()] = dirty
	}

	return repository.SaveErr
}

func (repository *MockDirtyRepository) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	repository.DeleteCalls++
	repository.LastContext = ctx
	repository.LastDeletedID = id

	delete(repository.Items, id)

	return repository.DeleteErr
}

func (repository *MockDirtyRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (Dirty, error) {
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

func (repository *MockDirtyRepository) FindBetween(
	ctx context.Context,
	participant participatables.Participatable,
	target participatables.Participatable,
) (Dirty, error) {
	repository.FindBetweenCalls++
	repository.LastContext = ctx
	repository.LastParticipant = participant
	repository.LastTarget = target

	if repository.FindBetweenErr != nil {
		return nil, repository.FindBetweenErr
	}

	if repository.FindBetweenValue != nil {
		return repository.FindBetweenValue, nil
	}

	for _, item := range repository.Items {
		if item.Participant().Identifier() == participant.Identifier() &&
			item.Target().Identifier() == target.Identifier() {
			return item, nil
		}
	}

	return nil, nil
}

func (repository *MockDirtyRepository) Find(
	ctx context.Context,
	index int,
	amount int,
) ([]Dirty, error) {
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

	items := repository.sorted()

	if index >= len(items) {
		return []Dirty{}, nil
	}

	end := index + amount
	if end > len(items) {
		end = len(items)
	}

	return items[index:end], nil
}

func (repository *MockDirtyRepository) FindAfter(
	ctx context.Context,
	cursor uuid.UUID,
	amount int,
) ([]Dirty, error) {
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

	items := repository.sorted()

	start := 0

	if cursor != uuid.Nil {
		for i, item := range items {
			if item.Identifier() == cursor {
				start = i + 1
				break
			}
		}
	}

	if start >= len(items) {
		return []Dirty{}, nil
	}

	end := start + amount
	if end > len(items) {
		end = len(items)
	}

	return items[start:end], nil
}

func (repository *MockDirtyRepository) Count(
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

func (repository *MockDirtyRepository) sorted() []Dirty {
	out := make([]Dirty, 0, len(repository.Items))

	for _, item := range repository.Items {
		out = append(out, item)
	}

	sort.Slice(out, func(left int, right int) bool {
		return out[left].Identifier().String() <
			out[right].Identifier().String()
	})

	return out
}
