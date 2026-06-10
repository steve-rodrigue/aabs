package scores

import (
	"time"

	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/scorables"
	"github.com/steve-rodrigue/aabs/services/saas/domain/scores/factors"
)

func NewMockScore(
	target scorables.Scorable,
	scoreType Type,
	value float64,
) Score {
	return &MockScore{
		id:         uuid.New(),
		target:     target,
		scoreType:  scoreType,
		value:      value,
		confidence: 1.0,
		factors:    []factors.Factor{},
	}
}

type MockScore struct {
	id           uuid.UUID
	scoreType    Type
	target       scorables.Scorable
	value        float64
	confidence   float64
	factors      []factors.Factor
	calculatedOn time.Time
}

func (score *MockScore) Identifier() uuid.UUID {
	return score.id
}

func (score *MockScore) Type() Type {
	return score.scoreType
}

func (score *MockScore) Target() scorables.Scorable {
	return score.target
}

func (score *MockScore) Value() float64 {
	return score.value
}

func (score *MockScore) Confidence() float64 {
	return score.confidence
}

func (score *MockScore) Factors() []factors.Factor {
	return score.factors
}

func (score *MockScore) CalculatedOn() time.Time {
	return score.calculatedOn
}

func NewMockScoreRepository() *MockScoreRepository {
	return &MockScoreRepository{}
}

type MockScoreRepository struct {
	SaveCalls int
	SaveErr   error

	FindLatestByTargetCalls int
	FindLatestByTargetErr   error
	FindLatestByTargetValue Score

	FindHistoryByTargetCalls int
	FindHistoryByTargetErr   error
	FindHistoryByTargetValue []Score

	FindLatestByTargetAndTypesCalls int
	FindLatestByTargetAndTypesErr   error
	FindLatestByTargetAndTypesValue []Score
}

func (repository *MockScoreRepository) Save(score Score) error {
	repository.SaveCalls++
	return repository.SaveErr
}

func (repository *MockScoreRepository) FindLatestByTarget(
	target scorables.Scorable,
	scoreType Type,
) (Score, error) {
	repository.FindLatestByTargetCalls++

	if repository.FindLatestByTargetErr != nil {
		return nil, repository.FindLatestByTargetErr
	}

	return repository.FindLatestByTargetValue, nil
}

func (repository *MockScoreRepository) FindHistoryByTarget(
	target scorables.Scorable,
	scoreType Type,
) ([]Score, error) {
	repository.FindHistoryByTargetCalls++

	if repository.FindHistoryByTargetErr != nil {
		return nil, repository.FindHistoryByTargetErr
	}

	return repository.FindHistoryByTargetValue, nil
}

func (repository *MockScoreRepository) FindLatestByTargetAndTypes(
	target scorables.Scorable,
	scoreTypes []Type,
) ([]Score, error) {
	repository.FindLatestByTargetAndTypesCalls++

	if repository.FindLatestByTargetAndTypesErr != nil {
		return nil, repository.FindLatestByTargetAndTypesErr
	}

	return repository.FindLatestByTargetAndTypesValue, nil
}

type MockScoreCalculator struct {
	CalculatorType Type

	CalculateCalls int
	CalculateErr   error
	CalculateValue Score
}

func NewMockScoreCalculator(
	scoreType Type,
) *MockScoreCalculator {
	return &MockScoreCalculator{
		CalculatorType: scoreType,
	}
}

func (calculator *MockScoreCalculator) Type() Type {
	return calculator.CalculatorType
}

func (calculator *MockScoreCalculator) Calculate(
	target scorables.Scorable,
) (Score, error) {
	calculator.CalculateCalls++

	if calculator.CalculateErr != nil {
		return nil, calculator.CalculateErr
	}

	return calculator.CalculateValue, nil
}
