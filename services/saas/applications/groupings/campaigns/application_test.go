package campaigns

import (
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/communities"
	domain_campaigns "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/campaigns"
	domain_participations "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
	"github.com/steve-rodrigue/aabs/services/saas/domain/platforms"
	domain_posts "github.com/steve-rodrigue/aabs/services/saas/domain/posts"
	"github.com/steve-rodrigue/aabs/services/saas/domain/users"
)

var errTest = errors.New("test error")

func TestFindByID(t *testing.T) {
	fixture := newApplicationFixture()

	id := uuid.New()
	campaign := domain_campaigns.NewMockCampaign("Campaign", "Description")

	fixture.repository.Items[id] = campaign

	result, err := fixture.application.FindByID(id)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindByIDCalls != 1 {
		t.Fatalf("expected 1 find by id call")
	}

	if result != campaign {
		t.Fatalf("expected campaign")
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

func TestFindAll(t *testing.T) {
	fixture := newApplicationFixture()

	campaign := domain_campaigns.NewMockCampaign("Campaign", "Description")
	fixture.repository.Items[campaign.Identifier()] = campaign

	result, err := fixture.application.FindAll()

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindAllCalls != 1 {
		t.Fatalf("expected 1 find all call")
	}

	if len(result) != 1 || result[0] != campaign {
		t.Fatalf("expected campaign result")
	}
}

func TestFindAllReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.repository.FindAllErr = errTest

	_, err := fixture.application.FindAll()

	if !errors.Is(err, errTest) {
		t.Fatalf("expected find all error, got %v", err)
	}
}

func TestFindCampaignsByUser(t *testing.T) {
	fixture := newApplicationFixture()

	user := users.NewMockUser("@user", "User")
	campaign := domain_campaigns.NewMockCampaign("Campaign", "Description")

	participation := domain_participations.NewMockParticipationBetween(
		user,
		campaign,
	)

	fixture.repository.Items[campaign.Identifier()] = campaign
	fixture.participations.FindByParticipantValue = []domain_participations.Participation{
		participation,
	}

	result, err := fixture.application.FindCampaignsByUser(user)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.participations.FindByParticipantCalls != 1 {
		t.Fatalf("expected 1 find by participant call")
	}

	if fixture.repository.FindByIDCalls != 1 {
		t.Fatalf("expected 1 campaign lookup")
	}

	if len(result) != 1 || result[0] != campaign {
		t.Fatalf("expected campaign result")
	}
}

func TestFindCampaignsByCommunity(t *testing.T) {
	fixture := newApplicationFixture()

	community := communities.NewMockCommunity("Community", "Text")
	campaign := domain_campaigns.NewMockCampaign("Campaign", "Description")

	participation := domain_participations.NewMockParticipationBetween(
		community,
		campaign,
	)

	fixture.repository.Items[campaign.Identifier()] = campaign
	fixture.participations.FindByParticipantValue = []domain_participations.Participation{
		participation,
	}

	result, err := fixture.application.FindCampaignsByCommunity(community)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.participations.FindByParticipantCalls != 1 {
		t.Fatalf("expected 1 find by participant call")
	}

	if len(result) != 1 || result[0] != campaign {
		t.Fatalf("expected campaign result")
	}
}

func TestFindCampaignsByPlatform(t *testing.T) {
	fixture := newApplicationFixture()

	platform := platforms.NewMockPlatform("Platform", "platform")
	campaign := domain_campaigns.NewMockCampaign("Campaign", "Description")

	participation := domain_participations.NewMockParticipationBetween(
		platform,
		campaign,
	)

	fixture.repository.Items[campaign.Identifier()] = campaign
	fixture.participations.FindByParticipantValue = []domain_participations.Participation{
		participation,
	}

	result, err := fixture.application.FindCampaignsByPlatform(platform)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.participations.FindByParticipantCalls != 1 {
		t.Fatalf("expected 1 find by participant call")
	}

	if len(result) != 1 || result[0] != campaign {
		t.Fatalf("expected campaign result")
	}
}

func TestFindCampaignsSkipsNonCampaignTargets(t *testing.T) {
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

	result, err := fixture.application.FindCampaignsByUser(user)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindByIDCalls != 0 {
		t.Fatalf("expected campaign repository not to be called")
	}

	if len(result) != 0 {
		t.Fatalf("expected no campaigns")
	}
}

func TestFindCampaignsReturnsParticipationError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.participations.FindByParticipantErr = errTest

	_, err := fixture.application.FindCampaignsByUser(
		users.NewMockUser("@user", "User"),
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected participation error, got %v", err)
	}
}

func TestFindCampaignsReturnsCampaignLookupError(t *testing.T) {
	fixture := newApplicationFixture()

	user := users.NewMockUser("@user", "User")
	campaign := domain_campaigns.NewMockCampaign("Campaign", "Description")

	fixture.participations.FindByParticipantValue = []domain_participations.Participation{
		domain_participations.NewMockParticipationBetween(user, campaign),
	}

	fixture.repository.FindByIDErr = errTest

	_, err := fixture.application.FindCampaignsByUser(user)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected campaign lookup error, got %v", err)
	}
}

func TestRebuildCampaigns(t *testing.T) {
	fixture := newApplicationFixture()

	firstPost := domain_posts.NewMockPost("one")
	secondPost := domain_posts.NewMockPost("two")
	campaign := domain_campaigns.NewMockCampaign("Campaign", "Description")

	fixture.posts.FindAllValue = []domain_posts.Post{
		firstPost,
		secondPost,
	}

	fixture.classifier.ClassifyValue = campaign

	err := fixture.application.RebuildCampaigns()

	if err != nil {
		t.Fatal(err)
	}

	if fixture.posts.FindAllCalls != 1 {
		t.Fatalf("expected posts find all")
	}

	if fixture.classifier.ClassifyCalls != 2 {
		t.Fatalf("expected 2 classify calls, got %d", fixture.classifier.ClassifyCalls)
	}

	if fixture.repository.SaveCalls != 2 {
		t.Fatalf("expected 2 campaign saves, got %d", fixture.repository.SaveCalls)
	}
}

func TestRebuildCampaignsReturnsPostsError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.posts.FindAllErr = errTest

	err := fixture.application.RebuildCampaigns()

	if !errors.Is(err, errTest) {
		t.Fatalf("expected posts error, got %v", err)
	}

	if fixture.classifier.ClassifyCalls != 0 {
		t.Fatalf("expected classifier not to be called")
	}
}

func TestRebuildCampaignsReturnsClassifierError(t *testing.T) {
	fixture := newApplicationFixture()

	fixture.posts.FindAllValue = []domain_posts.Post{
		domain_posts.NewMockPost("one"),
	}

	fixture.classifier.ClassifyErr = errTest

	err := fixture.application.RebuildCampaigns()

	if !errors.Is(err, errTest) {
		t.Fatalf("expected classifier error, got %v", err)
	}

	if fixture.repository.SaveCalls != 0 {
		t.Fatalf("expected campaign not to be saved")
	}
}

func TestRebuildCampaignsReturnsSaveError(t *testing.T) {
	fixture := newApplicationFixture()

	fixture.posts.FindAllValue = []domain_posts.Post{
		domain_posts.NewMockPost("one"),
	}

	fixture.classifier.ClassifyValue = domain_campaigns.NewMockCampaign(
		"Campaign",
		"Description",
	)

	fixture.repository.SaveErr = errTest

	err := fixture.application.RebuildCampaigns()

	if !errors.Is(err, errTest) {
		t.Fatalf("expected save error, got %v", err)
	}
}
