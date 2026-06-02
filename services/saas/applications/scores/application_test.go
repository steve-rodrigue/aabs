package scores

import (
	"errors"
	"testing"

	"github.com/google/uuid"

	domain_scores "github.com/steve-rodrigue/aabs/services/saas/domain/scores"
	"github.com/steve-rodrigue/aabs/services/saas/domain/scores/scorables"
)

var errTest = errors.New("test error")

func TestCalculate(t *testing.T) {
	fixture := newApplicationFixture()

	target := scorables.NewMockScorable(uuid.New(), scorables.UserKind)

	trustScore := domain_scores.NewMockScore(target, domain_scores.TrustType, 0.9)
	spamScore := domain_scores.NewMockScore(target, domain_scores.SpamType, 0.1)

	fixture.trustCalculator.CalculateValue = trustScore
	fixture.spamCalculator.CalculateValue = spamScore

	result, err := fixture.application.Calculate(target)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.trustCalculator.CalculateCalls != 1 {
		t.Fatalf("expected trust calculator to be called once")
	}

	if fixture.spamCalculator.CalculateCalls != 1 {
		t.Fatalf("expected spam calculator to be called once")
	}

	if fixture.repository.SaveCalls != 2 {
		t.Fatalf("expected 2 saved scores, got %d", fixture.repository.SaveCalls)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 scores, got %d", len(result))
	}

	if result[0] != trustScore {
		t.Fatalf("expected trust score")
	}

	if result[1] != spamScore {
		t.Fatalf("expected spam score")
	}
}

func TestCalculateReturnsCalculatorError(t *testing.T) {
	fixture := newApplicationFixture()

	target := scorables.NewMockScorable(uuid.New(), scorables.UserKind)
	fixture.trustCalculator.CalculateErr = errTest

	_, err := fixture.application.Calculate(target)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected calculator error, got %v", err)
	}

	if fixture.repository.SaveCalls != 0 {
		t.Fatalf("expected score not to be saved")
	}
}

func TestCalculateReturnsRepositorySaveError(t *testing.T) {
	fixture := newApplicationFixture()

	target := scorables.NewMockScorable(uuid.New(), scorables.UserKind)
	fixture.trustCalculator.CalculateValue = domain_scores.NewMockScore(
		target,
		domain_scores.TrustType,
		0.9,
	)
	fixture.repository.SaveErr = errTest

	_, err := fixture.application.Calculate(target)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected save error, got %v", err)
	}
}

func TestLatestScore(t *testing.T) {
	fixture := newApplicationFixture()

	id := uuid.New()
	target := scorables.NewMockScorable(id, scorables.UserKind)
	score := domain_scores.NewMockScore(target, domain_scores.TrustType, 0.9)

	fixture.scorableRepository.Items[id] = target
	fixture.repository.FindLatestByTargetValue = score

	result, err := fixture.application.LatestScore(id, domain_scores.TrustType)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.scorableRepository.FindByIDCalls != 1 {
		t.Fatalf("expected scorable find by id to be called")
	}

	if fixture.repository.FindLatestByTargetCalls != 1 {
		t.Fatalf("expected latest score lookup")
	}

	if result != score {
		t.Fatalf("expected latest score")
	}
}

func TestLatestScoreReturnsScorableError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.scorableRepository.FindByIDErr = errTest

	_, err := fixture.application.LatestScore(uuid.New(), domain_scores.TrustType)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected scorable error, got %v", err)
	}

	if fixture.repository.FindLatestByTargetCalls != 0 {
		t.Fatalf("expected score repository not to be called")
	}
}

func TestLatestScoreReturnsRepositoryError(t *testing.T) {
	fixture := newApplicationFixture()

	id := uuid.New()
	target := scorables.NewMockScorable(id, scorables.UserKind)

	fixture.scorableRepository.Items[id] = target
	fixture.repository.FindLatestByTargetErr = errTest

	_, err := fixture.application.LatestScore(id, domain_scores.TrustType)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected repository error, got %v", err)
	}
}

func TestScoreHistory(t *testing.T) {
	fixture := newApplicationFixture()

	id := uuid.New()
	target := scorables.NewMockScorable(id, scorables.UserKind)

	score := domain_scores.NewMockScore(target, domain_scores.TrustType, 0.9)

	fixture.scorableRepository.Items[id] = target
	fixture.repository.FindHistoryByTargetValue = []domain_scores.Score{
		score,
	}

	result, err := fixture.application.ScoreHistory(id, domain_scores.TrustType)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.scorableRepository.FindByIDCalls != 1 {
		t.Fatalf("expected scorable find by id")
	}

	if fixture.repository.FindHistoryByTargetCalls != 1 {
		t.Fatalf("expected history lookup")
	}

	if len(result) != 1 || result[0] != score {
		t.Fatalf("expected score history")
	}
}

func TestScoreHistoryReturnsScorableError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.scorableRepository.FindByIDErr = errTest

	_, err := fixture.application.ScoreHistory(uuid.New(), domain_scores.TrustType)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected scorable error, got %v", err)
	}

	if fixture.repository.FindHistoryByTargetCalls != 0 {
		t.Fatalf("expected history repository not to be called")
	}
}

func TestScoreHistoryReturnsRepositoryError(t *testing.T) {
	fixture := newApplicationFixture()

	id := uuid.New()
	target := scorables.NewMockScorable(id, scorables.UserKind)

	fixture.scorableRepository.Items[id] = target
	fixture.repository.FindHistoryByTargetErr = errTest

	_, err := fixture.application.ScoreHistory(id, domain_scores.TrustType)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected repository error, got %v", err)
	}
}

func TestRecalculateScores(t *testing.T) {
	fixture := newApplicationFixture()

	first := scorables.NewMockScorable(uuid.New(), scorables.UserKind)
	second := scorables.NewMockScorable(uuid.New(), scorables.PostKind)

	fixture.scorableRepository.FindAllValue = []scorables.Scorable{
		first,
		second,
	}

	fixture.trustCalculator.CalculateValue = domain_scores.NewMockScore(
		first,
		domain_scores.TrustType,
		0.9,
	)
	fixture.spamCalculator.CalculateValue = domain_scores.NewMockScore(
		first,
		domain_scores.SpamType,
		0.1,
	)

	err := fixture.application.RecalculateScores()

	if err != nil {
		t.Fatal(err)
	}

	if fixture.scorableRepository.FindAllCalls != 1 {
		t.Fatalf("expected scorable find all")
	}

	if fixture.trustCalculator.CalculateCalls != 2 {
		t.Fatalf("expected trust calculator called twice")
	}

	if fixture.spamCalculator.CalculateCalls != 2 {
		t.Fatalf("expected spam calculator called twice")
	}

	if fixture.repository.SaveCalls != 4 {
		t.Fatalf("expected 4 saved scores, got %d", fixture.repository.SaveCalls)
	}
}

func TestRecalculateScoresReturnsScorableRepositoryError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.scorableRepository.FindAllErr = errTest

	err := fixture.application.RecalculateScores()

	if !errors.Is(err, errTest) {
		t.Fatalf("expected scorable repository error, got %v", err)
	}

	if fixture.trustCalculator.CalculateCalls != 0 {
		t.Fatalf("expected calculator not to be called")
	}
}

func TestRecalculateScoresReturnsCalculateError(t *testing.T) {
	fixture := newApplicationFixture()

	fixture.scorableRepository.FindAllValue = []scorables.Scorable{
		scorables.NewMockScorable(uuid.New(), scorables.UserKind),
	}

	fixture.trustCalculator.CalculateErr = errTest

	err := fixture.application.RecalculateScores()

	if !errors.Is(err, errTest) {
		t.Fatalf("expected calculate error, got %v", err)
	}
}
