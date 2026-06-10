package evidences

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/participatables"
	domain_evidences "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/evidences"
)

var errTest = errors.New("test error")

func TestFindByID(t *testing.T) {
	fixture := newApplicationFixture()

	ctx := context.Background()
	id := uuid.New()
	evidence := domain_evidences.NewMockEvidence()

	fixture.repository.FindByIDValue = evidence

	result, err := fixture.application.FindByID(ctx, id)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindByIDCalls != 1 {
		t.Fatalf("expected 1 find by id call")
	}

	if fixture.repository.LastContext != ctx {
		t.Fatalf("expected context")
	}

	if fixture.repository.LastID != id {
		t.Fatalf("expected id")
	}

	if result != evidence {
		t.Fatalf("expected evidence result")
	}
}

func TestFindByIDReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.repository.FindByIDErr = errTest

	_, err := fixture.application.FindByID(
		context.Background(),
		uuid.New(),
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected find by id error, got %v", err)
	}
}

func TestFindByParticipation(t *testing.T) {
	fixture := newApplicationFixture()

	ctx := context.Background()
	participation := uuid.New()
	evidence := domain_evidences.NewMockEvidence()

	fixture.repository.FindByParticipationValue = []domain_evidences.Evidence{
		evidence,
	}

	result, err := fixture.application.FindByParticipation(
		ctx,
		participation,
	)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindByParticipationCalls != 1 {
		t.Fatalf("expected 1 find by participation call")
	}

	if fixture.repository.LastContext != ctx {
		t.Fatalf("expected context")
	}

	if fixture.repository.LastParticipation != participation {
		t.Fatalf("expected participation")
	}

	if len(result) != 1 || result[0] != evidence {
		t.Fatalf("expected evidence result")
	}
}

func TestFindByParticipationReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.repository.FindByParticipationErr = errTest

	_, err := fixture.application.FindByParticipation(
		context.Background(),
		uuid.New(),
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected find by participation error, got %v", err)
	}
}

func TestFindByPost(t *testing.T) {
	fixture := newApplicationFixture()

	ctx := context.Background()
	post := uuid.New()
	evidence := domain_evidences.NewMockEvidence()

	fixture.repository.FindByPostValue = []domain_evidences.Evidence{
		evidence,
	}

	result, err := fixture.application.FindByPost(
		ctx,
		post,
	)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindByPostCalls != 1 {
		t.Fatalf("expected 1 find by post call")
	}

	if fixture.repository.LastContext != ctx {
		t.Fatalf("expected context")
	}

	if fixture.repository.LastPost != post {
		t.Fatalf("expected post")
	}

	if len(result) != 1 || result[0] != evidence {
		t.Fatalf("expected evidence result")
	}
}

func TestFindByPostReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.repository.FindByPostErr = errTest

	_, err := fixture.application.FindByPost(
		context.Background(),
		uuid.New(),
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected find by post error, got %v", err)
	}
}

func TestFindByParticipant(t *testing.T) {
	fixture := newApplicationFixture()

	ctx := context.Background()

	participant := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.UserKind,
	)

	evidence := domain_evidences.NewMockEvidence()

	fixture.repository.FindByParticipantValue = []domain_evidences.Evidence{
		evidence,
	}

	result, err := fixture.application.FindByParticipant(
		ctx,
		participant,
	)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindByParticipantCalls != 1 {
		t.Fatalf("expected 1 find by participant call")
	}

	if fixture.repository.LastContext != ctx {
		t.Fatalf("expected context")
	}

	if fixture.repository.LastParticipant != participant {
		t.Fatalf("expected participant")
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

	_, err := fixture.application.FindByParticipant(
		context.Background(),
		participant,
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected find by participant error, got %v", err)
	}
}

func TestFindByTarget(t *testing.T) {
	fixture := newApplicationFixture()

	ctx := context.Background()

	target := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.CampaignKind,
	)

	evidence := domain_evidences.NewMockEvidence()

	fixture.repository.FindByTargetValue = []domain_evidences.Evidence{
		evidence,
	}

	result, err := fixture.application.FindByTarget(
		ctx,
		target,
	)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindByTargetCalls != 1 {
		t.Fatalf("expected 1 find by target call")
	}

	if fixture.repository.LastContext != ctx {
		t.Fatalf("expected context")
	}

	if fixture.repository.LastTarget != target {
		t.Fatalf("expected target")
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

	_, err := fixture.application.FindByTarget(
		context.Background(),
		target,
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected find by target error, got %v", err)
	}
}
