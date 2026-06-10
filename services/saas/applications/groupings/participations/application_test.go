package participations

import (
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/participatables"
	domain_participations "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations"
	domain_evidences "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/evidences"
)

var errTest = errors.New("test error")

func TestEvidences(t *testing.T) {
	fixture := newApplicationFixture()

	result := fixture.application.Evidences()

	if result != fixture.evidenceApplication {
		t.Fatalf("expected evidence application")
	}
}

func TestFindByID(t *testing.T) {
	fixture := newApplicationFixture()

	id := uuid.New()
	participation := domain_participations.NewMockParticipation()
	fixture.repository.FindByIDValue = participation

	result, err := fixture.application.FindByID(id)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindByIDCalls != 1 {
		t.Fatalf("expected 1 find by id call")
	}

	if result != participation {
		t.Fatalf("expected participation")
	}
}

func TestFindByParticipant(t *testing.T) {
	fixture := newApplicationFixture()

	participant := participatables.NewMockParticipatable(uuid.New(), participatables.UserKind)
	participation := domain_participations.NewMockParticipation()

	fixture.repository.FindByParticipantValue = []domain_participations.Participation{
		participation,
	}

	result, err := fixture.application.FindByParticipant(participant)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindByParticipantCalls != 1 {
		t.Fatalf("expected 1 find by participant call")
	}

	if len(result) != 1 || result[0] != participation {
		t.Fatalf("expected participation result")
	}
}

func TestFindByTarget(t *testing.T) {
	fixture := newApplicationFixture()

	target := participatables.NewMockParticipatable(uuid.New(), participatables.CampaignKind)
	participation := domain_participations.NewMockParticipation()

	fixture.repository.FindByTargetValue = []domain_participations.Participation{
		participation,
	}

	result, err := fixture.application.FindByTarget(target)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindByTargetCalls != 1 {
		t.Fatalf("expected 1 find by target call")
	}

	if len(result) != 1 || result[0] != participation {
		t.Fatalf("expected participation result")
	}
}

func TestFindBetween(t *testing.T) {
	fixture := newApplicationFixture()

	participant := participatables.NewMockParticipatable(uuid.New(), participatables.UserKind)
	target := participatables.NewMockParticipatable(uuid.New(), participatables.CampaignKind)
	participation := domain_participations.NewMockParticipation()

	fixture.repository.FindBetweenValue = participation

	result, err := fixture.application.FindBetween(participant, target)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindBetweenCalls != 1 {
		t.Fatalf("expected 1 find between call")
	}

	if result != participation {
		t.Fatalf("expected participation")
	}
}

func TestRebuildParticipations(t *testing.T) {
	fixture := newApplicationFixture()

	participant := participatables.NewMockParticipatable(uuid.New(), participatables.UserKind)
	target := participatables.NewMockParticipatable(uuid.New(), participatables.CampaignKind)

	participation := domain_participations.NewMockParticipation()
	evidence := domain_evidences.NewMockEvidence()

	fixture.participatableRepository.FindAllParticipantsValue = []participatables.Participatable{
		participant,
	}

	fixture.participatableRepository.FindAllTargetsValue = []participatables.Participatable{
		target,
	}

	fixture.calculator.Value = participation
	fixture.evidenceCalculator.Value = []domain_evidences.Evidence{
		evidence,
	}

	err := fixture.application.RebuildParticipations()

	if err != nil {
		t.Fatal(err)
	}

	if fixture.participatableRepository.FindAllParticipantsCalls != 1 {
		t.Fatalf("expected participants to be loaded")
	}

	if fixture.participatableRepository.FindAllTargetsCalls != 1 {
		t.Fatalf("expected targets to be loaded")
	}

	if fixture.calculator.CalculateCalls != 1 {
		t.Fatalf("expected participation calculator to be called once")
	}

	if fixture.repository.SaveCalls != 1 {
		t.Fatalf("expected participation to be saved once")
	}

	if fixture.evidenceCalculator.CalculateCalls != 1 {
		t.Fatalf("expected evidence calculator to be called once")
	}

	if fixture.evidenceRepository.SaveCalls != 1 {
		t.Fatalf("expected evidence to be saved once")
	}
}

func TestRebuildParticipationsSkipsSameIdentifier(t *testing.T) {
	fixture := newApplicationFixture()

	id := uuid.New()

	participant := participatables.NewMockParticipatable(id, participatables.UserKind)
	target := participatables.NewMockParticipatable(id, participatables.CampaignKind)

	fixture.participatableRepository.FindAllParticipantsValue = []participatables.Participatable{
		participant,
	}

	fixture.participatableRepository.FindAllTargetsValue = []participatables.Participatable{
		target,
	}

	err := fixture.application.RebuildParticipations()

	if err != nil {
		t.Fatal(err)
	}

	if fixture.calculator.CalculateCalls != 0 {
		t.Fatalf("expected calculator not to be called")
	}

	if fixture.repository.SaveCalls != 0 {
		t.Fatalf("expected participation not to be saved")
	}
}

func TestRebuildParticipationsReturnsParticipantsError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.participatableRepository.FindAllParticipantsErr = errTest

	err := fixture.application.RebuildParticipations()

	if !errors.Is(err, errTest) {
		t.Fatalf("expected participants error, got %v", err)
	}

	if fixture.participatableRepository.FindAllTargetsCalls != 0 {
		t.Fatalf("expected targets not to be loaded")
	}
}

func TestRebuildParticipationsReturnsTargetsError(t *testing.T) {
	fixture := newApplicationFixture()

	fixture.participatableRepository.FindAllParticipantsValue = []participatables.Participatable{
		participatables.NewMockParticipatable(uuid.New(), participatables.UserKind),
	}

	fixture.participatableRepository.FindAllTargetsErr = errTest

	err := fixture.application.RebuildParticipations()

	if !errors.Is(err, errTest) {
		t.Fatalf("expected targets error, got %v", err)
	}

	if fixture.calculator.CalculateCalls != 0 {
		t.Fatalf("expected calculator not to be called")
	}
}

func TestRebuildParticipationsReturnsCalculatorError(t *testing.T) {
	fixture := newApplicationFixture()

	fixture.participatableRepository.FindAllParticipantsValue = []participatables.Participatable{
		participatables.NewMockParticipatable(uuid.New(), participatables.UserKind),
	}

	fixture.participatableRepository.FindAllTargetsValue = []participatables.Participatable{
		participatables.NewMockParticipatable(uuid.New(), participatables.CampaignKind),
	}

	fixture.calculator.CalculateErr = errTest

	err := fixture.application.RebuildParticipations()

	if !errors.Is(err, errTest) {
		t.Fatalf("expected calculator error, got %v", err)
	}

	if fixture.repository.SaveCalls != 0 {
		t.Fatalf("expected participation not to be saved")
	}
}

func TestRebuildParticipationsReturnsSaveError(t *testing.T) {
	fixture := newApplicationFixture()

	fixture.participatableRepository.FindAllParticipantsValue = []participatables.Participatable{
		participatables.NewMockParticipatable(uuid.New(), participatables.UserKind),
	}

	fixture.participatableRepository.FindAllTargetsValue = []participatables.Participatable{
		participatables.NewMockParticipatable(uuid.New(), participatables.CampaignKind),
	}

	fixture.calculator.Value = domain_participations.NewMockParticipation()
	fixture.repository.SaveErr = errTest

	err := fixture.application.RebuildParticipations()

	if !errors.Is(err, errTest) {
		t.Fatalf("expected save error, got %v", err)
	}

	if fixture.evidenceCalculator.CalculateCalls != 0 {
		t.Fatalf("expected evidence calculator not to be called")
	}
}

func TestRebuildParticipationsReturnsEvidenceCalculatorError(t *testing.T) {
	fixture := newApplicationFixture()

	fixture.participatableRepository.FindAllParticipantsValue = []participatables.Participatable{
		participatables.NewMockParticipatable(uuid.New(), participatables.UserKind),
	}

	fixture.participatableRepository.FindAllTargetsValue = []participatables.Participatable{
		participatables.NewMockParticipatable(uuid.New(), participatables.CampaignKind),
	}

	fixture.calculator.Value = domain_participations.NewMockParticipation()
	fixture.evidenceCalculator.CalculateErr = errTest

	err := fixture.application.RebuildParticipations()

	if !errors.Is(err, errTest) {
		t.Fatalf("expected evidence calculator error, got %v", err)
	}

	if fixture.evidenceRepository.SaveCalls != 0 {
		t.Fatalf("expected evidence not to be saved")
	}
}

func TestRebuildParticipationsReturnsEvidenceSaveError(t *testing.T) {
	fixture := newApplicationFixture()

	fixture.participatableRepository.FindAllParticipantsValue = []participatables.Participatable{
		participatables.NewMockParticipatable(uuid.New(), participatables.UserKind),
	}

	fixture.participatableRepository.FindAllTargetsValue = []participatables.Participatable{
		participatables.NewMockParticipatable(uuid.New(), participatables.CampaignKind),
	}

	fixture.calculator.Value = domain_participations.NewMockParticipation()
	fixture.evidenceCalculator.Value = []domain_evidences.Evidence{
		domain_evidences.NewMockEvidence(),
	}
	fixture.evidenceRepository.SaveErr = errTest

	err := fixture.application.RebuildParticipations()

	if !errors.Is(err, errTest) {
		t.Fatalf("expected evidence save error, got %v", err)
	}
}
