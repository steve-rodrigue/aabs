package topics

import (
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/communities"
	domain_participations "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
	domain_topics "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/topics"
	domain_posts "github.com/steve-rodrigue/aabs/services/saas/domain/posts"
	"github.com/steve-rodrigue/aabs/services/saas/domain/users"
)

var errTest = errors.New("test error")

func TestFindByID(t *testing.T) {
	fixture := newApplicationFixture()

	topic := domain_topics.NewMockTopic("Topic", "Description")
	fixture.repository.Items[topic.Identifier()] = topic

	result, err := fixture.application.FindByID(topic.Identifier())

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindByIDCalls != 1 {
		t.Fatalf("expected 1 find by id call")
	}

	if result != topic {
		t.Fatalf("expected topic result")
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

	topic := domain_topics.NewMockTopic("Topic", "Description")
	fixture.repository.FindValue = []domain_topics.Topic{
		topic,
	}

	result, err := fixture.application.Find(0, 25)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindCalls != 1 {
		t.Fatalf("expected 1 find call")
	}

	if len(result) != 1 || result[0] != topic {
		t.Fatalf("expected topic result")
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
	topic := domain_topics.NewMockTopic("Topic", "Description")

	fixture.repository.FindAfterValue = []domain_topics.Topic{
		topic,
	}

	result, err := fixture.application.FindAfter(cursor, 25)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindAfterCalls != 1 {
		t.Fatalf("expected 1 find after call")
	}

	if len(result) != 1 || result[0] != topic {
		t.Fatalf("expected topic result")
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

func TestFindTopicsByUser(t *testing.T) {
	fixture := newApplicationFixture()

	user := users.NewMockUser("@user", "User")
	topic := domain_topics.NewMockTopic("Topic", "Description")

	participation := domain_participations.NewMockParticipationBetween(
		user,
		topic,
	)

	fixture.repository.Items[topic.Identifier()] = topic
	fixture.participations.FindByParticipantValue = []domain_participations.Participation{
		participation,
	}

	result, err := fixture.application.FindTopicsByUser(user)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.participations.FindByParticipantCalls != 1 {
		t.Fatalf("expected 1 find by participant call")
	}

	if fixture.repository.FindByIDCalls != 1 {
		t.Fatalf("expected 1 topic lookup")
	}

	if len(result) != 1 || result[0] != topic {
		t.Fatalf("expected topic result")
	}
}

func TestFindTopicsByCommunity(t *testing.T) {
	fixture := newApplicationFixture()

	community := communities.NewMockCommunity("Community", "Text")
	topic := domain_topics.NewMockTopic("Topic", "Description")

	participation := domain_participations.NewMockParticipationBetween(
		community,
		topic,
	)

	fixture.repository.Items[topic.Identifier()] = topic
	fixture.participations.FindByParticipantValue = []domain_participations.Participation{
		participation,
	}

	result, err := fixture.application.FindTopicsByCommunity(community)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.participations.FindByParticipantCalls != 1 {
		t.Fatalf("expected 1 find by participant call")
	}

	if len(result) != 1 || result[0] != topic {
		t.Fatalf("expected topic result")
	}
}

func TestFindTopicsSkipsNonTopicTargets(t *testing.T) {
	fixture := newApplicationFixture()

	user := users.NewMockUser("@user", "User")
	target := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.CampaignKind,
	)

	participation := domain_participations.NewMockParticipationBetween(
		user,
		target,
	)

	fixture.participations.FindByParticipantValue = []domain_participations.Participation{
		participation,
	}

	result, err := fixture.application.FindTopicsByUser(user)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindByIDCalls != 0 {
		t.Fatalf("expected topic repository not to be called")
	}

	if len(result) != 0 {
		t.Fatalf("expected no topics")
	}
}

func TestFindTopicsReturnsParticipationError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.participations.FindByParticipantErr = errTest

	_, err := fixture.application.FindTopicsByUser(
		users.NewMockUser("@user", "User"),
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected participation error, got %v", err)
	}
}

func TestFindTopicsReturnsTopicLookupError(t *testing.T) {
	fixture := newApplicationFixture()

	user := users.NewMockUser("@user", "User")
	topic := domain_topics.NewMockTopic("Topic", "Description")

	fixture.participations.FindByParticipantValue = []domain_participations.Participation{
		domain_participations.NewMockParticipationBetween(user, topic),
	}

	fixture.repository.FindByIDErr = errTest

	_, err := fixture.application.FindTopicsByUser(user)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected topic lookup error, got %v", err)
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

func TestRebuildTopics(t *testing.T) {
	fixture := newApplicationFixture()

	firstPost := domain_posts.NewMockPost("one")
	secondPost := domain_posts.NewMockPost("two")
	topic := domain_topics.NewMockTopic("Topic", "Description")

	fixture.posts.FindAfterValue = []domain_posts.Post{
		firstPost,
		secondPost,
	}

	fixture.builder.BuildValue = []domain_topics.Topic{
		topic,
	}

	err := fixture.application.RebuildTopics()

	if err != nil {
		t.Fatal(err)
	}

	if fixture.posts.FindAfterCalls != 2 {
		t.Fatalf("expected 2 posts find after calls, got %d", fixture.posts.FindAfterCalls)
	}

	if fixture.builder.BuildCalls != 1 {
		t.Fatalf("expected 1 build call")
	}

	if fixture.repository.SaveCalls != 1 {
		t.Fatalf("expected 1 topic save")
	}
}

func TestRebuildTopicsReturnsPostsError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.posts.FindAfterErr = errTest

	err := fixture.application.RebuildTopics()

	if !errors.Is(err, errTest) {
		t.Fatalf("expected posts error, got %v", err)
	}

	if fixture.builder.BuildCalls != 0 {
		t.Fatalf("expected builder not to be called")
	}
}

func TestRebuildTopicsReturnsBuilderError(t *testing.T) {
	fixture := newApplicationFixture()

	fixture.posts.FindAfterValue = []domain_posts.Post{
		domain_posts.NewMockPost("one"),
	}

	fixture.builder.BuildErr = errTest

	err := fixture.application.RebuildTopics()

	if !errors.Is(err, errTest) {
		t.Fatalf("expected builder error, got %v", err)
	}

	if fixture.repository.SaveCalls != 0 {
		t.Fatalf("expected topic not to be saved")
	}
}

func TestRebuildTopicsReturnsSaveError(t *testing.T) {
	fixture := newApplicationFixture()

	topic := domain_topics.NewMockTopic("Topic", "Description")

	fixture.posts.FindAfterValue = []domain_posts.Post{
		domain_posts.NewMockPost("one"),
	}

	fixture.builder.BuildValue = []domain_topics.Topic{
		topic,
	}

	fixture.repository.SaveErr = errTest

	err := fixture.application.RebuildTopics()

	if !errors.Is(err, errTest) {
		t.Fatalf("expected save error, got %v", err)
	}
}
