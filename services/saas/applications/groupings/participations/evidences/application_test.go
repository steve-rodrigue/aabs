package evidences

import (
	"errors"
	"testing"

	"github.com/google/uuid"

	domain_evidences "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/evidences"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
)

var errTest = errors.New("test error")

func TestFindByID(t *testing.T) {
	fixture := newApplicationFixture()

	id := uuid.New()
	evidence := domain_evidences.NewMockEvidence()

	fixture.repository.FindByIDValue = evidence

	result, err := fixture.application.FindByID(id)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindByIDCalls != 1 {
		t.Fatalf("expected 1 find by id call")
	}

	if result != evidence {
		t.Fatalf("expected evidence result")
	}
}

func TestFindByIDReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.repository.FindByIDErr = errTest

	_, err := fixture.application.FindByID(uuid.New())

	if !errors.Is(err, errTest) {
		t.Fatalf("expected find by id error, got %v", err)
	}
}

func TestFindByParticipation(t *testing.T) {
	fixture := newApplicationFixture()

	participation := uuid.New()
	evidence := domain_evidences.NewMockEvidence()

	fixture.repository.FindByParticipationValue = []domain_evidences.Evidence{
		evidence,
	}

	result, err := fixture.application.FindByParticipation(participation)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindByParticipationCalls != 1 {
		t.Fatalf("expected 1 find by participation call")
	}

	if len(result) != 1 || result[0] != evidence {
		t.Fatalf("expected evidence result")
	}
}

func TestFindByParticipationReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.repository.FindByParticipationErr = errTest

	_, err := fixture.application.FindByParticipation(uuid.New())

	if !errors.Is(err, errTest) {
		t.Fatalf("expected find by participation error, got %v", err)
	}
}

func TestFindByPost(t *testing.T) {
	fixture := newApplicationFixture()

	post := uuid.New()
	evidence := domain_evidences.NewMockEvidence()

	fixture.repository.FindByPostValue = []domain_evidences.Evidence{
		evidence,
	}

	result, err := fixture.application.FindByPost(post)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindByPostCalls != 1 {
		t.Fatalf("expected 1 find by post call")
	}

	if len(result) != 1 || result[0] != evidence {
		t.Fatalf("expected evidence result")
	}
}

func TestFindByPostReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.repository.FindByPostErr = errTest

	_, err := fixture.application.FindByPost(uuid.New())

	if !errors.Is(err, errTest) {
		t.Fatalf("expected find by post error, got %v", err)
	}
}

func TestFindByParticipant(t *testing.T) {
	fixture := newApplicationFixture()

	participant := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.UserKind,
	)

	evidence := domain_evidences.NewMockEvidence()

	fixture.repository.FindByParticipantValue = []domain_evidences.Evidence{
		evidence,
	}

	result, err := fixture.application.FindByParticipant(participant)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindByParticipantCalls != 1 {
		t.Fatalf("expected 1 find by participant call")
	}

	if len(result) != 1 || result[0] != evidence {
		t.Fatalf("expected evidence result")
	}
}

func TestFindByParticipantReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.repository.FindByParticipantErr = errTest

	participant := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.UserKind,
	)

	_, err := fixture.application.FindByParticipant(participant)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected find by participant error, got %v", err)
	}
}

func TestFindByTarget(t *testing.T) {
	fixture := newApplicationFixture()

	target := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.CampaignKind,
	)

	evidence := domain_evidences.NewMockEvidence()

	fixture.repository.FindByTargetValue = []domain_evidences.Evidence{
		evidence,
	}

	result, err := fixture.application.FindByTarget(target)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindByTargetCalls != 1 {
		t.Fatalf("expected 1 find by target call")
	}

	if len(result) != 1 || result[0] != evidence {
		t.Fatalf("expected evidence result")
	}
}

func TestFindByTargetReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.repository.FindByTargetErr = errTest

	target := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.CampaignKind,
	)

	_, err := fixture.application.FindByTarget(target)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected find by target error, got %v", err)
	}
}
