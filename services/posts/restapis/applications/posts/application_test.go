package posts

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	domain_posts "github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/users"
)

var errTest = errors.New("test error")

func TestSave(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	post := domain_posts.NewMockPost("hello")

	err := fixture.application.Save(ctx, post)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.service.SaveCalls != 1 {
		t.Fatalf("expected 1 save call, got %d", fixture.service.SaveCalls)
	}

	if fixture.service.LastContext != ctx {
		t.Fatalf("expected context to be passed")
	}

	if fixture.service.LastPost != post {
		t.Fatalf("expected post to be passed")
	}
}

func TestSaveReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	fixture.service.SaveErr = errTest

	err := fixture.application.Save(ctx, domain_posts.NewMockPost("hello"))

	if !errors.Is(err, errTest) {
		t.Fatalf("expected save error, got %v", err)
	}
}

func TestFindByID(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	id := uuid.New()
	post := domain_posts.NewMockPost("hello")

	fixture.repository.Items[id] = post

	result, err := fixture.application.FindByID(ctx, id)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindByIDCalls != 1 {
		t.Fatalf("expected 1 find by id call")
	}

	if fixture.repository.LastContext != ctx {
		t.Fatalf("expected context to be passed")
	}

	if fixture.repository.LastID != id {
		t.Fatalf("expected id to be passed")
	}

	if result != post {
		t.Fatalf("expected post result")
	}
}

func TestFindByIDReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	fixture.repository.FindByIDErr = errTest

	_, err := fixture.application.FindByID(ctx, uuid.New())

	if !errors.Is(err, errTest) {
		t.Fatalf("expected find by id error, got %v", err)
	}
}

func TestFind(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	post := domain_posts.NewMockPost("hello")
	fixture.repository.FindValue = []domain_posts.Post{
		post,
	}

	result, err := fixture.application.Find(ctx, 0, 25)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindCalls != 1 {
		t.Fatalf("expected 1 find call")
	}

	if fixture.repository.LastContext != ctx {
		t.Fatalf("expected context to be passed")
	}

	if fixture.repository.LastIndex != 0 {
		t.Fatalf("expected index 0")
	}

	if fixture.repository.LastAmount != 25 {
		t.Fatalf("expected amount 25")
	}

	if len(result) != 1 || result[0] != post {
		t.Fatalf("expected post result")
	}
}

func TestFindReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	fixture.repository.FindErr = errTest

	_, err := fixture.application.Find(ctx, 0, 25)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected find error, got %v", err)
	}
}

func TestFindAfter(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	cursor := uuid.New()
	post := domain_posts.NewMockPost("hello")

	fixture.repository.FindAfterValue = []domain_posts.Post{
		post,
	}

	result, err := fixture.application.FindAfter(ctx, cursor, 25)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindAfterCalls != 1 {
		t.Fatalf("expected 1 find after call")
	}

	if fixture.repository.LastContext != ctx {
		t.Fatalf("expected context to be passed")
	}

	if fixture.repository.LastCursor != cursor {
		t.Fatalf("expected cursor to be passed")
	}

	if fixture.repository.LastAmount != 25 {
		t.Fatalf("expected amount 25")
	}

	if len(result) != 1 || result[0] != post {
		t.Fatalf("expected post result")
	}
}

func TestFindAfterReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	fixture.repository.FindAfterErr = errTest

	_, err := fixture.application.FindAfter(ctx, uuid.New(), 25)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected find after error, got %v", err)
	}
}

func TestCount(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	fixture.repository.CountValue = 123

	result, err := fixture.application.Count(ctx)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.CountCalls != 1 {
		t.Fatalf("expected 1 count call")
	}

	if fixture.repository.LastContext != ctx {
		t.Fatalf("expected context to be passed")
	}

	if result != 123 {
		t.Fatalf("expected count 123, got %d", result)
	}
}

func TestCountReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	fixture.repository.CountErr = errTest

	_, err := fixture.application.Count(ctx)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected count error, got %v", err)
	}
}

func TestFindByCriteria(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	user := users.NewMockUser("@user", "User")
	post := domain_posts.NewMockPostWithUser("hello", user)

	criteria := domain_posts.Criteria{
		UserIDs: []uuid.UUID{user.Identifier()},
	}

	fixture.repository.FindByCriteriaValue = []domain_posts.Post{
		post,
	}

	result, err := fixture.application.FindByCriteria(ctx, criteria, 0, 25)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindByCriteriaCalls != 1 {
		t.Fatalf("expected 1 find by criteria call")
	}

	if fixture.repository.LastContext != ctx {
		t.Fatalf("expected context to be passed")
	}

	if fixture.repository.LastCriteria.UserIDs[0] != user.Identifier() {
		t.Fatalf("expected criteria to be passed")
	}

	if fixture.repository.LastIndex != 0 {
		t.Fatalf("expected index 0")
	}

	if fixture.repository.LastAmount != 25 {
		t.Fatalf("expected amount 25")
	}

	if len(result) != 1 || result[0] != post {
		t.Fatalf("expected post result")
	}
}

func TestFindByCriteriaReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	fixture.repository.FindByCriteriaErr = errTest

	_, err := fixture.application.FindByCriteria(
		ctx,
		domain_posts.Criteria{
			UserIDs: []uuid.UUID{uuid.New()},
		},
		0,
		25,
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected find by criteria error, got %v", err)
	}
}

func TestFindByCriteriaAfter(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	cursor := uuid.New()
	communityID := uuid.New()
	post := domain_posts.NewMockPostWithCommunities(
		"hello",
		[]uuid.UUID{communityID},
	)

	criteria := domain_posts.Criteria{
		CommunityIDs: []uuid.UUID{communityID},
	}

	fixture.repository.FindByCriteriaAfterValue = []domain_posts.Post{
		post,
	}

	result, err := fixture.application.FindByCriteriaAfter(
		ctx,
		criteria,
		cursor,
		25,
	)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindByCriteriaAfterCalls != 1 {
		t.Fatalf("expected 1 find by criteria after call")
	}

	if fixture.repository.LastContext != ctx {
		t.Fatalf("expected context to be passed")
	}

	if fixture.repository.LastCriteria.CommunityIDs[0] != communityID {
		t.Fatalf("expected criteria to be passed")
	}

	if fixture.repository.LastCursor != cursor {
		t.Fatalf("expected cursor to be passed")
	}

	if fixture.repository.LastAmount != 25 {
		t.Fatalf("expected amount 25")
	}

	if len(result) != 1 || result[0] != post {
		t.Fatalf("expected post result")
	}
}

func TestFindByCriteriaAfterReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	fixture.repository.FindByCriteriaAfterErr = errTest

	_, err := fixture.application.FindByCriteriaAfter(
		ctx,
		domain_posts.Criteria{
			PlatformIDs: []uuid.UUID{uuid.New()},
		},
		uuid.New(),
		25,
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected find by criteria after error, got %v", err)
	}
}

func TestCountByCriteria(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	platformID := uuid.New()

	criteria := domain_posts.Criteria{
		PlatformIDs: []uuid.UUID{platformID},
	}

	fixture.repository.CountByCriteriaValue = 123

	result, err := fixture.application.CountByCriteria(ctx, criteria)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.CountByCriteriaCalls != 1 {
		t.Fatalf("expected 1 count by criteria call")
	}

	if fixture.repository.LastContext != ctx {
		t.Fatalf("expected context to be passed")
	}

	if fixture.repository.LastCriteria.PlatformIDs[0] != platformID {
		t.Fatalf("expected criteria to be passed")
	}

	if result != 123 {
		t.Fatalf("expected count 123, got %d", result)
	}
}

func TestCountByCriteriaReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	fixture.repository.CountByCriteriaErr = errTest

	_, err := fixture.application.CountByCriteria(
		ctx,
		domain_posts.Criteria{
			UserIDs: []uuid.UUID{uuid.New()},
		},
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected count by criteria error, got %v", err)
	}
}
