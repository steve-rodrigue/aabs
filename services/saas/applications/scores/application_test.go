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

	target := scorables.NewMockScorable(uuid.New(), scorables.PostKind)

	trustScore := domain_scores.NewMockScore(
		target,
		domain_scores.TrustType,
		0.75,
	)

	spamScore := domain_scores.NewMockScore(
		target,
		domain_scores.SpamType,
		0.25,
	)

	fixture.calculators[0].CalculateValue = trustScore
	fixture.calculators[1].CalculateValue = spamScore

	result, err := fixture.application.Calculate(target)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.calculators[0].CalculateCalls != 1 {
		t.Fatalf("expected 1 trust calculate call")
	}

	if fixture.calculators[1].CalculateCalls != 1 {
		t.Fatalf("expected 1 spam calculate call")
	}

	if fixture.repository.SaveCalls != 2 {
		t.Fatalf("expected 2 save calls")
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 score results")
	}

	if result[0] != trustScore {
		t.Fatalf("expected trust score result")
	}

	if result[1] != spamScore {
		t.Fatalf("expected spam score result")
	}
}

func TestCalculateReturnsCalculatorError(t *testing.T) {
	fixture := newApplicationFixture()

	target := scorables.NewMockScorable(uuid.New(), scorables.PostKind)
	fixture.calculators[0].CalculateErr = errTest

	_, err := fixture.application.Calculate(target)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected calculator error, got %v", err)
	}

	if fixture.repository.SaveCalls != 0 {
		t.Fatalf("expected score not to be saved")
	}
}

func TestCalculateReturnsSaveError(t *testing.T) {
	fixture := newApplicationFixture()

	target := scorables.NewMockScorable(uuid.New(), scorables.PostKind)
	score := domain_scores.NewMockScore(target, domain_scores.TrustType, 0.75)

	fixture.calculators[0].CalculateValue = score
	fixture.repository.SaveErr = errTest

	_, err := fixture.application.Calculate(target)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected save error, got %v", err)
	}
}

func TestLatestScore(t *testing.T) {
	fixture := newApplicationFixture()

	id := uuid.New()
	target := scorables.NewMockScorable(id, scorables.PostKind)
	score := domain_scores.NewMockScore(target, domain_scores.TrustType, 0.75)

	fixture.scorables.Items[id] = target
	fixture.repository.FindLatestByTargetValue = score

	result, err := fixture.application.LatestScore(id, domain_scores.TrustType)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.scorables.FindByIDCalls != 1 {
		t.Fatalf("expected 1 scorable find by id call")
	}

	if fixture.repository.FindLatestByTargetCalls != 1 {
		t.Fatalf("expected 1 latest score lookup")
	}

	if result != score {
		t.Fatalf("expected latest score")
	}
}

func TestLatestScoreReturnsScorableError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.scorables.FindByIDErr = errTest

	_, err := fixture.application.LatestScore(uuid.New(), domain_scores.TrustType)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected scorable error, got %v", err)
	}
}

func TestLatestScoreReturnsRepositoryError(t *testing.T) {
	fixture := newApplicationFixture()

	id := uuid.New()
	target := scorables.NewMockScorable(id, scorables.PostKind)

	fixture.scorables.Items[id] = target
	fixture.repository.FindLatestByTargetErr = errTest

	_, err := fixture.application.LatestScore(id, domain_scores.TrustType)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected latest score error, got %v", err)
	}
}

func TestScoreHistory(t *testing.T) {
	fixture := newApplicationFixture()

	id := uuid.New()
	target := scorables.NewMockScorable(id, scorables.PostKind)
	score := domain_scores.NewMockScore(target, domain_scores.TrustType, 0.75)

	fixture.scorables.Items[id] = target
	fixture.repository.FindHistoryByTargetValue = []domain_scores.Score{
		score,
	}

	result, err := fixture.application.ScoreHistory(id, domain_scores.TrustType)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.scorables.FindByIDCalls != 1 {
		t.Fatalf("expected 1 scorable find by id call")
	}

	if fixture.repository.FindHistoryByTargetCalls != 1 {
		t.Fatalf("expected 1 history lookup")
	}

	if len(result) != 1 || result[0] != score {
		t.Fatalf("expected score history")
	}
}

func TestScoreHistoryReturnsScorableError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.scorables.FindByIDErr = errTest

	_, err := fixture.application.ScoreHistory(uuid.New(), domain_scores.TrustType)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected scorable error, got %v", err)
	}
}

func TestScoreHistoryReturnsRepositoryError(t *testing.T) {
	fixture := newApplicationFixture()

	id := uuid.New()
	target := scorables.NewMockScorable(id, scorables.PostKind)

	fixture.scorables.Items[id] = target
	fixture.repository.FindHistoryByTargetErr = errTest

	_, err := fixture.application.ScoreHistory(id, domain_scores.TrustType)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected history error, got %v", err)
	}
}

func TestRecalculateScores(t *testing.T) {
	fixture := newApplicationFixture()

	target := scorables.NewMockScorable(uuid.New(), scorables.PostKind)
	score := domain_scores.NewMockScore(target, domain_scores.TrustType, 0.75)

	fixture.scorables.FindAfterValue = []scorables.Scorable{
		target,
	}
	fixture.calculators[0].CalculateValue = score

	err := fixture.application.RecalculateScores()

	if err != nil {
		t.Fatal(err)
	}

	if fixture.scorables.FindAfterCalls != 2 {
		t.Fatalf("expected 2 find after calls, got %d", fixture.scorables.FindAfterCalls)
	}

	if fixture.calculators[0].CalculateCalls != 1 {
		t.Fatalf("expected 1 calculate call")
	}

	if fixture.repository.SaveCalls != 2 {
		t.Fatalf("expected 1 save call")
	}
}

func TestRecalculateScoresReturnsFindAfterError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.scorables.FindAfterErr = errTest

	err := fixture.application.RecalculateScores()

	if !errors.Is(err, errTest) {
		t.Fatalf("expected find after error, got %v", err)
	}
}

func TestRecalculateScoresReturnsCalculateError(t *testing.T) {
	fixture := newApplicationFixture()

	target := scorables.NewMockScorable(uuid.New(), scorables.PostKind)

	fixture.scorables.FindAfterValue = []scorables.Scorable{
		target,
	}
	fixture.calculators[0].CalculateErr = errTest

	err := fixture.application.RecalculateScores()

	if !errors.Is(err, errTest) {
		t.Fatalf("expected calculate error, got %v", err)
	}
}
