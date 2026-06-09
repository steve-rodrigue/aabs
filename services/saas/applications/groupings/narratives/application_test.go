package narratives

import (
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/communities"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/users"
	domain_narratives "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/narratives"
	domain_participations "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
)

var errTest = errors.New("test error")

func TestFindByID(t *testing.T) {
	fixture := newApplicationFixture()

	narrative := domain_narratives.NewMockNarrative("Narrative", "Description")
	fixture.repository.Items[narrative.Identifier()] = narrative

	result, err := fixture.application.FindByID(narrative.Identifier())

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindByIDCalls != 1 {
		t.Fatalf("expected 1 find by id call")
	}

	if result != narrative {
		t.Fatalf("expected narrative result")
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

func TestFind(t *testing.T) {
	fixture := newApplicationFixture()

	narrative := domain_narratives.NewMockNarrative("Narrative", "Description")
	fixture.repository.FindValue = []domain_narratives.Narrative{
		narrative,
	}

	result, err := fixture.application.Find(0, 25)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindCalls != 1 {
		t.Fatalf("expected 1 find call")
	}

	if len(result) != 1 || result[0] != narrative {
		t.Fatalf("expected narrative result")
	}
}

func TestFindReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.repository.FindErr = errTest

	_, err := fixture.application.Find(0, 25)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected find error, got %v", err)
	}
}

func TestFindAfter(t *testing.T) {
	fixture := newApplicationFixture()

	cursor := uuid.New()
	narrative := domain_narratives.NewMockNarrative("Narrative", "Description")

	fixture.repository.FindAfterValue = []domain_narratives.Narrative{
		narrative,
	}

	result, err := fixture.application.FindAfter(cursor, 25)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindAfterCalls != 1 {
		t.Fatalf("expected 1 find after call")
	}

	if len(result) != 1 || result[0] != narrative {
		t.Fatalf("expected narrative result")
	}
}

func TestFindAfterReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.repository.FindAfterErr = errTest

	_, err := fixture.application.FindAfter(uuid.New(), 25)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected find after error, got %v", err)
	}
}

func TestFindNarrativesByUser(t *testing.T) {
	fixture := newApplicationFixture()

	user := users.NewMockUser("@user", "User")
	narrative := domain_narratives.NewMockNarrative("Narrative", "Description")

	participation := domain_participations.NewMockParticipationBetween(
		user,
		narrative,
	)

	fixture.repository.Items[narrative.Identifier()] = narrative
	fixture.participations.FindByParticipantValue = []domain_participations.Participation{
		participation,
	}

	result, err := fixture.application.FindNarrativesByUser(user)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.participations.FindByParticipantCalls != 1 {
		t.Fatalf("expected 1 find by participant call")
	}

	if fixture.repository.FindByIDCalls != 1 {
		t.Fatalf("expected 1 narrative lookup")
	}

	if len(result) != 1 || result[0] != narrative {
		t.Fatalf("expected narrative result")
	}
}

func TestFindNarrativesByCommunity(t *testing.T) {
	fixture := newApplicationFixture()

	community := communities.NewMockCommunity("Community", "Text")
	narrative := domain_narratives.NewMockNarrative("Narrative", "Description")

	participation := domain_participations.NewMockParticipationBetween(
		community,
		narrative,
	)

	fixture.repository.Items[narrative.Identifier()] = narrative
	fixture.participations.FindByParticipantValue = []domain_participations.Participation{
		participation,
	}

	result, err := fixture.application.FindNarrativesByCommunity(community)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.participations.FindByParticipantCalls != 1 {
		t.Fatalf("expected 1 find by participant call")
	}

	if len(result) != 1 || result[0] != narrative {
		t.Fatalf("expected narrative result")
	}
}

func TestFindNarrativesSkipsNonNarrativeTargets(t *testing.T) {
	fixture := newApplicationFixture()

	user := users.NewMockUser("@user", "User")
	target := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.TopicKind,
	)

	participation := domain_participations.NewMockParticipationBetween(
		user,
		target,
	)

	fixture.participations.FindByParticipantValue = []domain_participations.Participation{
		participation,
	}

	result, err := fixture.application.FindNarrativesByUser(user)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindByIDCalls != 0 {
		t.Fatalf("expected narrative repository not to be called")
	}

	if len(result) != 0 {
		t.Fatalf("expected no narratives")
	}
}

func TestFindNarrativesReturnsParticipationError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.participations.FindByParticipantErr = errTest

	_, err := fixture.application.FindNarrativesByUser(
		users.NewMockUser("@user", "User"),
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected participation error, got %v", err)
	}
}

func TestFindNarrativesReturnsNarrativeLookupError(t *testing.T) {
	fixture := newApplicationFixture()

	user := users.NewMockUser("@user", "User")
	narrative := domain_narratives.NewMockNarrative("Narrative", "Description")

	fixture.participations.FindByParticipantValue = []domain_participations.Participation{
		domain_participations.NewMockParticipationBetween(user, narrative),
	}

	fixture.repository.FindByIDErr = errTest

	_, err := fixture.application.FindNarrativesByUser(user)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected narrative lookup error, got %v", err)
	}
}

func TestCount(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.repository.CountValue = 123

	result, err := fixture.application.Count()

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.CountCalls != 1 {
		t.Fatalf("expected 1 count call")
	}

	if result != 123 {
		t.Fatalf("expected count 123, got %d", result)
	}
}

func TestCountReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.repository.CountErr = errTest

	_, err := fixture.application.Count()

	if !errors.Is(err, errTest) {
		t.Fatalf("expected count error, got %v", err)
	}
}

func TestRebuildNarratives(t *testing.T) {
	fixture := newApplicationFixture()

	err := fixture.application.RebuildNarratives()

	if err != nil {
		t.Fatal(err)
	}
}
